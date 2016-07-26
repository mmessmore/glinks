package disk

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const TESTING bool = true

var TESTITR int = 1

type AccessStat struct {
	Count        int
	Merged       int
	Sectors      int
	Milliseconds int
}

type DiskData struct {
	Reads                  AccessStat
	Writes                 AccessStat
	IoInProgress           int
	IoMilliseconds         int
	IoWeightedMilliseconds int
}

type Data struct {
	Time  time.Time
	Disks map[string]DiskData
}

func Load() Data {
	disk_stat_file := "/proc/diskstats"

	file, err := os.Open(disk_stat_file)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_diskstats.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now(), Disks: make(map[string]DiskData)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := fields[2]
		data.Disks[name] = loadLine(fields[3:])
	}

	return data
}

func loadLine(fields []string) DiskData {
	data := DiskData{}

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

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// atoi converts a string to an integer and panic if something is awry.  This should not be a problem given that we
// are dealing with a very fixed format
func atoi(s string) int {
	val, err := strconv.Atoi(s)
	check(err)
	return val
}
