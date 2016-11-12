package twitter

import (
	"os"
	"testing"
)

const (
	consumerKey = "TWITTER_CONSUMER_KEY"
	secret      = "TWITTER_SECRET"
)

func TestNewFromEnv(t *testing.T) {
	cKey := os.Getenv(consumerKey)
	scr := os.Getenv(secret)

	os.Setenv(consumerKey, "test%")
	os.Setenv(secret, "ab11")

	defer os.Setenv(consumerKey, cKey)
	defer os.Setenv(secret, scr)

	c := NewFromEnv()

	expect := "dGVzdCUyNTphYjEx"
	if c.apiKey != expect {
		t.Errorf("Encoding api key failed. Expect: %s, actual: %s.", expect, c.apiKey)
	}
}

func TestGetAccessToken(t *testing.T) {
	if !existEnvValues() {
		t.Logf("GetAccessToken test skipped due to not setting environment values.")
		return
	}
	c := NewFromEnv()

	err := c.GetAccessToken()
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestSearchTweets(t *testing.T) {
	if !existEnvValues() {
		t.Logf("SearchTweets test skipped due to not setting environment values.")
	}
	c := NewFromEnv()

	if err := c.GetAccessToken(); err != nil {
		t.Errorf("%v\n", err)
	}

	query := "\"楽天銀行\" -RT"
	if err := c.SearchTweets(query); err != nil {
		t.Errorf("%v\n", err)
	}

}

func existEnvValues() bool {
	return os.Getenv(consumerKey) != "" && os.Getenv(secret) != ""
}
