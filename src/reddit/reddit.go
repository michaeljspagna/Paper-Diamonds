package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

//GetTickers used by outside code. Get list of most popular ticker symbols
func GetTickers() []string {

	//list of subreddits to use
	subreddits := []string{"wallstreetbets", "pennystocks", "RobinHoodPennyStocks"}
	data := redditAPICall(subreddits)
	tickers := formatTickerData(data)

	return tickers
}

func formatTickerData(data []postAPIResponse) []string {
	tickerList := make(map[string]int)
	for _, subreddit := range data {
		for _, post := range subreddit.Data.Children {

			titleTickers := parseStringToTickers(post.Data.Title)
			for _, ticker := range titleTickers {
				if tickerList[ticker] == 0 {
					tickerList[ticker] = 1
				} else {
					tickerList[ticker]++
				}
			}
			messageTickers := parseStringToTickers(post.Data.Message)
			for _, ticker := range messageTickers {
				if tickerList[ticker] == 0 {
					tickerList[ticker] = 1
				} else {
					tickerList[ticker]++
				}
			}

		}
	}
	fmt.Println(tickerList)

	return sortTickersByVolume(tickerList)
}

func parseStringToTickers(data string) []string {

	var tickers []string

	var tickIndexes []int
	for i, c := range data {
		if string(c) == "$" {
			tickIndexes = append(tickIndexes, i)
		}
	}

	for _, tickIndex := range tickIndexes {
		if tickIndex > -1 && unicode.IsLetter(rune(data[tickIndex+1])) {
			x := data[tickIndex+1:]
			var breakChar int
			for char := range x {
				if !unicode.IsLetter(rune(x[char])) {
					breakChar = char
					break
				}
			}
			final := x[:breakChar]
			final = strings.ToUpper(final)
			tickers = append(tickers, final)
		}
	}
	return tickers
}

func sortTickersByVolume(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	ranked := make([]string, len(values))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}

func redditAPICall(titles []string) []postAPIResponse {

	//number of posts per subreddit
	numOfPosts := 500

	var subreddits []postAPIResponse
	for i := 0; i < len(titles); i++ {

		//open client
		client := &http.Client{}
		var URL = "https://www.reddit.com/r/" + titles[i] + "/new.json?limit=" + strconv.Itoa(numOfPosts)
		req, err := http.NewRequest("GET", URL, nil)
		//set user header
		req.Header.Add("User-Agent", "golang:Paper-Diamonds:v0.0.0 (by /u/dfiu65)")
		res, err1 := client.Do(req)

		//stop on errors
		if err != nil {
			log.Fatal(err)
		}
		if err1 != nil {
			log.Fatal(err1)
		}

		//read data
		data, _ := ioutil.ReadAll(res.Body)

		s, _ := formatPosts([]byte(data))

		subreddits = append(subreddits, *s)
	}
	return subreddits
}

type post struct {
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	Message   string `json:"selftext"`
}

type postResponse struct {
	Kind string `json:"kind"`
	Data post   `json:"data"`
}

type apiData struct {
	Modhash  string         `json:"modhash"`
	Dist     int            `json:"dist"`
	Children []postResponse `json:"children"`
	After    string         `json:"after"`
	Before   string         `json:"before"`
}

type postAPIResponse struct {
	Kind string  `json:"kind"`
	Data apiData `json:"data"`
}

func formatPosts(body []byte) (*postAPIResponse, error) {
	var s = new(postAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}
