package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	var tickers []string
	for _, subreddit := range data {
		for _, title := range subreddit.Data.Children {
			tickIndex := strings.Index(title.Data.Title, "$")
			if tickIndex > -1 && unicode.IsLetter(rune(title.Data.Title[tickIndex+1])) {
				x := title.Data.Title[tickIndex+1:]
				var breakChar int
				for char := range x {
					if !unicode.IsLetter(rune(x[char])) {
						breakChar = char
						break
					}
				}
				final := x[:breakChar]
				tickers = append(tickers, final)
			}
		}
	}
	return tickers
}

func redditAPICall(titles []string) []postAPIResponse {

	//number of posts per subreddit
	numOfPosts := 50

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
