# brightcove-metadata-mapper
[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/brightcove-metadata-notifier/badge.svg?branch=master)](https://coveralls.io/github/Financial-Times/brightcove-metadata-notifier?branch=master)
[![Circle CI](https://circleci.com/gh/Financial-Times/brightcove-metadata-notifier/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/brightcove-metadata-notifier/tree/master)

Maps brightcove tags to appropriate TME IDs and sends the created metadata publish message to the cms-metadata-notifier.

##Installation

`go get github.com/Financial-Times/brightcove-metadata-notifier`

##Run the binary

```
go build
export MAPPING_URL="https://bertha.ig.ft.com/view/publish/gss/1UTzLQPgg_POaBZeOfAQoxY5YmjKxMqdP9O9oa_Fto9s/mappings"
export CMS_METADATA_NOTIFIER_ADDR="https://pub-xp-up.ft.com/__cms-notifier"
export CMS_METADATA_NOTIFIER_HOST_HEADER="cms-metadata-notifier"
export CMS_METADATA_NOTIFIER_AUTH="Basic dXB..."
./brightcove-metadata-notifier
```

## Endpoints

### /notify

The main operation. It maps brightcove tags to appropriate TME IDs and sends the created metadata publish message to the cms-metadata-notifier.

### POST /notify

Brightcove metadata.
* tags: the tags to be mapped

### /__reload

Reload the tags mappings loaded in application by querying the remote endpoint (set with MAPPING_URL).

### POST /__reload

Examples:
```
curl -X POST -H "Content-Type: application/json" localhost:8080/notify --data '{"uuid":"123", "tags":["brazil"]}'

curl -X POST -H "Content-Type: application/json" localhost:8080/__reload
```