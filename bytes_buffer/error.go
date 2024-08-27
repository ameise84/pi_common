package bytes_buffer

import "errors"

var (
	ErrInVirtualRead  = errors.New("buffer is in virtual read")
	ErrVirtualReadLen = errors.New("virtual read length must bigger then 0")
	ErrReserveZero    = errors.New("bytes_buffer reserve cap < 0")
	ErrReserveSmaller = errors.New("bytes_buffer reserve cap is small than using cap")
)
