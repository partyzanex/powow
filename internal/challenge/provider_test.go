package challenge

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/partyzanex/powow/pkg/proto"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/argon2"
)

func TestProvider_CreateTaskRequest(t *testing.T) {
	p := NewProvider(&Config{
		MinTime:      2,
		MinKeyLength: 32,
		MinThreads:   2,
		MinMemory:    1024 * 64,
		TTL:          time.Minute,
	})

	for i := 0; i < 1000; i++ {
		req, err := p.CreateTaskRequest()
		require.NoError(t, err)
		require.NotNil(t, req)
	}
}

func TestProvider_VerifyTaskResult(t *testing.T) {
	p := NewProvider(&Config{
		MinTime:      2,
		MinKeyLength: 32,
		MinThreads:   2,
		MinMemory:    1024 * 32,
		TTL:          time.Minute,
	})

	for i := 0; i < 10; i++ {
		req, err := p.CreateTaskRequest()
		require.NoError(t, err)
		require.NotNil(t, req)
	}

	req, err := p.CreateTaskRequest()
	require.NoError(t, err)
	require.NotNil(t, req)

	salt := make([]byte, 24)

	_, err = rand.Read(salt)
	require.NoError(t, err)

	hash := argon2.IDKey(req.ID[:], salt, req.Time, req.Memory, req.Threads, req.KeyLength)

	res := &proto.TaskResult{
		Salt: salt,
		Hash: hash,
	}

	ts := time.Now()
	err = p.VerifyTaskResult(req.ID, res)
	require.NoError(t, err)

	t.Log(time.Since(ts))
}
