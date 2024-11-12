package util

import (
	"time"

	"golang.org/x/exp/rand"
)

func Shuffle[T any](items []T) {
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	for i := range items {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}
