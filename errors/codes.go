package errors

// ErrorCode is a code for an error
type ErrorCode string

// Define the various error codes possible
const (
	SpeechNilAudio = "SpeechNilAudio"
	WAVCorruptFile = "WAVCorruptFile"
	AudioFileNilStream = "AudioFileNilStream"
	AudioFileNotWritableStream = "AudioFileNotWritableStream"
)

// errorMessages converts an error code to its corresponding message
var errorMessages = map[ErrorCode]string{
	SpeechNilAudio: "The audio file was nil. In order to convert a Speech object to Text, it must have a valid audio file. Usually, this means you created a Speech object that wasn't created using one of the Listen methods.",
	WAVCorruptFile: "The wav file was corrupted. The wav file sent in does not have an incorrect header. This could be due to an incorrect file passed in or the header parameters were not fully passed in.",
	AudioFileNilStream: "The audio stream was nil. Portaudio encountered an error in creating the file stream which may be due to an incorrectly formatted RIFF header. Check the file to make sure it was not corrupted or incomplete",
	AudioFileNotWritableStream: "The data could not be written into the stream. This may be due to the buffer not being completely written or the stream is not open. Check that the stream is open and that the buffer has been written",
}
