# schema-publisher

`schema-publisher` is a GitHub action for publishing schema updates to 
[Streamdal](https://streamdal.com).

This action will allow you to **avoid** having to manually update the schemas 
in the Streamdal console every time your schemas change.

<sub>**NOTE**: While this repo is intended for use with Github Actions, you can
use the `schema-publisher` tool to perform the schema update in any CI system.

## Example Flow

1. Your schemas are in a Github repo
2. Upon creating a PR in your schema repo, CI runs against the PR and tests the schema
3. Upon merging the PR, CI compiles the schemas and generates a _descriptor set_ file
   1. `$ ./protoc --descriptor_set_out=FILE ...`
4. The `schema-publisher` workflow is configured to point to the _descriptor set_
artifact that was created in the previous step
5. `schema-publisher` workflow updates the existing schema in Streamdal with the
new schema and includes the new schema version in the name

<sub>**NOTE**: If using `protobuf` with complex schemas, it is best to use
`descriptor_set` for the `input_type`. While Streamdal and the `schema-publisher`
tool support directory-based protobuf schemas, it is possible to run into errors
during schema updates if the source schemas have sophisticated directory 
structure.</sub>

## Config

| Flag          | Required | Description                                                         | 
|---------------|----------|---------------------------------------------------------------------|
| `api_token`   | **YES**  | API token used for Streamdal API (dashboard -> account -> security) |
| `schema_id`   | **YES**  | Schema ID in Streamdal (dashboard -> collection -> schema           |
| `input`       | **YES**  | File or directory with schema                                       |
| `input_type`  | No       | `descriptor_set` (default) OR `dir`                                 | 
| `schema_name` | No       | New "friendly" schema name (tip: set to current schema version)     |
| `schema_type` | No       | `protobuf` (default) OR `avro`                                      |
| `api_address` | No       | Override the Streamdal API endpoint                                 |
| `output`      | No       | Optionally write `dir` artifact (zip file) to a file                |
| `debug`       | No       | Display additional debug output when running workflow               |

## Example Github Workflow

Your `.github/workflows/foo.yaml` should look something like this:

```yaml
name: Bump version
on:
  push:
    branches:
      - master
jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      new_tag: ${{ steps.create_tag.outputs.new_tag }}
    steps:
      - uses: actions/checkout@master
      - name: Bump version and push tag
        uses: mathieudutour/github-tag-action@v4.5
        id: create_tag
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}

  streamdal_push:
    runs-on: ubuntu-latest
    name: Update protobuf schema in Streamdal
    needs: build
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Push new schemas
        uses: batchcorp/schema-publisher@latest
        id: publish
        with:
          api_token: '${{ secrets.STREAMDAL_API_TOKEN }}'
          schema_id: '${{ secrets.STREAMDAL_SCHEMA_ID }}'
          schema_name: '${{ github.repository }}: ${{ needs.build.outputs.new_tag }}'
          input_type: descriptor_set
          input: protoset.fds
          debug: true
```

