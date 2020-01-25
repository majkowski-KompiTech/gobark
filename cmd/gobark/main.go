package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const FMEDIA_EXECUTABLE="C:\\Users\\WerK\\go\\src\\gobark\\fmedia\\fmedia.exe"
const TMP_PATH ="clip.mp3"
const SAVED_PATH="saved"

func main() {
	_ = os.Remove(TMP_PATH)

	for {
		cmd := exec.Command(FMEDIA_EXECUTABLE, "--record", "--until=10", "--out="+TMP_PATH) // record 10 sec chunks then command terminates
		var out bytes.Buffer
		cmd.Stderr = &out // capture stdout

		log.Printf("recording...")
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		log.Printf("analyzing...")
		cleanStr := stripCtlFromBytes(out.String())

		peak := findPeak(&cleanStr)
		// peak: <0:10>

		if peak >= 5 { // remember to turn up mike sensitivity in Windows to max
			filename := filepath.Join(SAVED_PATH, fmt.Sprintf("%s%d.mp3", time.Now().Format(time.RFC3339), peak))
			filename = strings.Replace(filename, ":", ";", -1) // Windows does not like character : in filenames

			if err := os.Rename(TMP_PATH, filename); err != nil {
				log.Fatal(err)
			}

			log.Print("Clip saved: " + filename)
		} else {
			if err := os.Remove(TMP_PATH); err != nil {
				log.Fatal(err)
			}

			log.Print("Clip not saved")
		}
	}
}

// remove control chars from returned buffer, so we have history of all levels recorded, instead of it being overwritten when printed on terminal
// https://rosettacode.org/wiki/Strip_control_codes_and_extended_characters_from_a_string#Go
func stripCtlFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c != 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

func getVU(in *string, pos *int) int {
	// fmedia v1.14 (win-x64)Recording...  Source: int16 48000Hz stereo.  Press "s" to stop.0:00  [..........] -40.00dB / -40.00dB  0:00  [===.......] -27.56dB / -18.70dB  0:00  [=.........] -33.42dB / -18.70dB  0:00  [===.......] -27.03dB / -18.70dB  0:00  [==........] -28.65dB / -18.70dB

	for {
		// find potential start of next VU
		potStart := strings.Index((*in)[*pos:], "[")

		if potStart == -1 {
			// no more VUs
			return -1
		}

		vuPos := potStart + *pos // absolute start of next VU
		vuEnd := vuPos + 11 // get expected end of next VU

		// Skip something that is not a VU, just starts with [
		if string((*in)[vuEnd]) != "]" {
			// continue search at vuEnd + 1
			*pos = vuEnd + 1
			continue
		}

		value := 0
		vuPos += 1 // start search at first VU index

		// calculate VU value
		for i := 0 ; i < 10 ; i++ {
			s := string((*in)[vuPos + i])

			if s == "=" {
				value++
			} else if s == "." {
				// ok, do not incr value
			} else {
				// unexpected, continue with next VU
				// should not happen :)
				log.Printf("unexpected value in VU: %s", (*in)[vuPos:vuEnd])
				continue
			}
		}

		// let caller know where to continue search
		*pos = vuEnd + 1
		return value
	}
}

func findPeak(in *string) int {
	// lets just parse the VU meter instead of messing with float dB values
	pos := 0
	max := 0

	for {
		val := getVU(in, &pos)
		if val == -1 {
			break // no more VUs
		}

		if val > max {
			max = val
		}
	}

	return max
}
