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
x 2 33002
x 3 33003
x 4 33004
x 5 33005
x 6 33006
x 7 33007
x 8 33008
x 9 33009
x 10 33010
x 11 33011
x 12 33012
x 13 33013
x 0 33014
x 1 33015
x 2 33016
x 3 33017
x 4 33018
x 5 33019
x 6 33020
x 7 33021
x 8 33022
x 9 33023
x 10 33024
x 11 33025
x 12 33026
x 13 33027
x 0 33028
x 1 33029
x 2 33030
x 3 33031
x 4 33032
x 5 33033
x 6 33034
x 7 33035
x 8 33036
x 9 33037
x 10 33038
x 11 33039
x 12 33040
x 13 33041
x 0 33042
x 1 33043
x 2 33044
x 3 33045
x 4 33046
x 5 33047
x 6 33048
x 7 33049
x 8 33050
x 9 33051
x 10 33052
x 11 33053
x 12 33054
x 13 33055
x 0 33056
x 1 33057
x 2 33058
x 3 33059
x 4 33060
x 5 33061
x 6 33062
x 7 33063
x 8 33064
x 9 33065
x 10 33066
x 11 33067
x 12 33068
x 13 33069
x 0 33070
x 1 33071
x 2 33072
x 3 33073
x 4 33074
x 5 33075
x 6 33076
x 7 33077
x 8 33078
x 9 33079
x 10 33080
x 11 33081
x 12 33082
x 13 33083
x 0 33084
x 1 33085
x 2 33086
x 3 33087
x 4 33088
h 4 1.2.124.10
h 4 1.2.124.10
p 4 286 33088
x 5 33089
h 5 1.2.96.154
h 5 1.2.96.154
p 5 156 33089
x 6 33090
h 6 1.2.96.154
h 6 1.2.96.154
d 6 1.2.96.154
p 6 234 33090
x 0 33091
x 1 33092
x 2 33093
h 2 1.3.240.31
h 2 1.3.240.31
p 2 1079 33093
x 3 33094
h 3 1.3.240.50
h 3 1.3.240.50
p 3 275 33094
x 4 33095
x 5 33096
d 5 1.2.96.154
p 5 150 33096
x 0 33097
x 1 33098
x 2 33099
d 2 1.3.240.31
p 2 704 33099
x 3 33100
d 3 1.3.240.50
p 3 266 33100
x 4 33101
d 4 1.2.124.10
p 4 234 33101
x 5 33102
p 5 160 33102
x 0 33103
x 1 33104
h 1 1.3.240.85
h 1 1.3.240.85
p 1 237 33104
x 2 33105
h 2 1.3.240.33
h 2 1.3.240.33
p 2 1085 33105
x 3 33106
h 3 1.3.242.4
h 3 1.3.242.4
p 3 338 33106
x 4 33107
h 4 1.2.124.8
h 4 1.2.124.8
p 4 206 33107
x 5 33108
p 5 149 33108
x 0 33109
x 1 33110
h 1 1.3.240.81
h 1 1.3.240.81
p 1 279 33110
x 2 33111
p 2 24000 33111
x 3 33112
h 3 1.3.242.6
h 3 1.3.242.6
p 3 343 33112
x 4 33113
h 4 1.2.124.10
p 4 239 33113
x 5 33114
p 5 194 33114
x 0 33115
x 1 33116
h 1 1.3.240.85
d 1 1.3.240.85
p 1 281 33116
x 2 33117
p 2 1139 33117
x 3 33118
h 3 1.3.240.50
p 3 299 33118
x 4 33119
h 4 1.2.124.8
p 4 316 33119
x 5 33120
p 5 133 33120
x 0 33121
x 1 33122
h 1 1.3.240.83
h 1 1.3.240.83
p 1 316 33122
x 2 33123
h 2 1.3.240.7
h 2 1.3.240.7
p 2 325 33123
x 3 33124
h 3 1.3.242.12
h 3 1.3.242.12
p 3 246 33124
x 4 33125
h 4 1.2.124.10
p 4 276 33125
x 5 33126
p 5 196 33126
x 0 33127
x 1 33128
x 2 33129
h 2 1.3.240.21
h 2 1.3.240.21
p 2 207 33129
x 3 33130
h 3 1.3.242.6
p 3 257 33130
x 4 33131
x 5 33132
p 5 176 33132
x 0 33133
x 1 33134
p 1 245 33134
x 2 33135
h 2 1.3.240.19
h 2 1.3.240.19
p 2 276 33135
x 3 33136
x 4 33137
p 4 362 33137
x 5 33138
p 5 217 33138
x 0 33139
x 1 33140
h 1 1.3.240.81
p 1 159 33140
x 2 33141
x 3 33142
x 4 33143
p 4 304 33143
x 5 33144
p 5 197 33144
x 0 33145
x 1 33146
x 2 33147
h 2 1.3.240.33
p 2 8823 33147
x 3 33148
x 4 33149
p 4 241 33149
x 5 33150
p 5 197 33150
x 0 33151
x 1 33152
p 1 328 33152
x 2 33153
x 3 33154
x 4 33155
p 4 241 33155
x 5 33156
x 0 33157
x 1 33158
h 1 1.3.240.85
p 1 232 33158
x 2 33159
x 3 33160
x 4 33161
p 4 323 33161
x 5 33162
p 5 203 33162
x 0 33163
x 1 33164
h 1 1.3.240.83
p 1 233 33164
x 2 33165
x 3 33166
x 4 33167
p 4 225 33167
x 5 33168
p 5 201 33168
x 0 33169
x 1 33170
p 1 223 33170
x 2 33171
h 2 1.3.240.19
p 2 421 33171
x 3 33172
p 3 373 33172
x 4 33173
x 5 33174
x 0 33175
x 1 33176
h 1 1.3.240.81
p 1 262 33176
x 2 33177
h 2 1.3.240.21
p 2 209 33177
x 3 33178
p 3 317 33178
x 4 33179
p 4 226 33179
x 5 33180
p 5 155 33180
x 0 33181
x 1 33182
p 1 279 33182
x 2 33183
h 2 1.3.240.19
p 2 1372 33183
x 3 33184
h 3 1.3.240.50
p 3 203 33184
x 4 33185
h 4 1.2.124.8
p 4 282 33185
x 5 33186
p 5 139 33186
x 0 33187
x 1 33188
p 1 231 33188
x 2 33189
h 2 1.3.240.21
p 2 655 33189
x 3 33190
h 3 1.3.242.6
p 3 297 33190
x 4 33191
h 4 1.2.124.10
p 4 294 33191
x 5 33192
p 5 185 33192
x 0 33193
x 1 33194
h 1 1.3.240.85
p 1 188 33194
x 2 33195
h 2 1.3.240.5
h 2 1.3.240.5
p 2 1157 33195
x 3 33196
h 3 1.3.242.4
p 3 235 33196
x 4 33197
h 4 1.2.124.8
p 4 228 33197
x 5 33198
p 5 150 33198
x 0 33199
x 1 33200
h 1 1.3.240.81
p 1 204 33200
x 2 33201
h 2 1.3.240.33
p 2 713 33201
x 3 33202
p 3 311 33202
x 4 33203
h 4 1.2.124.10
p 4 299 33203
x 5 33204
x 0 33205
x 1 33206
h 1 1.3.240.85
p 1 244 33206
x 2 33207
x 3 33208
p 3 187 33208
x 4 33209
x 5 33210
x 0 33211
x 1 33212
p 1 349 33212
x 2 33213
h 2 1.3.240.19
p 2 780 33213
x 3 33214
h 3 1.3.242.6
p 3 303 33214
x 4 33215
h 4 1.2.124.8
p 4 296 33215
x 5 33216
x 0 33217
x 1 33218
h 1 1.3.240.83
p 1 301 33218
x 2 33219
x 3 33220
x 4 33221
x 5 33222
x 0 33223
x 1 33224
h 1 1.3.240.85
p 1 200 33224
x 2 33225
p 2 669 33225
x 3 33226
p 3 299 33226
x 4 33227
p 4 227 33227
x 5 33228
x 0 33229
x 1 33230
h 1 1.3.240.83
p 1 215 33230
x 2 33231
x 3 33232
h 3 1.3.242.4
p 3 267 33232
x 4 33233
h 4 1.2.124.10
p 4 284 33233
x 5 33234
p 5 225 33234
x 0 33235
x 1 33236
x 2 33237
h 2 1.3.240.33
p 2 1151 33237
x 3 33238
p 3 377 33238
x 4 33239
x 5 33240
p 5 220 33240
x 0 33241
x 1 33242
p 1 237 33242
x 2 33243
h 2 1.3.240.7
p 2 218 33243
x 3 33244
x 4 33245
h 4 1.2.124.8
p 4 233 33245
x 5 33246
p 5 180 33246
x 0 33247
x 1 33248
p 1 236 33248
x 2 33249
p 2 472 33249
x 3 33250
h 3 1.3.242.12
p 3 262 33250
x 4 33251
h 4 1.2.124.10
p 4 286 33251
x 5 33252
p 5 146 33252
x 0 33253
x 1 33254
h 1 1.3.240.81
p 1 133 33254
x 2 33255
h 2 1.3.240.21
p 2 729 33255
x 3 33256
h 3 1.3.242.4
p 3 265 33256
x 4 33257
x 5 33258
x 0 33259
x 1 33260
h 1 1.3.240.85
p 1 189 33260
x 2 33261
h 2 1.3.240.19
p 2 1214 33261
x 3 33262
h 3 1.3.242.6
p 3 283 33262
x 4 33263
h 4 1.2.124.8
p 4 210 33263
x 5 33264
p 5 185 33264
x 0 33265
x 1 33266
h 1 1.3.240.81
p 1 321 33266
x 2 33267
h 2 1.3.240.5
p 2 499 33267
x 3 33268
x 4 33269
h 4 1.2.124.10
p 4 319 33269
x 5 33270
p 5 188 33270
x 0 33271
x 1 33272
p 1 333 33272
x 2 33273
x 3 33274
x 4 33275
p 4 323 33275
x 5 33276
p 5 175 33276
x 0 33277
x 1 33278
p 1 227 33278
x 2 33279
p 2 898 33279
x 3 33280
p 3 337 33280
x 4 33281
x 5 33282
x 0 33283
x 1 33284
p 1 209 33284
x 2 33285
p 2 646 33285
x 3 33286
h 3 1.3.242.12
p 3 213 33286
x 4 33287
p 4 229 33287
x 5 33288
p 5 151 33288
x 0 33289
x 1 33290
p 1 256 33290
x 2 33291
p 2 279 33291
x 3 33292
x 4 33293
p 4 251 33293
x 5 33294
p 5 141 33294
x 0 33295
x 1 33296
h 1 1.3.240.85
p 1 294 33296
x 2 33297
x 3 33298
x 4 33299
h 4 1.2.124.8
p 4 223 33299
x 5 33300
p 5 128 33300
x 0 33301
x 1 33302
h 1 1.3.240.83
p 1 287 33302
x 2 33303
h 2 1.3.240.31
p 2 1147 33303
x 3 33304
p 3 288 33304
x 4 33305
h 4 1.2.124.10
p 4 249 33305
x 5 33306
p 5 134 33306
x 0 33307
x 1 33308
p 1 277 33308
x 2 33309
h 2 1.3.240.7
p 2 395 33309
x 3 33310
p 3 257 33310
x 4 33311
p 4 263 33311
x 5 33312
p 5 192 33312
x 0 33313
x 1 33314
p 1 251 33314
x 2 33315
p 2 15880 33315
x 3 33316
x 4 33317
h 4 1.2.124.8
p 4 207 33317
x 5 33318
p 5 195 33318
x 0 33319
x 1 33320
h 1 1.3.240.85
p 1 258 33320
x 2 33321
h 2 1.3.240.19
p 2 9687 33321
x 3 33322
h 3 1.3.242.6
p 3 245 33322
x 4 33323
p 4 277 33323
x 5 33324
x 0 33325
x 1 33326
p 1 283 33326
x 2 33327
x 3 33328
h 3 1.3.242.4
p 3 312 33328
x 4 33329
h 4 1.2.124.10
p 4 279 33329
x 5 33330
p 5 144 33330
x 0 33331
x 1 33332
x 2 33333
h 2 1.3.240.21
p 2 1321 33333
x 3 33334
h 3 1.3.242.6
p 3 231 33334
x 4 33335
x 5 33336
x 0 33337
x 1 33338
h 1 1.3.240.83
p 1 244 33338
x 2 33339
x 3 33340
x 4 33341
x 5 33342
x 0 33343
x 1 33344
h 1 1.3.240.85
p 1 230 33344
x 2 33345
x 3 33346
x 4 33347
h 4 1.2.124.8
p 4 214 33347
x 5 33348
p 5 175 33348
`

	m := &MTR{
		OutputRaw:   []byte(rawOutput),
		PacketsSent: 50,
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
	hop1, exists := m.Hops[1]
	if !exists {
		t.Fatal("Hop 1 not found")
	}

	expectedHop1IPs := map[string]bool{
		"1.3.240.81": false,
		"1.3.240.83": false,
		"1.3.240.85": false,
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
	hop2, exists := m.Hops[2]
	if !exists {
		t.Fatal("Hop 2 not found")
	}

	expectedHop2IPs := map[string]bool{
		"1.3.240.31": false,
		"1.3.240.33": false,
		"1.3.240.7":  false,
		"1.3.240.21": false,
		"1.3.240.19": false,
		"1.3.240.5":  false,
	}

	if len(hop2.Hosts) != len(expectedHop2IPs) {
		t.Errorf("Expected %d hosts at hop 2 (ECMP), got %d", len(expectedHop2IPs), len(hop2.Hosts))
	}

	// Verify hop 6 was filtered out (duplicate of hop 5)
	if _, exists := m.Hops[6]; exists {
		t.Error("Hop 6 should have been filtered out as duplicate of hop 5")
	}

	// Verify hop 5 exists and is the destination
	hop5, exists := m.Hops[5]
	if !exists {
		t.Fatal("Hop 5 not found")
	}

	if len(hop5.Hosts) == 0 {
		t.Error("Hop 5 should have at least one host")
	} else if hop5.Hosts[0].IP.String() != "1.2.96.154" {
		t.Errorf("Hop 5 should be destination 1.2.96.154, got %s", hop5.Hosts[0].IP.String())
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
