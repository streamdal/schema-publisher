package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

var (
	opts *Options
)

type Options struct {
	SchemaID    string
	SchemaType  string
	SchemaName  string
	APIToken    string
	RootDir     string
	RootMessage string
	APIAddress  string
	OutputFile  string
}

func init() {
	opts = &Options{}

	flag.StringVar(&opts.SchemaID, "schema-id", "", "Schema ID that will be updated")
	flag.StringVar(&opts.SchemaType, "schema-type", "", "protobuf, avro")
	flag.StringVar(&opts.SchemaName, "schema-name", "", "Specify schema release name")
	flag.StringVar(&opts.APIToken, "api-token", "", "Batch API token")
	flag.StringVar(&opts.RootDir, "root-dir", "", "Which directory to treat as root for schemas")
	flag.StringVar(&opts.RootMessage, "root-message", "", "Root message type")
	flag.StringVar(&opts.APIAddress, "api-address", "", "HTTP address for Batch API")
	flag.StringVar(&opts.OutputFile, "output", "", "Optional output file")

	flag.Parse()
}

func main() {
	if err := validateOptions(opts); err != nil {
		log.Fatalf("unable to validate options: %s", err)
	}

	apiClient, err := NewAPIClient(opts)
	if err != nil {
		log.Fatalf("unable to create new API client: %s", err)
	}

	sch, err := apiClient.GetSchema()
	if err != nil {
		log.Fatalf("unable to fetch schema '%s': %s", opts.SchemaID, err)
	}

	if sch.Type != opts.SchemaType {
		log.Fatalf("mismatching schema types: opts.SchemaType is '%s' while Batch schema id '%s' is '%s'",
			opts.SchemaType, opts.SchemaID, sch.Type)
	}

	archive, err := createZip(opts.RootDir)
	if err != nil {
		log.Fatalf("unable to create zip archive: '%s'", err)
	}

	if opts.OutputFile != "" {
		if err := os.WriteFile(opts.OutputFile, archive, 0600); err != nil {
			log.Fatalf("unable to write file: %s", err)
		}
	}

	if _, err := apiClient.UpdateSchema(archive); err != nil {
		log.Fatalf("unable to new schema upload: %s", err)
	}

	fmt.Printf("Schema updated!\nSchema ID: %s\nSchema name: %s\n", opts.SchemaID, opts.SchemaName)
}

func validateOptions(opts *Options) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}

	if opts.RootDir == "" {
		return errors.New("opts.RootDir cannot be empty")
	}

	// Check if root dir exists
	fi, err := os.Stat(opts.RootDir)
	if err != nil {
		return fmt.Errorf("unable to stat '%s': %s", opts.RootDir, err)
	}

	if !fi.IsDir() {
		return fmt.Errorf("'%s' is not a directory", opts.RootDir)
	}

	if opts.APIToken == "" {
		return errors.New("opts.APIToken cannot be empty")
	}

	if opts.SchemaType == "" {
		return errors.New("SchemaType cannot be empty")
	}

	if opts.SchemaType != "protobuf" {
		return errors.New("unsupported schema type")
	}

	if opts.SchemaName == "" {
		return errors.New("opts.SchemaName cannot be empty")
	}

	if opts.APIAddress == "" {
		return errors.New("opts.APIAddress cannot be empty")
	}

	if _, err := url.Parse(opts.APIAddress); err != nil {
		return fmt.Errorf("unable to parse opts.APIAddress as URL: %s", err)
	}

	return nil
}
