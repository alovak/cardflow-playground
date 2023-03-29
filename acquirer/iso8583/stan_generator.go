package iso8583

import (
	"fmt"
	"sync"
)

type stanGenerator struct {
	mu  sync.Mutex
	num int
}

func NewStanGenerator() *stanGenerator {
	return &stanGenerator{
		num: 1,
	}
}

func (g *stanGenerator) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	defer func() {
		g.num++
		if g.num > 999999 {
			g.num = 1
		}
	}()

	return fmt.Sprintf("%06d", g.num)
}
