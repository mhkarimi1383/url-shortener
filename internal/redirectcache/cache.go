package redirectcache

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/mhkarimi1383/url-shortener/internal/log"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
)

const (
	keyPrefix           = "{url-shortener}:cache:"
	generationKey       = keyPrefix + "generation"
	mutationKey         = keyPrefix + "mutation"
	mutationLockKey     = keyPrefix + "mutation-lock"
	redirectKeyPrefix   = keyPrefix + "redirect:"
	mutationLease       = 2 * time.Minute
	mutationWait        = 2 * time.Second
	mutationPoll        = 20 * time.Millisecond
	mutationRenewPeriod = mutationLease / 3
)

var (
	Default                *Cache
	ErrMutationUnavailable = errors.New("cache mutation coordination is unavailable")

	readScript = redis.NewScript(`
local generation = redis.call('GET', KEYS[1])
if not generation then generation = '0' end
if redis.call('EXISTS', KEYS[2]) == 1 then
  return {generation, false, 1}
end
local value = redis.call('GET', ARGV[1] .. generation .. ':' .. ARGV[2])
return {generation, value, 0}
`)
	setScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[2]) == 1 then return 0 end
local generation = redis.call('GET', KEYS[1])
if not generation then generation = '0' end
if generation ~= ARGV[1] then return 0 end
redis.call('SET', KEYS[3], ARGV[2], 'PX', ARGV[3])
return 1
`)
	beginMutationScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) ~= ARGV[1] then return false end
local active = redis.call('GET', KEYS[3])
if active and active ~= ARGV[1] then return false end
local generation = redis.call('INCR', KEYS[2])
redis.call('SET', KEYS[3], ARGV[1], 'PX', ARGV[2])
return generation
`)
	renewMutationScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) ~= ARGV[1] then return 0 end
if redis.call('GET', KEYS[2]) ~= ARGV[1] then return 0 end
redis.call('PEXPIRE', KEYS[1], ARGV[2])
redis.call('PEXPIRE', KEYS[2], ARGV[2])
return 1
`)
	finishMutationScript = redis.NewScript(`
redis.call('INCR', KEYS[3])
if redis.call('GET', KEYS[2]) == ARGV[1] then
  redis.call('DEL', KEYS[2])
end
if redis.call('GET', KEYS[1]) == ARGV[1] then
  redis.call('DEL', KEYS[1])
end
return 1
`)
	releaseLockScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) == ARGV[1] then
  redis.call('DEL', KEYS[1])
end
return 1
`)
)

type Entry struct {
	URLID    int64  `json:"url_id"`
	EntityID int64  `json:"entity_id,omitempty"`
	Target   string `json:"target"`
	Missing  bool   `json:"missing,omitempty"`
}

type Loader func() (Entry, bool, error)

type Cache struct {
	client        *redis.Client
	ttl           time.Duration
	loads         singleflight.Group
	disabledUntil atomic.Int64
}

type cacheLookup struct {
	generation string
	entry      Entry
	found      bool
	mutating   bool
}

type resolveResult struct {
	entry Entry
	found bool
}

func Init(cfg *configuration.Config) {
	if cfg.RedisAddress == "" {
		Default = nil
		return
	}

	opts := &redis.Options{
		Addr:         cfg.RedisAddress,
		Username:     cfg.RedisUsername,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDatabase,
		DialTimeout:  cfg.RedisDialTimeout,
		ReadTimeout:  cfg.RedisReadTimeout,
		WriteTimeout: cfg.RedisWriteTimeout,
		MaxRetries:   1,
	}
	if cfg.RedisTLS {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	Default = &Cache{client: redis.NewClient(opts), ttl: cfg.RedisCacheTTL}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RedisDialTimeout)
	defer cancel()
	if err := Default.client.Ping(ctx).Err(); err != nil {
		Default.markUnavailable()
		log.Logger.Warn("Redis is unavailable; redirect cache will fail open", zap.Error(err))
		return
	}
	log.Logger.Info("Redis redirect cache initialized", zap.String("address", cfg.RedisAddress))
}

func (c *Cache) Resolve(ctx context.Context, shortCode string, loader Loader) (Entry, bool, error) {
	if c == nil {
		return loader()
	}
	result, err, _ := c.loads.Do(shortCode, func() (any, error) {
		return c.resolve(ctx, shortCode, loader)
	})
	if err != nil {
		return Entry{}, false, err
	}
	resolved := result.(resolveResult)
	return resolved.entry, resolved.found, nil
}

func (c *Cache) resolve(ctx context.Context, shortCode string, loader Loader) (resolveResult, error) {
	for attempt := 0; attempt < 3; attempt++ {
		if c.unavailable() {
			entry, found, err := loader()
			return resolveResult{entry: entry, found: found}, err
		}

		lookup, err := c.lookup(ctx, shortCode)
		if err != nil {
			c.markUnavailable()
			entry, found, loadErr := loader()
			return resolveResult{entry: entry, found: found}, loadErr
		}
		if lookup.mutating {
			if err := waitForMutation(ctx, attempt); err != nil {
				entry, found, loadErr := loader()
				return resolveResult{entry: entry, found: found}, loadErr
			}
			continue
		}
		if lookup.found {
			return resolveResult{entry: lookup.entry, found: !lookup.entry.Missing}, nil
		}

		entry, found, err := loader()
		if err != nil {
			return resolveResult{}, err
		}
		stored, err := c.storeIfCurrent(ctx, shortCode, lookup.generation, entry, found)
		if err != nil {
			c.markUnavailable()
			return resolveResult{entry: entry, found: found}, nil
		}
		if stored {
			return resolveResult{entry: entry, found: found}, nil
		}
	}

	entry, found, err := loader()
	return resolveResult{entry: entry, found: found}, err
}

// Mutate serializes cache-affecting database mutations across every replica.
// Redirects remain fail-open, but mutations fail closed when Redis coordination
// is unavailable so another replica can never retain a stale redirect.
func (c *Cache) Mutate(ctx context.Context, operation func() error) error {
	if c == nil {
		return operation()
	}
	token, err := randomToken()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMutationUnavailable, err)
	}
	coordinationContext, cancel := context.WithTimeout(ctx, mutationWait)
	defer cancel()
	if err := c.acquireMutation(coordinationContext, token); err != nil {
		return fmt.Errorf("%w: %v", ErrMutationUnavailable, err)
	}

	leaseStop := make(chan struct{})
	leaseDone := make(chan struct{})
	go c.renewMutation(token, leaseStop, leaseDone)
	defer func() {
		close(leaseStop)
		<-leaseDone
		if err := c.finishMutation(token); err != nil {
			log.Logger.Error("Failed to finish distributed cache mutation; cache remains bypassed until the lease expires", zap.Error(err))
		}
	}()
	return operation()
}

func (c *Cache) lookup(ctx context.Context, shortCode string) (cacheLookup, error) {
	result, err := readScript.Run(ctx, c.client, []string{generationKey, mutationKey}, redirectKeyPrefix, shortCode).Slice()
	if err != nil {
		return cacheLookup{}, err
	}
	if len(result) != 3 {
		return cacheLookup{}, errors.New("invalid Redis cache lookup response")
	}
	generation, ok := result[0].(string)
	if !ok {
		return cacheLookup{}, errors.New("invalid Redis cache generation")
	}
	lookup := cacheLookup{generation: generation, mutating: redisInteger(result[2]) == 1}
	raw, ok := result[1].(string)
	if !ok || raw == "" {
		return lookup, nil
	}
	if err := json.Unmarshal([]byte(raw), &lookup.entry); err != nil {
		_ = c.client.Del(ctx, redirectKeyPrefix+generation+":"+shortCode).Err()
		return cacheLookup{}, err
	}
	if !lookup.entry.Missing && (lookup.entry.URLID == 0 || lookup.entry.Target == "") {
		_ = c.client.Del(ctx, redirectKeyPrefix+generation+":"+shortCode).Err()
		return cacheLookup{}, errors.New("invalid Redis redirect entry")
	}
	lookup.found = true
	return lookup, nil
}

func (c *Cache) storeIfCurrent(ctx context.Context, shortCode, generation string, entry Entry, found bool) (bool, error) {
	ttl := c.ttl
	if !found {
		entry = Entry{Missing: true}
		ttl = min(ttl, 30*time.Second)
	}
	raw, err := json.Marshal(entry)
	if err != nil {
		return false, err
	}
	key := redirectKeyPrefix + generation + ":" + shortCode
	result, err := setScript.Run(ctx, c.client, []string{generationKey, mutationKey, key}, generation, raw, ttl.Milliseconds()).Int64()
	return result == 1, err
}

func (c *Cache) acquireMutation(ctx context.Context, token string) error {
	ticker := time.NewTicker(mutationPoll)
	defer ticker.Stop()
	for {
		acquired, err := c.client.SetNX(ctx, mutationLockKey, token, mutationLease).Result()
		if err != nil {
			return err
		}
		if acquired {
			result, err := beginMutationScript.Run(ctx, c.client,
				[]string{mutationLockKey, generationKey, mutationKey},
				token, mutationLease.Milliseconds(),
			).Result()
			if err != nil || result == nil {
				_ = c.releaseLock(token)
				if err != nil {
					return err
				}
				return errors.New("lost distributed mutation lock")
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *Cache) renewMutation(token string, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	ticker := time.NewTicker(mutationRenewPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			result, err := renewMutationScript.Run(ctx, c.client,
				[]string{mutationLockKey, mutationKey},
				token, mutationLease.Milliseconds(),
			).Int64()
			cancel()
			if err != nil || result != 1 {
				log.Logger.Error("Failed to renew distributed cache mutation lease", zap.Error(err))
			}
		}
	}
}

func (c *Cache) finishMutation(token string) error {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, lastErr = finishMutationScript.Run(ctx, c.client, []string{mutationLockKey, mutationKey, generationKey}, token).Result()
		cancel()
		if lastErr == nil {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return lastErr
}

func (c *Cache) releaseLock(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := releaseLockScript.Run(ctx, c.client, []string{mutationLockKey}, token).Result()
	return err
}

func (c *Cache) Close() error {
	if c == nil {
		return nil
	}
	return c.client.Close()
}

func (c *Cache) unavailable() bool {
	return time.Now().UnixNano() < c.disabledUntil.Load()
}

func (c *Cache) markUnavailable() {
	c.disabledUntil.Store(time.Now().Add(time.Second).UnixNano())
}

func randomToken() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func redisInteger(value any) int64 {
	switch typed := value.(type) {
	case int64:
		return typed
	case string:
		parsed, _ := strconv.ParseInt(typed, 10, 64)
		return parsed
	default:
		return 0
	}
}

func waitForMutation(ctx context.Context, attempt int) error {
	delay := time.Duration(attempt+1) * 10 * time.Millisecond
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
