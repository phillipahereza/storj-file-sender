package codegen

import (
	"math/rand"
	"time"
	"strconv"
	"strings"
)

func Make(seed int64) string {
	s := rand.NewSource(time.Now().Unix() + seed)
	r := rand.New(s)

	number := strconv.Itoa(r.Intn(1000))
	animal := animals[r.Intn(len(animals))]
	colour := colours[r.Intn(len(colours))]
	adjective := adjectives[r.Intn(len(adjectives))]

	return strings.ToLower(strings.Replace(adjective +"-"+ colour +"-"+ animal +"-"+ number, " ", "-", -1))
}
