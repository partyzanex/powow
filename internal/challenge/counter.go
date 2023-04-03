package challenge

import (
	"sync"
	"time"
)

const div = 1000

type series struct {
	timestamp int64
	value     uint32
}

type counter struct {
	series   []series
	period   time.Duration
	interval time.Duration
	mu       sync.RWMutex
}

func newCounter(period, interval time.Duration) *counter {
	return &counter{
		series:   make([]series, period/time.Second),
		period:   period,
		interval: interval,
		mu:       sync.RWMutex{},
	}
}

func (c *counter) Add(t time.Time, count uint32) {
	m := t.UnixMilli()
	p := c.period.Milliseconds()
	i := (m / div) % (p / div)

	c.mu.Lock()
	defer c.mu.Unlock()

	if m-p > c.series[i].timestamp {
		c.series[i].timestamp = m
		c.series[i].value = 0
	} else {
		c.series[i].value += count
	}
}

func (c *counter) Value() uint32 {
	start := time.Now().UnixMilli() - c.period.Milliseconds()
	count := uint32(0)

	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.series {
		if c.series[i].timestamp >= start {
			count += c.series[i].value
		}
	}

	return count
}
