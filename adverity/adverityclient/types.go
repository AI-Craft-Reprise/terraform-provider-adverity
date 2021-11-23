package adverityclient

import (
	"net/http"
	"net/url"
)

type Client struct {
	token          string
	restURL        *url.URL
	requestsParams map[string]string
	httpClient     *http.Client
}

type errorString struct {
	s string
}

func (e errorString) Error() string {
	return e.s
}

type CreateWorkspaceConfig struct {
	DatalakeID string `json:"datalake_id,omitempty url:"datalake_id,omitempty"`
	Name       string `json:"name,omitempty" url:"name,omitempty"`
	ParentID   int    `json:"parent_id,omitempty url:"parent_id,omitempty"`
}

type UpdateWorkspaceConfig struct {
	DatalakeID string `json:"datalake_id,omitempty url:"datalake_id,omitempty"`
	ParentID   int    `json:"parent_id,omitempty url:"parent_id,omitempty"`
	StackSlug  string `json:"stack_slug,omitempty url:"stack_slug,omitempty"`
	Name       string `json:"name,omitempty" url:"name,omitempty"`
}

type DeleteWorkspaceConfig struct {
	StackSlug string `json:"stack_slug,omitempty url:"stack_slug,omitempty"`
}

type Parameters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ParametersListInt struct {
	Name  string `json:"name"`
	Value []int  `json:"value"`
}

type ParametersListStr struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type ConnectionConfig struct {
	Name       string        `json:"name"`
	Stack      int           `json:"stack"`
	Parameters []*Parameters `json:"parameters"`
}

type DatastreamConfig struct {
	Name              string               `json:"name"`
	Stack             int                  `json:"stack"`
	Parameters        []*Parameters        `json:"parameters"`
	ParametersListInt []*ParametersListInt `json:"parameters_int"`
	ParametersListStr []*ParametersListStr `json:"parameters_str"`
}

type DestinationConfig struct {
	Name          string `json:"name"`
	Stack         int    `json:"stack"`
	ProjectID     string `json:"project"`
	DatasetID     string `json:"dataset"`
	Auth          int    `json:"auth"`
	SchemaMapping bool   `json:"schema_mapping"`
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
	DefaultManageExtractNames bool   `json:"default_manage_extract_names"`
	Updated                   string `json:"updated"`
	Created                   string `json:"created"`
}

type Connection struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	MetadataSlack int    `json:"metadata_slack"`
	Stack         int    `json:"stack"`
	App           int    `json:"app"`
	User          int    `json:"user"`
	IsAuthorized  bool   `json:"is_authorized"`
}

type Destination struct {
	ID                      int    `json:"id"`
	LogoURL                 string `json:"logo_url"`
	IsSchemaMappingRequired bool   `json:"is_schema_mapping_required"`
	Name                    string `json:"name"`
	SchemaMapping           bool   `json:"schema_mapping"`
	ForceString             bool   `json:"force_string"`
	FormatHeaders           bool   `json:"format_headers"`
	ColumnNamesToLowercase  bool   `json:"column_names_to_lowercase"`
	Project                 string `json:"project"`
	Dataset                 string `json:"dataset"`
	HeadersFormatting       int    `json:"headers_formatting"`
	Stack                   int    `json:"stack"`
	Auth                    int    `json:"auth"`
}

type Datastream struct {
	ID                 int    `json:"id"`
	CronType           string `json:"cron_type"`
	CronInterval       int    `json:"cron_interval"`
	CronStartOfDay     string `json:"cron_start_of_day"`
	CronIntervalStart  int    `json:"cron_interval_start"`
	TimeRangePreset    int    `json:"time_range_preset"`
	DeltaType          int    `json:"delta_type"`
	DeltaInterval      int    `json:"delta_interval"`
	DeltaIntervalStart int    `json:"delta_interval_start"`
	DeltaStartOfDay    string `json:"delta_start_of_day"`
	Datatype           string `json:"datatype"`
	Creator            string `json:"creator"`
	DatastreamTypeID   int    `json:"datastream_type_id"`
	AbsoluteURL        string `json:"absolute_url"`
	Created            string `json:"created"`
	Updated            string `json:"updated"`
	Slug               string `json:"slug"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Enabled            bool   `json:"enabled"`
	Auth               int    `json:"auth"`
	Frequency          string `json:"frequency"`
	LastFetch          string `json:"last_fetch"`
	NextRun            string `json:"next_run"`
	OverviewURL        string `json:"overview_url"`
	StackID            int    `json:"stack_id"`
}

type StorageConfig struct {
	Name  string `json:"name"`
	Stack int    `json:"stack,omitempty"`
	URL   string `json:"url"`
	Auth  int    `json:"auth"`
}

type Storage struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Stack          int    `json:"stack"`
	URL            string `json:"url"`
	BackupExisting bool   `json:"backup_existing"`
	Auth           int    `json:"auth"`
}