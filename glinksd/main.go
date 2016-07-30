package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mmessmore/glinks/glinks"
)

func defaultRoute(c *gin.Context) {
	value := c.Param("stat")

	switch value {
	case "cpu":
		c.JSON(200, glinks.CpuLoad())
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

func main() {
	r := gin.Default()
	r.GET("/:stat", defaultRoute)

	r.Run()

}
