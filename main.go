package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	// Set the alarm time
	alarmTime := time.Now().Add(10 * time.Second)

	fmt.Println("Alarm set for:", alarmTime.Format("15:04:05"))

	// Wait for the alarm time
	for time.Now().Before(alarmTime) {
		time.Sleep(time.Second)
	}

	// Open the MP3 file
	f, err := os.Open("alarm.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Decode the MP3 file
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}

	// Initialize the audio player
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	// Play the audio
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing
	<-done
}
