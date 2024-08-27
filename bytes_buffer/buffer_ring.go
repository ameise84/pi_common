package bytes_buffer

import (
	"github.com/ameise84/pi_common/math"
	"github.com/ameise84/pi_common/str_conv"
	"io"
)

type ringBuffer struct {
	bytes       []byte //字节
	cap         int    //容量
	maxCap      int
	vrIndex     int  //虚读位置
	rIndex      int  //已读位置
	wIndex      int  //已写位置
	hasData     bool //是否有数据
	isInVirtual bool //是否在虚读状态
	autoGrow    bool //是否自动增长
}

func NewRingBuffer(cap int, maxCap int, autoGrow bool) RingBuffer {
	if maxCap != 0 && maxCap < cap {
		maxCap = cap
	}
	if !autoGrow && maxCap == 0 {
		maxCap = cap
	}

	return &ringBuffer{
		bytes:    make([]byte, cap),
		cap:      cap,
		maxCap:   maxCap,
		autoGrow: autoGrow,
	}
}

func (b *ringBuffer) IsCanGrow() bool {
	return b.autoGrow
}

func (b *ringBuffer) GetMaxCapacity() int {
	return b.maxCap
}

func (b *ringBuffer) GetCapacity() int {
	return b.cap
}

func (b *ringBuffer) GetDataSize() int {
	n := b.wIndex - b.rIndex
	if n == 0 {
		if b.hasData {
			return b.cap
		} else {
			return 0
		}
	}
	if n > 0 {
		return n
	}
	return b.cap + n
}

func (b *ringBuffer) getVirtualDataSize() int {
	n := b.wIndex - b.vrIndex
	if n == 0 {
		if b.hasData {
			return b.cap
		} else {
			return 0
		}
	}
	if n > 0 {
		return n
	}
	return b.cap + n
}

func (b *ringBuffer) GetEmptySize() int {
	return b.cap - b.GetDataSize()
}

func (b *ringBuffer) IsEmpty() bool {
	return !b.hasData
}

func (b *ringBuffer) Clean() {
	b.rIndex = 0
	b.vrIndex = 0
	b.wIndex = 0
	b.hasData = false
	b.isInVirtual = false
}

func (b *ringBuffer) Reserve(cap int) error {
	if cap <= 0 {
		return ErrReserveZero
	}
	if b.maxCap != 0 && cap > b.maxCap {
		cap = b.maxCap
	}
	if cap == b.cap {
		return nil
	}
	if cap < b.GetDataSize() {
		return ErrReserveSmaller
	}

	newBytes := make([]byte, cap)

	if b.hasData {
		if b.wIndex > b.rIndex {
			copy(newBytes, b.bytes[b.rIndex:b.wIndex])
			b.vrIndex -= b.rIndex
			b.wIndex -= b.rIndex
			b.rIndex = 0
		} else {
			tailSize := b.cap - b.rIndex
			copy(newBytes, b.bytes[b.rIndex:])
			copy(newBytes[tailSize:], b.bytes[:b.wIndex])

			if b.isInVirtual {
				if b.vrIndex > b.rIndex {
					b.vrIndex = b.vrIndex - b.rIndex
				} else {
					b.vrIndex = tailSize + b.vrIndex
				}
			} else {
				b.vrIndex = 0
			}
			b.wIndex = tailSize + b.wIndex
			b.rIndex = 0
		}
	}
	b.bytes = newBytes
	b.cap = cap
	return nil
}

func (b *ringBuffer) Write(p []byte) (int, error) {
	n := b.AppendSomeBytes(p)
	if n != len(p) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (b *ringBuffer) AssignString(s string) error {
	b.Clean()
	return b.AppendBytes(str_conv.ToBytes(s))
}

func (b *ringBuffer) AssignBytes(p []byte) error {
	b.Clean()
	return b.AppendBytes(p)
}

func (b *ringBuffer) AssignByte(p byte) error {
	return b.AssignBytes([]byte{p})
}

func (b *ringBuffer) AppendString(s string) error {
	return b.AppendBytes(str_conv.ToBytes(s))
}

func (b *ringBuffer) AppendBytes(p []byte) error {
	dataSize := len(p)
	if dataSize == 0 {
		return nil
	}
	usedSize := b.GetDataSize()
	toSize := usedSize + dataSize
	if b.maxCap != 0 && toSize > b.maxCap {
		return io.ErrShortWrite
	}
	if toSize > b.cap {
		if !b.autoGrow || b.cap == b.maxCap {
			return io.ErrShortWrite
		}

		toCap := math.CeilToPowerOfTwo(dataSize + b.cap)
		if err := b.Reserve(toCap); err != nil {
			return err
		}
	}

	n := copy(b.bytes[b.wIndex:], p)
	if n != dataSize {
		k := copy(b.bytes[:], p[n:])
		b.wIndex = k
	} else {
		b.hasData = true
		b.wIndex += n
	}
	return nil
}

func (b *ringBuffer) AppendByte(p byte) error {
	return b.AppendBytes([]byte{p})
}

func (b *ringBuffer) AppendSomeBytes(p []byte) int {
	writeSize := len(p)
	if writeSize == 0 {
		return 0
	}
	usedSize := b.GetDataSize()
	toSize := usedSize + writeSize
	if toSize > b.cap {
		writeSize = b.cap - usedSize
	}
	n := copy(b.bytes[b.wIndex:], p)
	writeSize -= n
	if writeSize > 0 {
		k := copy(b.bytes[:], p[n:])
		b.wIndex = k
	}
	return writeSize
}

func (b *ringBuffer) Peek() (first, end []byte, err error) {
	if !b.hasData {
		return nil, nil, io.EOF
	}
	if b.wIndex > b.rIndex {
		return b.bytes[b.rIndex:b.wIndex], nil, nil
	} else {
		return b.bytes[b.rIndex:], b.bytes[:b.wIndex], nil
	}
}

func (b *ringBuffer) PeekLen(outLen int) (first, end []byte, err error) {
	if !b.hasData {
		return nil, nil, io.EOF
	}
	if b.GetDataSize() < outLen {
		return nil, nil, io.ErrUnexpectedEOF
	}
	rTo := b.rIndex + outLen
	if rTo > b.cap {
		return b.bytes[b.rIndex:], b.bytes[:rTo-b.cap], nil
	} else {
		return b.bytes[b.rIndex:rTo], nil, nil
	}
}

func (b *ringBuffer) VirtualFetchLen(outLen int) (first, end []byte, err error) {
	if outLen == 0 {
		return nil, nil, ErrVirtualReadLen
	}
	if b.vrIndex == b.wIndex && b.isInVirtual {
		return nil, nil, io.EOF
	}
	if b.getVirtualDataSize() < outLen {
		return nil, nil, io.ErrUnexpectedEOF
	}

	vrTo := b.vrIndex + outLen
	if vrTo > b.cap {
		vrTo = vrTo - b.cap
		first = b.bytes[b.vrIndex:]
		end = b.bytes[:vrTo]
	} else {
		first = b.bytes[b.vrIndex:vrTo]
	}
	b.vrIndex = vrTo
	b.isInVirtual = true
	return
}

func (b *ringBuffer) VirtualFlush() {
	if b.isInVirtual {
		if b.vrIndex == b.wIndex {
			b.Clean()
		} else {
			b.rIndex = b.vrIndex
			b.isInVirtual = false
		}
	}
}

func (b *ringBuffer) VirtualReset() {
	b.vrIndex = b.rIndex
	b.isInVirtual = true
}

func (b *ringBuffer) Fetch() (first, end []byte, dataSize int, err error) {
	if b.isInVirtual {
		return nil, nil, 0, ErrInVirtualRead
	}

	if !b.hasData {
		return nil, nil, 0, io.EOF
	}
	dataSize = b.GetDataSize()
	if b.wIndex > b.rIndex {
		first = b.bytes[b.rIndex:b.wIndex]
	} else {
		first = b.bytes[b.rIndex:]
		end = b.bytes[:b.wIndex]
	}
	b.Clean()
	return
}

func (b *ringBuffer) FetchLen(outLen int) (first, end []byte, err error) {
	if b.isInVirtual {
		return nil, nil, ErrInVirtualRead
	}
	if !b.hasData {
		return nil, nil, io.EOF
	}
	if b.GetDataSize() < outLen {
		return nil, nil, io.ErrUnexpectedEOF
	}
	rTo := b.rIndex + outLen
	if rTo > b.cap {
		rTo = rTo - b.cap
		first = b.bytes[b.rIndex:]
		end = b.bytes[:rTo]
	} else {
		first = b.bytes[b.rIndex:rTo]
	}

	if rTo == b.wIndex {
		b.Clean()
	} else {
		b.rIndex = rTo
		b.vrIndex = rTo
	}
	return
}

func (b *ringBuffer) Copy(out []byte) (int, error) {
	if len(out) < b.GetDataSize() {
		return 0, io.ErrShortBuffer
	}
	first, end, dataSize, err := b.Fetch()
	if err != nil {
		return 0, err
	}
	_ = copy(out, first)
	if end != nil {
		_ = copy(out[len(first):], end)
	}
	return dataSize, nil
}

func (b *ringBuffer) CopyLen(out []byte, outLen int) error {
	if len(out) < outLen {
		return io.ErrShortBuffer
	}
	first, end, err := b.FetchLen(outLen)
	if err != nil {
		return err
	}
	copy(out, first)
	if end != nil {
		_ = copy(out[len(first):], end)
	}
	return nil
}

func (b *ringBuffer) CopyOut() ([]byte, int, error) {
	first, end, dataSize, err := b.Fetch()
	if err != nil {
		return nil, 0, err
	}
	out := make([]byte, dataSize)
	copy(out, first)
	if end != nil {
		_ = copy(out[len(first):], end)
	}
	return out, dataSize, nil
}

func (b *ringBuffer) CopyOutLen(outLen int) ([]byte, error) {
	first, end, err := b.FetchLen(outLen)
	if err != nil {
		return nil, err
	}
	out := make([]byte, outLen)
	copy(out, first)
	if end != nil {
		_ = copy(out[len(first):], end)
	}
	return out, nil
}

func (b *ringBuffer) Read(out []byte) (int, error) {
	maxReadSize := len(out)
	if maxReadSize == 0 {
		return 0, io.ErrShortBuffer
	}

	dataSize := b.GetDataSize()
	if dataSize < maxReadSize {
		maxReadSize = dataSize
	}

	err := b.CopyLen(out, maxReadSize)
	return maxReadSize, err
}
