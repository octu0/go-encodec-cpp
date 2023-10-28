package main

import (
	"encoding/gob"
	"flag"
	"os"
	"runtime"

	"github.com/go-audio/wav"
	"github.com/octu0/go-encodec-cpp"
)

type ECDLayout struct {
	BitDepth    int
	NumChannels int
	SampleRate  int
	AudioFormat int
	Bandwidth   int
	Data        []int32
}

func main() {
	var (
		model     = flag.String("model", "/path/to/ggml-model.bin", "model path")
		in        = flag.String("in", "/path/to/in.wav", "input path")
		out       = flag.String("out", "/path/to/out.ecd", "output path")
		gpuLayers = flag.Int("n-gpu-layers", 0, "number of GPU layers to use during computation")
		bandwidth = flag.Int("bandwidth", 6, "target bandwidth")
	)
	flag.Parse()

	ectx, err := encodec.LoadModel(*model, *gpuLayers)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(*in)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := wav.NewDecoder(f)
	buf, err := dec.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	ectx.SetTargetBandwidth(*bandwidth)

	data, free, err := ectx.CompressAudio(buf.AsFloat32Buffer().Data, runtime.NumCPU())
	if err != nil {
		panic(err)
	}
	defer free()

	w, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if err := gob.NewEncoder(w).Encode(ECDLayout{
		BitDepth:    32,
		NumChannels: 1,
		SampleRate:  24000,
		AudioFormat: 1,
		Bandwidth:   *bandwidth,
		Data:        data,
	}); err != nil {
		panic(err)
	}

	println("write to", *out)
}
