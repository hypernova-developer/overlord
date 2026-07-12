package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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
		fmt.Printf("[-] Network Error: %v\n", err)
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
	if contentType == "Movie" {
		finalURL = fmt.Sprintf("https://vidsrc.xyz/embed/movie/%s", match.ID)
	} else {
		fmt.Print("Enter Season: ")
		s, _ := reader.ReadString('\n')
		fmt.Print("Enter Episode: ")
		e, _ := reader.ReadString('\n')
		finalURL = fmt.Sprintf("https://vidsrc.xyz/embed/tv/%s/%s/%s", match.ID, strings.TrimSpace(s), strings.TrimSpace(e))
	}

	fmt.Printf("[*] Launching VLC via yt-dlp: %s\n", finalURL)
	
	// VLC'yi yt-dlp ile entegre bir şekilde çalıştırıyoruz.
	// --meta-title ile VLC'ye dosya değil, bir yayın olduğu bilgisini veriyoruz.
	cmd := exec.Command("vlc", "--meta-title", "Overlord Stream", finalURL)
	
	err = cmd.Start()
	if err != nil {
		fmt.Printf("[-] Error: %v. Make sure yt-dlp is installed and in your PATH.\n", err)
	}
}

