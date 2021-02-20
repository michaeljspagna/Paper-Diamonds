package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Function used by outside code. Get list of most popular ticker symbols
func getTickers() postAPIResponse {

	data := redditAPICall()

	s, _ := formatPosts([]byte(data))

	fmt.Println(s)
	return *s
}

func redditAPICall() []byte {

	subreddits := []string{"wallstreetbets"}

	client := &http.Client{}
	var URL = "https://www.reddit.com/r/" + subreddits[0] + "/new.json?limit=2"
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Add("User-Agent", "golang:Paper-Diamonds:v0.0.0 (by /u/dfiu65)")
	res, err1 := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if err1 != nil {
		log.Fatal(err1)
	}

	data, _ := ioutil.ReadAll(res.Body)

	return data

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
