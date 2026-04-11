package service

import (
	"bufio"
	"net"
	"strings"

	"github.com/hellinvoid/network-programming/common"
)

func ProxyRequest(in, out net.Conn) {
	defer in.Close()
	defer out.Close()

	r := bufio.NewReader(in)
	for {
		// Read a line
		msg, err := r.ReadString('\n')
		if err != nil {
			return
		}
		msg = alterBoguscoinAddress(msg)
		out.Write([]byte(msg))
	}
}

func alterBoguscoinAddress(msg string) string {
	const (
		MIN_BOGUSCOIN_ADDRESS_SIZE = 26
		MAX_BOGUSCOIN_ADDRESS_SIZE = 35
		TONY_BOGUSCOIN_ADDRESS     = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	)

	temp := msg[:len(msg)-1]

	// Split the message byt " " and check for each string if its a valid address
	for str := range strings.SplitSeq(temp, " ") {
		if isValidAddress(str) {
			// Replace with Tony's address if valid
			msg = strings.ReplaceAll(msg, str, TONY_BOGUSCOIN_ADDRESS)
		}
	}

	return msg
}

func isValidAddress(address string) bool {
	return len(address) >= 26 && len(address) <= 35 && common.IsAlnum(address) && address[0] == '7'
}
