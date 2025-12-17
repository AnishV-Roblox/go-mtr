// mtr is a package that processes and saves
// an `mtr --raw` call's output.
package mtr

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os/exec"
	"strconv"
)

type Host struct {
	IP   net.IP `json:"ip"`
	Name string `json:"hostname"`
}

type Hop struct {
	Hosts           []*Host `json:"hosts"`
	HopNumber       int     `json:"hop-number"`
	PacketMicrosecs []int   `json:"packet-times"`
	Sent            int     `json:"sent"`
	Received        int     `json:"received"`
	Dropped         int     `json:"dropped"`
	LostPercent     float64 `json:"lost-percent"`
	// All packet units are in microseconds
	Mean               float64 `json:"mean"`
	Best               int     `json:"best"`
	Worst              int     `json:"worst"`
	StandardDev        float64 `json:"standard-dev"`
	MeanJitter         float64 `json:"mean-jitter"`
	WorstJitter        int     `json:"worst-jitter"`
	InterarrivalJitter int     `json:"interarrival-jitter"` // calculated with rfc3550 A.8 shortcut
}

type MTR struct {
	Done        chan struct{}
	OutputRaw   []byte
	Error       error
	PacketsSent int
	Hops        []*Hop `json:"hops"`
}

// New runs mtr --raw -c reportCycles hostname args... Thus, you can add more arguments
// to the default required hostname and report cycles. The MTR call is signified as done
// when the MTR.Done chan closes. First wait for this, then check the MTR.Error field before
// looking at the output. Other than that, the fields and json tags document what everything
// means.
func New(reportCycles int, host string, args ...string) *MTR {
	m := &MTR{Done: make(chan struct{}), PacketsSent: reportCycles}
	args = append([]string{"--raw", "-c", strconv.Itoa(reportCycles), host}, args...)
	go func() {
		defer close(m.Done)
		m.OutputRaw, m.Error = exec.Command("mtr", args...).Output()
		if m.Error == nil {
			m.processOutput()
		}
	}()
	return m
}

func parseByteNum(input []byte) int {
	i := 0
	for _, v := range input {
		i = 10*i + int(v) - 48 // ascii 48 is `0`
	}
	return i
}

func parseHopNumber(line []byte) (num int, finalFieldIdx int) {
	finalFieldIdx = bytes.IndexByte(line[2:], ' ') + 2
	return parseByteNum(line[2:finalFieldIdx]), finalFieldIdx + 1 // `c ### <content>`
}

// Wait waits for the MTR call to complete and returns the error if any.
func (m *MTR) Wait() error {
	<-m.Done
	return m.Error
}

func (m *MTR) processOutput() {
	// h (host): host #, ip address
	// d (dns): host #, resolved dns name
	// p (packet): host #, microseconds
	defer func() {
		if x := recover(); x != nil {
			m.Error = fmt.Errorf("Unable to process output, error %v, meaning the mtr output is WHACK! Check the output to see what's up!", x)
		}
	}()

	// Track which hop index corresponds to which hop number
	hopNumToIdx := make(map[int]int)

	// Track which IPs we've seen at each hop for deduplication
	// Key format: "hopnum:ip"
	seenHostAtHop := make(map[string]bool)

	output := m.OutputRaw
	output = append(output, ' ') // tack on a space at the end so that `output = output[lineIdx+1:] doesn't panic on last newline
	for {
		lineIdx := bytes.IndexByte(output, '\n')
		if lineIdx == -1 {
			break
		}
		line := output[:lineIdx]
		output = output[lineIdx+1:]

		hopnum, finalFieldIdx := parseHopNumber(line)

		switch line[0] {
		case 'h':
			// Get or create the hop
			hopIdx, exists := hopNumToIdx[hopnum]
			if !exists {
				hop := &Hop{
					HopNumber: hopnum,
					Hosts:     []*Host{},
				}
				hopIdx = len(m.Hops)
				m.Hops = append(m.Hops, hop)
				hopNumToIdx[hopnum] = hopIdx
			}

			// Add this IP to the hop's hosts if not already seen
			ipStr := string(line[finalFieldIdx:])
			key := fmt.Sprintf("%d:%s", hopnum, ipStr)
			if !seenHostAtHop[key] {
				seenHostAtHop[key] = true
				m.Hops[hopIdx].Hosts = append(m.Hops[hopIdx].Hosts, &Host{
					IP: net.ParseIP(ipStr),
				})
			}

		case 'd':
			// Find the hop and update the hostname for the matching IP
			if hopIdx, ok := hopNumToIdx[hopnum]; ok {
				hostname := string(line[finalFieldIdx:])
				// Find the host with matching IP or hostname and update it
				for i := range m.Hops[hopIdx].Hosts {
					if m.Hops[hopIdx].Hosts[i].IP.String() == hostname ||
						m.Hops[hopIdx].Hosts[i].Name == "" {
						m.Hops[hopIdx].Hosts[i].Name = hostname
						break
					}
				}
			}

		case 'p':
			// Add packet data to the hop (aggregated across all IPs)
			if hopIdx, ok := hopNumToIdx[hopnum]; ok {
				m.Hops[hopIdx].PacketMicrosecs = append(
					m.Hops[hopIdx].PacketMicrosecs,
					parseByteNum(line[finalFieldIdx:]),
				)
			}
		}
	}
	m.filterDuplicateLastHops()

	m.processHops()
}

// filterDuplicateLastHops removes consecutive hops at the end of the path
// that have the same IP address. This handles MTR's behavior of reporting
// the destination multiple times at increasing TTL values.
func (m *MTR) filterDuplicateLastHops() {
	if len(m.Hops) == 0 {
		return
	}

	finalIdx := 0
	var previousIP string

	for idx, hop := range m.Hops {
		if len(hop.Hosts) > 0 {
			currentIP := hop.Hosts[0].IP.String()
			if currentIP != previousIP {
				previousIP = currentIP
				finalIdx = idx + 1
			}
		}
	}

	// Trim to the last hop where IP changed
	if finalIdx > 0 && finalIdx < len(m.Hops) {
		m.Hops = m.Hops[0:finalIdx]
	}
}

func (m *MTR) processHops() {
	for _, hop := range m.Hops {
		hop.Sent = m.PacketsSent
		hop.Received = len(hop.PacketMicrosecs)
		hop.Dropped = hop.Sent - hop.Received
		hop.LostPercent = float64(hop.Dropped) / float64(hop.Sent)
		if hop.Received == 0 {
			continue
		}
		totalPacketTime := 0
		best := 1<<31 - 1
		worst := 0
		jitters := make([]int, hop.Received)
		worstJitter := 0
		for i, packet := range hop.PacketMicrosecs {
			if i > 0 {
				newJitter := packet - hop.PacketMicrosecs[i-1]
				if newJitter < 0 {
					newJitter = -newJitter
				}
				if newJitter > worstJitter {
					worstJitter = newJitter
				}
				jitters[i] = newJitter
				hop.InterarrivalJitter += newJitter - ((hop.InterarrivalJitter + 8) >> 4) // rfc3550 A.8
			}
			totalPacketTime += packet
			if packet > worst {
				worst = packet
			}
			if packet < best {
				best = packet
			}
		}
		// mtr keeps a running average, so values may be different than a true average
		hop.Mean = float64(totalPacketTime) / float64(hop.Received)
		hop.WorstJitter = worstJitter
		hop.Best = best
		hop.Worst = worst
		sqrDiff := float64(0)
		jitterSum := 0
		for i, packet := range hop.PacketMicrosecs {
			diff := float64(packet) - hop.Mean
			sqrDiff += diff * diff
			jitterSum += jitters[i]
		}
		hop.MeanJitter = float64(jitterSum) / float64(hop.Received)
		hop.StandardDev = math.Sqrt(sqrDiff / float64(hop.Mean))
	}
}
