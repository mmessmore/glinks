// Package cpu parses, exposes, and does computation based on Linux /proc/stat
package cpu

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Allow for testing on non-Linux platforms like my laptop
const TESTING bool = true

var TESTITR int = 1

// CpuInfo encapsulates CPU statistics.
// For individual data points each field is in the unit of Jiffies. Total is the sum of all fields
// When returned by Delta, each field is a percentage of total Jiffies over the duration
type CpuInfo struct {
	User      int
	Nice      int
	System    int
	Idle      int
	IoWait    int
	Irq       int
	SoftIrq   int
	Steal     int
	Guest     int
	GuestNice int
	Total     int
}

// IntrInfo presents interrupts serviced since boot.  The total is for all interrupts, then Intrs are the values for
// each
type IntrInfo struct {
	Total int
	Intrs []int
}

// Data holds all values held by /proc/stats
type Data struct {
	Cpu          CpuInfo
	Processors   []CpuInfo
	Intr         IntrInfo
	Ctxt         int
	Btime        int
	Processes    int
	ProcsRunning int
	ProcsBlocked int
	SoftIrq      IntrInfo
	Time         time.Time
}

// Data holds all values held by /proc/stats
type Delta struct {
	Cpu          CpuInfo
	Processors   []CpuInfo
	Intr         IntrInfo
	Ctxt         int
	Btime        int
	Processes    int
	ProcsRunning int
	ProcsBlocked int
	SoftIrq      IntrInfo
	Duration     time.Duration
}

// loadIntr creates an IntrInfo from the line of text
func loadIntr(fields []string) IntrInfo {
	info := IntrInfo{
		Total: atoi(fields[1]),
		Intrs: make([]int, len(fields)-1, len(fields)),
	}
	for index, val := range fields[2:] {
		info.Intrs[index] = atoi(val)
	}
	return info

}

// loadCpu creates an CpuInfo from the line of text for either the total or individual CPU
func loadCpu(fields []string) CpuInfo {
	info := CpuInfo{Steal: 0, Guest: 0, GuestNice: 0}

	// compute total jiffies from sum of available fields
	total := 0
	for _, val := range fields[1:] {
		total += atoi(val)
	}

	// I'm assuming at least a 2.6.0 kernel
	info.Total = total
	info.User = atoi(fields[1])
	info.Nice = atoi(fields[2])
	info.System = atoi(fields[3])
	info.Idle = atoi(fields[4])
	info.IoWait = atoi(fields[5])
	info.Irq = atoi(fields[6])
	info.SoftIrq = atoi(fields[7])

	// supported by 2.6.11+
	if len(fields) > 8 {
		info.Steal = atoi(fields[8])
	}

	// supported by 2.6.24+
	if len(fields) > 9 {
		info.Guest = atoi(fields[9])
	}
	// supported by 2.6.33+
	if len(fields) > 10 {
		info.GuestNice = atoi(fields[10])
	}
	return info
}

// Load takes the data from /proc/stat and presents a usable data structure
func Load() Data {
	cpu_stat_file := "/proc/stat"

	file, err := os.Open(cpu_stat_file)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_stat.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := Data{Time: time.Now()}
	cpuCount := runtime.NumCPU()
	data.Processors = make([]CpuInfo, cpuCount, cpuCount*2)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		switch {
		case fields[0] == "cpu":
			data.Cpu = loadCpu(fields)
		case strings.HasPrefix(fields[0], "cpu"):
			cpu_num := atoi(fields[0][3:])
			data.Processors[cpu_num] = loadCpu(fields)
		case fields[0] == "intr":
			data.Intr = loadIntr(fields)
		case fields[0] == "softirq":
			data.SoftIrq = loadIntr(fields)
		case fields[0] == "ctxt":
			data.Ctxt = atoi(fields[1])
		case fields[0] == "btime":
			data.Btime = atoi(fields[1])
		case fields[0] == "processes":
			data.Processes = atoi(fields[1])
		case fields[0] == "procs_running":
			data.ProcsRunning = atoi(fields[1])
		case fields[0] == "procs_blocked":
			data.ProcsBlocked = atoi(fields[1])
		}
	}

	check(scanner.Err())

	return data
}

// Delta computes the difference between a first and second Data element
// second is assumed to be later in time than first
// CPU statistics are reported back as percentages of use rather than Jiffies
func Diff(first Data, second Data) Delta {
	diff := Delta{Duration: second.Time.Sub(first.Time)}
	diff.Cpu = deltaCpu(first.Cpu, second.Cpu)
	diff.Processors = make([]CpuInfo, len(second.Processors))
	for i := range second.Processors {
		_ = "breakpoint"
		diff.Processors[i] = deltaCpu(first.Processors[i], second.Processors[i])
	}
	diff.Intr = deltaIntr(first.Intr, second.Intr)
	diff.SoftIrq = deltaIntr(first.SoftIrq, second.SoftIrq)

	diff.Ctxt = second.Ctxt - first.Ctxt
	diff.Btime = second.Btime - first.Btime
	diff.Processes = second.Processes - first.Processes
	diff.ProcsRunning = second.ProcsRunning - first.ProcsRunning
	diff.ProcsBlocked = second.ProcsBlocked - first.ProcsBlocked

	return diff
}

// deltaCpu computes percentage of use over time for each field
func deltaCpu(first CpuInfo, second CpuInfo) CpuInfo {
	diff := CpuInfo{}
	diff.Total = second.Total - first.Total
	diff.User = percent(first.User, second.User, diff.Total)
	diff.Nice = percent(first.Nice, second.Nice, diff.Total)
	diff.System = percent(first.System, second.System, diff.Total)
	diff.Idle = percent(first.Idle, second.Idle, diff.Total)
	diff.IoWait = percent(first.IoWait, second.IoWait, diff.Total)
	diff.Irq = percent(first.Irq, second.Irq, diff.Total)
	diff.SoftIrq = percent(first.SoftIrq, second.SoftIrq, diff.Total)
	diff.Steal = percent(first.Steal, second.Steal, diff.Total)
	diff.Guest = percent(first.Guest, second.Guest, diff.Total)
	diff.GuestNice = percent(first.GuestNice, second.GuestNice, diff.Total)

	return diff
}

func percent(first int, second int, total int) int {
	return int(float32(second-first) / float32(total) * 100.0)
}

// deltaIntr just calculates the difference in each field between two points in time
func deltaIntr(first IntrInfo, second IntrInfo) IntrInfo {
	diff := IntrInfo{}
	diff.Total = second.Total - first.Total
	diff.Intrs = make([]int, len(second.Intrs))
	for i := range second.Intrs {
		diff.Intrs[i] = second.Intrs[i] - first.Intrs[i]
	}
	return diff
}

// check panics on failure
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
