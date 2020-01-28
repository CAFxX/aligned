package aligned

import (
	"fmt"
	"golang.org/x/sys/cpu"
	"unsafe"
)

const (
	cacheLineSize = unsafe.Sizeof(cpu.CacheLinePad{})
)

func malign(size uintptr) unsafe.Pointer {
	p, _ := _malign(size, cacheLineSize)
	return p
}

// _malign returns an aligned allocation of the specified size
// both size and align must be >0
func _malign(size, align uintptr) (unsafe.Pointer, []byte) {
	s := align * 2
	if size > align {
		s = align * (1 + (size+align-1)/align)
	}
	b := make([]byte, s)
	for i := range b[:uintptr(len(b))-size] {
		p := (unsafe.Pointer)(&b[i])
		if (uintptr)(p)%align == 0 {
			i := uintptr(i)
			return p, b[i : i+size]
		}
	}
	panic(fmt.Sprintf("mailgn(%d, %d): bug", size, align))
}

// Uint64 returns a cacheline-aligned uint64 that is guaranteed not
// to share the cacheline with any other data.
func Uint64() *uint64 {
	return (*uint64)(malign(unsafe.Sizeof(uint64(0))))
}

// Int64 returns a cacheline-aligned int64 that is guaranteed not
// to share the cacheline with any other data.
func Int64() *int64 {
	return (*int64)(malign(unsafe.Sizeof(int64(0))))
}

// Uint32 returns a cacheline-aligned uint32 that is guaranteed not
// to share the cacheline with any other data.
func Uint32() *uint32 {
	return (*uint32)(malign(unsafe.Sizeof(uint32(0))))
}

// Int32 returns a cacheline-aligned int32 that is guaranteed not
// to share the cacheline with any other data.
func Int32() *int32 {
	return (*int32)(malign(unsafe.Sizeof(int32(0))))
}

// Uintptr returns a cacheline-aligned uintptr that is guaranteed not
// to share the cacheline with any other data.
func Uintptr() *uintptr {
	return (*uintptr)(malign(unsafe.Sizeof(uintptr(0))))
}

// Uint128 returns a cacheline-aligned contiguous pair of uint64 that
// is guaranteed not to share the cacheline with any other data.
func Uint128() *[2]uint64 {
	return (*[2]uint64)(malign(unsafe.Sizeof([2]uint64{})))
}

// Int128 returns a cacheline-aligned contiguous pair of int64 that
// is guaranteed not to share the cacheline with any other data.
func Int128() *[2]int64 {
	return (*[2]int64)(malign(unsafe.Sizeof([2]int64{})))
}
