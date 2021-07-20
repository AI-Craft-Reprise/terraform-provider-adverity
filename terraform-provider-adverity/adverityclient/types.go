package adverityclient

import (
	"net/http"
	"net/url"
)


type Client struct {
// 	username, password, bearerToken string
	token string
	restURL                    *url.URL
	requestsParams             map[string]string
	httpClient                 *http.Client
}

type errorString struct {
	s string
}

func (e errorString) Error() string {
	return e.s
}


type CreateWorkspaceConfig struct {
	DatalakeID      string `json:"datalake_id,omitempty url:"datalake_id,omitempty"`
	Name      string `json:"name,omitempty" url:"name,omitempty"`
	ParentID      int `json:"parent_id,omitempty url:"parent_id,omitempty"`
}

type UpdateWorkspaceConfig struct {
	ParentID      int `json:"parent_id,omitempty url:"parent_id,omitempty"`
	StackSlug     string `json:"stack_slug,omitempty url:"stack_slug,omitempty"`
	Name      string `json:"name,omitempty" url:"name,omitempty"`
}

type DeleteWorkspaceConfig struct {
	StackSlug     string `json:"stack_slug,omitempty url:"stack_slug,omitempty"`
}

type ConnectionConfig struct {
	Name         string      `json:"name"`
	Stack        int         `json:"stack"`
	App          int         `json:"app"`
}


type Workspace struct {
	AddConnectionURL string      `json:"add_connection_url"`
	AddDatastreamURL string      `json:"add_datastream_url"`
	ChangeURL        string      `json:"change_url"`
	Datalake         string      `json:"datalake"`
	Destination      interface{} `json:"destination"`
	ExtractsURL      string      `json:"extracts_url"`
	IssuesURL        string      `json:"issues_url"`
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	OverviewURL      string      `json:"overview_url"`
	Parent           string      `json:"parent"`
	ParentID         int         `json:"parent_id"`
	Slug             string      `json:"slug"`
	URL              string      `json:"url"`
	Counts           struct {
		Connections int `json:"connections"`
		Datastreams int `json:"datastreams"`
	} `json:"counts"`
	Permissions struct {
		IsCreator           bool `json:"isCreator"`
		IsDatastreamManager bool `json:"isDatastreamManager"`
		IsViewer            bool `json:"isViewer"`
	} `json:"permissions"`
	DefaultManageExtractNames bool      `json:"default_manage_extract_names"`
	Updated                   string `json:"updated"`
	Created                   string `json:"created"`
}

type Connection struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Stack        int         `json:"stack"`
	App          int         `json:"app"`
	IsAuthorized bool        `json:"is_authorized"`
}