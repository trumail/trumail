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
	// Retrieve this servers IP address
	ip, err := v.ipify()
	if err != nil {
		return false, err
	}

	// Generate the Spamhaus query hostname
	revIP := strings.Join(reverse(strings.Split(ip, ".")), ".")

	for _, host := range blacklists {
		// Generate the blacklist url
		url := fmt.Sprintf("%s.%s", revIP, host)

		// Perform a host lookup and return true if the
		// host is found to blacklisted by Spamhaus
		addrs, _ := net.LookupHost(url)
		if len(addrs) > 0 {
			return true, nil
		}
	}
	return false, nil
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

var blacklists = []string{
	"zen.spamhaus.org",
	"xbl.spamhaus.org",
	"all.spamrats.com",
	"noptr.spamrats.com",
	"spam.spamrats.com",
	// "0spam-killlist.fusionzero.com",
	// "0spam.fusionzero.com",
	// "access.redhawk.org",
	// "all.rbl.jp",
	// "all.spam-rbl.fr",
	// "aspews.ext.sorbs.net",
	// "b.barracudacentral.org",
	// "backscatter.spameatingmonkey.net",
	// "badnets.spameatingmonkey.net",
	// "bb.barracudacentral.org",
	// "bl.drmx.org",
	// "bl.konstant.no",
	// "bl.nszones.com",
	// "bl.spamcannibal.org",
	// "bl.spameatingmonkey.net",
	// "bl.spamstinks.com",
	// "black.junkemailfilter.com",
	// "blackholes.five-ten-sg.com",
	// "blacklist.sci.kun.nl",
	// "blacklist.woody.ch",
	// "bogons.cymru.com",
	// "bsb.empty.us",
	// "bsb.spamlookup.net",
	// "cart00ney.surriel.com",
	// "cbl.abuseat.org",
	// "cbl.anti-spam.org.cn",
	// "cblless.anti-spam.org.cn",
	// "cblplus.anti-spam.org.cn",
	// "cdl.anti-spam.org.cn",
	// "cidr.bl.mcafee.com",
	// "combined.rbl.msrbl.net",
	// "db.wpbl.info",
	// "dev.null.dk",
	// "dialups.visi.com",
	// "dnsbl-0.uceprotect.net",
	// "dnsbl-1.uceprotect.net",
	// "dnsbl-2.uceprotect.net",
	// "dnsbl-3.uceprotect.net",
	// "dnsbl.anticaptcha.net",
	// "dnsbl.aspnet.hu",
	// "dnsbl.inps.de",
	// "dnsbl.justspam.org",
	// "dnsbl.kempt.net",
	// "dnsbl.madavi.de",
	// "dnsbl.rizon.net",
	// "dnsbl.rv-soft.info",
	// "dnsbl.rymsho.ru",
	// "dnsbl.sorbs.net",
	// "dnsbl.zapbl.net",
	// "dnsrbl.swinog.ch",
	// "dul.pacifier.net",
	// "dyn.nszones.com",
	// "dyna.spamrats.com",
	// "fnrbl.fast.net",
	// "fresh.spameatingmonkey.net",
	// "hostkarma.junkemailfilter.com",
	// "images.rbl.msrbl.net",
	// "ips.backscatterer.org",
	// "ix.dnsbl.manitu.net",
	// "korea.services.net",
	// "l2.bbfh.ext.sorbs.net",
	// "l3.bbfh.ext.sorbs.net",
	// "l4.bbfh.ext.sorbs.net",
	// "list.bbfh.org",
	// "list.blogspambl.com",
	// "mail-abuse.blacklist.jippg.org",
	// "netbl.spameatingmonkey.net",
	// "netscan.rbl.blockedservers.com",
	// "no-more-funn.moensted.dk",
	// "orvedb.aupads.org",
	// "pbl.spamhaus.org",
	// "phishing.rbl.msrbl.net",
	// "pofon.foobar.hu",
	// "psbl.surriel.com",
	// "rbl.abuse.ro",
	// "rbl.blockedservers.com",
	// "rbl.dns-servicios.com",
	// "rbl.efnet.org",
	// "rbl.efnetrbl.org",
	// "rbl.iprange.net",
	// "rbl.schulte.org",
	// "rbl.talkactive.net",
	// "rbl2.triumf.ca",
	// "rsbl.aupads.org",
	// "sbl-xbl.spamhaus.org",
	// "sbl.nszones.com",
	// "sbl.spamhaus.org",
	// "short.rbl.jp",
	// "spam.dnsbl.anonmails.de",
	// "spam.pedantic.org",
	// "spam.rbl.blockedservers.com",
	// "spam.rbl.msrbl.net",
	// "spamrbl.imp.ch",
	// "spamsources.fabel.dk",
	// "st.technovision.dk",
	// "tor.dan.me.uk",
	// "tor.dnsbl.sectoor.de",
	// "tor.efnet.org",
	// "torexit.dan.me.uk",
	// "truncate.gbudb.net",
	// "ubl.unsubscore.com",
	// "uribl.spameatingmonkey.net",
	// "urired.spameatingmonkey.net",
	// "virbl.dnsbl.bit.nl",
	// "virus.rbl.jp",
	// "virus.rbl.msrbl.net",
	// "vote.drbl.caravan.ru",
	// "vote.drbl.gremlin.ru",
	// "web.rbl.msrbl.net",
	// "work.drbl.caravan.ru",
	// "work.drbl.gremlin.ru",
	// "wormrbl.imp.ch",
}
