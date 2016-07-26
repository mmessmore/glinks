package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mmessmore/glinks/cpu"
	"github.com/mmessmore/glinks/disk"
	"github.com/mmessmore/glinks/entropy"
	"github.com/mmessmore/glinks/fh"
	"github.com/mmessmore/glinks/iface"
	"github.com/mmessmore/glinks/inode"
	"github.com/mmessmore/glinks/load"
	"github.com/mmessmore/glinks/mem"
	"github.com/mmessmore/glinks/pty"
	"github.com/mmessmore/glinks/uptime"
	"github.com/mmessmore/glinks/vmstat"
	"time"
)

func output(d interface{}) {
	value, err := json.MarshalIndent(d, "", "   ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", string(value))
	}
}

// run that shit dawg
func main() {

	var cpuFlag = flag.Bool("cpu", false, "Print CPU Metrics")
	var diskFlag = flag.Bool("disk", false, "Print Disk Metrics")
	var loadFlag = flag.Bool("load", false, "Print Load Metrics")
	var memFlag = flag.Bool("mem", false, "Print Memory Metrics")
	var ifaceFlag = flag.Bool("iface", false, "Print Network Interface Metrics")
	var fhFlag = flag.Bool("fh", false, "Print File Handle Metrics")
	var inodeFlag = flag.Bool("inode", false, "Print INode Metrics")
	var ptyFlag = flag.Bool("pty", false, "Print PTY Metrics")
	var entropyFlag = flag.Bool("entropy", false, "Print Entropy Metrics")
	var uptimeFlag = flag.Bool("uptime", false, "Print Uptime Metrics")
	var vmstatFlag = flag.Bool("vmstat", false, "Print VMStat Metrics")
	var allFlag = flag.Bool("all", false, "Print all Metrics")

	flag.Parse()

	if *cpuFlag || *allFlag {
		cpu_data_one := cpu.Load()
		time.Sleep(2 * time.Second)
		cpu_data_two := cpu.Load()
		output(cpu.Diff(cpu_data_one, cpu_data_two))
	}

	if *diskFlag || *allFlag {
		output(disk.Load())
	}

	if *loadFlag || *allFlag {
		output(load.Load())
	}

	if *memFlag || *allFlag {
		output(mem.Load())
	}

	if *ifaceFlag || *allFlag {
		iface_data_one := iface.Load()
		time.Sleep(2 * time.Second)
		iface_data_two := iface.Load()
		output(iface.Diff(iface_data_one, iface_data_two))
	}

	if *fhFlag || *allFlag {
		output(fh.Load())
	}

	if *inodeFlag || *allFlag {
		output(inode.Load())
	}

	if *ptyFlag || *allFlag {
		output(pty.Load())
	}

	if *entropyFlag || *allFlag {
		output(entropy.Load())
	}

	if *uptimeFlag || *allFlag {
		output(uptime.Load())
	}

	if *vmstatFlag || *allFlag {
		output(vmstat.Load())
	}
}
