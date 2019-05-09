package main

import "github.com/gin-gonic/gin"
import "net"
import "net/http"
import "os"
import "encoding/json"
import "os/exec"
import "fmt"
import "io/ioutil"
import "bytes"
import "bufio"

var dyndnsConfig = &Config{}

type Config struct {
    User string
    Password string
    Zone string
    Domains []string
    TTL string
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

func updateZone(zone string, domain string, recordType string, ttl string, ip string)string{

	f, err := ioutil.TempFile("/tmp", "dyndns")
    if err != nil {
        return err.Error()
    }

    defer os.Remove(f.Name())
    w := bufio.NewWriter(f)

    w.WriteString(fmt.Sprintf("server localhost\n"))
    w.WriteString(fmt.Sprintf("zone %s.\n", zone))
    w.WriteString(fmt.Sprintf("update delete %s.%s %s\n", domain, zone, recordType))
    w.WriteString(fmt.Sprintf("update add %s.%s %s %s %s\n", domain, zone, ttl, recordType, ip))
    w.WriteString("send\n")

    w.Flush()
    f.Close()


    cmd := exec.Command("/usr/bin/nsupdate", f.Name())
	var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err = cmd.Run()
    if err != nil {
        return err.Error() + ": " + stderr.String()
    }
    return out.String()
}


func main(){

    dyndnsConfig.parseConfig("/root/dyndnsConfig.json")

    router := gin.Default()

    authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
        dyndnsConfig.User : dyndnsConfig.Password,
    }))

	authorized.GET("/update", func(c *gin.Context) {

		domain := c.Query("domain")
		ip := c.Query("ip")

        if isDomainValid(domain, dyndnsConfig.Domains) {
            if (validateIpV4(ip) || validateIpV6(ip)) {
                err := updateZone(dyndnsConfig.Zone, domain, "A", dyndnsConfig.TTL, ip)
                if err != "" {
                c.String(http.StatusOK, "Set record for domain: %s to ip: %s", domain, ip)
                } else {
                   c.String(http.StatusBadRequest, "%s", err)
               }
             } else {
                c.String(http.StatusBadRequest, "ip: %s ist not in a valid format", ip)
            }
        } else {
            c.String(http.StatusBadRequest, "subdomain: %s not allowed", domain)
        }
	})

    router.Run(":8080")

}
