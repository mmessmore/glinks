package main

import (
	"github.com/gin-gonic/gin"
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
)

func defaultRoute(c *gin.Context) {
	value := c.Param("stat")

	switch value {
	case "cpu":
		c.JSON(200, cpu.Load())
	case "disk":
		c.JSON(200, disk.Load())
	case "entropy":
		c.JSON(200, entropy.Load())
	case "fh":
		c.JSON(200, fh.Load())
	case "iface":
		c.JSON(200, iface.Load())
	case "inode":
		c.JSON(200, inode.Load())
	case "load":
		c.JSON(200, load.Load())
	case "mem":
		c.JSON(200, mem.Load())
	case "pty":
		c.JSON(200, pty.Load())
	case "uptime":
		c.JSON(200, uptime.Load())
	case "vmstat":
		c.JSON(200, vmstat.Load())
	default:
		c.AbortWithStatus(404)
	}
}

func main() {
	r := gin.Default()
	r.GET("/:stat", defaultRoute)

	r.Run()

}
