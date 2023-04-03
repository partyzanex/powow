package challenge

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCounter_Add(t *testing.T) {
	cnt := newCounter(time.Minute/10, time.Second/1000)
	require.NotNil(t, cnt)

	for i := 0; i < 100; i++ {
		cnt.Add(time.Now(), 1)
	}

	v := cnt.Value()
	assert.True(t, v >= 90 && v <= 100, v)
}

func BenchmarkCounter_Add(b *testing.B) {
	cnt := newCounter(time.Minute, time.Second)
	require.NotNil(b, cnt)

	for i := 0; i < b.N; i++ {
		cnt.Add(time.Now(), 1)
	}
}

func BenchmarkCounter_Value(b *testing.B) {
	cnt := newCounter(time.Minute, time.Second)
	require.NotNil(b, cnt)

	go func() {
		for {
			cnt.Add(time.Now(), 1)
			time.Sleep(time.Millisecond * 10)
		}
	}()

	for i := 0; i < b.N; i++ {
		_ = cnt.Value()
	}
}
