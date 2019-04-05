// Package gophish provides a convenient wrapper around the Phish.Net API.
package gophish

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	dateFormat = "2006-01-02"
	baseUrl    = "https://api.phish.net/v3"
	// By default, we will limit queries to 120/min.
	defaultQueryRate = time.Minute / 120
)

// Client is a http.Client wrapper that handles rate limiting as well as
// provides convenience methods for various Phish.Net API endpoints.
type Client struct {
	apiKey    string
	queryRate time.Duration
	client    *http.Client
	baseUrl   string
	throttle  <-chan time.Time
}

// ClientOpt is a option configuring the Client.
type ClientOpt func(*Client)

// WithQueryRate sets the rate limiting that should be used when making requests
// with the client.
func WithQueryRate(rate time.Duration) ClientOpt {
	return func(c *Client) {
		c.queryRate = rate
	}
}

// WithTimeout sets the default timeout for making HTTP requests.
func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithBaseUrl sets the baseurl that should be used for requests. The default is https://api.phish.net/v3.
func WithBaseUrl(url string) ClientOpt {
	return func(c *Client) {
		c.baseUrl = url
	}
}

// NewClient returns a new Client for making requests to the Phish.Net API.
func NewClient(apiKey string, opts ...ClientOpt) *Client {
	c := &Client{
		apiKey:    apiKey,
		queryRate: defaultQueryRate,
		client:    new(http.Client),
		baseUrl:   baseUrl,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.throttle = time.Tick(c.queryRate)
	return c
}

// ParseDate parses the date string into a time.Time for use in the Phish.Net
// API.
func ParseDate(date string) (time.Time, error) {
	return time.Parse(dateFormat, date)
}

// FormatDate converts a time.Time into the correct format that the Phish.Net API expects.
func FormatDate(date time.Time) string {
	return date.Format(dateFormat)
}

func (c *Client) do(method, apiPath string, req, resp interface{}) error {
	vals, err := query.Values(req)
	if err != nil {
		return errors.Wrapf(err, "unable to convert %v to URL values", req)
	}
	vals.Add("apikey", c.apiKey)
	urlString := c.baseUrl + apiPath + "?" + vals.Encode()
	hreq, err := http.NewRequest(method, urlString, nil)
	if err != nil {
		return errors.Wrapf(err, "unable to create http request for url %s", urlString)
	}
	res, err := c.Do(hreq)
	if err != nil {
		return errors.Wrap(err, "unable to perform the API lookup")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read body")
	}
	if res.StatusCode != 200 {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return errors.Wrap(err, "unable to unmarshal json error")
		}
		return errResp
	}
	return errors.Wrap(json.Unmarshal(body, resp), "unable to umarshal json response")
}

// Do sends a throttled request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	<-c.throttle
	return c.client.Do(req)
}

type ErrorResponseBody struct {
	Message string                 `json:"message"`
	Body    map[string]interface{} `json:"body"`
}

// ErrorResponse is the standard error format for API errors.
type ErrorResponse struct {
	ErrorCode int                `json:"error"`
	Response  *ErrorResponseBody `json:"response"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("API Error: %v", e)
}

type ResponseHeader struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// ShowsQueryRequest is the request type for the shows/query endpoint.
type ShowsQueryRequest struct {
	ShowIds     []string `url:"showids,omitempty"`
	Year        int      `url:"year,omitempty"`
	Month       int      `url:"month,omitempty"`
	Day         int      `url:"day,omitempty"`
	VenueId     int      `url:"venueid,omitempty"`
	TourId      int      `url:"tourid,omitempty"`
	Country     string   `url:"country,omitempty"`
	City        string   `url:"city,omitempty"`
	State       string   `url:"state,omitempty"`
	ShowdateGt  string   `url:"showdate_gt,omitempty"`
	ShowdateGte string   `url:"showdate_gte,omitempty"`
	ShowdateLt  string   `url:"showdate_lt,omitempty"`
	ShowdateLte string   `url:"showdate_lte,omitempty"`
	ShowyearGt  int      `url:"showyear_gt,omitempty"`
	ShowyearGte int      `url:"showyear_gte,omitempty"`
	ShowyearLt  int      `url:"showyear_lt,omitempty"`
	ShowyearLte int      `url:"showyear_lte,omitempty"`
	Limit       int      `url:"limit,omitempty"`
	Order       string   `url:"order,omitempty"`
}

type Show struct {
	ShowId       int    `json:"showid"`
	ShowDate     string `json:"showdate"`
	ArtistId     int    `json:"artistid"`
	BilledAs     string `json:"billed_as"`
	Link         string `json:"link"`
	Location     string `json:"location"`
	Venue        string `json:"venue"`
	SetlistNotes string `json:"setlistnotes"`
	VenueId      int    `json:"venueid"`
	TourId       int    `json:"tourid"`
	TourName     string `json:"tourname"`
	TourWhen     string `json:"tourwhen"`
	ArtistLink   string `json:"artistlink"`
}

type ShowsQueryResponse struct {
	*ResponseHeader
	Response *ShowsQueryResponseBody `json:"response"`
}

type ShowsQueryResponseBody struct {
	Count int     `json:"count"`
	Data  []*Show `json:"data"`
}

func (c *Client) ShowsQuery(req *ShowsQueryRequest) (*ShowsQueryResponse, error) {
	var resp ShowsQueryResponse
	if err := c.do(http.MethodPost, "/shows/query", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type SetlistsGetRequest struct {
	ShowId   int    `url:"showid,omitempty"`
	ShowDate string `url:"showdate,omitempty"`
}

type SetlistsRecentRequest struct {
	Limit int `url:"limit,omitempty"`
}

type SetlistsResponse struct {
	*ResponseHeader
	Response *SetlistsResponseBody `json:"response"`
}

type SetlistsResponseBody struct {
	Count int        `json:"count"`
	Data  []*Setlist `json:"data"`
}

type Setlist struct {
	ShowId       int    `json:"showid"`
	ShowDate     string `json:"showdate"`
	ShortDate    string `json:"short_date"`
	LongDate     string `json:"long_date"`
	RelativeDate string `json:"relative_date"`
	Url          string `json:"url"`
	GapChart     string `json:"gapchart"`
	Artist       string `json:"artist"`
	ArtistId     int    `json:"artistid"`
	VenueId      int    `json:"venueid"`
	Venue        string `json:"venue"`
	Location     string `json:"location"`
	SetlistData  string `json:"setlistdata"`
	SetlistNotes string `json:"setlistnotes"`
	Rating       string `json:"rating"`
}

func (c *Client) SetlistsGet(req *SetlistsGetRequest) (*SetlistsResponse, error) {
	var resp SetlistsResponse
	if err := c.do(http.MethodGet, "/setlists/get", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetlistsLatest() (*SetlistsResponse, error) {
	var resp SetlistsResponse
	if err := c.do(http.MethodGet, "/setlists/latest", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetlistsRecent(req *SetlistsRecentRequest) (*SetlistsResponse, error) {
	var resp SetlistsResponse
	if err := c.do(http.MethodGet, "/setlists/recent", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetlistsTiph() (*SetlistsResponse, error) {
	var resp SetlistsResponse
	if err := c.do(http.MethodGet, "/setlists/tiph", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetlistsRandom() (*SetlistsResponse, error) {
	var resp SetlistsResponse
	if err := c.do(http.MethodGet, "/setlists/random", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
