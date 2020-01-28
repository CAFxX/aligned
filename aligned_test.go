package aligned

import (
	"fmt"
	"math/rand"
	"testing"
	"unsafe"
)

func get(t int) (unsafe.Pointer, string) {
	switch t {
	case 0:
		return unsafe.Pointer(Uint64()), "uint64"
	case 1:
		return unsafe.Pointer(Int64()), "int64"
	case 2:
		return unsafe.Pointer(Uint32()), "uint32"
	case 3:
		return unsafe.Pointer(Int32()), "int32"
	case 4:
		return unsafe.Pointer(Uintptr()), "uintptr"
	case 5:
		return unsafe.Pointer(Uint128()), "[2]uint64"
	case 6:
		return unsafe.Pointer(Int128()), "[2]int64"
	default:
		panic("unknown type")
	}
}

func TestAligned(t *testing.T) {
	for i := 0; i < 1<<20; i++ {
		p, s := get(rand.Intn(7))
		if p == nil {
			t.Fatalf("nil %s", s)
		}

		if unaligned(p, cacheLineSize) {
			t.Fatalf("unaligned %s", s)
		}

		e := (*[cacheLineSize]byte)(p)
		*e = [cacheLineSize]byte{}
	}
}

func TestMalign4(t *testing.T) {
	testMalign(t, 4, cacheLineSize)
}

func TestMalign8(t *testing.T) {
	testMalign(t, 8, cacheLineSize)
}

func TestMalignBigSize(t *testing.T) {
	for size := uintptr(1); size <= 256; size++ {
		for align := uintptr(1); align <= 1024; align *= 2 {
			testMalign(t, size, align)
		}
	}
}

func TestMalignPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_malign(0, 0)
}

func testMalign(t *testing.T, size, align uintptr) {
	t.Run(fmt.Sprintf("%d/%d", size, align), func(t *testing.T) {
		p, b := _malign(size, align)
		if uintptr(len(b)) != size {
			t.Fatalf("wrong length %d", len(b))
		}
		if unaligned(p, align) {
			t.Fatal("unaligned p")
		}
		if unsafe.Pointer(&b[0]) != p {
			t.Fatal("wrong p")
		}
	})
}

func unaligned(p unsafe.Pointer, align uintptr) bool {
	return (uintptr)(p)%align != 0
}
