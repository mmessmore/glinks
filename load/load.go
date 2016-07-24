package load

import (
	"os"
	"log"
	"fmt"
	"time"
	"bufio"
	"strings"
	"encoding/json"
	"strconv"
)

const TESTING bool = true
var TESTITR int = 1

type Data struct {
	FiveMin float32
	TenMin float32
	FifteenMin float32
	RunableEnt int
	TotalEnt int
	LastPid int
	Time time.Time
}

func Load() (Data, string) {
	loadAvgFile := "/proc/loadavg"

	file, err := os.Open(loadAvgFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_loadavg.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
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

	perdy, err := json.MarshalIndent(data, "", "    ")
	check(err)
	return data, string(perdy)
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
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

func atof(s string) float32 {
	val, err := strconv.ParseFloat(s, 32)
	check(err)
	return float32(val)
}
