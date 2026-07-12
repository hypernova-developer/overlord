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

// IMDbSearchResult represents the official unofficial IMDb suggestion API structure
type IMDbSearchResult struct {
	D []struct {
		ID    string `json:"id"`    // IMDb ID (tt1234567)
		Title string `json:"l"`     // Label/Title
		Year  int    `json:"y"`     // Year
		Q     string `json:"q"`     // Type (feature, TV series vb.)
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
		fmt.Println("[-] Error: Query cannot be empty.")
		return
	}

	fmt.Printf("[*] Querying IMDb Servers for '%s'...\n", query)

	// IMDb's zero-protection suggestion endpoint
	firstChar := string(strings.ToLower(query)[0])
	apiURL := fmt.Sprintf("https://v3.sg.media-imdb.com/suggestion/%s/%s.json", firstChar, url.QueryEscape(strings.ToLower(query)))

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("[-] Network Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var searchData IMDbSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchData); err != nil {
		fmt.Printf("[-] Failed to parse data from IMDb: %v\n", err)
		return
	}

	if len(searchData.D) == 0 {
		fmt.Println("[-] No results found on IMDb!")
		return
	}

	// Get the first item
	match := searchData.D[0]
	
	// Format the content type beautifully
	contentType := "Movie"
	if match.Q == "TV series" || match.Q == "TV mini-series" {
		contentType = "TV Show"
	}

	fmt.Printf("[+] Found on IMDb: %s (%d) [%s]\n", match.Title, match.Year, contentType)

	var finalURL string
	if contentType == "Movie" {
		finalURL = fmt.Sprintf("https://vidsrc.xyz/embed/movie/%s", match.ID)
	} else {
		fmt.Print("Enter Season Number: ")
		season, _ := reader.ReadString('\n')
		season = strings.TrimSpace(season)

		fmt.Print("Enter Episode Number: ")
		episode, _ := reader.ReadString('\n')
		episode = strings.TrimSpace(episode)

		finalURL = fmt.Sprintf("https://vidsrc.xyz/embed/tv/%s/%s/%s", match.ID, season, episode)
	}

	fmt.Printf("[*] Launching VLC with Stream: %s\n", finalURL)
	cmd := exec.Command("vlc", finalURL)
	
	err = cmd.Start()
	if err != nil {
		fmt.Printf("[-] Error launching VLC: %v\n", err)
		return
	}

	fmt.Println("[+] Overlord task completed successfully. Stay Tuned.")
}

