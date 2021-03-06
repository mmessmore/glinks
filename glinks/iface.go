package glinks

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type IfaceInfo struct {
	RxBytes      int
	RxPackets    int
	RxErrs       int
	RxDrop       int
	RxFifo       int
	RxFrame      int
	RxCompressed int
	RxMulticast  int
	TxBytes      int
	TxPackets    int
	TxErrs       int
	TxDrop       int
	TxFifo       int
	TxColls      int
	TxCarrier    int
	TxCompressed int
}

type IfaceData struct {
	Interfaces map[string]IfaceInfo
	Time       time.Time
}

func (d *IfaceData) SampleTime() int64 {
	return d.Time.Unix()
}

type IfaceDelta struct {
	Interfaces map[string]IfaceInfo
	Duration   time.Duration
	Time       time.Time
}

func (d IfaceDelta) SampleTime() int64 {
	return d.Time.Unix()
}

func IfaceLoad() IfaceData {
	ifaceFile := "/proc/net/dev"

	file, err := os.Open(ifaceFile)

	// for testing on non-Linux OSes I have an example copied off a linux host
	if TESTING {
		if os.IsNotExist(err) {
			log.Print("Falling back to dummy data")
			file, err = os.Open(fmt.Sprintf("dummy_data/proc_net_dev.%d", TESTITR))
			TESTITR += 1
		}
	}
	defer file.Close()
	check(err)

	data := IfaceData{Time: time.Now(), Interfaces: make(map[string]IfaceInfo)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		// skip the headers on the text table
		if !strings.Contains(fields[0], ":") {
			continue
		}
		iface := strings.Trim(fields[0], " ")
		iface = iface[0 : len(iface)-1]
		data.Interfaces[iface] = loadIfaceLine(fields)
	}

	return data
}

func IfaceDiff(first IfaceData, second IfaceData) IfaceDelta {
	delta := IfaceDelta{
		Duration:   second.Time.Sub(first.Time),
		Interfaces: make(map[string]IfaceInfo),
		Time:       second.Time,
	}
	for k := range second.Interfaces {
		delta.Interfaces[k] = diffIfaceLine(first.Interfaces[k], second.Interfaces[k])
	}
	if TESTING {
		TESTITR = 1
	}
	return delta
}

func loadIfaceLine(fields []string) IfaceInfo {
	info := IfaceInfo{}
	info.RxBytes = atoi(fields[1])
	info.RxPackets = atoi(fields[2])
	info.RxErrs = atoi(fields[3])
	info.RxDrop = atoi(fields[4])
	info.RxFifo = atoi(fields[5])
	info.RxFrame = atoi(fields[6])
	info.RxCompressed = atoi(fields[7])
	info.RxMulticast = atoi(fields[8])
	info.TxBytes = atoi(fields[9])
	info.TxPackets = atoi(fields[10])
	info.TxErrs = atoi(fields[11])
	info.TxDrop = atoi(fields[12])
	info.TxFifo = atoi(fields[13])
	info.TxColls = atoi(fields[14])
	info.TxCarrier = atoi(fields[15])
	info.TxCompressed = atoi(fields[16])
	return info
}

func diffIfaceLine(first IfaceInfo, second IfaceInfo) IfaceInfo {
	info := IfaceInfo{}
	info.RxBytes = second.RxBytes - first.RxBytes
	info.RxPackets = second.RxPackets - first.RxPackets
	info.RxErrs = second.RxErrs - first.RxErrs
	info.RxDrop = second.RxDrop - first.RxDrop
	info.RxFifo = second.RxFifo - first.RxFifo
	info.RxFrame = second.RxFrame - first.RxFrame
	info.RxCompressed = second.RxCompressed - first.RxCompressed
	info.RxMulticast = second.RxMulticast - first.RxMulticast
	info.TxBytes = second.TxBytes - first.TxBytes
	info.TxPackets = second.TxPackets - first.TxPackets
	info.TxErrs = second.TxErrs - first.TxErrs
	info.TxDrop = second.TxDrop - first.TxDrop
	info.TxFifo = second.TxFifo - first.TxFifo
	info.TxColls = second.TxColls - first.TxColls
	info.TxCarrier = second.TxCarrier - first.TxCarrier
	info.TxCompressed = second.TxCompressed - first.TxCompressed
	return info
}
