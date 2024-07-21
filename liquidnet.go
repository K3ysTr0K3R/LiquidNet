package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

func generateRandomIP() net.IP {
	ip := make(net.IP, 4)
	rand.Seed(time.Now().UnixNano())
	rand.Read(ip)
	ip[0] = byte(rand.Intn(256))
	ip[1] = byte(rand.Intn(256))
	ip[2] = byte(rand.Intn(256))
	ip[3] = byte(rand.Intn(256))
	return ip
}

func isPortOpen(ip net.IP, port int) bool {
	address := fmt.Sprintf("%s:%d", ip.String(), port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func worker(ctx context.Context, portsOpen chan<- string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ip := generateRandomIP()
			if isPortOpen(ip, port) {
				select {
				case portsOpen <- ip.String():
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func main() {
	numIPs := flag.Int("n", 1, "Number of IP addresses to gather")
	port := flag.Int("p", 80, "Port to check if open")
	numThreads := flag.Int("t", 10, "Number of threads to use")

	flag.Parse()

	portsOpen := make(chan string, *numIPs)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < *numThreads; i++ {
		wg.Add(1)
		go worker(ctx, portsOpen, *port, &wg)
	}

	for i := 0; i < *numIPs; i++ {
		fmt.Println(<-portsOpen)
	}

	cancel()
	wg.Wait()
	close(portsOpen)
}
