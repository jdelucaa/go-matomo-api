package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	API = "API"
)

type ApiClient struct {
	apiUrl    *url.URL
	authToken string
	client    *http.Client
	UserAgent string

	Sites *SitesService
}

const (
	userAgent = "go-matomo-api/api/" + LibraryVersion
)

func (c *ApiClient) SetApiUrl(urlStr string) error {
	if urlStr == "" {
		return ErrApiUrlCannotBeEmpty
	}
	var err error
	c.apiUrl, err = url.Parse(urlStr)
	return err
}

func (c *ApiClient) SetAuthToken(authToken string) error {
	if authToken == "" {
		return ErrTokenAuthCannotBeEmpty
	}
	c.authToken = authToken

	return nil
}

func NewClient(httpClient *http.Client, apiUrl string, authToken string) (*ApiClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	c := &ApiClient{}
	if err := c.SetApiUrl(apiUrl); err != nil {
		return nil, err
	}
	if err := c.SetAuthToken(authToken); err != nil {
		return nil, err
	}

	c.client = httpClient
	c.UserAgent = userAgent
	c.Sites = &SitesService{client: c}

	return c, nil
}

// StardardOpt represents the opt present in all queries
type StandardReqOpt struct {
	Module    string `url:"module"`
	Method    string `url:"method"`
	Format    string `url:"format"`
	AuthToken string `url:"token_auth"`
}

func (c *ApiClient) newRequest(module string, method string, opt interface{}) (*http.Request, error) {
	var u = *c.apiUrl
	standardReqOpt := &StandardReqOpt{
		Module:    module,
		Method:    method,
		Format:    "JSON",
		AuthToken: c.authToken,
	}

	sq, err := query.Values(standardReqOpt)
	if err != nil {
		return nil, err
	}
	u.RawQuery = sq.Encode()

	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = u.RawQuery + "&" + q.Encode()
	}

	req := &http.Request{
		Method:     "POST",
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// Response is a Matomo API response. This wraps the standard http.Response
// returned from Matomo API and provides convenient access to things like errors
type Response struct {
	*http.Response
}

// do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *ApiClient) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	if v != nil && response.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}
