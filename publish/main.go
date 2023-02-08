package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

const (
	DefaultAPIAddress      = "https://api.streamdal.com"
	InputTypeDescriptorSet = "descriptor_set"
	InputTypeProtosArchive = "protos_archive"
	InputTypeDir           = "dir"
	SchemaTypeProtobuf     = "protobuf"
	SchemaTypeAvro         = "avro"
)

var (
	opts *Options
)

type Options struct {
	SchemaID   string
	SchemaType string
	SchemaName string
	APIToken   string
	APIAddress string
	Input      string
	InputType  string
	OutputFile string
}

func init() {
	opts = &Options{}

	flag.StringVar(&opts.SchemaID, "schema-id",
		envar("STREAMDAL_SCHEMA_ID", ""),
		"Schema ID that will be updated")

	flag.StringVar(&opts.SchemaType, "schema-type",
		envar("STREAMDAL_SCHEMA_TYPE", SchemaTypeProtobuf),
		fmt.Sprintf("%s OR %s", SchemaTypeProtobuf, SchemaTypeAvro))

	flag.StringVar(&opts.SchemaName, "schema-name",
		envar("STREAMDAL_SCHEMA_NAME", ""),
		"Updated schema name")

	flag.StringVar(&opts.APIToken, "api-token",
		envar("STREAMDAL_API_TOKEN", ""),
		"Streamdal API token (dashboard -> account -> security)")

	flag.StringVar(&opts.InputType, "input-type",
		envar("STREAMDAL_INPUT_TYPE", InputTypeDescriptorSet),
		fmt.Sprintf("Type of data input (valid '%s', '%s')", InputTypeDescriptorSet, InputTypeDir))

	flag.StringVar(&opts.Input, "input",
		envar("STREAMDAL_INPUT", ""),
		"Input file (descriptor set) or directory")

	flag.StringVar(&opts.APIAddress, "api-address",
		envar("STREAMDAL_API_ADDRESS", DefaultAPIAddress),
		"HTTP address for Streamdal API")

	flag.StringVar(&opts.OutputFile, "output",
		envar("STREAMDAL_OUTPUT", ""),
		"Optional output file (only used with 'dir' -input-type)")

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
		log.Fatalf("unable to fetch existing schema id '%s': %s", opts.SchemaID, err)
	}

	if sch.Type != opts.SchemaType {
		log.Fatalf("mismatching schema types: opts.SchemaType is '%s' while Batch schema id '%s' is '%s'",
			opts.SchemaType, opts.SchemaID, sch.Type)
	}
	archive, err := generateArchive(opts)
	if err != nil {
		log.Fatalf("unable to generate archive: '%s'", err)
	}

	if opts.InputType == InputTypeProtosArchive && opts.OutputFile != "" {
		if err := os.WriteFile(opts.OutputFile, archive, 0600); err != nil {
			log.Fatalf("unable to write file: %s", err)
		}
	}

	if _, err := apiClient.UpdateSchema(opts, archive); err != nil {
		log.Fatalf("unable to new schema upload: %s", err)
	}

	fmt.Printf("Schema updated!\nSchema ID: %s\nSchema name: %s\n", opts.SchemaID, opts.SchemaName)
}

func envar(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func validateOptions(opts *Options) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}

	if opts.Input == "" {
		return errors.New("-input must be set")
	}

	switch opts.InputType {
	case InputTypeDir:
		opts.InputType = InputTypeProtosArchive

		fallthrough
	case InputTypeProtosArchive:
		fi, err := os.Stat(opts.Input)
		if err != nil {
			return fmt.Errorf("unable to stat '%s': %s", opts.Input, err)
		}

		if !fi.IsDir() {
			return fmt.Errorf("'%s' is not a directory", opts.Input)
		}
	case InputTypeDescriptorSet:
		if _, err := os.Stat(opts.Input); err != nil {
			return fmt.Errorf("unable to stat descriptor set file '%s': %s", opts.Input, err)
		}
	default:
		return fmt.Errorf("unrecognized -input-type '%s'", opts.InputType)
	}

	if opts.APIToken == "" {
		return errors.New("-api-token cannot be empty")
	}

	if opts.SchemaType != "protobuf" {
		return errors.New("-schema-type not supported")
	}

	if opts.APIAddress == "" {
		opts.APIAddress = DefaultAPIAddress
	}

	if _, err := url.Parse(opts.APIAddress); err != nil {
		return fmt.Errorf("unable to parse opts.APIAddress as URL: %s", err)
	}

	return nil
}

func generateArchive(opts *Options) ([]byte, error) {
	var (
		archive []byte
		err     error
	)

	if opts.InputType == InputTypeProtosArchive {
		archive, err = createZip(opts.Input)
	} else if opts.InputType == InputTypeDescriptorSet {
		archive, err = os.ReadFile(opts.Input)
	} else {
		return nil, errors.New("unsupported input type")
	}

	if err != nil {
		return nil, fmt.Errorf("unable to generate archive: %s", err)
	}

	return []byte(base64.StdEncoding.EncodeToString(archive)), nil
}
