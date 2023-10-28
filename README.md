# `go-encodec-cpp`

[![MIT License](https://img.shields.io/github/license/octu0/go-encodec-cpp)](https://github.com/octu0/go-encodec-cpp/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/go-encodec-cpp?status.svg)](https://godoc.org/github.com/octu0/go-encodec-cpp)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/go-encodec-cpp)](https://goreportcard.com/report/github.com/octu0/go-encodec-cpp)
[![Releases](https://img.shields.io/github/v/release/octu0/go-encodec-cpp)](https://github.com/octu0/go-encodec-cpp/releases)

Go binding for [encodec.cpp](https://github.com/PABannier/encodec.cpp)  

encodec.cpp provides a C++ interface, so I create a bridge called `cencodec` (c-encodec-cpp) that can be handled by cgo to call encodec.cpp from Go.

## Installation

install `libcencodec.(so|dylib)`

```
$ git clone https://github.com/octu0/go-encodec-cpp.git
$ cd go-encodec-cpp/c-encodec-cpp
$ make install
```

## Example

Compress

```go
imort (
	"os"
	"runtime"

	"github.com/go-audio/wav"
	"github.com/octu0/go-encodec-cpp"
)

func main() {
	gpuLayers := 0
	targetBandwidth := 6

	in, _ := os.Open("/path/to/input.wav")
	out, _ := os.Create("/path/to/output.compress")

	wd := wav.NewDecoder(in)
	buf, _ := wd.FullPCMBuffer()

	ectx, _ := encodec.LoadModel("/path/to/ggml-model.bin", gpuLayers)
	defer ectx.Release()
	ectx.SetTargetBandwidth(targetBandwidth)

	compressedData, free, _ := ectx.CompressAudio(buf.AsFloat32Buffer().Data, runtime.NumCPU())
	defer free()

	gob.NewEncoder(w).Encode(ECDC{
		BitDepth:    32,
		NumChannels: 1,
		SampleRate:  24000,
		AudioFormat: 1,
		Bandwidth:   targetBandwidth,
		Data:        compressedData,
	})
}
```

Decompress

```go
imort (
	"os"
	"runtime"

	"github.com/go-audio/wav"
	"github.com/octu0/go-encodec-cpp"
)

func main() {
	gpuLayers := 0

	in, _ := os.Open("/path/to/input.compress")
	out, _ := os.Create("/path/to/output.wav")

	ecdc := ECDC{}
	gob.NewDecoder(in).Decode(&ecdc)

	ectx, _ := encodec.LoadModel("/path/to/ggml-model.bin", gpuLayers)
	defer ectx.Release()
	ectx.SetTargetBandwidth(ecdc.Bandwidth)

	decompressData, free, _ := ectx.DecompressAudio(ecdc.Data, runtime.NumCPU())
	defer free()

	we := wav.NewEncoder(out, ecdc.SampleRate, 16, ecdc.NumChannels, ecdc.AudioFormat)
	for _, d := range Decompress {
		d *= 0x7fff
		switch {
		case 0x7fff < d:
			d = 0x7fff
		case d < -0x8000:
			d = -0x8000
		case 0 < d:
			d += 0.5
		case d < 0:
			d -= -0.5
		}
		we.WriteFrame(int16(d))
	}
}
```

See [_example](https://github.com/octu0/go-encodec-cpp/tree/master/_example) for more details

## License

MIT, see LICENSE file for details.
