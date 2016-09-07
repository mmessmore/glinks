package http

import (
	"fmt"
	"os"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/mmessmore/glinks/glinks"
)

var db storm.DB

func cpuDelta() glinks.CpuDelta {
	var old = glinks.CpuData{}
	//db.Get("cpu", uuid_key, &old)
	new := glinks.CpuLoad()
	return glinks.CpuDiff(old, new)
}

func defaultRoute(c *gin.Context) {
	value := c.Param("stat")

	switch value {
	case "cpu":
		delta := cpuDelta()
		c.JSON(200, delta)
	case "disk":
		c.JSON(200, glinks.DiskLoad())
	case "entropy":
		c.JSON(200, glinks.EntropyLoad())
	case "fh":
		c.JSON(200, glinks.FhLoad())
	case "iface":
		c.JSON(200, glinks.IfaceLoad())
	case "inode":
		c.JSON(200, glinks.InodeLoad())
	case "load":
		c.JSON(200, glinks.LoadLoad())
	case "mem":
		c.JSON(200, glinks.MemLoad())
	case "pty":
		c.JSON(200, glinks.PtyLoad())
	case "uptime":
		c.JSON(200, glinks.UptimeLoad())
	case "vmstat":
		c.JSON(200, glinks.VmstatLoad())
	default:
		c.AbortWithStatus(404)
	}
}

func Httpd(dbPath string) {
	db, err := storm.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I can't open the database file: %s\n", dbPath)
		panic(err)
	}
	defer db.Close()
	r := gin.Default()
	r.GET("/:stat", defaultRoute)

	r.Run()
}
