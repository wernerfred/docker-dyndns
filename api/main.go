package main

import "github.com/gin-gonic/gin"
import "net"
import "net/http"
import "os"
import "encoding/json"

var dyndnsConfig = &Config{}

type Config struct {
    User string
    Password string
    Domains []string
}

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

func (conf *Config) parseConfig(pathToConfig string) {
	file, err := os.Open(pathToConfig)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
}

func isDomainValid(domain string, domains []string) bool {
    for _, cur := range domains {
        if cur == domain {
            return true
        }
    }
    return false
}

func main(){

    dyndnsConfig.parseConfig("/tmp/dyndnsConfig.json")

    router := gin.Default()

    authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
        dyndnsConfig.User : dyndnsConfig.Password,
    }))

	authorized.GET("/update", func(c *gin.Context) {

		domain := c.Query("domain")
		ip := c.Query("ip")

        if isDomainValid(domain, dyndnsConfig.Domains) {
            if (validateIpV4(ip) || validateIpV6(ip)) {
                c.String(http.StatusOK, "domain: %s; ip: %s", domain, ip)
             } else {
                c.String(http.StatusBadRequest, "ip: %s ist not in a valid format", ip)
            }
        } else {
            c.String(http.StatusBadRequest, "subdomain: %s not allowed", domain)
        }
	})

    router.Run(":8080")

}
