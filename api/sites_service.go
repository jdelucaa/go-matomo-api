package api

import (
	"bytes"
	"net/http"
)

// Site represents a Site resource
type Site struct {
	ID    string `json:"idSite,omitempty"`
	Name  string `json:"siteName"`
	Value int    `json:"value,omitempty"`
}

// Site represents a Site resource
type SiteOptions struct {
	ID   string `url:"idSite,omitempty"`
	Name string `url:"siteName"`
}

// Site represents a Site resource
type GetSitesOptions struct {
	Pattern *string `url:"pattern"`
}

type SitesService struct {
	client *apiClient
}

const (
	GetSiteFromID        = "SitesManager.getSiteFromId"
	AddSite              = "SitesManager.addSite"
	UpdateSite           = "SitesManager.updateSite"
	DeleteSite           = "SitesManager.deleteSite"
	GetPatternMatchSites = "SitesManager.getPatternMatchSites"
)

// GetSiteByID retrieves a site by ID
func (p *SitesService) GetSiteByID(idSite string) (*Site, *Response, error) {
	siteOpt := &SiteOptions{
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
func (p *SitesService) CreateSite(opt *SiteOptions) (*Site, *Response, error) {
	req, _ := p.client.newRequest(API, AddSite, opt)

	var createdSite Site

	resp, err := p.client.do(req, &createdSite)
	if err != nil {
		return nil, resp, err
	}
	return &createdSite, resp, err
}

// CreateSite updates a Site
func (p *SitesService) UpdateSite(opt *SiteOptions) (*Site, *Response, error) {
	req, _ := p.client.newRequest(API, UpdateSite, opt)

	var updatedSite Site

	resp, err := p.client.do(req, &updatedSite)
	if err != nil {
		return nil, resp, err
	}
	return &updatedSite, resp, err
}

// DeleteSite deletes the given Site
func (p *SitesService) DeleteSite(idSite string) (bool, *Response, error) {
	siteOpt := &SiteOptions{
		ID: idSite,
	}
	req, err := p.client.newRequest(API, DeleteSite, siteOpt)
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

// GetSites retrieves sites by pattern
func (p *SitesService) GetSites(opt *GetSitesOptions) (*[]Site, *Response, error) {
	req, err := p.client.newRequest(API, GetPatternMatchSites, opt)
	if err != nil {
		return nil, nil, err
	}

	var sites []Site

	resp, err := p.client.do(req, &sites)
	if err != nil {
		return nil, resp, err
	}
	return &sites, resp, err
}
