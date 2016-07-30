package glinks

import (
	"bufio"
	"log"
	"os"
	"time"
)

type EntropyData struct {
	Available int
	Time      time.Time
}

func (d EntropyData) SampleTime() int64 {
	return d.Time.Unix()
}

func EntropyLoad() EntropyData {
	entropyFile := "/proc/sys/kernel/random/entropy_avail"
	file, err := os.Open(entropyFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_sys_kernel_random_entropy_avail.1")
		}
	}
	defer file.Close()
	check(err)

	data := EntropyData{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data.Available = atoi(scanner.Text())
	}

	return data
}
