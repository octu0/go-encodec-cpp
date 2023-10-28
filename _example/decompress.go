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
		in        = flag.String("in", "/path/to/in.ecd", "input path")
		out       = flag.String("out", "/path/to/out.wav", "output path")
		gpuLayers = flag.Int("n-gpu-layers", 0, "number of GPU layers to use during computation")
	)
	flag.Parse()

	ectx, err := encodec.LoadModel(*model, *gpuLayers)
	if err != nil {
		panic(err)
	}

	r, err := os.Open(*in)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	ecd := ECDLayout{}
	if err := gob.NewDecoder(r).Decode(&ecd); err != nil {
		panic(err)
	}

	ectx.SetTargetBandwidth(ecd.Bandwidth)

	data, free, err := ectx.DecompressAudio(ecd.Data, runtime.NumCPU())
	if err != nil {
		panic(err)
	}
	defer free()

	w, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	enc := wav.NewEncoder(w, ecd.SampleRate, 16, ecd.NumChannels, ecd.AudioFormat)
	defer enc.Close()

	for _, d := range data {
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
		enc.WriteFrame(int16(d))
	}
	println("write to", *out)
}
