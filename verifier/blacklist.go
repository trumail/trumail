package verifier

import (
	"fmt"
	"net"
	"strings"
)

// ipifyBaseURL is used to retrieve our servers IP
const ipifyBaseURL = "https://api.ipify.org"

// Blacklisted is a parent blacklist checking method that
// returns true if our IP is blacklisted in any of the
// monitored blacklisting services
func (v *Verifier) Blacklisted() bool {
	return v.dnsBlacklisted(blacklists) ||
		v.matchBlacklisted("support@me.com", "proofpoint") ||
		v.matchBlacklisted("support@orange.fr", "cloudmark")
}

// dnsBlacklisted takes a list of dns blacklist addresses
// and checks each
func (v *Verifier) dnsBlacklisted(lists []string) bool {
	// Retrieve this servers IP address
	ip, err := v.client.GetString(ipifyBaseURL)
	if err != nil {
		return false
	}

	// Generate the blacklist query subdomain
	revIPArr := reverse(strings.Split(ip, "."))
	revIP := strings.Join(revIPArr, ".")

	for _, host := range lists {
		// Generate the blacklist url
		url := fmt.Sprintf("%s.%s", revIP, host)

		// Perform a host lookup and return true if found
		if addrs, _ := net.LookupHost(url); len(addrs) > 0 {
			return true
		}
	}
	return false
}

// matchBlacklisted returns a boolean value that defines
// whether or not our sending IP is blacklisted on the passed
// emails mail server using the passed selector string
func (v *Verifier) matchBlacklisted(email, selector string) bool {
	// Parse the address passed
	a, err := ParseAddress(email)
	if err != nil {
		return false
	}

	// Attempts to form an SMTP Connection to the proofpoint
	// protected mailserver
	deliverabler, err := NewDeliverabler(a.Domain, v.hostname, v.sourceAddr)
	if err != nil {
		// If the error confirms we are blocked with the selector
		le := parseRCPTErr(err)
		if le != nil && le.Message == ErrBlocked &&
			insContains(le.Details, selector) {
			return true
		}
		return false
	}

	// Checks deliverability of an arbitrary proofpoint protected
	// address
	if err := deliverabler.IsDeliverable(a.Address, 5); err != nil {
		// If the error confirms we are blocked with the selector
		le := parseRCPTErr(err)
		if le != nil && le.Message == ErrBlocked &&
			insContains(le.Details, selector) {
			return true
		}
		return false
	}
	return false
}

// reverse reverses the order of the passed in string slice,
// returning a new one of the opposite order
func reverse(in []string) []string {
	if len(in) == 0 {
		return in
	}
	return append(reverse(in[1:]), in[0])
}

// blacklists contains all the blacklist hostnames we want to check
var blacklists = []string{
	"zen.spamhaus.org",
	"xbl.spamhaus.org",
	"pbl.spamhaus.org",
	"sbl-xbl.spamhaus.org",
	"sbl.spamhaus.org",
	"all.spamrats.com",
	"noptr.spamrats.com",
	"spam.spamrats.com",
	"dyna.spamrats.com",
	"bl.spamstinks.com",
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
	// "sbl.nszones.com",
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
