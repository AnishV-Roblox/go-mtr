package mtr

import (
	"fmt"
	"strings"
	"testing"
)

func TestECMPProcessing(t *testing.T) {
	// Full MTR raw output with ECMP
	rawOutput := `x 0 33000
x 1 33001
h 1 1.2.3.87
h 1 1.2.3.87
p 1 267 33001
x 2 33002
h 2 1.2.3.45
h 2 1.2.3.45
p 2 213 33002
x 3 33003
h 3 1.2.5.6
h 3 1.2.5.6
p 3 286 33003
x 4 33004
h 4 1.6.9.10
h 4 1.6.9.10
p 4 201 33004
x 5 33005
h 5 1.6.5.154
h 5 1.6.5.154
p 5 98 33005
x 6 33006
h 6 1.6.5.154
h 6 1.6.5.154
d 6 1.6.5.154
p 6 80 33006
x 0 33007
x 1 33008
h 1 1.2.3.81
h 1 1.2.3.81
p 1 191 33008
x 2 33009
h 2 1.2.3.7
h 2 1.2.3.7
p 2 1928 33009
x 3 33010
h 3 1.2.5.12
h 3 1.2.5.12
p 3 149 33010
x 4 33011
d 4 1.6.9.10
p 4 233 33011
x 5 33012
d 5 1.6.5.154
p 5 88 33012
x 0 33013
x 1 33014
h 1 1.2.3.83
h 1 1.2.3.83
p 1 171 33014
x 2 33015
d 2 1.2.3.7
p 2 804 33015
x 3 33016
h 3 1.2.3.50
h 3 1.2.3.50
p 3 229 33016
x 4 33017
p 4 201 33017
x 5 33018
p 5 126 33018
x 0 33019
x 1 33020
h 1 1.2.3.81
d 1 1.2.3.81
p 1 217 33020
x 2 33021
h 2 1.2.3.19
h 2 1.2.3.19
p 2 18393 33021
x 3 33022
h 3 1.2.5.12
d 3 1.2.5.12
p 3 245 33022
x 4 33023
p 4 216 33023
x 5 33024
p 5 142 33024
x 0 33025
x 1 33026
h 1 1.2.3.83
p 1 252 33026
x 2 33027
h 2 1.2.3.5
h 2 1.2.3.5
p 2 640 33027
x 3 33028
h 3 1.2.5.4
h 3 1.2.5.4
p 3 307 33028
x 4 33029
p 4 325 33029
x 5 33030
p 5 238 33030
x 0 33031
x 1 33032
h 1 1.2.3.81
p 1 274 33032
x 2 33033
h 2 1.2.3.19
p 2 1200 33033
x 3 33034
h 3 1.2.5.12
p 3 237 33034
x 4 33035
p 4 244 33035
x 5 33036
p 5 168 33036
x 0 33037
x 1 33038
p 1 285 33038
x 2 33039
h 2 1.2.3.5
p 2 1251 33039
x 3 33040
p 3 171 33040
x 4 33041
p 4 305 33041
x 5 33042
p 5 202 33042
x 0 33043
x 1 33044
p 1 191 33044
x 2 33045
h 2 1.2.3.45
p 2 239 33045
x 3 33046
h 3 1.2.3.50
p 3 184 33046
x 4 33047
p 4 232 33047
x 5 33048
p 5 133 33048
x 0 33049
x 1 33050
h 1 1.2.3.83
p 1 172 33050
x 2 33051
h 2 1.2.3.7
p 2 408 33051
x 3 33052
p 3 162 33052
x 4 33053
p 4 203 33053
x 5 33054
p 5 130 33054
x 0 33055
x 1 33056
p 1 206 33056
x 2 33057
h 2 1.2.3.45
p 2 286 33057
x 3 33058
p 3 326 33058
x 4 33059
p 4 200 33059
x 5 33060
p 5 225 33060
`

	m := &MTR{
		OutputRaw:   []byte(rawOutput),
		PacketsSent: 10,
	}

	// Process the output
	m.processOutput()

	if m.Error != nil {
		t.Fatalf("processOutput failed: %v", m.Error)
	}

	// Print out what we found
	separator := strings.Repeat("=", 70)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("MTR Test Results - ECMP Network\n")
	fmt.Printf("%s\n\n", separator)
	fmt.Printf("Total hops found: %d\n", len(m.Hops))
	fmt.Printf("Packets sent per hop: %d\n\n", m.PacketsSent)

	// Print each hop with MTR-style output
	for _, hop := range m.Hops {
		fmt.Printf("Hop %d: Loss=%.1f%% Snt=%d Rcv=%d Avg=%.1fµs Best=%dµs Worst=%dµs\n",
			hop.HopNumber,
			hop.LostPercent*100,
			hop.Sent,
			hop.Received,
			hop.Mean,
			hop.Best,
			hop.Worst,
		)

		// Print all IPs at this hop (ECMP)
		fmt.Printf("  Hosts (%d):\n", len(hop.Hosts))
		for _, host := range hop.Hosts {
			hostname := host.Name
			if hostname == "" {
				hostname = "(no DNS)"
			}
			fmt.Printf("    - %-15s  %s\n", host.IP.String(), hostname)
		}
		fmt.Println()
	}

	// Verify hop 1 has ECMP (3 different IPs)
	var hop1 *Hop
	for _, hop := range m.Hops {
		if hop.HopNumber == 1 {
			hop1 = hop
			break
		}
	}

	if hop1 == nil {
		t.Fatal("Hop 1 not found")
	}

	expectedHop1IPs := map[string]bool{
		"1.2.3.87": false,
		"1.2.3.81": false,
		"1.2.3.83": false,
	}

	if len(hop1.Hosts) != len(expectedHop1IPs) {
		t.Errorf("Expected %d hosts at hop 1 (ECMP), got %d", len(expectedHop1IPs), len(hop1.Hosts))
	}

	for _, host := range hop1.Hosts {
		ipStr := host.IP.String()
		if _, expected := expectedHop1IPs[ipStr]; expected {
			expectedHop1IPs[ipStr] = true
		} else {
			t.Errorf("Unexpected IP at hop 1: %s", ipStr)
		}
	}

	for ip, found := range expectedHop1IPs {
		if !found {
			t.Errorf("Expected IP %s not found at hop 1", ip)
		}
	}

	// Verify statistics are aggregated at hop level
	if hop1.Received != len(hop1.PacketMicrosecs) {
		t.Errorf("Hop 1 Received (%d) doesn't match packet count (%d)",
			hop1.Received, len(hop1.PacketMicrosecs))
	}

	// Verify loss calculation
	expectedLoss := float64(hop1.Dropped) / float64(hop1.Sent)
	if hop1.LostPercent != expectedLoss {
		t.Errorf("Hop 1 loss percent incorrect: got %.2f, expected %.2f",
			hop1.LostPercent, expectedLoss)
	}

	// Verify hop 2 also has ECMP
	var hop2 *Hop
	for _, hop := range m.Hops {
		if hop.HopNumber == 2 {
			hop2 = hop
			break
		}
	}

	if hop2 == nil {
		t.Fatal("Hop 2 not found")
	}

	expectedHop2IPs := []string{
		"1.2.3.45",
		"1.2.3.7",
		"1.2.3.19",
		"1.2.3.5",
	}

	if len(hop2.Hosts) != len(expectedHop2IPs) {
		t.Errorf("Expected %d hosts at hop 2 (ECMP), got %d", len(expectedHop2IPs), len(hop2.Hosts))
	}

	fmt.Printf("%s\n", separator)
	fmt.Printf("Test Summary:\n")
	fmt.Printf("  ✓ Parsed %d hops\n", len(m.Hops))
	fmt.Printf("  ✓ Hop 1 has %d ECMP paths\n", len(hop1.Hosts))
	fmt.Printf("  ✓ Hop 2 has %d ECMP paths\n", len(hop2.Hosts))
	fmt.Printf("  ✓ Statistics aggregated at hop level\n")
	fmt.Printf("  ✓ Loss calculation: Hop 1 = %.1f%%, Hop 2 = %.1f%%\n",
		hop1.LostPercent*100, hop2.LostPercent*100)
	fmt.Printf("%s\n\n", separator)
}
