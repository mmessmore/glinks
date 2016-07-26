package uptime

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

type Data struct {
	Uptime float32
	Idle   float32
	Time   time.Time
}

func Load() Data {
	uptimeFile := "/proc/uptime"
	file, err := os.Open(uptimeFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_uptime.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Uptime = atof(fields[0])
		data.Idle = atof(fields[1])

	}

	return data
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		log.Panic(e)
	}
}

func atof(s string) float32 {
	val, err := strconv.ParseFloat(s, 32)
	check(err)
	return float32(val)
}
