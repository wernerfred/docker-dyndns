package main

import "github.com/gin-gonic/gin"
import "net/http"

func main(){

    router := gin.Default()

	router.GET("/update", func(c *gin.Context) {

		domain := c.Query("domain")
		ip := c.Query("ip")

        c.String(http.StatusOK, "domain: %s; ip: %s", domain, ip)
	})
    router.Run(":8080")

}
