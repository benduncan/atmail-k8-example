package main

// EXAMPLE atmail-rbl service for Kubernetes tutorial

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// JSON response for query
type Resp struct {
	RblServer string
	IsMatch   bool
	IPs       []net.IP
}

func main() {

	// Setup our gin routes
	r := gin.Default()

	// Bind to the specified port
	port := os.Getenv("API_PORT")

	// Defualt to port 8001
	if port == "" {
		port = "8001"
	}

	// Load which RBL servers to query, otherwise fallback to the defaults
	rblHost := os.Getenv("RBL_DNS_LOOKUP")

	if rblHost == "" {
		rblHost = "zen.spamhaus.org,bl.score.senderscore.com,b.barracudacentral.org,bl.spamcop.net"
	}

	// The 'crux' of the micro-service, allow to query an IP via a HTTP get
	r.GET("/query/:ip", func(c *gin.Context) {

		ip := net.ParseIP(c.Param("ip"))

		// RBL servers require the IP to be reversed as per the in-addr.arpa format
		reverseIP, err := ReverseIPAddress(ip)

		// Return an error if the input is malformed
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rblServers := strings.Split(rblHost, ",")

		// Build our JSON response
		var jsonResp []Resp
		jsonResp = make([]Resp, len(rblServers))

		// Loop through each query
		for i := range rblServers {

			// Query the reversed IP to the RBL host
			query := fmt.Sprintf("%s.%s.", reverseIP, rblServers[i])

			// Run the DNS query
			ips, err := net.LookupIP(query)

			// Append the output to our JSON object
			jsonResp[i].RblServer = rblServers[i]
			jsonResp[i].IPs = ips

			// Flag is the IP is matched on an RBL server
			if err != nil {
				jsonResp[i].IsMatch = false
			} else {
				jsonResp[i].IsMatch = true
			}

		}

		c.JSON(http.StatusOK, jsonResp)

	})

	// Health for K8 and container checks
	r.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, "OK")
	})

	// Run the gin service on our specified port
	r.Run(fmt.Sprintf(":%s", port))

}

// ReverseIPAddress Reverse an IP address for lookups
func ReverseIPAddress(ip net.IP) (reverseIP string, err error) {

	if ip.To4() != nil {
		// Split into slice by dot .
		addressSlice := strings.Split(ip.String(), ".")
		reverseSlice := []string{}

		for i := range addressSlice {
			octet := addressSlice[len(addressSlice)-1-i]
			reverseSlice = append(reverseSlice, octet)
		}

		return strings.Join(reverseSlice, "."), err

	} else {
		return "", errors.New("Invalid IPv4 address")
	}

}
