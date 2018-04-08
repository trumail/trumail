package spamhaus

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

const (
	spamhausHost = "zen.spamhaus.org"
	ipifyURL     = "https://api.ipify.org"
)

// Blocked returns a boolean value defining whether or
// not our sending IP is blocked by Spamhaus
func Blocked() (bool, error) {
	// Retrieve the servers IP address
	ip, err := IPify()
	if err != nil {
		return false, err
	}

	// Generate the Spamhaus query hostname
	revIP := strings.Join(reverse(strings.Split(ip, ".")), ".")
	url := revIP + "." + spamhausHost

	// Perform a host lookup and return true if the
	// host is found to blocked by Spamhaus
	addrs, err := net.LookupHost(url)
	if err != nil {
		return false, nil
	}
	return len(addrs) > 0, nil
}

// IPify performs a GET request to IPify to retrieve
// our IP address
func IPify() (string, error) {
	bytes, err := get(ipifyURL)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// get performs a get request using the passed URL
func get(url string) ([]byte, error) {
	// Perform the GET request on the passed URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to retrieve IP from api.ipify.org")
		return nil, err
	}
	defer resp.Body.Close()

	// Decode and return the bytes
	return ioutil.ReadAll(resp.Body)
}

// reverse reverses the order of the passed in string
// slice, returning a new one of the opposite order
func reverse(in []string) []string {
	if len(in) == 0 {
		return in
	}
	return append(reverse(in[1:]), in[0])
}
