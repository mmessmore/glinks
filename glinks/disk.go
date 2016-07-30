package glinks

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type AccessStat struct {
	Count        int
	Merged       int
	Sectors      int
	Milliseconds int
}

type DiskInfo struct {
	Reads                  AccessStat
	Writes                 AccessStat
	IoInProgress           int
	IoMilliseconds         int
	IoWeightedMilliseconds int
}

type DiskData struct {
	Time  time.Time
	Disks map[string]DiskInfo
}

func (d DiskData) SampleTime() int64 {
	return d.Time.Unix()
}

func DiskLoad() DiskData {
	disk_stat_file := "/proc/diskstats"

	file, err := os.Open(disk_stat_file)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_diskstats.1")
		}
	}
	defer file.Close()
	check(err)

	data := DiskData{Time: time.Now(), Disks: make(map[string]DiskInfo)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := fields[2]
		data.Disks[name] = loadDiskLine(fields[3:])
	}

	return data
}

func loadDiskLine(fields []string) DiskInfo {
	data := DiskInfo{}

	data.Reads = AccessStat{
		Count:        atoi(fields[0]),
		Merged:       atoi(fields[1]),
		Sectors:      atoi(fields[2]),
		Milliseconds: atoi(fields[3]),
	}
	data.Writes = AccessStat{
		Count:        atoi(fields[4]),
		Merged:       atoi(fields[5]),
		Sectors:      atoi(fields[6]),
		Milliseconds: atoi(fields[7]),
	}
	data.IoInProgress = atoi(fields[8])
	data.IoMilliseconds = atoi(fields[9])
	data.IoWeightedMilliseconds = atoi(fields[10])

	return data
}
