---
layout: "powerdns"
page_title: "PowerDNS: powerdns_zone_metadata"
sidebar_current: "docs-powerdns-resource-zone-metadata"
description: |-
  Provides a PowerDNS zone metadata resource.
---

# powerdns\_zone\_metadata

Provides a PowerDNS zone metadata resource.

## Example Usage

All possible zone metadata can be found [here](https://doc.powerdns.com/authoritative/domainmetadata.html#).

### ALLOW-AXFR-FROM example
For the v1 API (PowerDNS version 4):

```hcl
# Add ALLOW-AXFR-FROM metadata to the zone
resource "powerdns_zone_metadata" "foobar" {
  zone     = "example.com."
  kind     = "ALLOW-AXFR-FROM"
  metadata = ["AUTO-NS", "10.0.0.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The name of zone to contain this metadata.
* `kind` - (Required) The kind of the metadata.
* `metadata` - (Required) A string list of metadata.

### Attribute Reference

The id of the resource is a composite of the zone name and metadata kind, joined by a separator - `:::`.

For example, metadata in zone `foo.test.com.` of kind `ALLOW-AXFR-FROM` will be represented with the following `id`: `foo.test.com.:::ALLOW-AXFR-FROM`

### Importing

An existing record can be imported into this resource by supplying both the zone name and metadata kind it belongs to.
If the kind or zone is not found, or if the record is of a different type or in a different zone, an error will be returned.

For example:

```
$ terraform import powerdns_zone_metadata.test-a test.com.:::AXFR-SOURCE
```

For more information on how to use terraform's `import` command, please refer to terraform's [core documentation](https://www.terraform.io/docs/import/index.html#currently-state-only).

