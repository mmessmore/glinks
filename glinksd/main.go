package main

import "github.com/gin-gonic/gin"

func defaultRoute(c *gin.Context) {

}

func main() {
	r := gin.Default()
	r.GET("/:stat", defaultRoute)

	r.Run()

}
