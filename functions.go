package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// Convert cookies to token
func convertCookiesToken() string {

	url := "https://open.spotify.com/get_access_token?reason=transport&productType=web_player"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Host", "open.spotify.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Cookie", "sp_dc="+session)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	accessToken := gjson.Get(string(body), "accessToken")
	if strings.Contains(string(body), `"isAnonymous": false`) {
		return ""
	}
	return accessToken.String()
}

// Find according Spotify lyrics from current playing song and trackId
func getLyric(token string, trackId string, currentProgress int64) string {
	url := fmt.Sprintf("https://spclient.wg.spotify.com/lyrics/v1/track/%s", trackId)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Host", "open.spotify.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if string(body) == "" {
		return ""
	}
	// Append all lines in a list according to progress time
	var listMS []int
	var lines = gjson.Get(string(body), "lines").Array()
	for _, line := range lines {
		var TimeSecond = line.Get("time").Int()
		listMS = append(listMS, int(TimeSecond))
	}

	number := findClosest(listMS, int(currentProgress))

	// Discord API doesn't support high frequence rate to change status, so we need to merge lines
	var word = lines[number].Get("words.0.string").String()
	if number+1 < len(lines) {
		word = word + ", " + lines[number+1].Get("words.0.string").String()
	}
	return word
}
func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

// Depreceated function to find nearest number according to lyrics..
//	func findNearestNumber(numbers []int, goal int) []int {
//		n := []int{}
//		for i, num := range numbers {
//			if int(i+1) < len(numbers) {
//				if math.Abs(float64(numbers[i+1]-goal)) < math.Abs(float64(num-goal)) {
//					n = append(n, numbers[i+1])
//				} else {
//					n = append(n, num)
//				}
//			}
//		}
//		if len(n) > 1 {
//			return findNearestNumber(n, goal)
//		} else {
//			return n
//		}
//	}

func findClosest(arr []int, x int) int {
	closestIndex := 0
	closestValue := abs(arr[0] - x)
	for index, val := range arr {
		currentValue := abs(val - x)
		if currentValue < closestValue {
			closestValue = currentValue
			closestIndex = index
		}
	}
	return closestIndex
}

// Get information about current playing track and return result as JSON
func getCurrentPlaying(token string) gjson.Result {
	url := "https://api.spotify.com/v1/me/player"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return gjson.Result{}
	}
	req.Header.Add("Host", "open.spotify.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return gjson.Result{}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return gjson.Result{}
	}
	return gjson.Parse(string(body))
}
