package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mmessmore/glinks/glinks"
	"os"
	"time"
)

func output(d glinks.Data, is_graphite bool, prefix string) {
	if is_graphite {
		metrics := glinks.FromStruct(d, prefix)
		for _, m := range metrics {
			fmt.Println(m.String())
		}
	} else {
		value, err := json.MarshalIndent(d, "", "   ")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n", string(value))
		}
	}
}

// run that shit dawg
func main() {

	hostname, err := os.Hostname()

	if err != nil {
		fmt.Println("I can't determine hostname.  I don't even know who I am")
		panic(err)
	}

	var cpuFlag = flag.Bool("cpu", false, "Output CPU Metrics")
	var diskFlag = flag.Bool("disk", false, "Output Disk Metrics")
	var loadFlag = flag.Bool("load", false, "Output Load Metrics")
	var memFlag = flag.Bool("mem", false, "Output Memory Metrics")
	var ifaceFlag = flag.Bool("iface", false, "Output Network Interface Metrics")
	var fhFlag = flag.Bool("fh", false, "Output File Handle Metrics")
	var inodeFlag = flag.Bool("inode", false, "Output INode Metrics")
	var ptyFlag = flag.Bool("pty", false, "Output PTY Metrics")
	var entropyFlag = flag.Bool("entropy", false, "Output Entropy Metrics")
	var uptimeFlag = flag.Bool("uptime", false, "Output Uptime Metrics")
	var vmstatFlag = flag.Bool("vmstat", false, "Output VMStat Metrics")
	var allFlag = flag.Bool("all", false, "Output all Metrics")
	var is_graphite = flag.Bool("graphite", false, "Use graphite output instead of JSON")
	var prefix = flag.String("graphite-prefix", hostname, "Prefix for graphite values")
	var count = flag.Int("count", 2, "number of iterations to run")
	var delay = flag.Int("delay", 1, "sleep between iterations")

	flag.Parse()

	lastCpu := glinks.CpuData{}
	lastIface := glinks.IfaceData{}

	for i := 1; i <= *count; i++ {
		if *cpuFlag || *allFlag {
			if lastCpu.Time == (time.Time{}) {
				lastCpu = glinks.CpuLoad()
			} else {
				currentCpu := glinks.CpuLoad()
				output(glinks.CpuDiff(lastCpu, currentCpu), *is_graphite, makePrefix(*prefix, "cpu"))
				lastCpu = currentCpu
			}
		}

		if *ifaceFlag || *allFlag {
			if lastIface.Time == (time.Time{}) {
				lastIface = glinks.IfaceLoad()
			} else {
				currentIface := glinks.IfaceLoad()
				output(glinks.IfaceDiff(lastIface, currentIface),
					*is_graphite,
					makePrefix(*prefix, "network_interface"))
				lastIface = currentIface
			}
		}

		if *diskFlag || *allFlag {
			output(glinks.DiskLoad(), *is_graphite, makePrefix(*prefix, "disk"))
		}

		if *loadFlag || *allFlag {
			output(glinks.LoadLoad(), *is_graphite, makePrefix(*prefix, "load"))
		}

		if *memFlag || *allFlag {
			output(glinks.MemLoad(), *is_graphite, makePrefix(*prefix, "memory"))
		}

		if *fhFlag || *allFlag {
			output(glinks.FhLoad(), *is_graphite, makePrefix(*prefix, "file_handles"))
		}

		if *inodeFlag || *allFlag {
			output(glinks.InodeLoad(), *is_graphite, makePrefix(*prefix, "inodes"))
		}

		if *ptyFlag || *allFlag {
			output(glinks.PtyLoad(), *is_graphite, makePrefix(*prefix, "pty"))
		}

		if *entropyFlag || *allFlag {
			output(glinks.EntropyLoad(), *is_graphite, makePrefix(*prefix, "entropy"))
		}

		if *uptimeFlag || *allFlag {
			output(glinks.UptimeLoad(), *is_graphite, makePrefix(*prefix, "uptime"))
		}

		if *vmstatFlag || *allFlag {
			output(glinks.VmstatLoad(), *is_graphite, makePrefix(*prefix, "vmstat"))
		}
		// sleep till the next go 'round
		time.Sleep(time.Duration(*delay) * time.Second)
	}
}

func makePrefix(prefix string, postfix string) string {
	return fmt.Sprintf("%s.%s", prefix, postfix)
}
