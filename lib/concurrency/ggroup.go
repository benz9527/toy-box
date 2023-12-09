package concurrency

import (
	"context"
	"sync"
)

type GGroup struct {
	wg sync.WaitGroup
}

func (g *GGroup) Start(fn func()) {
	g.wg.Add(1)
	go func() {
		defer func() {
			g.wg.Done()
		}()
		fn()
	}()
}

func (g *GGroup) Wait() {
	g.wg.Wait()
}

func (g *GGroup) StartWithChannel(stopC GenericWaitChannel[struct{}], fn func(stopC GenericWaitChannel[struct{}])) {
	g.Start(func() {
		fn(stopC)
	})
}

func (g *GGroup) StartWithContext(ctx context.Context, fn func(ctx context.Context)) {
	g.Start(func() {
		fn(ctx)
	})
}
