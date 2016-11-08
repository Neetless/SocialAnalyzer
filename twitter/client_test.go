package twitter

import (
	"os"
	"testing"
)

func TestNewFromEnv(t *testing.T) {
	os.Setenv("TWITTER_CONSUMER_KEY", "test%")
	os.Setenv("TWITTER_SECRET", "ab11")
	c := NewFromEnv()

	expect := "dGVzdCUyNTphYjEx"
	if c.apiKey != expect {
		t.Errorf("Encoding api key failed. Expect: %s, actual: %s.", expect, c.apiKey)
	}
}

func TestGetAccessToken(t *testing.T) {
	os.Setenv("TWITTER_CONSUMER_KEY", "test%")
	os.Setenv("TWITTER_SECRET", "ab11")
	c := NewFromEnv()

	err := c.GetAccessToken()
	if err != nil {
		t.Errorf("%v\n", err)
	}
}
