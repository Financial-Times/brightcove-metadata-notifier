# brightcove-metadata-mapper
[![Circle CI](https://circleci.com/gh/Financial-Times/brightcove-metadata-notifier/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/brightcove-metadata-notifier/tree/master)

Maps brightcove tags to appropriate TME IDs and sends the created metadata publish msg to the cms-metadata-notifier

##Run the binary

```
go build
export MAPPING_URL="https://bertha.ig.ft.com/view/publish/gss/1UTzLQPgg_POaBZeOfAQoxY5YmjKxMqdP9O9oa_Fto9s/mappings"
export CMS_METADATA_NOTIFIER_ADDR="https://pub-xp-up.ft.com/__cms-notifier"
export CMS_METADATA_NOTIFIER_HOST_HEADER="cms-metadata-notifier"
export CMS_METADATA_NOTIFIER_AUTH="Basic dXB..."
./brightcove-metadata-notifier
```
