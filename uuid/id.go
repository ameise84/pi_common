package uuid

import (
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
)

// github.com/rs/xid : 一定程度上避免了时钟回拨,但没有完全解决,和V7类似,带有时序性
// shortuuid:基于v4的改版

func New() string {
	return NewXid()
}

func NewXid() string {
	return xid.New().String()
}

func NewShortV4() string {
	return shortuuid.NewWithEncoder(enc58)
}

func NewShortV7() string {
	s, _ := uuid.NewV7()
	return enc58.Encode(s)
}
