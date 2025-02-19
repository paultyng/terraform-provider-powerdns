---
layout: "powerdns"
page_title: "PowerDNS: powerdns_record_soa"
sidebar_current: "docs-powerdns-resource-record-soa"
description: |-
  Provides a PowerDNS SOA record resource.
---

# powerdns\_record

Provides a PowerDNS SOA record resource.
This is offers an alternative to managing the SOA record of a zone through the generic `powerdns_record`.

## Example Usage

### A record example

```hcl
resource "powerdns_record_soa" "soa-example_com" {
  zone    = "example.com."
  name    = "example.com."
  type    = "SOA"
  ttl     = 8600
  mnme    = "ns1.example.com."
  rname   = "hostmaster.example.com."
  serial  = 0
  refresh = 7200
  retry   = 600
  expire  = 1209600
  minimum = 6400

}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The name of zone to contain this record.
* `name` - (Required) The name of the record, typically the same as `zone`.
* `type` - (Required) The record type, must be `SOA`.
* `ttl` - (Required) The TTL of the record.
* `mname` - (Required) SOA MNAME.
* `rname` - (Required) SOA RNAME.
* `serial` - (Required) SOA SERIAL - it will only be used for creation of the zone, subsequent changes are expected to use SOA-EDIT-API if a serial number update is desired. Set to `0` if no specific starting serial number is desired, or if the zone to manage already exists.
* `refresh` - (Required) SOA REFRESH.
* `retry` - (Required) SOA RETRY.
* `expire` - (Required) SOA EXPIRE.
* `minimum` - (Required) SOA MINIMUM.

### Attribute Reference

The id of the resource is a composite of the record name and record type, joined by a separator - `:::`.

For example, record `example.com.` of type `SOA` will be represented with the following `id`: `example.com.:::SOA`
