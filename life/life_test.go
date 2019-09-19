package life

import (
	"strconv"
	"testing"
	"time"
)

func runGame(w, h, nw, nh, iter int) time.Duration {
	img := NewTestCreator()
	g, err := NewGame("random", h, w, nh, nw, iter, img)
	if err != nil {
		panic("")
	}
	start := time.Now()
	g.Start()
	return time.Since(start)
}

func Benchmark_1(b *testing.B) {
	parts := [...]int{1, 2, 4, 8, 500}
	dims := [...]int{500, 250, 125, 62, 1}
	iter := 100

	for i := 0; i < len(parts); i++ {
		t := runGame(dims[i], dims[i], parts[i], parts[i], iter)
		b.Log(strconv.Itoa(dims[i]) + "( * " + strconv.Itoa(parts[i]) + " parts): " + t.String())
	}
}
