package proto

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type Kind uint8

const (
	KindTaskRequest Kind = iota + 1
	KindTaskResult
	KindTaskError
	KindQuote

	enumKindTaskRequest = "TaskRequest"
	enumKindTaskResult  = "TaskResult"
	enumKindTaskError   = "TaskError"
	enumKindQuote       = "Quote"
)

func (k Kind) String() string {
	switch k {
	case KindTaskRequest:
		return enumKindTaskRequest
	case KindTaskResult:
		return enumKindTaskResult
	case KindTaskError:
		return enumKindTaskError
	case KindQuote:
		return enumKindQuote
	default:
		return strconv.Itoa(int(k))
	}
}

type Message struct {
	Kind Kind `json:"kind"`
	Data any  `json:"data"`
}

func (m *Message) GetTaskRequest() *TaskRequest {
	return getDataAs[*TaskRequest](m.Data)
}

func (m *Message) GetTaskResult() *TaskResult {
	return getDataAs[*TaskResult](m.Data)
}

func (m *Message) GetTaskError() *TaskError {
	return getDataAs[*TaskError](m.Data)
}

func (m *Message) GetQuote() *Quote {
	return getDataAs[*Quote](m.Data)
}

func getDataAs[T any](data any) (t T) {
	if data == nil {
		return t
	}

	t, _ = data.(T)

	return t
}

func (m *Message) UnmarshalJSON(src []byte) error {
	raw := new(message)

	err := json.Unmarshal(src, raw)
	if err != nil {
		return err
	}

	var data any

	switch raw.Kind {
	case KindQuote:
		data = new(Quote)
	case KindTaskRequest:
		data = new(TaskRequest)
	case KindTaskResult:
		data = new(TaskResult)
	case KindTaskError:
		data = new(TaskError)
	default:
		return errors.Errorf("unknown message kind: %v", raw.Kind)
	}

	err = json.Unmarshal(raw.Data, data)
	if err != nil {
		return err
	}

	m.Kind = raw.Kind
	m.Data = data

	return nil
}

type message struct {
	Kind Kind            `json:"kind"`
	Data json.RawMessage `json:"data"`
}
