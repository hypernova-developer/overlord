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
	fmt.Println("👑 OVERLORD v3.2.5-BETA — The Ultimate Stream Controller by hypernova-developer")

	var query string
	survey.AskOne(&survey.Input{Message: "Search for movie or series:"}, &query)

	apiURL := fmt.Sprintf("https://v3.sg.media-imdb.com/suggestion/%s/%s.json", string(strings.ToLower(query)[0]), url.QueryEscape(strings.ToLower(query)))
	resp, _ := http.Get(apiURL)
	defer resp.Body.Close()

	var searchData IMDbSearchResult
	json.NewDecoder(resp.Body).Decode(&searchData)

	if len(searchData.D) == 0 {
		fmt.Println("[-] No results found.")
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

	if !isTV {
		launch("https://vidsrc.to/embed/movie/" + selectedID)
	} else {
		var s, e string
		survey.AskOne(&survey.Input{Message: "Season (e.g., 1):"}, &s)
		survey.AskOne(&survey.Input{Message: "Episode (e.g., 1):"}, &e)
		launch(fmt.Sprintf("https://vidsrc.to/embed/tv/%s/%s/%s", selectedID, s, e))
	}
}

func launch(url string) {
	fmt.Printf("[+] Launching stream: %s\n", url)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("open", url)
	}
	cmd.Start()
}

