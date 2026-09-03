package visits

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"xorm.io/xorm"

	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/internal/log"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
)

var Default *Accumulator

type counter struct {
	Count         int64
	LastVisitedAt time.Time
}

type Accumulator struct {
	operationMu sync.RWMutex
	flushMu     sync.Mutex
	mu          sync.Mutex
	urls        map[int64]counter
	entities    map[int64]counter
	interval    time.Duration
	maxEntries  int
	dropped     atomic.Int64
	flushSignal chan struct{}
	signalSent  atomic.Bool
}

func Init(cfg *configuration.Config) {
	Default = &Accumulator{
		urls:        make(map[int64]counter),
		entities:    make(map[int64]counter),
		interval:    cfg.VisitFlushInterval,
		maxEntries:  cfg.VisitBufferMaxEntries,
		flushSignal: make(chan struct{}, 1),
	}
}

func (a *Accumulator) Increment(urlID, entityID int64, visitedAt time.Time) {
	if a == nil || urlID == 0 {
		return
	}
	a.operationMu.RLock()
	defer a.operationMu.RUnlock()
	a.mu.Lock()
	defer a.mu.Unlock()

	value, exists := a.urls[urlID]
	if !exists && len(a.urls) >= a.maxEntries {
		a.dropped.Add(1)
		return
	}
	value.Count++
	value.LastVisitedAt = visitedAt
	a.urls[urlID] = value
	if len(a.urls) >= a.maxEntries*8/10 && a.signalSent.CompareAndSwap(false, true) {
		select {
		case a.flushSignal <- struct{}{}:
		default:
		}
	}

	if entityID != 0 {
		entityValue := a.entities[entityID]
		entityValue.Count++
		entityValue.LastVisitedAt = visitedAt
		a.entities[entityID] = entityValue
	}
}

func (a *Accumulator) Run(ctx context.Context) {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.flushAndReport()
		case <-a.flushSignal:
			a.flushAndReport()
		case <-ctx.Done():
			a.flushAndReport()
			return
		}
	}
}

func (a *Accumulator) SynchronizeCleanup(operation func() error) error {
	a.operationMu.Lock()
	defer a.operationMu.Unlock()
	if err := a.flush(); err != nil {
		return err
	}
	return operation()
}

func (a *Accumulator) flushAndReport() {
	if err := a.flush(); err != nil {
		log.Logger.Error("Failed to flush visit statistics", zap.Error(err))
	}
	if dropped := a.dropped.Swap(0); dropped > 0 {
		log.Logger.Warn("Visit statistics were dropped to preserve redirect availability",
			zap.Int64("dropped", dropped),
		)
	}
}

func (a *Accumulator) flush() error {
	a.flushMu.Lock()
	defer a.flushMu.Unlock()
	urls, entities := a.snapshot()
	if len(urls) == 0 && len(entities) == 0 {
		return nil
	}

	session := database.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		a.merge(urls, entities)
		return err
	}
	err := updateCounters(session, "url", urls)
	if err == nil {
		err = updateCounters(session, "entity", entities)
	}
	if err != nil {
		_ = session.Rollback()
		a.merge(urls, entities)
		return err
	}
	if err := session.Commit(); err != nil {
		a.merge(urls, entities)
		return err
	}
	return nil
}

func (a *Accumulator) snapshot() (map[int64]counter, map[int64]counter) {
	a.mu.Lock()
	defer a.mu.Unlock()
	urls, entities := a.urls, a.entities
	a.urls = make(map[int64]counter)
	a.entities = make(map[int64]counter)
	a.signalSent.Store(false)
	return urls, entities
}

func (a *Accumulator) merge(urls, entities map[int64]counter) {
	a.mu.Lock()
	defer a.mu.Unlock()
	mergeCounters(a.urls, urls)
	mergeCounters(a.entities, entities)
}

func mergeCounters(destination, source map[int64]counter) {
	for id, value := range source {
		current := destination[id]
		current.Count += value.Count
		if value.LastVisitedAt.After(current.LastVisitedAt) {
			current.LastVisitedAt = value.LastVisitedAt
		}
		destination[id] = current
	}
}

func updateCounters(session *xorm.Session, table string, counters map[int64]counter) error {
	if len(counters) == 0 {
		return nil
	}
	if configuration.CurrentConfig.DatabaseEngine != "pgx" {
		for id, value := range counters {
			if _, err := session.Exec(
				fmt.Sprintf("UPDATE %s SET visit_count = visit_count + ?, last_visited_at = ? WHERE id = ?", table),
				value.Count, value.LastVisitedAt, id,
			); err != nil {
				return err
			}
		}
		return nil
	}

	const maxRowsPerStatement = 10000
	values := make([]string, 0, min(len(counters), maxRowsPerStatement))
	args := make([]any, 0, min(len(counters), maxRowsPerStatement)*3)
	for id, value := range counters {
		values = append(values, "(CAST(? AS BIGINT), CAST(? AS BIGINT), CAST(? AS TIMESTAMPTZ))")
		args = append(args, id, value.Count, value.LastVisitedAt)
		if len(values) == maxRowsPerStatement {
			if err := executePostgresBatch(session, table, values, args); err != nil {
				return err
			}
			values = values[:0]
			args = args[:0]
		}
	}
	return executePostgresBatch(session, table, values, args)
}

func executePostgresBatch(session *xorm.Session, table string, values []string, args []any) error {
	if len(values) == 0 {
		return nil
	}
	query := fmt.Sprintf(`UPDATE %s AS target
SET visit_count = target.visit_count + pending.delta,
    last_visited_at = GREATEST(COALESCE(target.last_visited_at, pending.visited_at), pending.visited_at)
FROM (VALUES %s) AS pending(id, delta, visited_at)
WHERE target.id = pending.id`, table, strings.Join(values, ","))
	queryArgs := make([]any, 0, len(args)+1)
	queryArgs = append(queryArgs, query)
	queryArgs = append(queryArgs, args...)
	_, err := session.Exec(queryArgs...)
	return err
}
