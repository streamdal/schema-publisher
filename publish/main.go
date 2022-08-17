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
	SchemaID          string
	SchemaType        string
	SchemaName        string
	APIToken          string
	RootDir           string
	APIAddress        string
	OutputFile        string
	ArtifactType      string
	DescriptorSetPath string
}

func init() {
	opts = &Options{}

	flag.StringVar(&opts.SchemaID, "schema-id", envar("BATCH_SCHEMA_ID", ""), "Schema ID that will be updated")
	flag.StringVar(&opts.SchemaType, "schema-type", envar("BATCH_SCHEMA_TYPE", ""), "'protobuf' only currently")
	flag.StringVar(&opts.SchemaName, "schema-name", envar("BATCH_SCHEMA_NAME", ""), "Specify schema release name")
	flag.StringVar(&opts.APIToken, "api-token", envar("BATCH_API_TOKEN", ""), "Batch API token")
	flag.StringVar(&opts.RootDir, "root-dir", envar("BATCH_ROOT_DIR", ""), "Which directory to treat as root for schemas")
	flag.StringVar(&opts.APIAddress, "api-address", envar("BATCH_API_ADDRESS", ""), "HTTP address for Batch API")
	flag.StringVar(&opts.OutputFile, "output", envar("BATCH_OUTPUT_FILE", ""), "Optional output file")
	flag.StringVar(&opts.ArtifactType, "artifact-type", envar("BATCH_ARTIFACT_TYPE", "protos_archive"), "Type of artifact to publish (protos_archive, descriptor_set)")
	flag.StringVar(&opts.DescriptorSetPath, "descriptor-set-path", envar("BATCH_DESCRIPTOR_SET_PATH", ""), "Path to descriptor set file if using --artifact-type=descriptor_set")

	flag.Parse()
}

func envar(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
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

	var archive []byte

	if opts.ArtifactType == "protos_archive" {
		archive, err = createZip(opts.RootDir)
		if err != nil {
			log.Fatalf("unable to create zip archive: '%s'", err)
		}
	} else if opts.ArtifactType == "descriptor_set" {
		// Verify file exists
		if _, err := os.Stat(opts.DescriptorSetPath); err != nil {
			log.Fatalf("could not find descriptor set file '%s'", opts.DescriptorSetPath)
		}

		archive, err = os.ReadFile(opts.DescriptorSetPath)
		if err != nil {
			log.Fatalf("unable to read descriptor set file '%s': %s", opts.DescriptorSetPath, err)
		}
	} else {
		log.Fatalf("unknown artifact type: '%s'", opts.ArtifactType)
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

	if opts.ArtifactType == "protos_archive" {
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
	} else if opts.ArtifactType == "descriptor_set" {
		if opts.DescriptorSetPath == "" {
			return errors.New("opts.DescriptorSetPath cannot be empty when using --artifact-type=descriptor_set")
		}
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

	if opts.ArtifactType == "" {
		return errors.New("opts.ArtifactType cannot be empty")
	}

	if _, err := url.Parse(opts.APIAddress); err != nil {
		return fmt.Errorf("unable to parse opts.APIAddress as URL: %s", err)
	}

	return nil
}
