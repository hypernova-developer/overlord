package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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

	var query string
	survey.AskOne(&survey.Input{Message: "Enter Movie or TV Show Name:"}, &query)

	apiURL := fmt.Sprintf("https://v3.sg.media-imdb.com/suggestion/%s/%s.json", string(strings.ToLower(query)[0]), url.QueryEscape(strings.ToLower(query)))
	resp, _ := http.Get(apiURL)
	defer resp.Body.Close()

	var searchData IMDbSearchResult
	json.NewDecoder(resp.Body).Decode(&searchData)

	if len(searchData.D) == 0 {
		fmt.Println("[-] No results found!")
		return
	}

	var options []string
	for _, item := range searchData.D {
		options = append(options, fmt.Sprintf("%s (%d) [%s]", item.Title, item.Year, item.Q))
	}

	var choice string
	survey.AskOne(&survey.Select{Message: "Select Title:", Options: options}, &choice)

	var selectedID string
	var isTV bool
	for _, item := range searchData.D {
		if strings.Contains(choice, item.Title) {
			selectedID = item.ID
			if strings.Contains(item.Q, "TV") {
				isTV = true
			}
			break
		}
	}

	const streamSite = "https://vidsrc.to"
	var finalURL string

	if !isTV {
		finalURL = fmt.Sprintf("%s/embed/movie/%s", streamSite, selectedID)
	} else {
		var s, e string
		survey.AskOne(&survey.Input{Message: "Season Number:"}, &s)
		survey.AskOne(&survey.Input{Message: "Episode Number:"}, &e)
		finalURL = fmt.Sprintf("%s/embed/tv/%s/%s/%s", streamSite, selectedID, s, e)
	}

	fmt.Printf("[*] Launching: %s\n", finalURL)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", finalURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", finalURL)
	case "darwin":
		cmd = exec.Command("open", finalURL)
	}
	cmd.Start()
}

