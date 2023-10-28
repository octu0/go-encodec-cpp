package encodec

/*
#cgo CFLAGS: -I${SRCDIR}/c-encodec-cpp -I/usr/local/c-encodec-cpp/include -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -L${SRCDIR}/c-encodec-cpp -L/usr/local/lib -L/usr/lib
#cgo LDFLAGS: -lcencodec -lencodec
#cgo LDFLAGS: -lm -ldl -lstdc++

#include <stdlib.h>
#include "cencodec.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

var (
	ErrLoadModel            = errors.New("failed to cencodec_load_model")
	ErrCompressAudio        = errors.New("failed to cencodec_compress_audio")
	ErrFetchCompressAudio   = errors.New("failed to cencodec_get_compress_codes")
	ErrDecompressAudio      = errors.New("failed to cencodec_get_decompress_audio")
	ErrFetchDecompressAudio = errors.New("failed to cencodec_get_decompress_audio")
)

type ReleaseFunc func()

type noCopy struct{}

type EncodecContext struct {
	noCopy noCopy

	ptr  unsafe.Pointer
	once sync.Once
}

func (e *EncodecContext) Release() {
	e.once.Do(func() {
		runtime.SetFinalizer(e, nil)
		C.cencodec_free(e.ptr)
	})
}

func (e *EncodecContext) SetTargetBandwidth(bandwidth int) {
	C.cencodec_set_target_bandwidth(e.ptr, C.int(bandwidth))
}

func (e *EncodecContext) CompressAudio(audio []float32, nThreads int) ([]int32, ReleaseFunc, error) {
	ret := C.cencodec_compress_audio(
		e.ptr,
		(*C.float)(unsafe.Pointer(&audio[0])),
		C.int(len(audio)),
		C.int(nThreads),
	)
	if ret != C.int(0) {
		return nil, nil, ErrCompressAudio
	}

	data := C.cencodec_get_compress_codes(e.ptr)
	if data == nil {
		return nil, nil, ErrFetchCompressAudio
	}

	free := func() {
		C.cencodec_compressed_free(data)
	}
	once := new(sync.Once)
	codes := unsafe.Slice((*int32)(unsafe.Pointer(data.codes)), int(data.codes_len))
	return codes, func() { once.Do(free) }, nil
}

func (e *EncodecContext) DecompressAudio(codes []int32, nThreads int) ([]float32, ReleaseFunc, error) {
	ret := C.cencodec_decompress_audio(
		e.ptr,
		(*C.int32_t)(unsafe.Pointer(&codes[0])),
		C.int(len(codes)),
		C.int(nThreads),
	)
	if ret != C.int(0) {
		return nil, nil, ErrDecompressAudio
	}

	data := C.cencodec_get_decompress_audio(e.ptr)
	if data == nil {
		return nil, nil, ErrFetchDecompressAudio
	}

	free := func() {
		C.cencodec_decompressed_free(data)
	}
	once := new(sync.Once)
	audio := unsafe.Slice((*float32)(unsafe.Pointer(data.audio)), int(data.audio_len))
	return audio, func() { once.Do(free) }, nil
}

func LoadModel(path string, nGPULayers int) (*EncodecContext, error) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	ptr := unsafe.Pointer(C.cencodec_load_model(cs, C.int(nGPULayers)))
	if ptr == nil {
		return nil, ErrLoadModel
	}

	ctx := &EncodecContext{
		noCopy: noCopy{},
		ptr:    ptr,
		once:   sync.Once{},
	}
	runtime.SetFinalizer(ctx, finalizeEncodeContext)
	return ctx, nil
}

func finalizeEncodeContext(e *EncodecContext) {
	e.Release()
}
