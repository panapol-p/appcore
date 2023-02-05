package appcore_utils

import (
	"log"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker interface {
	Use(name string) *gobreaker.CircuitBreaker
}

type Breaker struct {
	goBreaker map[string]*gobreaker.CircuitBreaker
}

func NewGoBreaker() *Breaker {
	return &Breaker{
		goBreaker: make(map[string]*gobreaker.CircuitBreaker),
	}
}

func (b *Breaker) Use(name string) *gobreaker.CircuitBreaker {
	if _, ok := b.goBreaker[name]; !ok {
		log.Println("set goBreaker command", name)
		var st gobreaker.Settings
		st.Name = name
		st.Interval = 5 * time.Second
		st.Timeout = 15 * time.Second

		st.ReadyToTrip = func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		}

		b.goBreaker[name] = gobreaker.NewCircuitBreaker(st)
	}
	return b.goBreaker[name]
}
