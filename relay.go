package main

// JSAB's Pia traffic normally goes directly between consoles. An opt-in UDP
// relay is available for a measured pair that cannot punch through NAT. Both
// players must be named in relay_allow; an absent file never enables anyone.

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	nex "github.com/NextendoNetwork/nextendo-nex"
)

const relayBase = 39000
const relaySpan = 16

func readRelayAllow(path string) (allowed, forced []uint64) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.TrimSpace(strings.SplitN(scan.Text(), "#", 2)[0])
		force := strings.HasPrefix(line, "!")
		line = strings.TrimPrefix(line, "!")
		pid, err := strconv.ParseUint(line, 10, 64)
		if err != nil || pid == 0 {
			continue
		}
		allowed = append(allowed, pid)
		if force {
			forced = append(forced, pid)
		}
	}
	return allowed, forced
}

func startRelayWatcher() {
	go func() {
		var previous string
		for {
			hostBytes, _ := os.ReadFile("/app/relay_on")
			host := strings.TrimSpace(string(hostBytes))
			if net.ParseIP(host) == nil {
				host = ""
			}
			allowBytes, _ := os.ReadFile("/app/relay_allow")
			state := host + "|" + string(allowBytes)
			if state != previous {
				previous = state
				allowed, forced := readRelayAllow("/app/relay_allow")
				nex.SetRelayVolontaires(allowed, forced, false)
				if host == "" || len(allowed) < 2 {
					nex.SetPairRelay("", 0, 0)
					fmt.Printf("[JSAB relay] off; %d player(s) allowed\n", len(allowed))
				} else {
					nex.SetPairRelay(host, relayBase, relaySpan)
					fmt.Printf("[JSAB relay] on at %s:%d-%d; %d player(s) allowed, %d forced\n",
						host, relayBase, relayBase+relaySpan-1, len(allowed), len(forced))
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
