package pty

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const TESTING bool = true

var TESTITR int = 1

type Data struct {
	Number int
	Max    int
	Time   time.Time
}

func (*Data) isSerializable() bool {
	return true
}

func Load() Data {
	maxPtyFile := "/proc/sys/kernel/pty/max"
	file, err := os.Open(maxPtyFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_sys_kernel_pty_max.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data.Max = atoi(scanner.Text())
	}

	nrPtyFile := "/proc/sys/kernel/pty/nr"
	file, err = os.Open(nrPtyFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_sys_kernel_pty_nr.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		data.Number = atoi(scanner.Text())
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
