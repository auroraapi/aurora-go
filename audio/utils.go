package audio

import (
	"encoding/binary"
	"math"
)

// rms is a helper function used to calculate the RMS. This is called 
// by TrimSilent which uses RMS to determine whether the sample of audio
// is silent
func rms(sampleSize int, audioData []byte) float64 {
	sum := 0.0
	for i := 0; i < len(audioData); i += sampleSize {
		val := binary.LittleEndian.Uint64(audioData[i:(i + sampleSize)])
		sum += float64(val * val)
	}

	return math.Sqrt(sum / (float64(len(audioData)/sampleSize)))
}