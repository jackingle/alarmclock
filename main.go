package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	filename := "alarm.mp3"
	alarmTime := time.Now().Add(10 * time.Second)

	// Set up snooze channel and goroutine
	snooze := make(chan bool)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(input) == "s" {
				snooze <- true
			}
		}
	}()

	for {
		now := time.Now()
		if now.After(alarmTime) {
			alarm(filename, snooze)
			break
		}
		fmt.Printf("\rCurrent time: %s. Waiting for alarm at %s...", now.Format("15:04:05"), alarmTime.Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}
func alarm(filename string, snooze chan bool) {
	// Load MP3 file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s, format, _ := mp3.Decode(f)
	defer s.Close()

	// Initialize speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Create streamer
	done := make(chan bool)
	streamer := beep.Seq(s, beep.Callback(func() {
		done <- true
	}))

	// Play streamer
	speaker.Play(streamer)

	// Wait for streamer to finish or for snooze to be pressed
	select {
	case <-done:
		fmt.Println("Alarm stopped")
	case <-time.After(5 * time.Minute):
		fmt.Println("Snooze time elapsed")
	case <-snooze:
		fmt.Println("Alarm snoozed")
	}

	// Stop playing
	speaker.Lock()
	speaker.Clear()
	speaker.Unlock()
}
