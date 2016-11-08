package twitter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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
