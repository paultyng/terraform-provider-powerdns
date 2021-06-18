---
layout: "powerdns"
page_title: "PowerDNS: powerdns_record"
sidebar_current: "docs-powerdns-resource-record"
description: |-
  Provides a PowerDNS record resource.
---

# powerdns\_record

Provides a PowerDNS record resource.

## Example Usage

Note that PowerDNS may internally lowercase certain records (e.g. CNAME and AAAA), which may lead to resources being marked for a change in every single plan/apply.

### A record example
For the v1 API (PowerDNS version 4):

```hcl
# Add A record to the zone
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "www.example.com."
  type    = "A"
  ttl     = 300
  records = ["192.168.0.11"]
}
```

### PTR record example
An example creating PTR record:

```hcl
# Add PTR record to the zone
resource "powerdns_record" "foobar" {
  zone    = "0.10.in-addr.arpa."
  name    = "10.0.0.10.in-addr.arpa."
  type    = "PTR"
  ttl     = 300
  records = ["www.example.com."]
}
```

### Record with comments
An example creating a record with comments:

```hcl
# Add a record with comments
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "www.example.com."
  type = "A"
  ttl = 60
  records = [ "1.1.1.1" ]

  comment { 
    content = "Example comment #1"
    account = "Example account #1"
  }

  comment { 
    content = "Example comment #2"
    account = "Example account #2"
  }
}
```

### MX record example
The following example shows, how to setup MX record with a priority of `10`.
Please note that priority is not set as other `powerdns_record` properties; rather, it's part of the string that goes into `records` list.

```hcl
# Add MX record to the zone with priority 10
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "example.com."
  type    = "MX"
  ttl     = 300
  records = ["10 mail1.example.com"]
}
```

### Multiple values for the same resource record
Sometimes there is a need to have multiple values for the same DNS resource record, e.g. two IP addresses in A record, or `spf` and `DKIM` values in a TXT record.

Following examples show how to set mutliple values for A, TXT, and MX record types:

```hcl
# Add two A records for www.example.com
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "www.example.com."
  type    = "A"
  ttl     = 300
  records = ["192.168.0.11", "192.126.0.12"]
}
```

Similarly for the two TXT records:

```hcl
resource "powerdns_record" "foobar_txt_spf_and_dkim" {
  zone    = "example.com."
  name    = "example.com."
  type    = "TXT"
  ttl     = 60
  records = ["\"v=spf1 mx -all\"", "\"v=DKIM1 ;k=rsa; s=email; p=Msdsdfsdfsdfsdfsdfsdfsdfsdfsdfsfdfsdfsdfsdfds\""]
}
```

For two MX records with pritotities `10` and `20` respectively, the terraform code would look like this:

```hcl
# Add two MX records to the zone with priorities 10 and 20
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "example.com."
  type    = "MX"
  ttl     = 300
  records = ["10 mail1.example.com", "20 mail2.example.com"]
}
```

### Automatically set PTR record for A/AAAA records

!> **Deprecation warning:** _set_ptr_ feature is set to be deprecated in PowerDNS v4.3.0

PowerDNS API v4.2.0 offers a feature to automatically create corresponding PTR record for the A/AAAA record.
Existing PTR records with the same name are replaced. If no matching reverse zone is found, resource creation will fail.
You can use `powerdns_zone` resource to create the reverse zone.


!> **Warning:** Using _set_ptr:true_  will not automatically remove the PTR record when A/AAAA record is deleted. You should create PTR zone using `powerdns_zone` and manage PTR records using `powerdns_record`, rather than using _set_ptr_. With upcoming _set_ptr_ deprecation, this will be the only way of maintaining PTR records **via this provider**.

Here is an example of creating A record along with corresponding PTR record:

```hcl
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "www.example.com"
  type    = "A"
  ttl     = 300
  set_ptr = true
  records = ["192.168.0.11"]
}
```

For the legacy API (PowerDNS version 3.4):

```hcl
# Add a record to the zone
resource "powerdns_record" "foobar" {
  zone    = "example.com."
  name    = "www.example.com."
  type    = "A"
  ttl     = 300
  records = ["192.168.0.11"]
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The name of zone to contain this record.
* `name` - (Required) The name of the record.
* `type` - (Required) The record type.
* `ttl` - (Required) The TTL of the record.
* `records` - (Required) A string list of records.
* `set_ptr` - (Optional) [**_Deprecated in PowerDNS 4.3.0_**] A boolean (true/false), determining whether API server should automatically create PTR record in the matching reverse zone. Existing PTR records are replaced. If no matching reverse zone, an error is thrown.
* `comment` - (Optional) A comment about an RRSet.
    *  `content` - (Required) The content of the comment.
    *  `account` - (Required) The account of the comment.

### Attribute Reference

The id of the resource is a composite of the record name and record type, joined by a separator - `:::`.

For example, record `foo.test.com.` of type `A` will be represented with the following `id`: `foo.test.com.:::A`

### Importing

An existing record can be imported into this resource by supplying both the record id and zone name it belongs to.
If the record or zone is not found, or if the record is of a different type or in a different zone, an error will be returned.

For example:

```
$ terraform import powerdns_record.test-a '{"zone": "test.com.", "id": "foo.test.com.:::A"}'
```

For more information on how to use terraform's `import` command, please refer to terraform's [core documentation](https://www.terraform.io/docs/import/index.html#currently-state-only).

