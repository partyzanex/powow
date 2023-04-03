package challenge

import (
	"bytes"
	"crypto/rand"
	"math"
	"time"

	"github.com/partyzanex/powow/pkg/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

type Cache[K comparable, V any] interface {
	Get(key K) (value V, err error)
	Set(key K, value V, ttl time.Duration)
	TTL(key K) (time.Time, error)
	Remove(key K)
}

type Counter interface {
	Add(t time.Time, count uint32)
	Value() uint32
}

type Config struct {
	MinTime      uint32
	MinKeyLength uint32
	MinThreads   uint8
	MinMemory    uint32
	TTL          time.Duration
}

type Provider struct {
	config  *Config
	counter Counter
	cache   Cache[proto.ID, *proto.TaskRequest]
}

func NewProvider(config *Config) *Provider {
	return &Provider{
		config:  config,
		counter: newCounter(time.Minute*10, time.Second),
		cache:   newCache[proto.ID, *proto.TaskRequest](),
	}
}

func (p *Provider) CreateTaskRequest() (*proto.TaskRequest, error) {
	p.counter.Add(time.Now(), 1)

	var id proto.ID

	_, err := rand.Read(id[:])
	if err != nil {
		return nil, errors.Wrap(err, "cannot create token")
	}

	taskRequest := &proto.TaskRequest{
		ID:        id,
		Time:      p.getTime(),
		Memory:    p.getMemory(),
		Threads:   p.getThreads(),
		KeyLength: p.getKeyLength(),
	}

	p.cache.Set(id, taskRequest, p.config.TTL)

	return taskRequest, nil
}

func getOrder(n uint32) uint32 {
	return uint32(math.Floor(math.Log10(float64(n))) + 1)
}

func (p *Provider) getTime() uint32 {
	return p.config.MinTime + getOrder(p.counter.Value())
}

func (p *Provider) getKeyLength() uint32 {
	return p.config.MinKeyLength
}

func (p *Provider) getThreads() uint8 {
	return p.config.MinThreads
}

func (p *Provider) getMemory() uint32 {
	const kib = 1024

	order := getOrder(p.counter.Value())

	return p.config.MinMemory + (order*order*order)*kib
}

func (p *Provider) VerifyTaskResult(id proto.ID, result *proto.TaskResult) error {
	request, err := p.cache.Get(id)
	if err != nil {
		return errors.Wrap(err, "cannot get task by token")
	}

	hash := argon2.IDKey(
		request.ID[:],
		result.Salt,
		request.Time,
		request.Memory,
		request.Threads,
		request.KeyLength,
	)

	if !bytes.Equal(result.Hash, hash) {
		return errors.New("wrong task result")
	}

	return nil
}

func (p *Provider) GetDeadline() time.Time {
	return time.Now().Add(p.config.TTL)
}
