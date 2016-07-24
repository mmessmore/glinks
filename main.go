package main

import (
	"fmt"
	"time"
	"github.com/mmessmore/glinks/cpu"
	"github.com/mmessmore/glinks/disk"
	"github.com/mmessmore/glinks/load"
	"github.com/mmessmore/glinks/mem"
)


// run that shit dawg
func main() {
	data_one, _ := cpu.Load()
	time.Sleep(2 * time.Second)
	data_two, _ := cpu.Load()

	_, json_delta := cpu.Delta(data_one, data_two)
	fmt.Println("CPU DELTA")
	fmt.Print(json_delta)

	_, json_disk := disk.Load()
	fmt.Print(json_disk)

	_, json_load := load.Load()
	fmt.Print(json_load)

	_, json_mem := mem.Load()
	fmt.Print(json_mem)
}
