package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mmessmore/glinks/cli"
	"github.com/satori/go.uuid"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func decode_checks(cliChecks string) []string {
	possible_checks := []string{"cpu", "disk", "load", "mem", "iface", "fh", "inode", "pty", "entropy",
		"uptime", "vmstat"}
	checks := make([]string, 0, 12)
	for _, check := range strings.Split(cliChecks, ",") {
		if check == "all" {
			return possible_checks
		} else if contains(possible_checks, check) {
			checks = append(checks, check)
		} else {
			fmt.Printf("Invalid check: %s", check)
			os.Exit(22) // invalid argument
		}
	}
	return checks
}

// run that shit dawg
func main() {

	hostname, err := os.Hostname()

	if err != nil {
		fmt.Println("I can't determine hostname.  I don't even know who I am")
		panic(err)
	}

	var cliChecks = flag.String("checks", "all", "Comma separated list of checks or 'all' for all")
	var is_graphite = flag.Bool("graphite", false, "Use graphite output instead of JSON")
	var prefix = flag.String("graphite-prefix", hostname, "Prefix for graphite values")
	var count = flag.Int("count", 1, "number of iterations to run")
	var delay = flag.Int("delay", 1, "sleep between iterations")
	var dbPath = flag.String("db-path", "/tmp/glinks.db", "database path")
	var dbId = flag.String("db-id", uuid.NewV4().String(), "ID string of database session, often a UUID")

	flag.Parse()

	checks := decode_checks(*cliChecks)

	// If we're an httpd daemon, go do that
	cli.Cli(*dbPath, *dbId, checks, *is_graphite, *prefix, *count, *delay)
}
