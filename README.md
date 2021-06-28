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

  batch_push:
    runs-on: ubuntu-latest
    name: Update protobuf schema in Batch
    needs: build
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Push new schemas
        uses: batchcorp/schema-publisher@9b53662cfca8785253b91f160f4eca5faceb6f37
        id: publish
        with:
          api_token: '${{ secrets.BATCH_API_TOKEN }}'
          schema_id: '${{ secrets.BATCH_SCHEMA_ID }}'
          schema_name: '${{ github.repository }}: ${{ needs.build.outputs.new_tag }}'
          schema_type: protobuf
          root_message: events.Message
          root_dir: events
          api_address: 'https://api.dev.batch.sh'
```

