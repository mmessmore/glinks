package glinks

import (
	"log"
	"strconv"
)

// Allow for testing on non-Linux platforms like my laptop
const TESTING bool = true

var TESTITR int = 1

type Data interface {
	SampleTime() int64
}

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

func atof(s string) float32 {
	val, err := strconv.ParseFloat(s, 32)
	check(err)
	return float32(val)
}
