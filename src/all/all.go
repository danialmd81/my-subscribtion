package all

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var protocols = []string{"hysteria", "ss", "trojan", "vless", "vmess", "other"}

func testConnection(protocol, entry string) bool {

	if protocol == "other" {
		return true
	}

	// Try to parse as URL
	u, err := url.Parse(entry)
	if err != nil || u.Host == "" {
		// Try to parse as host:port
		hostPort := entry
		if !strings.Contains(hostPort, ":") {
			return false
		}
		dialer := net.Dialer{Timeout: 500 * time.Millisecond}
		conn, err := dialer.Dial("tcp", hostPort)
		if err == nil {
			conn.Close()
			return true
		}
		return false
	}

	hostPort := u.Host
	if !strings.Contains(hostPort, ":") {
		// Try to get port from scheme
		switch u.Scheme {
		case "ss":
			hostPort += ":8388"
		case "trojan":
			hostPort += ":443"
		case "vless", "vmess":
			hostPort += ":443"
		case "hysteria":
			hostPort += ":443"
		default:
			return false
		}
	}
	// Use TCP for all except hysteria (which may use UDP, but we use TCP for reachability)
	dialer := net.Dialer{Timeout: 500 * time.Millisecond}
	conn, err := dialer.Dial("tcp", hostPort)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func mergeAndTest() error {
	baseDirs := []string{"subs/", "telegram/"} // Directories to read input files from
	outputDir := "all/"                        // Output directory

	for _, proto := range protocols {
		fmt.Printf("[INFO] Processing protocol: %s\n", proto)
		entriesMap := make(map[string]struct{}) // Deduplicate entries
		totalEntries := 0
		for _, dir := range baseDirs {
			filePath := filepath.Join(dir, proto+".txt") // Input file path
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("[WARN] Missing file: %s\n", filePath)
				continue // Skip if file is missing
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text()) // Read and trim each line
				if line != "" {
					entriesMap[line] = struct{}{} // Add to map if not empty
					totalEntries++
				}
			}
			file.Close()
		}
		fmt.Printf("[INFO] Total entries: %d\n", totalEntries)

		goodEntries := []string{} // List of reachable entries
		for entry := range entriesMap {
			if testConnection(proto, entry) { // Test reachability
				goodEntries = append(goodEntries, entry)
			}
		}

		outPath := filepath.Join(outputDir, proto+".txt") // Output file path
		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create output file for %s: %w", proto, err)
		}
		for _, entry := range goodEntries {
			_, _ = outFile.WriteString(entry + "\n") // Write each good entry
		}
		outFile.Close()
		fmt.Printf("[SUMMARY] %s: %d/%d good entries saved to %s\n", proto, len(goodEntries), totalEntries, outPath)
	}
	return nil
}

func Run() {
	err := mergeAndTest()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
