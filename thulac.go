package thulacgo

/*
#cgo CXXFLAGS: -I./cppthulac -DLOGGING_LEVEL=LL_WARNING -O3 -Wall
#include <stdlib.h>
#include "thulac.h"
*/
import "C"
import (
	"unsafe"
)

type Thulacgo struct {
	lac C.Thulac
}

func NewThulacgo(modelpath string, userpath string, justseg bool, t2s bool, ufilter bool, separator byte) (*Thulacgo) {
	mpath, upath := C.CString(modelpath), C.CString(userpath)
	defer C.free(unsafe.Pointer(mpath))
	defer C.free(unsafe.Pointer(upath))

	_t2s := C.int(0)
	if t2s {
		_t2s = C.int(1)
	}
	_justseg := C.int(0)
	if justseg {
		_justseg = C.int(1)
	}
	_ufilter := C.int(0)
	if ufilter {
		_ufilter = C.int(1)
	}
	sep := C.char(separator)
	lac := C.NewThulac(mpath, upath, _justseg, _t2s, _ufilter, sep)
	return &Thulacgo{
		lac,
	}
}

func (self *Thulacgo) Deinit() {
	C.Deinit(self.lac)
}

func (self *Thulacgo) Seg(text string) string {
	input := C.CString(text)
	defer C.free(unsafe.Pointer(input))
	return C.GoString(C.Seg(self.lac, input))
}

func GoStrings(argc C.int, argv **C.char) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

func (self *Thulacgo) SegToSlice(text string) []string {
	input := C.CString(text)
	defer C.free(unsafe.Pointer(input))

	var output **C.char
	var size C.int
	C.SegToSlice(self.lac, input, &output, &size)
	defer C.free(unsafe.Pointer(output))

	return GoStrings(size, output)
}
