package vmstat

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
	Stats map[string]int
	Time  time.Time
}

func Load() Data {
	vmstatFile := "/proc/vmstat"

	file, err := os.Open(vmstatFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_vmstat.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now(), Stats: make(map[string]int)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Stats[fields[0]] = atoi(fields[1])
	}

	return data

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
