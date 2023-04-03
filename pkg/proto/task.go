package proto

import "encoding/hex"

type ID [32]byte

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

type TaskRequest struct {
	ID        ID     `json:"id"`
	Time      uint32 `json:"time"`
	Memory    uint32 `json:"memory"`
	Threads   uint8  `json:"threads"`
	KeyLength uint32 `json:"key_length"`
}

type TaskResult struct {
	Salt []byte `json:"salt"`
	Hash []byte `json:"hash"`
}

type TaskError struct {
	Details string `json:"details"`
}

func (e *TaskError) Error() string {
	return e.Details
}

func NewTaskError(err error) TaskError {
	return TaskError{
		Details: err.Error(),
	}
}
