package glinks

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type InodeData struct {
	Number int
	Free   int
	Time   time.Time
}

func (d InodeData) SampleTime() int64 {
	return d.Time.Unix()
}

func InodeLoad() InodeData {
	inodeFile := "/proc/sys/fs/inode-nr"

	file, err := os.Open(inodeFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open("dummy_data/proc_sys_fs_inode-nr.1")
		}
	}
	defer file.Close()
	check(err)

	data := InodeData{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		data.Number = atoi(fields[0])
		data.Free = atoi(fields[1])
	}

	return data
}
