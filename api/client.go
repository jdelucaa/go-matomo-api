package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type apiClient struct {
	apiUrl    *url.URL
	authToken string
	client    *http.Client
	UserAgent string

	Sites *SitesService
}

const (
	userAgent = "go-matomo-api/api/" + LibraryVersion
)

func (c *apiClient) SetApiUrl(urlStr string) error {
	if urlStr == "" {
		return ErrApiUrlCannotBeEmpty
	}
	// Make sure the given URL ends with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.apiUrl, err = url.Parse(urlStr)
	return err
}

func (c *apiClient) SetAuthToken(authToken string) error {
	if authToken == "" {
		return ErrTokenAuthCannotBeEmpty
	}
	c.authToken = authToken

	return nil
}

func newClient(httpClient *http.Client, apiUrl string, authToken string) (*apiClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	c := &apiClient{}
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

func (c *apiClient) newRequest(module string, method string, opt interface{}) (*http.Request, error) {
	var u url.URL
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
