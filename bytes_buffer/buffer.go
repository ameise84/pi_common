package bytes_buffer

import "io"

type base interface {
	IsCanGrow() bool     // 是否可增长
	GetMaxCapacity() int // 获取可扩展到最大容量
	GetCapacity() int    // 获取总空间大小
	GetDataSize() int    // 获取占用的数据空间大小
	GetEmptySize() int   // 获取空余总空间大小
	IsEmpty() bool       // 是否空数据
}

type PeekReader interface {
	Peek() ([]byte, error) //查看,但不取出数据
	PeekLen(outLen int) ([]byte, error)
}

type VirtualReader interface {
	VirtualFetchLen(outLen int) ([]byte, error) // 虚读指定长度的数据
	VirtualFlush()                              // 将虚读刷新为实读
	VirtualReset()                              // 恢复虚读
}

type FetchReader interface {
	Fetch() ([]byte, int, error) //取出数据,但不做数据拷贝,有数据被覆盖风险
	FetchLen(outLen int) ([]byte, error)
}

type CopyReader interface {
	Copy([]byte) (int, error)              // 读出所有数据,如何out长度不足,返回错误
	CopyLen([]byte, int) error             // 读取outLen长度的数据,如果数据不足 outLen 返回 error
	CopyOut() ([]byte, int, error)         // 读出所有数据,返回读出数据长度
	CopyOutLen(outLen int) ([]byte, error) // 读取outLen长度的数据,如果数据不足 outLen 返回 error
}

type PeekRReader interface {
	Peek() (first, end []byte, err error) //查看,但不取出数据
	PeekLen(outLen int) (first, end []byte, err error)
}

type VirtualRReader interface {
	VirtualFetchLen(outLen int) (first, end []byte, err error) // 虚读指定长度的数据
	VirtualFlush()                                             // 将虚读刷新为实读
	VirtualReset()                                             // 恢复虚读
}

type FetchRReader interface {
	Fetch() (first, end []byte, outLen int, err error) //取出数据,但不做数据拷贝,有数据被覆盖风险
	FetchLen(outLen int) (first, end []byte, err error)
}

type Reader interface {
	base
	PeekReader
	VirtualReader
	FetchReader
	CopyReader
	io.Reader //read 读出最大可读长度的数据,返回读出数据长度
}

type RReader interface {
	base
	PeekRReader
	VirtualRReader
	FetchRReader
	CopyReader
	io.Reader
}

type Writer interface {
	base
	io.Writer //写入最长可写入的数据,如果数据没有写入完,将返回错误
	Clean()
	Reserve(int) error
	AssignString(string) error
	AssignBytes([]byte) error
	AssignByte(byte) error
	AppendString(string) error
	AppendBytes([]byte) error // 将p写入 ringBuffer,如果 ringBuffer 不足,返回错误
	AppendByte(byte) error
	AppendSomeBytes([]byte) int //写入最长可写入的数据
}

type UnsafeWriter interface {
	Writer
	GetTailEmptyBytes() ([]byte, int) //与AddLen配合使用,用于需要转写数据的场景但又只有[]byte接口的地方
	AddLen(int)                       //增加数据长度, 与GetTailEmptyBytes配合使用
	ResetLen(int)                     //重置数据长度, 与GetTailEmptyBytes配合使用
}

type ShiftBuffer interface {
	Reader
	UnsafeWriter
}

type RingBuffer interface {
	RReader
	Writer
}
