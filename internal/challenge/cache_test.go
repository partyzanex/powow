package challenge

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := newCache[[20]byte, string]()

	var key [20]byte

	_, _ = rand.Read(key[:])

	want := hex.EncodeToString(key[:])
	wantTTL := time.Now().Add(time.Millisecond * 100)

	c.Set(key, want, time.Millisecond*100)

	got, err := c.Get(key)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	gotTTL, err := c.TTL(key)
	require.NoError(t, err)
	require.Equal(t, wantTTL.UnixMilli(), gotTTL.UnixMilli())

	time.Sleep(time.Millisecond * 110)

	got, err = c.Get(key)
	require.Error(t, err)
	require.Empty(t, got)

	c.Set(key, want, time.Second)

	got, err = c.Get(key)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	c.Remove(key)

	_, err = c.TTL(key)
	require.Error(t, err)

	got, err = c.Get(key)
	require.Error(t, err)
	require.Empty(t, got)
}
