package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	REMOTE_ADDRESS = "127.0.0.1:6969"
)

func main() {
	remote, err := net.ResolveUDPAddr("udp", REMOTE_ADDRESS)
	if err != nil {
		panic(err)
	}
	udp, err := net.DialUDP("udp", nil, remote)
	if err != nil {
		panic(err)
	}

	go func() {
		buf := make([]byte, 1000)
		for {
			n, err := udp.Read(buf)
			if err != nil {
				continue
			}
			fmt.Println(">", string(buf[:n]))
		}
	}()

	r := bufio.NewReader(os.Stdin)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			continue
		}
		str := string(line)
		str = strings.ReplaceAll(str, `\n`, "\n")
		udp.Write([]byte(str))
	}
}
