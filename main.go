package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/akamensky/argparse"
)

/**
	read subdomain wordlist and return a string list
**/

func readFile(fileName string) []string {
	wordlist, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal("Could not read wordlist")
	}

	lines := strings.Split(string(wordlist), "\n")

	return lines
}

/**
	make http request to target with subdomain wordlist
**/

func makeRequest(wordlist []string, domain string, ssl *bool, writeOut *bool, fileName *string, port *string) {

	bad_domain_counter := 0
	protocol := "http"
	var out []string

	if *ssl == true {
		protocol = "https"
	}

	for i := 0; i < len(wordlist); i++ {
		resp, err := http.Get(protocol + "://" + wordlist[i] + "." + domain + ":" + *port + "/")
		if err != nil && strings.Contains(err.Error(), "no such host") {
			bad_domain_counter += 1
			continue
		} else if resp.Status == "200 OK" {
			// fmt.Println(string("\033[32m"), protocol, "://"+wordlist[i], ".", domain, ":", *port, "/", ": 200 OK")
			result := fmt.Sprintf("\033[32m%s%s%s%s%s%s%s%s%s", protocol, "://", wordlist[i], ".", domain, ":", *port, "/", " | 200 OK")
			fmt.Println(result)
			if *writeOut == true {
				out = append(out, protocol+"://"+wordlist[i]+"."+domain+":"+*port+"/")
			}
		}
	}

	if *writeOut == true {
		writeOutPutToFile(out, fileName)
	}

	fmt.Println(string("\033[31m"), "Bad domains: ", bad_domain_counter)
}

func writeOutPutToFile(outList []string, fileName *string) {

	file, err := os.Create(*fileName)
	if err != nil {
		log.Fatal("Could not create file: ", fileName)
	}

	//close created file
	defer file.Close()

	for i := 0; i < len(outList); i++ {
		_, err := file.WriteString(outList[i] + "\n")
		if err != nil {
			log.Fatalf("Could not write datat to file %s", err)
		}
	}

	fmt.Println(string("\033[33m"), "Wrote data to:", *fileName)
}

func main() {
	fmt.Println("Welcome to my little web fuzzer")

	parser := argparse.NewParser("FUZZTHINGS", "This is a simple tool i'm building to learn Golang")
	domain := parser.String("d", "domain", &argparse.Options{Required: true, Help: "domain to fuzz"})
	subdomainsList := parser.String("w", "wordlist", &argparse.Options{Required: true, Help: "subdomains file"})
	fileName := parser.String("f", "filename", &argparse.Options{Required: false, Help: "File name to write to"})
	port := parser.String("p", "port", &argparse.Options{Required: false, Help: "Specify port for web server"})

	var ssl *bool = parser.Flag("s", "ssl", &argparse.Options{Required: false, Help: "Add flag if it's https"})
	var writeOut *bool = parser.Flag("o", "out", &argparse.Options{Required: false, Help: "Write output to file"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
	}

	if len(*fileName) == 0 {
		*fileName = "output.txt"
	}

	if len(*port) == 0 {
		*port = "80"
	} else if *ssl == true {
		*port = "443"
	}

	subdomains := readFile(*subdomainsList)
	makeRequest(subdomains, *domain, ssl, writeOut, fileName, port)
}
