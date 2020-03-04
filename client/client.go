package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultHost      = "https://scrapbox.io"
	DefaultUserAgent = "Scrapgox/0.1.0"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Token      string
	UserAgent  string
}

func (c *Client) buildRequest(method, path string, body io.Reader) (*http.Request, error) {
	baseURL := *c.URL
	u := fmt.Sprintf("%s/%s", baseURL.String(), path)

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)
	if len(c.Token) != 0 {
		req.Header.Set("Cookie", "connect.sid="+c.Token)
	}

	return req, nil
}

func buildPath(project string, skip, limit int, query string) string {
	escapedQuery := url.QueryEscape(query)
	params := fmt.Sprintf("skip=%d&limit=%d&q=%s", skip, limit, escapedQuery)
	if len(query) == 0 {
		return fmt.Sprintf("api/pages/%s?%s", project, params)
	} else {
		return fmt.Sprintf("api/pages/%s/search/query?%s", project, params)
	}
}

type Response struct {
	ProjectName string `json:"projectName"`
	SearchQuery string `json:"searchQuery"`
	Limit       int    `json:"limit"`
	Count       int    `json:"count"`
	Pages       []struct {
		ID           string   `json:"id"`
		Title        string   `json:"title"`
		Image        string   `json:"image"`
		Descriptions []string `json:"descriptions"`
		User         struct {
			ID string `json:"id"`
		} `json:"user"`
		Pin             int         `json:"pin"`
		Views           int         `json:"views"`
		Linked          int         `json:"linked"`
		Created         int         `json:"created"`
		Updated         int         `json:"updated"`
		Accessed        int         `json:"accessed"`
		SnapshotCreated interface{} `json:"snapshotCreated"`
		Snipet          []string    `json:"snipet"`
	} `json:"pages"`
	ExistsExactTitleMatch bool `json:"existsExactTitleMatch"`
	Query                 struct {
		Words    []string      `json:"words"`
		Excludes []interface{} `json:"excludes"`
	} `json:"query"`
}

type Page struct {
	Title string
}

func (r Response) getPages() []*Page {
	pages := make([]*Page, 0, 30)
	if len(r.Pages) == 0 {
		return pages
	} else {
		for _, v := range r.Pages {
			pages = append(pages, &Page{v.Title})
		}
		return pages
	}
}

func decodeResponse(body io.Reader, value *Response) error {
	return json.NewDecoder(body).Decode(value)
}

func NewClient(url *url.URL, token string, userAgent string) (*Client, error) {
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{},
		Token:      token,
		UserAgent:  userAgent,
	}, nil
}

func (c *Client) GetPages(project string, query string) ([]*Page, error) {
	path := buildPath(project, 0, 30, query)
	req, err := c.buildRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("http status is %q", resp.Status))
	}

	var body Response
	decodeResponse(resp.Body, &body)

	pages := body.getPages()
	return pages, nil
}
