package errors

// ErrorCode is a code for an error
type ErrorCode string

// Define the various error codes possible
const (
	SpeechNilAudio = "SpeechNilAudio"
	WAVCorruptFile = "WAVCorruptFile"
)

// errorMessages converts an error code to its corresponding message
var errorMessages = map[ErrorCode]string{
	SpeechNilAudio: "The audio file was nil. In order to convert a Speech object to Text, it must have a valid audio file. Usually, this means you created a Speech object that wasn't created using one of the Listen methods.",
	WAVCorruptFile: "The wav file was corrupted. The wav file sent in does not have an incorrect header. This could be due to an incorrect file passed in or the header parameters were not fully passed in.",
}
