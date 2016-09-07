package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/mmessmore/glinks/glinks"
)

var db *storm.DB

type storedCPU struct {
	ID        string
	CpuMetric glinks.CpuData
}
type storedIface struct {
	ID          string
	IfaceMetric glinks.IfaceData
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Awww snap!")
		panic(err)
	}
}

func output(d glinks.Data, is_graphite bool, prefix string) {
	if is_graphite {
		metrics := glinks.FromStruct(d, prefix)
		for _, m := range metrics {
			fmt.Println(m.String())
		}
	} else {
		value, err := json.MarshalIndent(d, "", "   ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", string(value))
		}
	}
}

func makePrefix(prefix string, postfix string) string {
	return fmt.Sprintf("%s.%s", prefix, postfix)
}

func do_check(dbID string, check string, is_graphite bool, prefix string) {
	switch check {
	case "cpu":
		var stored storedCPU
		newCPU := glinks.CpuLoad()

		err := db.One("ID", dbID, &stored)
		if err == storm.ErrNotFound {
			fmt.Fprintf(os.Stderr, "New ID: %s\n", dbID)
		} else if err != nil {
			checkErr(err)
		} else {
			output(glinks.CpuDiff(stored.CpuMetric, newCPU), is_graphite, makePrefix(prefix, "cpu"))
		}
		toStore := storedCPU{ID: dbID, CpuMetric: newCPU}
		err = db.Save(&toStore)
		checkErr(err)
	case "iface":
		var stored storedIface
		newIface := glinks.IfaceLoad()

		err := db.One("ID", dbID, &stored)
		if err == storm.ErrNotFound {
			fmt.Fprintf(os.Stderr, "New ID: %s\n", dbID)
		} else if err != nil {
			checkErr(err)
		} else {
			output(glinks.IfaceDiff(stored.IfaceMetric, newIface), is_graphite, makePrefix(prefix, "iface"))
		}
		toStore := storedIface{ID: dbID, IfaceMetric: newIface}
		err = db.Save(&toStore)
		checkErr(err)
	case "disk":
		output(glinks.DiskLoad(), is_graphite, makePrefix(prefix, "disk"))
	case "load":
		output(glinks.LoadLoad(), is_graphite, makePrefix(prefix, "load"))
	case "mem":
		output(glinks.MemLoad(), is_graphite, makePrefix(prefix, "memory"))
	case "fh":
		output(glinks.FhLoad(), is_graphite, makePrefix(prefix, "file_handles"))
	case "inode":
		output(glinks.InodeLoad(), is_graphite, makePrefix(prefix, "inodes"))
	case "pty":
		output(glinks.PtyLoad(), is_graphite, makePrefix(prefix, "pty"))
	case "entropy":
		output(glinks.EntropyLoad(), is_graphite, makePrefix(prefix, "entropy"))
	case "uptime":
		output(glinks.UptimeLoad(), is_graphite, makePrefix(prefix, "uptime"))
	case "vmstat":
		output(glinks.VmstatLoad(), is_graphite, makePrefix(prefix, "vmstat"))
	}
}

func Cli(dbPath string, dbID string, checks []string, is_graphite bool, prefix string, count int, delay int) {
	var err error
	db, err = storm.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I can't open the database file: %s\n", dbPath)
		panic(err)
	}
	defer db.Close()
	for i := 1; i <= count; i++ {
		for _, check := range checks {
			do_check(dbID, check, is_graphite, prefix)
		}
		// sleep till the next go 'round
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
