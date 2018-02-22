package audio

import (
	"encoding/binary"
	"math"

	"github.com/gordonklaus/portaudio"
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

	return math.Sqrt(sum / (float64(len(audioData) / sampleSize)))
}

// isSilent determines whether an audio slice is silent or not
func isSilent(audio []int16) bool {
	var max int16 = audio[0]
	for _, value := range audio {
		if value > max {
			max = value
		}
	}
	return max < SilentThresh
}

// recordResponse is a response emitted on the channel returns from the
// record function.
type recordResponse struct {
	// Data contains the raw data converted from the n-bit sample that
	// is read from PortAudio
	Data []byte
	// Samples contains the raw samples read from PortAudio (do not use --
	// this is overwritten in every message sent on the channel. Its use
	// is purely internal)
	Samples []int16
	// Error if an error occurred
	Error error
}

// record accesses the underlying audio hardware and reads data from it
// based on the given parameters. It returns a channel of slices which is
// closed when the recording session has finished. Each slice is raw WAV
// data, converted from 16-bit samples to an equivalent byte array (each
// sample is explicitly split into two bytes, but is still the same data --
// you can think of data as reinterpret_cast<char*>(int16array)
func record(length float64, silenceLen float64) chan *recordResponse {
	// we'll send the data in 2048-byte sized arrays. We'll allow buffering
	// up to 1000 of these, so up to 2MB of buffering. This is so that the
	// user doesn't cause stuttering audio if they can't consume the data fast
	// enough.
	ch := make(chan *recordResponse, 1000)

	// internal processing channel, works the same way as `ch`
	prch := make(chan *recordResponse, 1000)

	// Goroutine for reading data from portaudio
	go func() {
		// close the channel when we're done
		defer close(prch)

		// initialize the underlying APIs for audio transmission
		portaudio.Initialize()
		defer portaudio.Terminate()

		// we read in 16-bit samples, even though the output data is in bytes
		buf := make([]int16, BufSize)
		stream, err := portaudio.OpenDefaultStream(NumChannels, 0, SampleRate, BufSize, buf)
		if err != nil {
			prch <- &recordResponse{nil, nil, err}
			return
		}

		defer stream.Close()
		defer stream.Stop()
		if err := stream.Start(); err != nil {
			prch <- &recordResponse{nil, nil, err}
			return
		}

		// discard silence at the beginning of the recording. Why waste time with it?
		for {
			stream.Read()
			if !isSilent(buf) {
				prch <- &recordResponse{nil, buf, nil}
				break
			}
		}

		// read data until the specified amount of silence or until the
		// specified amount of length
		dataLen := 0
		silentFor := 0.0
		for {
			if err := stream.Read(); err != nil {
				prch <- &recordResponse{nil, nil, err}
				return
			}

			dataLen += BufSize
			prch <- &recordResponse{nil, buf, nil}

			if isSilent(buf) {
				silentFor += float64(BufSize) / SampleRate
			} else {
				silentFor = 0.0
			}

			if length == 0 && silentFor > silenceLen {
				break
			}

			if length > 0 && dataLen > int(length*SampleRate) {
				break
			}
		}
	}()

	// Goroutine for converting audio and sending to the user
	go func() {
		// close the channel when we're done receiving all audio samples
		defer close(ch)
		for res := range prch {
			// check if port audio threw an error
			if res.Error != nil {
				ch <- res
				return
			}

			// convert each element from a 16-bit value to two 8-bit values that
			// are equivalent in little-endian form.
			d := res.Samples
			res.Data = make([]byte, 2*len(d))
			for i, j := 0, 0; i < len(d); i, j = i+1, j+2 {
				res.Data[j] = byte(d[i] & 0xFF)
				res.Data[j+1] = byte((d[i] >> 8) & 0xFF)
			}

			ch <- res
		}
	}()

	return ch
}
