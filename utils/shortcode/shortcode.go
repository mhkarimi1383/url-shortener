package shortcode

import (
	"math/rand"
	"strconv"
	"time"
)

var random *rand.Rand

func Generate(id int64, timestamp time.Time) string {
	i := timestamp.UTC().UnixNano() + id
	random.Seed(i)
	return randomString(id)
}

func randomString(id int64) string {
	n, _ := strconv.ParseInt(
		strconv.Itoa(
			randInt(
				10000,
			),
		)+strconv.FormatInt(id, 10),
		10,
		64,
	)
	return strconv.FormatInt(n, 36)
}

func randInt(max int) int {
	return random.Intn(max)
}

func init() {
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}
