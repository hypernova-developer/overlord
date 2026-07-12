package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
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

	var query string
	survey.AskOne(&survey.Input{Message: "Search:"}, &query)

	apiURL := fmt.Sprintf("https://v3.sg.media-imdb.com/suggestion/%s/%s.json", string(strings.ToLower(query)[0]), url.QueryEscape(strings.ToLower(query)))
	resp, _ := http.Get(apiURL)
	defer resp.Body.Close()

	var searchData IMDbSearchResult
	json.NewDecoder(resp.Body).Decode(&searchData)

	var options []string
	for _, item := range searchData.D {
		options = append(options, fmt.Sprintf("%s (%d) [%s]", item.Title, item.Year, item.Q))
	}

	var choice string
	survey.AskOne(&survey.Select{Message: "Select:", Options: options}, &choice)

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

	finalURL := fmt.Sprintf("https://vidsrc.to/embed/movie/%s", selectedID)
	if isTV {
		var seasons []string
		for i := 1; i <= 20; i++ {
			seasons = append(seasons, strconv.Itoa(i))
		}

		var s string
		survey.AskOne(&survey.Select{Message: "Select Season:", Options: seasons}, &s)

		var episodes []string
		for i := 1; i <= 50; i++ {
			episodes = append(episodes, strconv.Itoa(i))
		}

		var e string
		survey.AskOne(&survey.Select{Message: "Select Episode:", Options: episodes}, &e)

		finalURL = fmt.Sprintf("https://vidsrc.to/embed/tv/%s/%s/%s", selectedID, s, e)
	}

	cmd := exec.Command("xdg-open", finalURL)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", finalURL)
	}
	cmd.Start()
}

