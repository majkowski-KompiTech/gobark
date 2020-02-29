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
)

func main() {
	triggerData, err := ioutil.ReadFile(TRIGGER_PATH)
	if err != nil {
		panic(err)
	}

	log.Printf("Ready")

	for {
		newerTriggerData, err := ioutil.ReadFile(TRIGGER_PATH)
		if err != nil {
			panic(err)
		}

		if bytes.Compare(triggerData, newerTriggerData) != 0 {
			log.Printf("Triggered")

			if err := exec.Command(FMEDIA_EXECUTABLE, PAYLOAD_PATH).Run(); err != nil {
				panic(err)
			}

			triggerData = newerTriggerData

			log.Printf("Ready")
		}

		time.Sleep(time.Second)
	}
}
