package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type IMDbSearchResult struct {
	D []struct {
		ID    string `json:"id"`
		Title string `json:"l"`
		Year  int    `json:"y"`
		Q     string `json:"q"`
	} `json:"d"`
}

func main() {
	fmt.Println("👑 OVERLORD v1.0 — The Ultimate CLI Stream Overlord")
	fmt.Println("--------------------------------------------------")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Movie or TV Show Name: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		return
	}

	fmt.Printf("[*] Querying IMDb Servers for '%s'...\n", query)
	firstChar := string(strings.ToLower(query)[0])
	apiURL := fmt.Sprintf("https://v3.sg.media-imdb.com/suggestion/%s/%s.json", firstChar, url.QueryEscape(strings.ToLower(query)))

	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var searchData IMDbSearchResult
	json.NewDecoder(resp.Body).Decode(&searchData)

	if len(searchData.D) == 0 {
		fmt.Println("[-] No results found!")
		return
	}

	match := searchData.D[0]
	contentType := "Movie"
	if strings.Contains(match.Q, "TV") {
		contentType = "TV Show"
	}

	fmt.Printf("[+] Found: %s (%d) [%s]\n", match.Title, match.Year, contentType)

	var finalURL string
	const streamSite = "https://vidsrc.to"

	if contentType == "Movie" {
		finalURL = fmt.Sprintf("%s/embed/movie/%s", streamSite, match.ID)
	} else {
		fmt.Print("Enter Season: ")
		s, _ := reader.ReadString('\n')
		fmt.Print("Enter Episode: ")
		e, _ := reader.ReadString('\n')
		finalURL = fmt.Sprintf("%s/embed/tv/%s/%s/%s", streamSite, match.ID, strings.TrimSpace(s), strings.TrimSpace(e))
	}

	fmt.Printf("[*] Launching Browser with: %s\n", finalURL)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", finalURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", finalURL)
	case "darwin":
		cmd = exec.Command("open", finalURL)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("[-] Error: %v\n", err)
	} else {
		fmt.Println("[+] Overlord task completed successfully.")
	}
}

