package twitter

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type OauthResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

// Client is http client for twitter api.
// This contains authentication information.
type Client struct {
	apiKey string
	token  string
	c      *http.Client
}

// NewFromEnv create new Client. Keys for API authentication
// are obtained from following environment value.
// - TWITTER_CONSUMER_KEY
// - TWITTER_SECRET
//
func NewFromEnv() *Client {
	c := &Client{}
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	secret := os.Getenv("TWITTER_SECRET")
	c.apiKey = base64.StdEncoding.EncodeToString([]byte(urlEncode(consumerKey) + ":" + urlEncode(secret)))
	c.c = &http.Client{}
	return c
}

// GetAccessToken get twitter access token via application-only
// authentication service provided by twitter.
//
func (c *Client) GetAccessToken() error {
	if c.apiKey == "" {
		return errors.New("No API key set for twitter access.")
	}

	urlStr := "https://api.twitter.com/oauth2/token"
	v := url.Values{}
	v.Add("grant_type", "client_credentials")
	body := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+c.apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset-UTF-8")

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(fmt.Sprintf("%d:%s\n", res.StatusCode, res.Status))

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println(string(respBody))

	var authToken OauthResponse
	if err := json.Unmarshal(respBody, &authToken); err != nil {
		return err
	}
	log.Println(authToken.AccessToken)

	c.token = authToken.AccessToken

	return nil
}

// SearchTweets get tweets with searching by given query string.
//
func (c *Client) SearchTweets(query string) error {
	urlStr := "https://api.twitter.com/1.1/search/tweets.json"

	v := url.Values{}
	v.Add("q", query)
	v.Add("count", "100")

	req, err := http.NewRequest("GET", urlStr, nil)
	// Get request doesn't pass query to URL.
	// query values have to be set directly
	req.URL.RawQuery = v.Encode()

	if err != nil {
		return err
	}

	log.Println(c.token)
	req.Header.Add("Authorization", "Bearer "+c.token)
	log.Println(req)

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(fmt.Sprintf("%d:%s\n", res.StatusCode, res.Status))
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))
	return nil
}

func urlEncode(s string) string {
	v := url.Values{}
	v.Set("", s)
	return v.Encode()[1:]
}
