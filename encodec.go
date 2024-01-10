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
	"sync/atomic"
	"unsafe"
)

var (
	ErrLoadModel            = errors.New("failed to cencodec_load_model")
	ErrCompressAudio        = errors.New("failed to cencodec_compress_audio")
	ErrFetchCompressAudio   = errors.New("failed to cencodec_get_compress_codes")
	ErrDecompressAudio      = errors.New("failed to cencodec_get_decompress_audio")
	ErrFetchDecompressAudio = errors.New("failed to cencodec_get_decompress_audio")
)

type noCopy struct{}

type EncodecContext struct {
	noCopy noCopy

	ptr unsafe.Pointer
	r   int32
}

func (e *EncodecContext) Release() {
	if atomic.CompareAndSwapInt32(&e.r, 0, 1) {
		runtime.SetFinalizer(e, nil)
		C.cencodec_free(e.ptr)
	}
}

func (e *EncodecContext) SetTargetBandwidth(bandwidth int) {
	C.cencodec_set_target_bandwidth(e.ptr, C.int(bandwidth))
}

func (e *EncodecContext) CompressAudio(audio []float32, nThreads int) ([]int32, error) {
	ret := C.cencodec_compress_audio(
		e.ptr,
		(*C.float)(unsafe.Pointer(&audio[0])),
		C.int(len(audio)),
		C.int(nThreads),
	)
	if ret != C.int(0) {
		return nil, ErrCompressAudio
	}

	data := C.cencodec_get_compress_codes(e.ptr)
	if data == nil {
		return nil, ErrFetchCompressAudio
	}

	tmp := unsafe.Slice((*int32)(unsafe.Pointer(data.codes)), int(data.codes_len))
	codes := make([]int32, int(data.codes_len))
	copy(codes, tmp)
	C.cencodec_compressed_free(data)
	return codes, nil
}

func (e *EncodecContext) DecompressAudio(codes []int32, nThreads int) ([]float32, error) {
	ret := C.cencodec_decompress_audio(
		e.ptr,
		(*C.int32_t)(unsafe.Pointer(&codes[0])),
		C.int(len(codes)),
		C.int(nThreads),
	)
	if ret != C.int(0) {
		return nil, ErrDecompressAudio
	}

	data := C.cencodec_get_decompress_audio(e.ptr)
	if data == nil {
		return nil, ErrFetchDecompressAudio
	}

	tmp := unsafe.Slice((*float32)(unsafe.Pointer(data.audio)), int(data.audio_len))
	audio := make([]float32, int(data.audio_len))
	copy(audio, tmp)
	C.cencodec_decompressed_free(data)
	return audio, nil
}

func LoadModel(path string, nGPULayers int) (*EncodecContext, error) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	ptr := C.cencodec_load_model(cs, C.int(nGPULayers))
	if ptr == nil {
		return nil, ErrLoadModel
	}

	ctx := &EncodecContext{
		noCopy: noCopy{},
		ptr:    ptr,
		r:      0,
	}
	runtime.SetFinalizer(ctx, finalizeEncodeContext)
	return ctx, nil
}

func finalizeEncodeContext(e *EncodecContext) {
	e.Release()
}
