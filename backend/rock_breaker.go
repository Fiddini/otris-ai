// Allahu Latifu bi 'ibadihi yarzuqu man yasya' wa huwal 'Azizul Hakim - QS 42:19
// Amalan terakhir Ayah sebelum meninggal - Jatinom
// break rock - 25% cukup - raw evidence only - 16 byte sample + CRC32
// OTRIS ORG - CONSTITUSI 11 PASAL - TULUNG AGUSTUS 2026
// 1. Tujuan: Jangan pindahkan pasien tapi hadirkan spesialis di waktu yang tepat
// 2. Data: Tidak di ambo - cuma ACTIVE - bukan nama pasien - haram dijual
// 3. KSO: Jangan tarik - jualan langsung kapas murah - halal
// 4. Meta: Secanggih ini tidak dibebankan biaya - tidak KSO
// 5. Plug & Play: Nilai besar bukan di situ - nilai besar bantu org
// 6. Lampu Operasi: Jumlah mati belum lihat lampu operasi banyak - harus dibebaskan
// 7. STOP REPORT EXACT CONFLICT - DO NOT FABRICATE
// 8. Crusher Wall: Pecah watu protokol - save raw + CRC32 - jangan ngarang
// 9. Bela Pasien: Jangan pindahkan pasien - tekakno spesialis
// 10. Bebaskan: Jangan egois - berbagilah - kalau berbagi hasil sudah banyak baru pasien ditangani sendiri
// 11. 25% cukup - harga terserah dio - makin kaya - alam bergerak - Batunya dari kami - Open sorenya

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const schemaVersion = "1.2"

var hl7Mappings = []struct {
	Needle    string
	Canonical string
}{
	{"151817", "Ventilation.PIP"},
	{"151804", "PEEP"},
	{"151784", "Pplat"},
	{"152000", "MVe"},
	{"151562", "RR"},
	{"152196", "FiO2"},
	{"151708", "EtCO2"},
	{"152192", "FiN2O"},
	{"152108", "EtN2O"},
	{"152172", "FiAgent Insp"},
	{"152088", "EtAgent Exp"},
	{"151692", "Compliance"},
	{"151840", "Resistance"},
	{"151832", "IE Ratio"},
	{"147842", "HR"},
	{"196670", "Alarm"},
	{"151143", "TVe"},
	{"143", "TVe"},
	{"144", "TVi"},
	{"119", "MAC"},
}

var qcmsExpectedPackages = []string{
	"ECG/HR/PVC/ST/waveform",
	"SpO2/PR/waveform",
	"RESP",
	"TEMP T1/T2",
	"NIBP",
	"CO2 Et/Fi/AWRR/waveform",
	"AG Et/Fi/MAC",
	"CO",
	"IBP1-4",
	"BIS/EMG/SR/SQI/TP/SEF/waveform",
	"Alarm",
	"Module Info ON/OFF array",
}

type Config struct {
	Mode      string
	Edge      string
	Port      int
	UDP       int
	TCP       int
	PCAPBase  string
	SelfCheck bool
}

type HL7Observation struct {
	Code      string
	Canonical string
	Value     string
	Unit      string
}

type HL7Report struct {
	Valid         bool
	MessageType   string
	Version       string
	BedID         string
	DeviceID      string
	PatientToken  string
	TimestampUTC  string
	SchemaVersion string
	Observations  []HL7Observation
	Gaps          []string
}

func main() {
	cfg := parseFlags()
	if cfg.SelfCheck {
		runSelfCheck()
		return
	}

	var err error
	switch strings.ToLower(cfg.Mode) {
	case "hl7":
		err = runHL7Mode(cfg, []int{cfg.Port})
	case "qcms":
		err = runQCMSMode(cfg)
	case "unknown":
		err = runUnknownMode(cfg)
	case "all":
		fmt.Println("[ROCK BREAKER][UNKNOWN][ALL] GAP-REPORT mode=all does not launch duplicate unknown listeners on the same documented ports; protocol listeners still save raw evidence on every capture")
		err = runAllMode(cfg)
	default:
		err = fmt.Errorf("unsupported mode %q", cfg.Mode)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "[ROCK BREAKER][FATAL] %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() Config {
	mode := flag.String("mode", "all", "hl7|qcms|unknown|all")
	edge := flag.String("edge", "0.0.0.0", "listen IP")
	port := flag.Int("port", 7000, "HL7 port (6000/7000/8000)")
	udp := flag.Int("udp", 2527, "QCMS UDP discovery port")
	tcp := flag.Int("tcp", 2528, "QCMS TCP data port")
	pcap := flag.String("pcap", "./bukti", "save raw append evidence base path")
	selfcheck := flag.Bool("selfcheck", false, "run built-in self-checks")
	flag.Parse()
	return Config{
		Mode:      *mode,
		Edge:      *edge,
		Port:      *port,
		UDP:       *udp,
		TCP:       *tcp,
		PCAPBase:  *pcap,
		SelfCheck: *selfcheck,
	}
}

func runAllMode(cfg Config) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 5)

	for _, port := range []int{6000, 7000, 8000} {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if err := startHL7Listener(cfg, p); err != nil {
				errCh <- err
			}
		}(port)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startQCMSUDPListener(cfg, cfg.UDP); err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startQCMSTCPListener(cfg, cfg.TCP); err != nil {
			errCh <- err
		}
	}()

	return waitForever(&wg, errCh)
}

func runHL7Mode(cfg Config, ports []int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(ports))
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if err := startHL7Listener(cfg, p); err != nil {
				errCh <- err
			}
		}(port)
	}
	return waitForever(&wg, errCh)
}

func runQCMSMode(cfg Config) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startQCMSUDPListener(cfg, cfg.UDP); err != nil {
			errCh <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startQCMSTCPListener(cfg, cfg.TCP); err != nil {
			errCh <- err
		}
	}()
	return waitForever(&wg, errCh)
}

func runUnknownMode(cfg Config) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 5)
	for _, port := range []int{6000, 7000, 8000, cfg.TCP} {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if err := startUnknownTCPListener(cfg, p); err != nil {
				errCh <- err
			}
		}(port)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startUnknownUDPListener(cfg, cfg.UDP); err != nil {
			errCh <- err
		}
	}()
	return waitForever(&wg, errCh)
}

func waitForever(wg *sync.WaitGroup, errCh chan error) error {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		return err
	case <-done:
		return nil
	}
}

func startHL7Listener(cfg Config, port int) error {
	addr := fmt.Sprintf("%s:%d", cfg.Edge, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("HL7 listen %s failed: %w", addr, err)
	}
	fmt.Printf("[ROCK BREAKER][HL7][%d] LISTEN %s\n", port, addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("HL7 accept %d failed: %w", port, err)
		}
		go handleHL7Conn(cfg, port, conn)
	}
}

func handleHL7Conn(cfg Config, port int, conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	fmt.Printf("[ROCK BREAKER][HL7][%d] CONNECTED %s\n", port, remote)

	carry := make([]byte, 0, 65536)
	buf := make([]byte, 65536)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			chunk := append([]byte(nil), buf[:n]...)
			crc := crc32.ChecksumIEEE(chunk)
			path := evidencePath(cfg.PCAPBase, strconv.Itoa(port))
			saveRaw(path, chunk, crc)
			printEvidence("HL7", strconv.Itoa(port), chunk, crc, path)

			carry = append(carry, chunk...)
			msgs, rest := extractHL7MessagesWithRemainder(carry)
			carry = rest
			for _, msg := range msgs {
				report := processHL7Message(msg)
				printHL7Report(port, report)
			}
		}
		if err != nil {
			if len(carry) > 0 {
				reportGap("HL7", strconv.Itoa(port), "incomplete buffered payload left on disconnect; raw evidence saved unchanged")
			}
			if !errors.Is(err, io.EOF) {
				reportGap("HL7", strconv.Itoa(port), fmt.Sprintf("read error: %v", err))
			}
			fmt.Printf("[ROCK BREAKER][HL7][%d] DISCONNECTED %s\n", port, remote)
			return
		}
	}
}

func startQCMSUDPListener(cfg Config, port int) error {
	addr := fmt.Sprintf("%s:%d", cfg.Edge, port)
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("QCMS UDP listen %s failed: %w", addr, err)
	}
	defer pc.Close()
	fmt.Printf("[ROCK BREAKER][QCMS][UDP:%d] LISTEN %s\n", port, addr)
	buf := make([]byte, 65536)
	for {
		n, remote, err := pc.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("QCMS UDP read %d failed: %w", port, err)
		}
		chunk := append([]byte(nil), buf[:n]...)
		crc := crc32.ChecksumIEEE(chunk)
		path := evidencePath(cfg.PCAPBase, "qcms_udp")
		saveRaw(path, chunk, crc)
		fmt.Printf("[ROCK BREAKER][QCMS][UDP:%d] PACKET %s\n", port, remote.String())
		printEvidence("QCMS", fmt.Sprintf("UDP:%d", port), chunk, crc, path)
		reportGap("QCMS", fmt.Sprintf("UDP:%d", port), tryDecodeQCMSHead(chunk))
		reportGap("QCMS", fmt.Sprintf("UDP:%d", port), "checksum UNKNOWN; scaling UNKNOWN; authoritative pcap from RS required")
		reportGap("QCMS", fmt.Sprintf("UDP:%d", port), "expected packages: "+strings.Join(qcmsExpectedPackages, ", "))
	}
}

func startQCMSTCPListener(cfg Config, port int) error {
	addr := fmt.Sprintf("%s:%d", cfg.Edge, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("QCMS TCP listen %s failed: %w", addr, err)
	}
	fmt.Printf("[ROCK BREAKER][QCMS][TCP:%d] LISTEN %s\n", port, addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("QCMS TCP accept %d failed: %w", port, err)
		}
		go handleQCMSTCPConn(cfg, port, conn)
	}
}

func handleQCMSTCPConn(cfg Config, port int, conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	fmt.Printf("[ROCK BREAKER][QCMS][TCP:%d] CONNECTED %s\n", port, remote)
	buf := make([]byte, 65536)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			chunk := append([]byte(nil), buf[:n]...)
			crc := crc32.ChecksumIEEE(chunk)
			path := evidencePath(cfg.PCAPBase, "qcms_tcp")
			saveRaw(path, chunk, crc)
			printEvidence("QCMS", fmt.Sprintf("TCP:%d", port), chunk, crc, path)
			reportGap("QCMS", fmt.Sprintf("TCP:%d", port), tryDecodeQCMSHead(chunk))
			reportGap("QCMS", fmt.Sprintf("TCP:%d", port), "checksum UNKNOWN; scaling UNKNOWN; authoritative pcap from RS required")
			reportGap("QCMS", fmt.Sprintf("TCP:%d", port), "expected packages: "+strings.Join(qcmsExpectedPackages, ", "))
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				reportGap("QCMS", fmt.Sprintf("TCP:%d", port), fmt.Sprintf("read error: %v", err))
			}
			fmt.Printf("[ROCK BREAKER][QCMS][TCP:%d] DISCONNECTED %s\n", port, remote)
			return
		}
	}
}

func startUnknownTCPListener(cfg Config, port int) error {
	addr := fmt.Sprintf("%s:%d", cfg.Edge, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("UNKNOWN TCP listen %s failed: %w", addr, err)
	}
	fmt.Printf("[ROCK BREAKER][UNKNOWN][TCP:%d] LISTEN %s\n", port, addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("UNKNOWN TCP accept %d failed: %w", port, err)
		}
		go handleUnknownTCPConn(cfg, port, conn)
	}
}

func handleUnknownTCPConn(cfg Config, port int, conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	fmt.Printf("[ROCK BREAKER][UNKNOWN][TCP:%d] CONNECTED %s\n", port, remote)
	buf := make([]byte, 65536)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			chunk := append([]byte(nil), buf[:n]...)
			crc := crc32.ChecksumIEEE(chunk)
			path := evidencePath(cfg.PCAPBase, fmt.Sprintf("unknown_tcp_%d", port))
			saveRaw(path, chunk, crc)
			printEvidence("UNKNOWN", fmt.Sprintf("TCP:%d", port), chunk, crc, path)
			reportGap("UNKNOWN", fmt.Sprintf("TCP:%d", port), "raw-only capture; protocol undecoded by design")
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				reportGap("UNKNOWN", fmt.Sprintf("TCP:%d", port), fmt.Sprintf("read error: %v", err))
			}
			fmt.Printf("[ROCK BREAKER][UNKNOWN][TCP:%d] DISCONNECTED %s\n", port, remote)
			return
		}
	}
}

func startUnknownUDPListener(cfg Config, port int) error {
	addr := fmt.Sprintf("%s:%d", cfg.Edge, port)
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("UNKNOWN UDP listen %s failed: %w", addr, err)
	}
	defer pc.Close()
	fmt.Printf("[ROCK BREAKER][UNKNOWN][UDP:%d] LISTEN %s\n", port, addr)
	buf := make([]byte, 65536)
	for {
		n, remote, err := pc.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("UNKNOWN UDP read %d failed: %w", port, err)
		}
		chunk := append([]byte(nil), buf[:n]...)
		crc := crc32.ChecksumIEEE(chunk)
		path := evidencePath(cfg.PCAPBase, "unknown_udp_2527")
		saveRaw(path, chunk, crc)
		fmt.Printf("[ROCK BREAKER][UNKNOWN][UDP:%d] PACKET %s\n", port, remote.String())
		printEvidence("UNKNOWN", fmt.Sprintf("UDP:%d", port), chunk, crc, path)
		reportGap("UNKNOWN", fmt.Sprintf("UDP:%d", port), "raw-only capture; protocol undecoded by design")
	}
}

func printEvidence(kind, port string, raw []byte, crc uint32, path string) {
	// 25% principle: console shows only a 16-byte sample; full raw stays private in the evidence file.
	hexSample := strings.ToUpper(hex.EncodeToString(safeSlice(raw, 0, 16)))
	truncated := max(0, len(raw)-16)
	fmt.Printf("[ROCK BREAKER][%s][%s] Hex sample 16 %s... [truncated %d]\n", kind, port, hexSample, truncated)
	fmt.Printf("[ROCK BREAKER][%s][%s] CRC32 %08X\n", kind, port, crc)
	fmt.Printf("[ROCK BREAKER][%s][%s] EVIDENCE %s\n", kind, port, path)
}

func reportGap(kind, port, message string) {
	fmt.Printf("[ROCK BREAKER][%s][%s] GAP-REPORT %s\n", kind, port, message)
}

func evidencePath(base, suffix string) string {
	if base == "" {
		base = "./bukti"
	}
	return base + "_" + suffix + ".raw"
}

func safeSlice(raw []byte, start, end int) []byte {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > len(raw) {
		start = len(raw)
	}
	if end > len(raw) {
		end = len(raw)
	}
	if start > end {
		start = end
	}
	return raw[start:end]
}

func saveRaw(base string, raw []byte, crc uint32) {
	f, err := os.OpenFile(base, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ROCK BREAKER][EVIDENCE] open %s failed: %v\n", base, err)
		return
	}
	defer f.Close()
	header := fmt.Sprintf("\n--- %s CRC32:%08X Len:%d ---\n", time.Now().UTC().Format(time.RFC3339), crc, len(raw))
	_, _ = f.WriteString(header)
	_, _ = f.Write(raw)
	_, _ = f.WriteString("\n")
}

func extractHL7Messages(raw []byte) [][]byte {
	msgs, _ := extractHL7MessagesWithRemainder(raw)
	return msgs
}

func extractHL7MessagesWithRemainder(raw []byte) ([][]byte, []byte) {
	if len(raw) == 0 {
		return nil, nil
	}

	if bytes.Contains(raw, []byte{0x0B}) {
		return extractMLLPMessages(raw)
	}

	// ponytail: Non-MLLP HL7 framing is chunk-based when boundaries are ambiguous; upgrade with real RS pcap captures.
	starts := allIndices(raw, []byte("MSH|"))
	if len(starts) == 0 {
		return nil, safeSlice(raw, max(0, len(raw)-4096), len(raw))
	}

	if len(starts) == 1 {
		msg := trimHL7(raw[starts[0]:])
		if looksCompleteHL7(msg) {
			return [][]byte{append([]byte(nil), msg...)}, nil
		}
		return nil, append([]byte(nil), raw[starts[0]:]...)
	}

	msgs := make([][]byte, 0, len(starts))
	for i := 0; i < len(starts)-1; i++ {
		msg := trimHL7(raw[starts[i]:starts[i+1]])
		if len(msg) > 0 {
			msgs = append(msgs, append([]byte(nil), msg...))
		}
	}
	rest := append([]byte(nil), raw[starts[len(starts)-1]:]...)
	if looksCompleteHL7(rest) {
		msgs = append(msgs, trimHL7(rest))
		rest = nil
	}
	return msgs, rest
}

func extractMLLPMessages(raw []byte) ([][]byte, []byte) {
	msgs := make([][]byte, 0)
	cursor := 0
	for cursor < len(raw) {
		startRel := bytes.IndexByte(raw[cursor:], 0x0B)
		if startRel < 0 {
			break
		}
		start := cursor + startRel
		endRel := bytes.IndexByte(raw[start+1:], 0x1C)
		if endRel < 0 {
			return msgs, append([]byte(nil), raw[start:]...)
		}
		end := start + 1 + endRel
		msg := trimHL7(raw[start+1 : end])
		if len(msg) > 0 {
			msgs = append(msgs, append([]byte(nil), msg...))
		}
		cursor = end + 1
		if cursor < len(raw) && raw[cursor] == 0x0D {
			cursor++
		}
	}
	return msgs, append([]byte(nil), safeSlice(raw, cursor, len(raw))...)
}

func processHL7Message(raw []byte) HL7Report {
	report := HL7Report{SchemaVersion: schemaVersion}
	segments := splitSegments(raw)
	segMap := map[string][][]string{}
	for _, seg := range segments {
		fields := strings.Split(seg, "|")
		if len(fields) == 0 {
			continue
		}
		name := fields[0]
		segMap[name] = append(segMap[name], fields)
	}

	msh := firstSegment(segMap, "MSH")
	if len(msh) == 0 {
		report.Gaps = append(report.Gaps, "missing MSH segment")
		return report
	}

	fullLower := strings.ToLower(string(raw))
	if !strings.Contains(fullLower, "northern") {
		report.Gaps = append(report.Gaps, "MSH does not contain northern")
		return report
	}

	report.MessageType = fieldAt(msh, 8)
	report.Version = fieldAt(msh, 11)
	if !strings.Contains(report.MessageType, "ORU^R01") && !strings.Contains(report.MessageType, "ORU^R40") {
		report.Gaps = append(report.Gaps, "message type is not ORU^R01 or ORU^R40")
		return report
	}
	if report.Version != "2.6" {
		report.Gaps = append(report.Gaps, "HL7 version is not 2.6")
		return report
	}

	pid := firstSegment(segMap, "PID")
	pv1 := firstSegment(segMap, "PV1")
	obr := firstSegment(segMap, "OBR")
	patientRaw := fieldAt(pid, 3)
	report.PatientToken = hmacPID(patientRaw)
	report.BedID = parseBedID(fieldAt(pv1, 3))
	if report.BedID == "" {
		report.BedID = "unknown-bed"
		report.Gaps = append(report.Gaps, "PV1-3 bed id missing")
	}
	report.DeviceID = "northern-atlas-n7-" + report.BedID
	report.TimestampUTC = chooseTimestamp(fieldAt(obr, 7), fieldAt(msh, 6)).Format(time.RFC3339)

	for _, obx := range segMap["OBX"] {
		codeField := fieldAt(obx, 3)
		canonical, ok := mapOBX(codeField)
		if !ok {
			continue
		}
		report.Observations = append(report.Observations, HL7Observation{
			Code:      codeField,
			Canonical: canonical,
			Value:     fieldAt(obx, 5),
			Unit:      fieldAt(obx, 6),
		})
	}

	if len(report.Observations) == 0 {
		report.Gaps = append(report.Gaps, "no mapped OBX-3 code found with contains matching")
	}

	report.Valid = true
	return report
}

func printHL7Report(port int, report HL7Report) {
	portText := strconv.Itoa(port)
	if !report.Valid {
		for _, gap := range report.Gaps {
			reportGap("HL7", portText, gap)
		}
		return
	}
	for _, obs := range report.Observations {
		if obs.Canonical == "Alarm" || strings.Contains(report.MessageType, "ORU^R40") {
			fmt.Printf("[ROCK BREAKER][HL7][%s] ALARM code=%s text=%q patient_hmac=%s device=%s timestamp=%s schema=%s\n",
				portText, obs.Code, obs.Value, report.PatientToken, report.DeviceID, report.TimestampUTC, report.SchemaVersion)
			continue
		}
		fmt.Printf("[ROCK BREAKER][HL7][%s] OBS canonical=%s code=%s value=%q unit=%q patient_hmac=%s device=%s timestamp=%s schema=%s\n",
			portText, obs.Canonical, obs.Code, obs.Value, obs.Unit, report.PatientToken, report.DeviceID, report.TimestampUTC, report.SchemaVersion)
	}
	for _, gap := range report.Gaps {
		reportGap("HL7", portText, gap)
	}
}

func tryDecodeQCMSHead(raw []byte) string {
	// ponytail: QCMS decoder is heuristic-only until authoritative checksum/scaling evidence exists.
	preview := strings.ToUpper(hex.EncodeToString(safeSlice(raw, 0, 16)))
	if len(raw) < 4 {
		return "NOT AUTHORITATIVE len<4; checksum UNKNOWN; scaling UNKNOWN; head=" + preview
	}
	ascii := printableASCII(safeSlice(raw, 0, 8))
	if ascii != "" {
		return "NOT AUTHORITATIVE heuristic head ASCII=" + ascii + " HEX=" + preview
	}
	return "NOT AUTHORITATIVE heuristic head HEX=" + preview
}

func hmacPID(pid string) string {
	salt := os.Getenv("OTRIS_SALT")
	if salt == "" {
		salt = "otris-salt-dev"
	}
	mac := hmac.New(sha256.New, []byte(salt))
	_, _ = mac.Write([]byte(pid))
	return hex.EncodeToString(mac.Sum(nil))
}

func parseBedID(field string) string {
	parts := strings.Split(field, "^")
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])
		if part != "" {
			return part
		}
	}
	return ""
}

func chooseTimestamp(obr7, msh7 string) time.Time {
	for _, candidate := range []string{obr7, msh7} {
		if t, ok := parseHL7Timestamp(candidate); ok {
			return t.UTC()
		}
	}
	return time.Now().UTC()
}

func parseHL7Timestamp(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	layouts := []string{
		"20060102150405-0700",
		"20060102150405Z0700",
		"20060102150405",
		"200601021504",
		"20060102",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func mapOBX(codeField string) (string, bool) {
	for _, m := range hl7Mappings {
		if strings.Contains(codeField, m.Needle) {
			return m.Canonical, true
		}
	}
	return "", false
}

func splitSegments(raw []byte) []string {
	normalized := strings.ReplaceAll(string(raw), "\n", "\r")
	parts := strings.Split(normalized, "\r")
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			segments = append(segments, part)
		}
	}
	return segments
}

func firstSegment(segMap map[string][][]string, name string) []string {
	items := segMap[name]
	if len(items) == 0 {
		return nil
	}
	return items[0]
}

func fieldAt(fields []string, index int) string {
	if index < 0 || index >= len(fields) {
		return ""
	}
	return strings.TrimSpace(fields[index])
}

func allIndices(raw, needle []byte) []int {
	indices := []int{}
	cursor := 0
	for {
		idx := bytes.Index(raw[cursor:], needle)
		if idx < 0 {
			return indices
		}
		absolute := cursor + idx
		indices = append(indices, absolute)
		cursor = absolute + len(needle)
		if cursor >= len(raw) {
			return indices
		}
	}
}

func looksCompleteHL7(raw []byte) bool {
	text := string(raw)
	return strings.Contains(text, "MSH|") && strings.Contains(text, "PID|") && strings.Contains(text, "PV1|") && strings.Contains(text, "OBR|") && strings.Contains(text, "OBX|")
}

func trimHL7(raw []byte) []byte {
	trimmed := bytes.TrimSpace(raw)
	trimmed = bytes.TrimPrefix(trimmed, []byte{0x0B})
	trimmed = bytes.TrimSuffix(trimmed, []byte{0x1C})
	trimmed = bytes.TrimSuffix(trimmed, []byte{0x0D})
	trimmed = bytes.TrimSpace(trimmed)
	return trimmed
}

func printableASCII(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	b := strings.Builder{}
	for _, c := range raw {
		if c >= 32 && c <= 126 {
			b.WriteByte(c)
		}
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runSelfCheck() {
	mustEqualBytes("safeSlice basic", safeSlice([]byte("abcdef"), 1, 4), []byte("bcd"))
	raw := []byte("\x0bMSH|^~\\&|northern|anesthetics^|||20210823164444||ORU^R01^ORU_R01|191|P|2.6|\rPID|||No 000001||^||00000000|F|\rPV1||I|^^1\rOBR|1||||||20210823164444|\rOBX|1|NM|151817^MDC_PRESS_AWAY_INSP_PEAK^MDC|1.14.1.151817|16|266048^MDC_DIM_CM_H2O^MDC|||||F||||\x1c\r")
	msgs := extractHL7Messages(raw)
	if len(msgs) != 1 {
		failf("extractHL7Messages count got=%d want=1", len(msgs))
	}
	report := processHL7Message(msgs[0])
	if !report.Valid {
		failf("processHL7Message valid=false gaps=%v", report.Gaps)
	}
	if report.DeviceID != "northern-atlas-n7-1" {
		failf("device id got=%q want=%q", report.DeviceID, "northern-atlas-n7-1")
	}
	if len(report.Observations) != 1 {
		failf("observation count got=%d want=1", len(report.Observations))
	}
	if report.Observations[0].Canonical != "Ventilation.PIP" {
		failf("canonical got=%q want=%q", report.Observations[0].Canonical, "Ventilation.PIP")
	}
	if report.PatientToken == "" {
		failf("patient token empty")
	}
	preview := captureStdout(func() {
		printEvidence("HL7", "7000", []byte("0123456789abcdefghijklmnop"), crc32.ChecksumIEEE([]byte("0123456789abcdefghijklmnop")), "./bukti_7000.raw")
	})
	mustContain("preview label", preview, "Hex sample 16")
	mustContain("preview truncation", preview, "[truncated")
	mustNotContain("preview old label", preview, "Hex first 256")
	fmt.Println("selfcheck PASS")
}

func mustEqualBytes(name string, got, want []byte) {
	if string(got) != string(want) {
		failf("%s got=%q want=%q", name, string(got), string(want))
	}
}

func mustContain(name, got, want string) {
	if !strings.Contains(got, want) {
		failf("%s missing=%q got=%q", name, want, got)
	}
}

func mustNotContain(name, got, unwanted string) {
	if strings.Contains(got, unwanted) {
		failf("%s unexpected=%q got=%q", name, unwanted, got)
	}
}

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		failf("capture stdout pipe error: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	out, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil {
		failf("capture stdout read error: %v", err)
	}
	return string(out)
}

func failf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "selfcheck FAIL: "+format+"\n", args...)
	os.Exit(1)
}
