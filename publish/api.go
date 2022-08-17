package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type SchemaArchiveType string

const (
	GetAccountEndpoint = "/v1/account"
	GetSchemaEndpoint  = "/v1/schema"

	ProtosArchive        SchemaArchiveType = "protos_archive"
	DescriptorSetArchive SchemaArchiveType = "descriptor_set"
)

type APIClient struct {
	Options *Options

	Client *http.Client
}

type Schema struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	RootType   string       `json:"root_type"`
	Type       string       `json:"type"`
	TeamId     string       `json:"team_id"`
	Shared     bool         `json:"shared"`
	Archived   bool         `json:"archived"`
	ProtoFiles []*Protofile `json:"proto_files"`
	InsertedAt time.Time    `json:"inserted_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type SchemaUpdateRequest struct {
	SchemaID          string            `json:"schema_id"`
	Name              string            `json:"name"`
	SchemaArchive     string            `json:"schema_archive"`
	SchemaArchiveType SchemaArchiveType `json:"schema_archive_type"`
}

type Protofile struct {
	Id         string    `json:"id"`
	FileName   string    `json:"file_name"`
	Contents   string    `json:"contents"`
	SchemaId   string    `json:"schema_id"`
	InsertedAt time.Time `json:"inserted_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewAPIClient(opts *Options) (*APIClient, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Test that we have a good token
	req, err := http.NewRequest("GET", opts.APIAddress+GetAccountEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to generate new HTTP request: %s", err)
	}

	req.Header.Add("Authorization", "Bearer "+opts.APIToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform HTTP request: %s", err)
	}

	if resp.StatusCode != 200 {
		output, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("received a non-200 response (%d) and cannot read response body: '%s'",
				resp.StatusCode, err)
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("received a non-200 response (%d); output: %s",
			resp.StatusCode, string(output))
	}

	return &APIClient{
		Options: opts,
		Client:  client,
	}, nil
}

func (a *APIClient) GetSchema() (*Schema, error) {
	fullGetSchemaEndpoint := opts.APIAddress + GetSchemaEndpoint + "/" + a.Options.SchemaID

	req, err := http.NewRequest("GET", fullGetSchemaEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to generate new HTTP request: %s", err)
	}

	req.Header.Add("Authorization", "Bearer "+opts.APIToken)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform HTTP request: %s", err)
	}

	if resp.StatusCode != 200 {
		output, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("received a non-200 response (%d) and cannot read response body: '%s'",
				resp.StatusCode, err)
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("received a non-200 response (%d); output: %s",
			resp.StatusCode, string(output))
	}

	// Got a good resp
	sch := &Schema{}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %s", err)
	}

	if err := json.Unmarshal(data, sch); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response into a schema: %s", err)
	}

	return sch, nil
}

func (a *APIClient) UpdateSchema(archive []byte) (*Schema, error) {
	fullGetSchemaEndpoint := opts.APIAddress + GetSchemaEndpoint

	putRequest := &SchemaUpdateRequest{
		SchemaID:          opts.SchemaID,
		Name:              opts.SchemaName,
		SchemaArchive:     base64.StdEncoding.EncodeToString(archive),
		SchemaArchiveType: SchemaArchiveType(opts.ArtifactType),
	}

	requestData, err := json.Marshal(putRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal schema update request: %s", err)
	}

	req, err := http.NewRequest("PUT", fullGetSchemaEndpoint, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, fmt.Errorf("unable to generate new HTTP request: %s", err)
	}

	req.Header.Add("Authorization", "Bearer "+opts.APIToken)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform HTTP request: %s", err)
	}

	if resp.StatusCode != 200 {
		output, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("received a non-200 response (%d) and cannot read response body: '%s'",
				resp.StatusCode, err)
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("received a non-200 response (%d); output: %s",
			resp.StatusCode, string(output))
	}

	sch := &Schema{}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %s", err)
	}

	if err := json.Unmarshal(data, sch); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response into a schema: %s", err)
	}

	return sch, nil
}
