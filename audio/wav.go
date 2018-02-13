package audio

import (
	"github.com/nkansal96/aurora-go/errors"
)

type WAV struct {
	RiffHeader *riffHeader
}

type riffHeader struct {
	ChunkID   [4]byte
	ChunkSize [4]byte
	Format    [4]byte
}

func calculateDataLength(b []byte) ([]byte) {
	i := 4
	for i < len(b) && b[i-4] != 'd' || b[i-3] != 'a' || b[i-2] != 't' || b[i-1] != 'a' {
		i++
	}

	dataLen := len(b) - i
	return dataLen
}

func newRiffHeader(chunkSize uint32) *riffHeader{
	riffHeaderLen := 12
	wav := make([]byte, riffHeaderLen)

	// RIFF header
	wav[0] = 'R'
	wav[1] = 'I'
	wav[2] = 'F'
	wav[3] = 'F'
	// chunk size
	wav[4] = byte((chunkSize >> 0) & 0xFF)
	wav[5] = byte((chunkSize >> 8) & 0xFF)
	wav[6] = byte((chunkSize >> 16) & 0xFF)
	wav[7] = byte((chunkSize >> 24) & 0xFF)
	// Format (WAVE)
	wav[8] = 'W'
	wav[9] = 'A'
	wav[10] = 'V'
	wav[11] = 'E'

	return &riffHeader{ ChunkID: wav[0:3], ChunkSize: wav[4:7], Format: wav[8:11] }
}

func MakeValidWAV(b []byte) ([]byte, *errors.Error) *WAV {
	dataLen := calculateDataLength(b)

	if dataLen <= 0 {
		return nil, errors.FromErrorCodeWithInfo(errors.TTSResponseInvalidAudioData, fmt.Sprintf("Received bytes with length %d and header length of %d", len(b), i))
	}

	headerLen := 44
	chunkSize := (dataLen + headerLen - 8)

	riffHeader := newRiffHeader(chunkSize)
	
	//Currently only creates the header
	return WAV{ RiffHeader: riffHeader }
}
