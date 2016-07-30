package glinks

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type VmstatData struct {
	Stats map[string]int
	Time  time.Time
}

func (d VmstatData) SampleTime() int64 {
	return d.Time.Unix()
}

func VmstatLoad() VmstatData {
	vmstatFile := "/proc/vmstat"

	file, err := os.Open(vmstatFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_vmstat.1")
		}
	}
	defer file.Close()
	check(err)

	data := VmstatData{Time: time.Now(), Stats: make(map[string]int)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Stats[fields[0]] = atoi(fields[1])
	}

	return data

}
