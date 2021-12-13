# go-matomo-api

This is a Go client library to interact with the [Matomo APIs](https://developer.matomo.org/api-reference) - in early development.

## Supported APIs

- [x] SitesManager
  - [x] getSiteFromId
  - [x] addSite
  - [x] updateSite
  - [x] deleteSite
  - [x] getPatternMatchSites

Only the `name` attribute is supported for Sites at the moment.