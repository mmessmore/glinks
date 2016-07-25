package fh

import (
	"bufio"
	"encoding/json"
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
	Used int
	//Free int
	Max  int
	Time time.Time
}

func Load() (Data, string) {
	fhFile := "/proc/sys/fs/file-nr"

	file, err := os.Open(fhFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_sys_fs_file-nr.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Used = atoi(fields[0])
		//data.Free = atoi(fields[1])
		data.Max = atoi(fields[2])
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
