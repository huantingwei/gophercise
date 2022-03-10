package main

import (
	"flag"
	"fmt"
	"net"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type result struct {
	ip     string
	port   int
	isOpen bool
}

func getAddrs(host string) []string {
	addrs, err := net.LookupHost(host)
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Printf("Scanning %s\n", addr)
	}
	return addrs
}

func scanPort(host string, sPort, ePort, timeout int, c chan result) bool {
	for p := sPort; p <= ePort; p++ {
		dest := net.JoinHostPort(host, strconv.Itoa(p))
		conn, err := net.DialTimeout("tcp", dest, time.Duration(timeout)*time.Second)
		if err != nil {
			// fmt.Printf("[I] %-16s port %-5d is closed\n", host, p)
			c <- result{host, p, false}
		} else {
			// fmt.Printf("[I] %-16s port %-5d is open\n", host, p)
			c <- result{host, p, true}
			// only close the connection if it's established
			conn.Close()
		}
	}

	return true
}

func main() {
	var host = flag.String("a", "google.com", "hostname to scan")
	var portRange = flag.String("p", "20-200", "port range to scan")
	var timeout = flag.Int("t", 3, "timeout before cancel connection(seconds)")
	var numThread = flag.Int("n", runtime.NumCPU()-1, "number of threads")
	flag.Parse()

	ports := strings.Split(*portRange, "-")
	sPort, _ := strconv.Atoi(ports[0])
	ePort, _ := strconv.Atoi(ports[1])

	addrs := getAddrs(*host)
	openPorts := make(chan result, (ePort-sPort+1)*len(addrs))
	numTask := (ePort - sPort + 1) / *numThread

	for _, a := range addrs {
		for i := sPort; i <= ePort; i = i + numTask {
			port := i
			host := a
			if port+numTask <= ePort {
				go func() {
					scanPort(host, port, port+numTask, *timeout, openPorts)
				}()
			} else {
				go func() {
					scanPort(host, port, ePort, *timeout, openPorts)
				}()
				break
			}
		}
	}

	reports := make([]string, (ePort-sPort+1)*len(addrs))

	for i := 0; i < (ePort-sPort+1)*len(addrs); i++ {
		p := <-openPorts
		if p.isOpen {
			reports[i] = fmt.Sprintf("%-16s port %-5d is open\n", p.ip, p.port)
			// fmt.Printf("%-16s port %-5d is open\n", p.ip, p.port)
		} else {
			reports[i] = fmt.Sprintf("%-16s port %-5d is closed\n", p.ip, p.port)
			// fmt.Printf("%-16s port %-5d is closed\n", p.ip, p.port)
		}
	}
	sort.Strings(reports)
	for _, l := range reports {
		fmt.Print(l)
	}

	fmt.Printf("\n====================\nFinished!\n")
}
