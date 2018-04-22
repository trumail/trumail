package verifier

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"golang.org/x/sync/errgroup"
)

// Blacklisted is a parent blacklist checking method that
// returns true if our IP is blacklisted in any of the
// monitored blacklisting services
func (v *Verifier) Blacklisted() error {
	var g errgroup.Group
	g.Go(func() error { return v.dnsBlacklisted(blacklists) })
	g.Go(func() error { return v.matchBlacklisted("support@me.com", "proofpoint") }) // Proofpoint
	// g.Go(func() error { return v.matchBlacklisted("support@bath.ac.uk", "outlook.com") })        // Outlook
	g.Go(func() error { return v.matchBlacklisted("support@orange.fr", "cloudmark") })           // Cloudmark
	g.Go(func() error { return v.matchBlacklisted("support@subaru.com.au", "trend micro rbl") }) // Trend Micro RBL+
	return g.Wait()
}

// dnsBlacklisted takes a list of dns blacklist addresses
// and checks each
func (v *Verifier) dnsBlacklisted(lists []string) error {
	// Retrieve this servers IP address
	ip, err := v.client.GetString("https://api.ipify.org")
	if err != nil {
		return nil
	}

	// Generate the blacklist query subdomain
	revIPArr := reverse(strings.Split(ip, "."))
	revIP := strings.Join(revIPArr, ".")

	// Perform a DNS lookup on all the lists
	for _, host := range lists {
		// Generate the blacklist url
		url := fmt.Sprintf("%s.%s", revIP, host)

		// Perform a host lookup and return true if found
		if addrs, _ := net.LookupHost(url); len(addrs) > 0 {
			return errors.New("Blocked by " + host)
		}
	}
	return nil
}

// matchBlacklisted returns a boolean value that defines
// whether or not our sending IP is blacklisted on the
// passed emails mail server using the passed selector
// string
func (v *Verifier) matchBlacklisted(email, selector string) error {
	// Perform a lookup on the email
	if _, err := v.Verify(email); err != nil {
		// If the error confirms we are blocked with the selector
		le := parseSMTPError(err)
		if le != nil && le.Message == ErrBlocked &&
			insContains(le.Details, selector) {
			return errors.New("Blocked by " + selector)
		}
		return nil
	}
	return nil
}

// reverse reverses the order of the passed in string
// slice, returning a new one of the opposite order
func reverse(in []string) []string {
	if len(in) == 0 {
		return in
	}
	return append(reverse(in[1:]), in[0])
}

// blacklists contains all the blacklist hostnames we
// want to check
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
	// "bl.drmx.org",
	// "bl.konstant.no",
	// "bl.nszones.com",
	// "bl.spamcannibal.org",
	// "bl.spameatingmonkey.net",
	// "black.junkemailfilter.com",
	// "rbl.abuse.ro",
	// "bl.spamstinks.com",
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
