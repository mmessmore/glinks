package glinks

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type UptimeData struct {
	Uptime float32
	Idle   float32
	Time   time.Time
}

func (d UptimeData) SampleTime() int64 {
	return d.Time.Unix()
}

func UptimeLoad() UptimeData {
	uptimeFile := "/proc/uptime"
	file, err := os.Open(uptimeFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_uptime.1")
		}
	}
	defer file.Close()
	check(err)

	data := UptimeData{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Uptime = atof(fields[0])
		data.Idle = atof(fields[1])

	}

	return data
}
