package audio

import (
	"encoding/binary"
	"math"
)

// getRMS is a helper function used to calculate the RMS. This is called 
// by TrimSilent which uses RMS to determine whether the sample of audio
// is silent
func GetRMS(sampleSize uint16, audioData []byte) float64 {
	sum := 0.0
	for i := uint16(0); i < uint16(len(audioData)); i += sampleSize {
		val := binary.LittleEndian.Uint64(audioData[i:(i + sampleSize)])
		sum += float64(val * val)
	}

	return math.Sqrt(sum / (float64(len(audioData)/int(sampleSize))))
}