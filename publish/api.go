package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	GetAccountEndpoint = "/v1/account"
	GetSchemaEndpoint  = "/v1/schema"
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
	SchemaID          string `json:"schema_id"`
	Name              string `json:"name"`
	SchemaArchive     string `json:"schema_archive"`
	SchemaArchiveType string `json:"schema_archive_type"`
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

func (a *APIClient) UpdateSchema(opts *Options, archive []byte) (*Schema, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}

	fullGetSchemaEndpoint := opts.APIAddress + GetSchemaEndpoint

	putRequest := &SchemaUpdateRequest{
		SchemaID:          opts.SchemaID,
		Name:              opts.SchemaName,
		SchemaArchiveType: opts.InputType,
		SchemaArchive:     string(archive),
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
