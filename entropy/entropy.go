package entropy

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
	Available int
	Time      time.Time
}

func Load() Data {
	entropyFile := "/proc/sys/kernel/random/entropy_avail"
	file, err := os.Open(entropyFile)
	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_sys_kernel_random_entropy_avail.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data.Available = atoi(scanner.Text())
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
