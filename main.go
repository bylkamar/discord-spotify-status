package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/tidwall/gjson"
)

var (
	session string
	token   string
)

func main() {
	dat, err := os.ReadFile("./session.json")
	if err != nil {
		return
	}
	session = gjson.Get(string(dat), "session").String()
	token = gjson.Get(string(dat), "discord_token").String()
	var token = convertCookiesToken()

	if token == "" {
		pterm.Error.Println("You are not connected to Spotify (SESSION TOKEN NOT FOUND OR INVALID)\nPlease save session cookie in session.json and restart the program.")
		return
	} else {
		pterm.Success.Println("You are connected to Spotify. Waiting for a song to play...")
	}

	var wordChange string
	var currentTrackId string
	var lastPlayedTrackId string

	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				pterm.Warning.Println("Refreshing Spotify session...")
				token = convertCookiesToken()
				if token == "" {
					pterm.Error.Println("You are not connected to Spotify (SESSION TOKEN NOT FOUND OR INVALID)\nPlease save session cookie in session.json and restart the program.")
					os.Exit(1)
					return
				} else {
					pterm.Success.Println("You are connected to Spotify. Waiting for a song to play...")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pterm.Warning.Println("Exiting...")
		os.Exit(0)
	}()

	for {

		// Get current playing song
		var UserView = getCurrentPlaying(token)
		// Check if user is playing a song
		lastPlayedTrackId = currentTrackId
		if UserView.Get("is_playing").Bool() {
			currentTrackId = UserView.Get("item.id").String()
			var word = getLyric(token, currentTrackId, UserView.Get("progress_ms").Int())
			if word != wordChange && word != "" {
				if currentTrackId != lastPlayedTrackId {
					pterm.Success.Println("New song playing: " + UserView.Get("item.name").String())
				}
				var timeDuration, _ = time.ParseDuration(fmt.Sprintf("%.0f", UserView.Get("progress_ms").Float()) + "ms")
				var endDuration, _ = time.ParseDuration(fmt.Sprintf("%.0f", UserView.Get("item.duration_ms").Float()) + "ms")

				// Format lyrics according to progression and end time
				wordChange = fmt.Sprintf("[%s:%s] %s", timeDuration.Round(time.Second).String(), endDuration.Round(time.Second).String(), word)

				changeStatus(wordChange)
			} else {
				if currentTrackId != lastPlayedTrackId {
					pterm.Success.Println("[No Lyrics Founds] New song playing: " + UserView.Get("item.name").String())
				}
			}

		}
		// Wait 3 seconds before next request to avoid rate limit. (Discord don't allow to change status too fast. Average 12 seconds to update status in real time view)
		time.Sleep(4000 * time.Millisecond)
	}
}

// Change discord status according to lyrics
func changeStatus(text string) {

	url := "https://discordapp.com/api/v8/users/@me/settings"
	method := "PATCH"

	payload := strings.NewReader(fmt.Sprintf(`{"custom_status": {"text": "%s","emoji_id": null}}`, text))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "fr,fr-FR;q=0.9")
	req.Header.Add("authorization", token)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://canary.discord.com")
	req.Header.Add("referer", "https://canary.discord.com/channels/@me")
	req.Header.Add("sec-ch-ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/116.0")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}
