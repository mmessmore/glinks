package main

import (
	"fmt"
	"github.com/mmessmore/glinks/cpu"
	"time"
)


// run that shit dawg
func main() {
	data_one, _ := cpu.Load()
	time.Sleep(2 * time.Second)
	data_two, _ := cpu.Load()

	_, json_delta := cpu.Delta(data_one, data_two)
	fmt.Println("CPU DELTA")
	fmt.Print(json_delta)

}
