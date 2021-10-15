package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiClient struct {
	apiUrl    *url.URL
	authToken string
	client    *http.Client
	UserAgent string

	Sites *SitesService
}

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

func newClient(httpClient *http.Client, d *schema.ResourceData, userAgent string) (*apiClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	c := &apiClient{}
	if err := c.SetApiUrl(d.Get("api_url").(string)); err != nil {
		return nil, err
	}
	if err := c.SetAuthToken(d.Get("auth_token").(string)); err != nil {
		return nil, err
	}

	c.client = httpClient
	c.UserAgent = userAgent
	c.Sites = &SitesService{client: c}

	return c, nil
}

func (c *apiClient) newRequest(endpoint string, method string, userAgent string, opt interface{}) (*http.Request, error) {
	var u url.URL
	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	req.Header.Set("User-Agent", userAgent)

	if method == "POST" || method == "PUT" {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}
