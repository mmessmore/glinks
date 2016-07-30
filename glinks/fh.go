package glinks

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type FhData struct {
	Used int
	//Free int
	Max  int
	Time time.Time
}

func (d FhData) SampleTime() int64 {
	return d.Time.Unix()
}

func FhLoad() FhData {
	fhFile := "/proc/sys/fs/file-nr"

	file, err := os.Open(fhFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_sys_fs_file-nr.1")
		}
	}
	defer file.Close()
	check(err)

	data := FhData{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Used = atoi(fields[0])
		//data.Free = atoi(fields[1])
		data.Max = atoi(fields[2])
	}

	return data
}
