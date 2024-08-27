package bytes_buffer

import (
	"github.com/ameise84/pi_common/math"
	"github.com/ameise84/pi_common/str_conv"
	"io"
)

// buffer 实现了类似于ringBuffer的效果,在使用上,比ringBuffer方便(不存在返回2段的情况),在效率上略微低于ringBuffer(存在拷贝偏移数据)

type shiftBuffer struct {
	bytes       []byte //字节
	cap         int    //容量
	maxCap      int
	vrIndex     int  //虚读位置
	rIndex      int  //已读位置
	wIndex      int  //已写位置
	isInVirtual bool //是否在虚读状态
	autoGrow    bool //是否自动增长
}

func NewShiftBuffer(cap int, maxCap int, autoGrow bool) ShiftBuffer {
	if maxCap != 0 && maxCap < cap {
		maxCap = cap
	}
	if !autoGrow && maxCap == 0 {
		maxCap = cap
	}
	return &shiftBuffer{
		bytes:    make([]byte, cap),
		cap:      cap,
		maxCap:   maxCap,
		autoGrow: autoGrow,
	}
}

func Warp(b []byte) Reader {
	return &shiftBuffer{
		bytes:    b,
		cap:      cap(b),
		wIndex:   len(b),
		autoGrow: false,
	}
}

func (b *shiftBuffer) IsCanGrow() bool {
	return b.autoGrow
}

func (b *shiftBuffer) GetMaxCapacity() int {
	return b.maxCap
}

func (b *shiftBuffer) GetCapacity() int {
	return b.cap
}

func (b *shiftBuffer) GetDataSize() int {
	return b.wIndex - b.rIndex
}

func (b *shiftBuffer) GetEmptySize() int {
	return b.cap - b.GetDataSize()
}

func (b *shiftBuffer) IsEmpty() bool {
	return b.wIndex == b.rIndex
}

func (b *shiftBuffer) Clean() {
	b.rIndex = 0
	b.vrIndex = 0
	b.wIndex = 0
	b.isInVirtual = false
}

func (b *shiftBuffer) Reserve(cap int) error {
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
	if b.wIndex > b.rIndex {
		copy(newBytes, b.bytes[b.rIndex:b.wIndex])
	}
	b.bytes = newBytes
	b.cap = cap
	b.vrIndex -= b.rIndex
	b.wIndex -= b.rIndex
	b.rIndex = 0
	return nil
}

func (b *shiftBuffer) Write(p []byte) (int, error) {
	n := b.AppendSomeBytes(p)
	if n != len(p) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (b *shiftBuffer) AssignString(s string) error {
	b.Clean()
	return b.AppendBytes(str_conv.ToBytes(s))
}

func (b *shiftBuffer) AssignBytes(p []byte) error {
	b.Clean()
	return b.AppendBytes(p)
}

func (b *shiftBuffer) AssignByte(p byte) error {
	b.Clean()
	return b.AppendByte(p)
}

func (b *shiftBuffer) AppendString(s string) error {
	return b.AppendBytes(str_conv.ToBytes(s))
}

func (b *shiftBuffer) AppendBytes(p []byte) error {
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

	if b.wIndex+dataSize > b.cap {
		b.moveHeadZeroRIndex()
	}
	n := copy(b.bytes[b.wIndex:], p)
	b.wIndex += n
	return nil
}

func (b *shiftBuffer) AppendByte(p byte) error {
	return b.AppendBytes([]byte{p})
}

func (b *shiftBuffer) AppendSomeBytes(p []byte) int {
	tailSize := b.cap - b.wIndex
	if b.rIndex == 0 && tailSize == 0 {
		return 0
	}
	dataSize := len(p)
	if tailSize < dataSize && b.GetEmptySize() >= dataSize {
		b.moveHeadZeroRIndex()
	}
	n := copy(b.bytes[b.wIndex:], p)
	b.wIndex += n
	return n
}

func (b *shiftBuffer) Peek() ([]byte, error) {
	if b.wIndex == b.rIndex {
		return nil, io.EOF
	}
	return b.bytes[b.rIndex:b.wIndex], nil
}

func (b *shiftBuffer) PeekLen(outLen int) ([]byte, error) {
	if b.wIndex == b.rIndex {
		return nil, io.EOF
	}
	toIdx := b.rIndex + outLen
	if toIdx < b.wIndex {
		return nil, io.ErrUnexpectedEOF
	}
	return b.bytes[b.rIndex:toIdx], nil
}

func (b *shiftBuffer) VirtualFetchLen(outLen int) (out []byte, err error) {
	if outLen == 0 {
		return nil, ErrVirtualReadLen
	}
	if b.vrIndex == b.wIndex {
		return nil, io.EOF
	}
	vrTo := b.vrIndex + outLen
	if vrTo > b.wIndex {
		return nil, io.ErrUnexpectedEOF
	}
	out = b.bytes[b.vrIndex:vrTo]
	b.vrIndex = vrTo
	b.isInVirtual = true
	return
}

func (b *shiftBuffer) VirtualFlush() {
	if b.isInVirtual {
		if b.vrIndex == b.wIndex {
			b.Clean()
		} else {
			b.rIndex = b.vrIndex
			b.tryMoveHeadZeroRIndex()
			b.isInVirtual = false
		}
	}
}

func (b *shiftBuffer) VirtualReset() {
	b.vrIndex = b.rIndex
	b.isInVirtual = false
}

func (b *shiftBuffer) Fetch() (out []byte, dataSize int, err error) {
	if b.isInVirtual {
		return nil, 0, ErrInVirtualRead
	}
	if b.wIndex == b.rIndex {
		return nil, 0, io.EOF
	}
	dataSize = b.wIndex - b.rIndex
	out = b.bytes[b.rIndex:b.wIndex]
	b.Clean()
	return
}

func (b *shiftBuffer) FetchLen(outLen int) (out []byte, err error) {
	if b.isInVirtual {
		return nil, ErrInVirtualRead
	}
	if b.rIndex == b.wIndex {
		return nil, io.EOF
	}
	rTo := b.rIndex + outLen
	if rTo > b.wIndex {
		return nil, io.ErrUnexpectedEOF
	}
	out = b.bytes[b.rIndex:rTo]
	if rTo == b.wIndex {
		b.Clean()
	} else {
		b.rIndex = rTo
		b.vrIndex = rTo
		b.tryMoveHeadZeroRIndex()
	}
	return
}

func (b *shiftBuffer) Copy(out []byte) (int, error) {
	if len(out) < b.GetDataSize() {
		return 0, io.ErrShortBuffer
	}
	src, dataSize, err := b.Fetch()
	if err != nil {
		return 0, err
	}
	_ = copy(out, src)
	return dataSize, nil
}

func (b *shiftBuffer) CopyLen(out []byte, outLen int) error {
	if len(out) < outLen {
		return io.ErrShortBuffer
	}
	src, err := b.FetchLen(outLen)
	if err != nil {
		return err
	}
	copy(out, src)
	return nil
}

func (b *shiftBuffer) CopyOut() ([]byte, int, error) {
	src, dataSize, err := b.Fetch()
	if err != nil {
		return nil, 0, err
	}
	out := make([]byte, dataSize)
	copy(out, src)
	return out, dataSize, nil
}

func (b *shiftBuffer) CopyOutLen(outLen int) ([]byte, error) {
	src, err := b.FetchLen(outLen)
	if err != nil {
		return nil, err
	}
	out := make([]byte, outLen)
	copy(out, src)
	return out, nil
}

func (b *shiftBuffer) Read(out []byte) (int, error) {
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

func (b *shiftBuffer) GetTailEmptyBytes() ([]byte, int) {
	size := b.cap - b.wIndex
	if size == 0 {
		b.moveHeadZeroRIndex()
		size = b.cap - b.wIndex
	}
	return b.bytes[b.wIndex:], size
}

func (b *shiftBuffer) AddLen(len int) {
	b.wIndex += len
}

func (b *shiftBuffer) ResetLen(len int) {
	b.wIndex = b.rIndex + len
}

func (b *shiftBuffer) tryMoveHeadZeroRIndex() {
	dataSize := b.GetDataSize()
	if dataSize < 256 && b.rIndex*2 > b.cap && (b.wIndex-b.rIndex)*16 <= b.cap {
		b.moveHeadZeroRIndex()
	}
}

func (b *shiftBuffer) moveHeadZeroRIndex() {
	if b.wIndex != b.rIndex && b.rIndex != 0 {
		copy(b.bytes, b.bytes[b.rIndex:b.wIndex])
		b.wIndex -= b.rIndex
		b.vrIndex -= b.rIndex
		b.rIndex = 0
	}
}
