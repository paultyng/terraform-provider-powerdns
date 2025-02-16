locals {
	zones = var.zones
	nameservers = var.nameservers
	nameservers_records = flatten([ for r in var.records : [ for rd in r.records : rd ] if r.type == "NS" ])
}

resource "powerdns_zone" "zone" {
	for_each = toset(local.zones)
	name = each.value
	kind = "Native"
	nameservers = length(var.nameservers) == 0 ? local.nameservers_records : var.nameservers
	lifecycle {
		ignore_changes = [
			# https://github.com/pan-net/terraform-provider-powerdns/issues/63
			# users of the module are expected to use NS records for tracking nameservers
			nameservers,
		]
	}
}

locals {
	records_expanded = {
		for i, record in var.records : join("-", compact([
			lower(record.type),
			try(lower(record.name), ""),
			])) => {
			type = record.type
			name = try(record.name, "")
			ttl  = try(record.ttl, null)
			idx  = i
		}
	}

	records_by_name = {
		for product in setproduct(local.zones, keys(local.records_expanded)) : "${product[1]}-${product[0]}" => {
			zone = powerdns_zone.zone[product[0]].name
			type = local.records_expanded[product[1]].type
			name = local.records_expanded[product[1]].name
			ttl  = local.records_expanded[product[1]].ttl
			idx  = local.records_expanded[product[1]].idx
		}
	}

	records = local.records_by_name
}

resource "powerdns_record" "record" {
	for_each = local.records
	name = each.value.name == "" ? each.value.zone : join(".", [each.value.name, each.value.zone])
	zone = each.value.zone
	type = each.value.type
	ttl  = each.value.ttl
	records = can(var.records[each.value.idx].records) ? [for r in var.records[each.value.idx].records :
		each.value.type == "TXT" && length(regexall("(\\\"\\\")", r)) == 0 ?
		format("\"%s\"", r) : r
	] : null
}
