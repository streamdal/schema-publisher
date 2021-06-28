# schema-publisher

`schema-publisher` is a GitHub action for publishing protobuf schemas to [Batch](https://batch.sh).

This action should be used if you are using either `protobuf` or `avro` schemas
for your collection and want to **avoid** having to manually update the schemas
in the Batch console, every time you make an update.

Your `.github/workflows/foo.yaml` should look something like this:

```yaml
name: Bump version
on:
  push:
    branches:
      - master
jobs:
  batch_push:
    runs-on: ubuntu-latest
    name: A job to push a new schema revision to Batch
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Push new schemas
        uses: batchcorp/schema-publisher
        id: publish
        with:
          # Required
          api_token: ${{ secrets.BATCH_API_TOKEN }}

          # Required
          schema_id: 'batch-schema-id'

          # Required
          schema_type: protobuf

          # Required
          root_message: events.Message

          # Optional (defaults to the below format)
          schema_name_format: "Protobuf: ${{ steps.publish.outputs.git_tag }}"

          # Optional (defaults to ./)
          root_dir: ./events
```

