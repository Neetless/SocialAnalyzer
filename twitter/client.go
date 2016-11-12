package twitter

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SearchMetadata struct {
	CompletedIn   float32 `json:"completed_in"`
	MaxId         int64   `json:"max_id"`
	MaxIdString   string  `json:"max_id_str"`
	Query         string  `json:"query"`
	RefreshUrl    string  `json:"refresh_url"`
	Count         int     `json:"count"`
	SinceId       int64   `json:"since_id"`
	SinceIdString string  `json:"since_id_str"`
	NextResults   string  `json:"next_results"`
}

type SearchResponse struct {
	Statuses []Tweet        `json:"statuses"`
	Metadata SearchMetadata `json:"search_metadata"`
}

type UrlEntity struct {
	Urls []struct {
		Indices      []int
		Url          string
		Display_url  string
		Expanded_url string
	}
}

type Entities struct {
	Hashtags []struct {
		Indices []int
		Text    string
	}
	Urls []struct {
		Indices      []int
		Url          string
		Display_url  string
		Expanded_url string
	}
	Url           UrlEntity
	User_mentions []struct {
		Name        string
		Indices     []int
		Screen_name string
		Id          int64
		Id_str      string
	}
	Media []EntityMedia
}

type EntityMedia struct {
	Id                   int64
	Id_str               string
	Media_url            string
	Media_url_https      string
	Url                  string
	Display_url          string
	Expanded_url         string
	Sizes                MediaSizes
	Source_status_id     int64
	Source_status_id_str string
	Type                 string
	Indices              []int
	VideoInfo            VideoInfo `json:"video_info"`
}

type MediaSizes struct {
	Medium MediaSize
	Thumb  MediaSize
	Small  MediaSize
	Large  MediaSize
}

type MediaSize struct {
	W      int
	H      int
	Resize string
}

type VideoInfo struct {
	AspectRatio    []int     `json:"aspect_ratio"`
	DurationMillis int64     `json:"duration_millis"`
	Variants       []Variant `json:"variants"`
}

type Variant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	Url         string `json:"url"`
}

type User struct {
	ContributorsEnabled            bool     `json:"contributors_enabled"`
	CreatedAt                      string   `json:"created_at"`
	DefaultProfile                 bool     `json:"default_profile"`
	DefaultProfileImage            bool     `json:"default_profile_image"`
	Description                    string   `json:"description"`
	Entities                       Entities `json:"entities"`
	FavouritesCount                int      `json:"favourites_count"`
	FollowRequestSent              bool     `json:"follow_request_sent"`
	FollowersCount                 int      `json:"followers_count"`
	Following                      bool     `json:"following"`
	FriendsCount                   int      `json:"friends_count"`
	GeoEnabled                     bool     `json:"geo_enabled"`
	Id                             int64    `json:"id"`
	IdStr                          string   `json:"id_str"`
	IsTranslator                   bool     `json:"is_translator"`
	Lang                           string   `json:"lang"` // BCP-47 code of user defined language
	ListedCount                    int64    `json:"listed_count"`
	Location                       string   `json:"location"` // User defined location
	Name                           string   `json:"name"`
	Notifications                  bool     `json:"notifications"`
	ProfileBackgroundColor         string   `json:"profile_background_color"`
	ProfileBackgroundImageURL      string   `json:"profile_background_image_url"`
	ProfileBackgroundImageUrlHttps string   `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool     `json:"profile_background_tile"`
	ProfileBannerURL               string   `json:"profile_banner_url"`
	ProfileImageURL                string   `json:"profile_image_url"`
	ProfileImageUrlHttps           string   `json:"profile_image_url_https"`
	ProfileLinkColor               string   `json:"profile_link_color"`
	ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string   `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool     `json:"profile_use_background_image"`
	Protected                      bool     `json:"protected"`
	ScreenName                     string   `json:"screen_name"`
	ShowAllInlineMedia             bool     `json:"show_all_inline_media"`
	Status                         *Tweet   `json:"status"` // Only included if the user is a friend
	StatusesCount                  int64    `json:"statuses_count"`
	TimeZone                       string   `json:"time_zone"`
	URL                            string   `json:"url"` // From UTC in seconds
	UtcOffset                      int      `json:"utc_offset"`
	Verified                       bool     `json:"verified"`
	WithheldInCountries            []string `json:"withheld_in_countries"`
	WithheldScope                  string   `json:"withheld_scope"`
}

type Place struct {
	Attributes  map[string]string `json:"attributes"`
	BoundingBox struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"bounding_box"`
	ContainedWithin []struct {
		Attributes  map[string]string `json:"attributes"`
		BoundingBox struct {
			Coordinates [][][]float64 `json:"coordinates"`
			Type        string        `json:"type"`
		} `json:"bounding_box"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		FullName    string `json:"full_name"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		PlaceType   string `json:"place_type"`
		URL         string `json:"url"`
	} `json:"contained_within"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	FullName    string `json:"full_name"`
	Geometry    struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"geometry"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	PlaceType string   `json:"place_type"`
	Polylines []string `json:"polylines"`
	URL       string   `json:"url"`
}

type Tweet struct {
	//Contributors         []Contributor          `json:"contributors"` // Not yet generally available to all, so hard to test
	//Coordinates          *Coordinates           `json:"coordinates"`
	CreatedAt            string                 `json:"created_at"`
	Entities             Entities               `json:"entities"`
	ExtendedEntities     Entities               `json:"extended_entities"`
	FavoriteCount        int                    `json:"favorite_count"`
	Favorited            bool                   `json:"favorited"`
	FilterLevel          string                 `json:"filter_level"`
	Id                   int64                  `json:"id"`
	IdStr                string                 `json:"id_str"`
	InReplyToScreenName  string                 `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64                  `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string                 `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64                  `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string                 `json:"in_reply_to_user_id_str"`
	Lang                 string                 `json:"lang"`
	Place                Place                  `json:"place"`
	QuotedStatusID       int64                  `json:"quoted_status_id"`
	QuotedStatusIdStr    string                 `json:"quoted_status_id_str"`
	QuotedStatus         *Tweet                 `json:"quoted_status"`
	PossiblySensitive    bool                   `json:"possibly_sensitive"`
	RetweetCount         int                    `json:"retweet_count"`
	Retweeted            bool                   `json:"retweeted"`
	RetweetedStatus      *Tweet                 `json:"retweeted_status"`
	Source               string                 `json:"source"`
	Scopes               map[string]interface{} `json:"scopes"`
	Text                 string                 `json:"text"`
	Truncated            bool                   `json:"truncated"`
	User                 User                   `json:"user"`
	WithheldCopyright    bool                   `json:"withheld_copyright"`
	WithheldInCountries  []string               `json:"withheld_in_countries"`
	WithheldScope        string                 `json:"withheld_scope"`

	//Geo is deprecated
	//Geo                  interface{} `json:"geo"`
}

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
	v.Add("result_type", "recent")

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
	var s SearchResponse
	if err := json.Unmarshal(respBody, &s); err != nil {
		return err
	}
	s.toTsv("test.tsv")
	return nil
}

func (s *SearchResponse) show() {
	for _, status := range s.Statuses {
		createdAt := status.CreatedAt
		idStr := status.IdStr
		text := status.Text
		userIdStr := status.User.IdStr
		fmt.Printf("%s\t%s\t%s\t%s\n", createdAt, idStr, text, userIdStr)
	}
}

func (s SearchResponse) toCsv(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	s.toWriter(w, ",")
	w.Flush()
	return nil
}

func (s SearchResponse) toTsv(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	s.toWriter(w, ",\t")
	w.Flush()
	return nil
}

func (s SearchResponse) toWriter(w io.Writer, sep string) {
	for _, t := range s.Statuses {
		strs := []string{
			addQuote(t.IdStr),
			addQuote(""),
			addQuote(strings.Replace(t.Text, "\n", " ", -1)),
		}
		line := strings.Join(strs, sep)
		fmt.Fprint(w, line+"\n")
	}
}

func addQuote(str string) string {
	return "\"" + str + "\""
}

func urlEncode(s string) string {
	v := url.Values{}
	v.Set("", s)
	return v.Encode()[1:]
}
