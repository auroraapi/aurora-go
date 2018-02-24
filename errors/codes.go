package errors

// ErrorCode is a code for an error
type ErrorCode string

// Define the various error codes possible
const (
	SpeechNilAudio = "SpeechNilAudio"
)

// errorMessages converts and error code to its corresponding message
var errorMessages = map[ErrorCode]string{
	SpeechNilAudio: "The audio file was nil. In order to convert a Speech object to Text, it must have a valid audio file. Usually, this means you created a Speech object that wasn't created using one of the Listen methods.",
}
