package api

import (
	"bytes"
	"net/http"
)

// Site represents a Site resource
type Site struct {
	ID       string `json:"idSite,omitempty"`
	Name     string `json:"siteName"`
	Timezone string `json:"timezone"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
}

type SitesService struct {
	client *apiClient
}

const (
	GetSiteFromID = "SitesManager.getSiteFromId"
	AddSite       = "SitesManager.addSite"
	DeleteSite    = "SitesManager.deleteSite"
)

// GetSiteByID retrieves a site by ID
func (p *SitesService) GetSiteByID(idSite string) (*Site, *Response, error) {
	siteOpt := &Site{
		ID: idSite,
	}
	req, err := p.client.newRequest(API, GetSiteFromID, siteOpt)
	if err != nil {
		return nil, nil, err
	}

	var site Site

	resp, err := p.client.do(req, &site)
	if err != nil {
		return nil, resp, err
	}
	if site.ID != idSite {
		return nil, resp, ErrNotFound
	}
	return &site, resp, err
}

// CreateSite creates a Site
func (p *SitesService) CreateSite(name string) (*Site, *Response, error) {
	site := &Site{
		Name: name,
	}
	req, _ := p.client.newRequest(API, AddSite, site)

	var createdSite Site

	resp, err := p.client.do(req, &createdSite)
	if err != nil {
		return nil, resp, err
	}
	return &createdSite, resp, err
}

// DeleteSite deletes the given Site
func (p *SitesService) DeleteSite(site Site) (bool, *Response, error) {
	req, err := p.client.newRequest(API, DeleteSite, site)
	if err != nil {
		return false, nil, err
	}

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		return false, resp, err
	}
	return true, resp, nil
}
