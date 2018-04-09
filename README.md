# Aurora Golang SDK


## Overview

Aurora is the enterprise end-to-end speech solution. This Golang SDK will allow you to quickly and easily use the Aurora service to integrate voice capabilities into your application.

The SDK is currently in a pre-alpha release phase. Bugs and limited functionality should be expected.

## Installation

**The Recommended Golang version is 1.9+**

The Go SDK currently does not bundle the necessary system headers and binaries to interact with audio hardware in a cross-platform manner. For this reason, before using the SDK, you need to install `PortAudio`. The Go binding we use needs to link the headers from PortAudio, so you'll also need `pkg-config`.

### macOS

```
$ brew install portaudio pkg-config
$ go get -u github.com/auroraapi/aurora-go
```

### Linux

```
$ sudo apt-get install libportaudio-dev pkg-config
$ go get -u github.com/auroraapi/aurora-go
```

This will install `PortAudio` and `pkg-config`. Use `yum` if your distribution uses `RPM`-based packages. If your distribution does not have `PortAudio` in its repository, install [PortAudio via source](http://www.portaudio.com/download.html).


## Basic Usage

First, make sure you have an account with [Aurora](http://dashboard.auroraapi.com) and have created an Application.

### Text to Speech (TTS)

```go
package main

import (
  "github.com/auroraapi/aurora-go"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Create a Text object and query the TTS service
  speech, err := aurora.NewText("Hello world").Speech()
  if err != nil {
    return
  }

  // Play the resulting audio...
  speech.Audio.Play()

  // ...or save it into a file
  speech.Audio.WriteToFile("test.wav")
}
```

### Speech to Text (STT)

#### Convert a WAV file to Speech

```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
  "github.com/auroraapi/aurora-go/audio"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Load a WAV file
  audio, err := audio.NewFileFromFileName("test.wav")
  if err != nil {
    return
  }

  speech := aurora.NewSpeech(audio)
  text, err := speech.Text()
  if err != nil {
    return
  }

  fmt.Printf("Transcribed: %s\n", text.Text)
}
```

#### Convert a previous Text API call to Speech
```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Call the TTS API to convert "Hello world" to speech
  speech, err := aurora.NewText("Hello world").Speech()
  if err != nil {
    return
  }

  // Convert the generated speech back to text
  text, err := speech.Text()
  if err != nil {
    return
  }

  fmt.Println(text.Text) // "hello world"
}
```

#### Listen for a specified amount of time
```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Create listen parameters. You should call this method so that the default
  // values get set. Then override them with whatever you want
  params := aurora.NewListenParams()
  params.Length = 3.0

  // Listen for 3 seconds
  speech, err := aurora.Listen(params)
  if err != nil {
    return
  }

  // Convert the recorded speech to text
  text, err := speech.Text()
  if err != nil {
    return
  }

  fmt.Println(text.Text)
}
```

#### Listen for an unspecified amount of time

Calling this API will start listening and will automatically stop listening after a certain amount of silence (default is 0.5 seconds).
```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Create listen parameters. You should call this method so that the default
  // values get set. Then override them with whatever you want
  params := aurora.NewListenParams()
  params.SilenceLen = 1.0

  // Listen until 1 second of silence
  speech, err := aurora.Listen(params)
  if err != nil {
    return
  }

  // Convert the recorded speech to text
  text, err := speech.Text()
  if err != nil {
    return
  }

  fmt.Println(text.Text)
}
```

#### Continuously listen

Continuously listen and retrieve speech segments. Note: you can do anything with these speech segments, but here we'll convert them to text. Just like the previous example, these segments are demarcated by silence (0.5 seconds by default) and can be changed by setting the `SilenceLen` parameter. Additionally, you can make these segments fixed length (as in the example before the previous) by setting the `Length` parameter.

```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

// this callback is passed to ContinuouslyListen. It is called every time a
// Speech object is available. The return value specifies whether or not we
// should continue to listen (true if so, false otherwise)
func listenCallback(s *aurora.Speech, err error) bool {
  if err != nil {
    // returning false in this function will stop listening
    // and quit ContinuouslyListen
    return false
  }

  // convert detected speech to text
  text, err := s.Text()
  if err != nil {
    return false
  }

  fmt.Println(text.Text)

  // Continue listening
  return true
}

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Create listen parameters. You should call this method so that the default
  // values get set. Then override them with whatever you want
  params := aurora.NewListenParams()

  // Continuously listen and convert to speech (blocks) with default params
  aurora.ContinuouslyListen(params, listenCallback)

  // Reduce the amount of silence in between speech segments
  params.SilenceLen = 0.5
  aurora.ContinuouslyListen(params, listenCallback)

  // Fixed-length speech segments of 3 seconds (overrides SilenceLen parameter)
  params.Length = 3.0
  aurora.ContinuouslyListen(params, listenCallback)
}
```

#### Listen and Transcribe

If you already know that you wanted the recorded speech to be converted to text, you can do it in one step, reducing the amount of code you need to write and also reducing latency. Using the `ListenAndTranscribe` method, the audio that is recorded automatically starts uploading as soon as you call the method and transcription begins. When the audio recording ends, you get back the final transcription.

```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

// this callback is passed to ContinuouslyListenAndTranscribe. It is called
// every time a Text object is available. The return value specifies whether
// or not we should continue to listen (true if so, false otherwise)
func listenCallback(t *aurora.Text, err error) bool {
  if err != nil {
    // returning false in this function will stop listening
    // and quit ContinuouslyListen
    return false
  }

  // Print and continue listening
  fmt.Println(t.Text)
  return true
}

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Create listen parameters. You should call this method so that the default
  // values get set. Then override them with whatever you want
  params := aurora.NewListenParams()

  // Listen and transcribe once
  t, err := aurora.ListenAndTranscribe(params)
  listenCallback(t, err)

  // Continuously listen. while recording, this method also streams the data
  // to the backend. Once recording is finished, a transcript is almost
  // instantly available. The callback here receives an *aurora.Text (as opposed
  // to the *aurora.Speech object in regular ContinuouslyListen).
  aurora.ContinuouslyListenAndTranscribe(params, listenCallback)
}
```

#### Listen and echo example

```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

func listenCallback(t *aurora.Text, err error) bool {
  if err != nil {
    return false
  }

  // Perform STT on the transcribed text
  s, err := t.Speech()
  if err != nil {
    return false
  }

  // Speak and continue listening
  s.Audio.Play()
  return true
}

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  params := aurora.NewListenParams()
  aurora.ContinuouslyListenAndTranscribe(params, listenCallback)
}
```

### Interpret (Language Understanding)

The interpret service allows you to take any Aurora `Text` object and understand the user's intent and extract additional query information. Interpret can only be called on `Text` objects and return `Interpret` objects after completion. To convert a user's speech into and `Interpret` object, it must be converted to text first.

#### Basic example

```go
package main

import (
  "fmt"
  "github.com/auroraapi/aurora-go"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Create a Text object
  text := aurora.NewText("What is the time in Los Angeles?")

  // Call the interpret service
  i, err := text.Interpret()
  if err != nil {
    return
  }

  // Print the detected intent (string) and entities (map[string]string)
  fmt.Printf("Intent:   %s\nEntities: %v\n", i.Intent, i.Entities)

  // This should print:
  // Intent:   time
  // Entities: map[location: los angeles]
}
```

#### User query example

```go
package main

import (
  "bufio"
  "fmt"
  "os"
  "github.com/auroraapi/aurora-go"
)

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  // Read line-by-line from stdin
  r := bufio.NewReader(os.Stdin)
  for {
    t, _ := r.ReadString('\n')

    // Interpret what the user type
    i, err := aurora.NewText(t).Interpret()
    if err != nil {
      break
    }

    // Print out the intent and entities
    fmt.Printf("%s %v\n", i.Intent, i.Entities)
  }
}
```

#### Smart Lamp

This example shows how easy it is to voice-enable a smart lamp. It responds to queries in the form of "turn on the lights" or "turn off the lamp". You define what `object` you're listening for (so that you can ignore queries like "turn on the music").

```go
package main

import (
  "github.com/auroraapi/aurora-go"
)

// handle what the user said
func listenCallback(t *aurora.Text, err error) bool {
  if err != nil {
    return true
  }
  i, err := t.Interpret()
  if err != nil {
    return true
  }

  intent := i.Intent
  object := i.Entities["object"]
  validWords := []string{ "light", "lights", "lamp" }

  for _, word := range validWords {
    if object == word {
      if intent == "turn_on" {
        // turn on the lamp
      } else if intent == "turn_off" {
        // turn off the lamp
      }

      break
    }
  }
  return true
}

func main() {
  // Set your application settings
  aurora.Config.AppID = "YOUR_APP_ID"
  aurora.Config.AppToken = "YOUR_APP_TOKEN"
  aurora.Config.DeviceID = "YOUR_DEVICE_ID"

  params := aurora.NewListenParams()
  aurora.ContinuouslyListenAndTranscribe(params, listenCallback)
}
```
