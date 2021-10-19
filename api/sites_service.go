package api

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
