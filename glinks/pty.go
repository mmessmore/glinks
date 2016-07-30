package glinks

import (
	"bufio"
	"log"
	"os"
	"time"
)

type PtyData struct {
	Number int
	Max    int
	Time   time.Time
}

func (d PtyData) SampleTime() int64 {
	return d.Time.Unix()
}

func PtyLoad() PtyData {
	maxPtyFile := "/proc/sys/kernel/pty/max"
	file, err := os.Open(maxPtyFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_sys_kernel_pty_max.1")
		}
	}
	defer file.Close()
	check(err)

	data := PtyData{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data.Max = atoi(scanner.Text())
	}

	nrPtyFile := "/proc/sys/kernel/pty/nr"
	file, err = os.Open(nrPtyFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_sys_kernel_pty_nr.1")
		}
	}
	defer file.Close()
	check(err)

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		data.Number = atoi(scanner.Text())
	}

	return data
}
