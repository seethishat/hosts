package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var debug = false

// OkayChars is a string of ASCII chars that are acceptable as part of a domain name.
var OkayChars = "abcdefghijklmnopqrstuvwxyz0123456789.-_"

// validTLD ensures a domain ends with a valid Top Level Domain (.com, .net, .cc, etc.).
// It takes a domain name (example.com) and a slice of TLDs (.com, .net, .cc) as argmuents
// and returns true or false.
func validTLD(domain string, tlds []string) bool {
	cleanDomain := strings.ToLower(strings.TrimSpace(domain))
	for _, tld := range tlds {
		if strings.HasSuffix(cleanDomain, tld) {
			return true
		}
	}
	return false
}

// loadTLDs takes a file with a list of TLDS from IANA, normalizes them and returns them as a slice.
// https://data.iana.org/TLD/tlds-alpha-by-domain.txt
// The data file is updated periodically and has this format:
// # Version 2019012900, Last Updated Tue Jan 29 07:07:01 2019 UTC
// AAA
// AARP
// ...
// ZW
func loadTLDs(fileName string) []string {
	tlds := make([]string, 0)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tld := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if !strings.HasPrefix(tld, "#") {
			tlds = append(tlds, "."+tld)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if debug {
		log.Printf("Loaded %d TLDs.", len(tlds))
	}

	return tlds
}

func loadDomains(fileName string) []string {
	domains := make([]string, 0)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.ToLower(strings.TrimSpace(scanner.Text()))
		domains = append(domains, domain)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if debug {
		log.Printf("Loaded %d Domains.", len(domains))
	}

	return domains
}

// validDomain ensures a domain is valid in the context of a DNS CNAME resource record.
func validDomain(domain string) bool {
	cleanDomain := strings.ToLower(strings.TrimSpace(domain))

	if len(cleanDomain) < 4 || len(cleanDomain) > 253 {
		return false
	}

	if strings.HasPrefix(cleanDomain, "-") || strings.HasPrefix(cleanDomain, "_") || strings.HasPrefix(cleanDomain, ".") {
		return false
	}

	if strings.HasSuffix(cleanDomain, "-") || strings.HasSuffix(cleanDomain, "_") || strings.HasSuffix(cleanDomain, ".") {
		return false
	}

	if !strings.ContainsAny(cleanDomain, ".") {
		return false
	}

	for _, r := range cleanDomain {
		if !strings.Contains(OkayChars, string(r)) {
			return false
		}
	}

	labels := strings.Split(cleanDomain, ".")
	for _, label := range labels {
		if len(label) > 63 || len(label) < 1 {
			return false
		}
	}

	return true
}

func main() {
	var helpFlag = flag.Bool("help", false, "show this help message.")
	var debugFlag = flag.Bool("debug", false, "enable debug logging.")
	var tlds = flag.String("tlds", "tlds-alpha-by-domain.txt", "the path to the IANA TLD list.")
	var domains = flag.String("domains", "", "the path to the list of domains.")

	flag.Parse()

	if *debugFlag {
		debug = true
	}

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	t := loadTLDs(*tlds)
	d := loadDomains(*domains)

	for _, domain := range d {
		fmt.Printf("Valid TLD? %t, %s\n", validTLD(domain, t), domain)
		fmt.Printf("Valid Domain? %t, %s\n", validDomain(domain), domain)
	}
}
