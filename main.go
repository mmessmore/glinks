package main

import (
	"fmt"
	"time"
	"github.com/mmessmore/glinks/cpu"
	"github.com/mmessmore/glinks/disk"
	"github.com/mmessmore/glinks/load"
	"github.com/mmessmore/glinks/mem"
	"github.com/mmessmore/glinks/iface"
	"flag"
)


// run that shit dawg
func main() {

	var cpuFlag = flag.Bool("cpu", false, "Run CPU Metrics")
	var diskFlag = flag.Bool("disk", false, "Run Disk Metrics")
	var loadFlag = flag.Bool("load", false, "Run Load Metrics")
	var memFlag = flag.Bool("mem", false, "Run Memory Metrics")
	var ifaceFlag = flag.Bool("iface", false, "Run Network Interface Metrics")
	var allFlag = flag.Bool("all", false, "Run all Metrics")

	flag.Parse()

	if *cpuFlag || *allFlag {
		cpu_data_one, _ := cpu.Load()
		time.Sleep(2 * time.Second)
		cpu_data_two, _ := cpu.Load()
		_, json_delta := cpu.Diff(cpu_data_one, cpu_data_two)
		fmt.Println("CPU DELTA")
		fmt.Print(json_delta)
		fmt.Println("")
	}

	if *diskFlag || *allFlag {
		_, json_disk := disk.Load()
		fmt.Print(json_disk)
		fmt.Println("")
	}

	if *loadFlag || *allFlag {
		_, json_load := load.Load()
		fmt.Print(json_load)
		fmt.Println("")
	}

	if *memFlag || *allFlag {
		_, json_mem := mem.Load()
		fmt.Print(json_mem)
		fmt.Println("")
	}

	if *ifaceFlag || *allFlag {
		iface_data_one, _ := iface.Load()
		time.Sleep(2 * time.Second)
		iface_data_two, _ := iface.Load()
		_, json_iface := iface.Diff(iface_data_one, iface_data_two)
		fmt.Print(json_iface)
		fmt.Println("")
	}
}
