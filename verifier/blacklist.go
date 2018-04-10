package verifier

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

// Blacklisted is a parent blacklist checking method that
// returns true if our IP is blacklisted in any of the
// monitored blacklisting services
func (v *Verifier) Blacklisted() bool {
	// Check Proofpoint blacklist
	pb, err := v.proofpointBlacklisted()
	if err == nil && pb {
		return true
	}

	// Check Spamhaus blacklist
	sb, err := v.spamhausBlacklisted()
	if err == nil && sb {
		return true
	}
	return false
}

// proofpointBlacklisted returns a boolean value that defines
// whether or not our sending IP is blacklisted by Proofpoint
func (v *Verifier) proofpointBlacklisted() (bool, error) {
	// Attempts to form an SMTP Connection to the proofpoint
	// protected mailserver
	deliverabler, err := NewDeliverabler("me.com", v.hostname,
		v.sourceAddr, v.client.Timeout)
	if err != nil {
		// If the error confirms we are blocked by proofpoint return true
		basicErr, detailErr := parseRCPTErr(err)
		if basicErr == ErrBlocked && insContains(detailErr.Error(), "proofpoint") {
			return true, nil
		}
		return false, err
	}

	// Checks deliverability of an arbitrary proofpoint protected
	// address
	if err := deliverabler.IsDeliverable("support@me.com", 5); err != nil {
		// If the error confirms we are blocked by proofpoint return true
		basicErr, detailErr := parseRCPTErr(err)
		if basicErr == ErrBlocked && insContains(detailErr.Error(), "proofpoint") {
			return true, nil
		}
	}
	return false, nil
}

// spamhausBlacklisted returns a boolean value defining whether or
// not our sending IP is blacklisted by Spamhaus
func (v *Verifier) spamhausBlacklisted() (bool, error) {
	// Retrieve the servers IP address
	ip, err := v.ipify()
	if err != nil {
		return false, err
	}

	// Generate the Spamhaus query hostname
	revIP := strings.Join(reverse(strings.Split(ip, ".")), ".")
	url := fmt.Sprintf("%s.zen.spamhaus.org", revIP)

	// Perform a host lookup and return true if the
	// host is found to blacklisted by Spamhaus
	addrs, err := net.LookupHost(url)
	if err != nil {
		return false, nil
	}
	return len(addrs) > 0, nil
}

// ipify performs a GET request to IPify to retrieve
// our IP address
func (v *Verifier) ipify() (string, error) {
	bytes, err := v.get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// get performs a get request using the passed URL
func (v *Verifier) get(url string) ([]byte, error) {
	// Perform the GET request on the passed URL
	resp, err := v.client.Get(url)
	if err != nil {
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
