package mem

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"os"
	"bufio"
	"strings"
)

const TESTING bool = true
var TESTITR int = 1

type Data struct {
	MemTotal int
	MemFree int
	MemAvailable int
	Buffers int
	Cached int
	SwapCached int
	Active int
	Inactive int
	ActiveAnon int
	InactiveAnon int
	ActiveFile int
	InactiveFile int
	Unevictable int
	Mlocked int
	SwapTotal int
	SwapFree int
	Dirty int
	Writeback int
	AnonPages int
	Mapped int
	Shmem int
	Slab int
	SReclaimable int
	SUnreclaim int
	KernelStack int
	PageTables int
	NFS_Unstable int
	Bounce int
	WritebackTmp int
	CommitLimit int
	Committed_AS int
	VmallocTotal int
	VmallocUsed int
	VmallocChunk int
	HardwareCorrupted int
	AnonHugePages int
	HugePages_Total int
	HugePages_Free int
	HugePages_Rsvd int
	HugePages_Surp int
	Hugepagesize int
	DirectMap4k int
	DirectMap2M int
	Time time.Time
}

func Load() (Data, string) {
	memInfoFile := "/proc/meminfo"

	file, err := os.Open(memInfoFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_meminfo.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		switch fields[0] {
		case "MemTotal:":
			data.MemTotal = atoi(fields[1])
		case "MemFree:":
			data.MemFree = atoi(fields[1])
		case "MemAvailable:":
			data.MemAvailable = atoi(fields[1])
		case "Buffers:":
			data.Buffers = atoi(fields[1])
		case "Cached:":
			data.Cached = atoi(fields[1])
		case "SwapCached:":
			data.SwapCached = atoi(fields[1])
		case "Active:":
			data.Active = atoi(fields[1])
		case "Inactive:":
			data.Inactive = atoi(fields[1])
		case "Active(anon):":
			data.ActiveAnon = atoi(fields[1])
		case "Inactive(anon):":
			data.InactiveAnon = atoi(fields[1])
		case "Active(file):":
			data.ActiveFile = atoi(fields[1])
		case "Inactive(file):":
			data.InactiveFile = atoi(fields[1])
		case "Unevictable:":
			data.Unevictable = atoi(fields[1])
		case "Mlocked:":
			data.Mlocked = atoi(fields[1])
		case "SwapTotal:":
			data.SwapTotal = atoi(fields[1])
		case "SwapFree:":
			data.SwapFree = atoi(fields[1])
		case "Dirty:":
			data.Dirty = atoi(fields[1])
		case "Writeback:":
			data.Writeback = atoi(fields[1])
		case "AnonPages:":
			data.AnonPages = atoi(fields[1])
		case "Mapped:":
			data.Mapped = atoi(fields[1])
		case "Shmem:":
			data.Shmem = atoi(fields[1])
		case "Slab:":
			data.Slab = atoi(fields[1])
		case "SReclaimable:":
			data.SReclaimable = atoi(fields[1])
		case "SUnreclaim:":
			data.SUnreclaim = atoi(fields[1])
		case "KernelStack:":
			data.KernelStack = atoi(fields[1])
		case "PageTables:":
			data.PageTables = atoi(fields[1])
		case "NFS_Unstable:":
			data.NFS_Unstable = atoi(fields[1])
		case "Bounce:":
			data.Bounce = atoi(fields[1])
		case "WritebackTmp:":
			data.WritebackTmp = atoi(fields[1])
		case "CommitLimit:":
			data.CommitLimit = atoi(fields[1])
		case "Committed_AS:":
			data.Committed_AS = atoi(fields[1])
		case "VmallocTotal:":
			data.VmallocTotal = atoi(fields[1])
		case "VmallocUsed:":
			data.VmallocUsed = atoi(fields[1])
		case "VmallocChunk:":
			data.VmallocChunk = atoi(fields[1])
		case "HardwareCorrupted:":
			data.HardwareCorrupted = atoi(fields[1])
		case "AnonHugePages:":
			data.AnonHugePages = atoi(fields[1])
		case "HugePages_Total:":
			data.HugePages_Total = atoi(fields[1])
		case "HugePages_Free:":
			data.HugePages_Free = atoi(fields[1])
		case "HugePages_Rsvd:":
			data.HugePages_Rsvd = atoi(fields[1])
		case "HugePages_Surp:":
			data.HugePages_Surp = atoi(fields[1])
		case "Hugepagesize:":
			data.Hugepagesize = atoi(fields[1])
		case "DirectMap4k:":
			data.DirectMap4k = atoi(fields[1])
		case "DirectMap2M:":
			data.DirectMap2M = atoi(fields[1])
		}

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
