package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

const (
	FMEDIA_EXECUTABLE = "C:\\Users\\WerK\\go\\src\\gobark\\fmedia\\fmedia.exe"
	TRIGGER_PATH      = "C:\\Users\\WerK\\Dropbox\\trigger.txt"
	PAYLOAD_PATH      = "C:\\Users\\WerK\\go\\src\\gobark\\payload.mp3"
	READY_STRING      = "READY"
)

func main() {
	if err := ioutil.WriteFile(TRIGGER_PATH, []byte(READY_STRING), 0644); err != nil {
		panic(err)
	}
	log.Printf("Ready")
	for {
		triggerData, err := ioutil.ReadFile(TRIGGER_PATH)
		if err != nil {
			panic(err)
		}
		if bytes.Compare(triggerData, []byte(READY_STRING)) != 0 {
			log.Printf("Triggered")
			if err := exec.Command(FMEDIA_EXECUTABLE, PAYLOAD_PATH).Run(); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(TRIGGER_PATH, []byte(READY_STRING), 0644); err != nil {
				panic(err)
			}
			log.Printf("Ready")
		}
		time.Sleep(time.Second)
	}
}
