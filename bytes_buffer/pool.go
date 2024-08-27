package bytes_buffer

import (
	"sync"
)

type Pool[T any] interface {
	Get() T
	Put(T)
}

func NewShiftBufferPool(cap int, maxCap int, autoGrow bool) Pool[ShiftBuffer] {
	return &shiftBufferPool{
		sync.Pool{
			New: func() any {
				b := NewShiftBuffer(cap, maxCap, autoGrow)
				return b
			},
		},
	}
}

type shiftBufferPool struct {
	sync.Pool
}

func (p *shiftBufferPool) Get() ShiftBuffer {
	b := p.Pool.Get().(*shiftBuffer)
	b.Clean()
	return b
}

func (p *shiftBufferPool) Put(b ShiftBuffer) {
	p.Pool.Put(b)
}

func NewRingBufferPool(cap int, maxCap int, autoGrow bool) Pool[RingBuffer] {
	return &ringBufferPool{
		sync.Pool{
			New: func() any {
				b := NewRingBuffer(cap, maxCap, autoGrow)
				return b
			},
		},
	}
}

type ringBufferPool struct {
	sync.Pool
}

func (p *ringBufferPool) Get() RingBuffer {
	b := p.Pool.Get().(*ringBuffer)
	b.Clean()
	return b
}

func (p *ringBufferPool) Put(b RingBuffer) {
	p.Pool.Put(b)
}
