package glinks

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type LoadData struct {
	FiveMin    float32
	TenMin     float32
	FifteenMin float32
	RunableEnt int
	TotalEnt   int
	LastPid    int
	Time       time.Time
}

func (d LoadData) SampleTime() int64 {
	return d.Time.Unix()
}

func LoadLoad() LoadData {
	loadAvgFile := "/proc/loadavg"

	file, err := os.Open(loadAvgFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_loadavg.1")
		}
	}
	defer file.Close()
	check(err)

	data := LoadData{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.FiveMin = atof(fields[0])
		data.TenMin = atof(fields[1])
		data.FifteenMin = atof(fields[2])

		subfields := strings.Split(fields[3], "/")
		data.RunableEnt = atoi(subfields[0])
		data.TotalEnt = atoi(subfields[1])

		data.LastPid = atoi(fields[4])
	}

	return data
}
