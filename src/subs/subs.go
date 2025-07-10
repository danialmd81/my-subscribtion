package subs

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ParseBase64SubsLink decodes a base64-encoded subscription string and separates configs by protocol type.
func ParseBase64SubsLink(b64 string) map[string]string {
	configs := map[string]string{
		"ss":       "",
		"vmess":    "",
		"trojan":   "",
		"vless":    "",
		"hysteria": "",
		"other":    "",
	}

	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return configs // or handle error as needed
	}

	lines := strings.Split(string(decoded), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "ss://"):
			configs["ss"] += line + "\n"
		case strings.HasPrefix(line, "vmess://"):
			configs["vmess"] += line + "\n"
		case strings.HasPrefix(line, "trojan://"):
			configs["trojan"] += line + "\n"
		case strings.HasPrefix(line, "vless://"):
			configs["vless"] += line + "\n"
		case strings.HasPrefix(line, "hysteria://"), strings.HasPrefix(line, "hy2://"):
			configs["hysteria"] += line + "\n"
		default:
			if line != "" {
				configs["other"] += line + "\n"
			}
		}
	}
	return configs
}

// Run reads subscription links from subs.txt, downloads and parses each, then saves separated configs to files by protocol.
func Run() {
	subsFile := "subs/subs.txt" // Input file with subscription links
	outputDir := "subs/"        // Output directory for config files

	fmt.Println("[INFO] Reading subscription links from:", subsFile)
	data, err := os.ReadFile(subsFile)
	if err != nil {
		fmt.Println("[ERROR] Error reading subs.txt:", err)
		return
	}
	links := strings.Split(string(data), "\n")
	configs := map[string]string{
		"ss":       "",
		"vmess":    "",
		"trojan":   "",
		"vless":    "",
		"hysteria": "",
		"other":    "",
	}
	for _, link := range links {
		link = strings.TrimSpace(link)
		if link == "" || strings.HasPrefix(link, "//") {
			continue // Skip empty or commented lines
		}
		fmt.Println("[INFO] Downloading:", link)
		resp, err := http.Get(link)
		if err != nil {
			fmt.Println("[ERROR] Error downloading:", link, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("[ERROR] Error reading body:", link, err)
			continue
		}
		fmt.Println("[INFO] Parsing configs from:", link)
		parsed := ParseBase64SubsLink(string(body)) // Parse and separate configs
		for k, v := range parsed {
			configs[k] += v // Merge configs by type
		}
	}
	fmt.Println("[INFO] Saving separated configs to files in:", outputDir)
	// Save each config type to a file
	for k, v := range configs {
		if v == "" {
			continue // Skip empty config types
		}
		filePath := outputDir + k + ".txt"
		err := os.WriteFile(filePath, []byte(strings.TrimSpace(v)), 0644)
		if err != nil {
			fmt.Println("[ERROR] Error writing", filePath, err)
		} else {
			fmt.Println("[INFO] Saved:", filePath)
		}
	}
	fmt.Println("[INFO] All done!")
}
