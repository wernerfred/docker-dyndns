package main

import "github.com/gin-gonic/gin"
import "net"
import "net/http"


func validateIpV4(ipV4 string) bool {
    v4addr := net.ParseIP(ipV4)
    if v4addr == nil {
        return false
    }
    return (v4addr.To4() != nil)
}


func validateIpV6(ipV6 string) bool {
    v6addr := net.ParseIP(ipV6)
    if v6addr == nil {
        return false
    }
    return (v6addr.To16() != nil)
}

func main(){

    router := gin.Default()

    authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
        "user": "test",
    }))

	authorized.GET("/update", func(c *gin.Context) {

		domain := c.Query("domain")
		ip := c.Query("ip")

        if (validateIpV4(ip) || validateIpV6(ip)) {
            c.String(http.StatusOK, "domain: %s; ip: %s", domain, ip)
        } else {
            c.String(http.StatusBadRequest, "ip: %s ist not in a valid format", ip)
        }
	})

    router.Run(":8080")

}
