locals {
	zones = var.zones
	nameservers = var.nameservers
	nameservers_records_data = flatten([ for r in var.records : [ for rd in r.records : rd ] if r.type == "NS" ])
	non_soa_records = [ for r in var.records : r if r.type != "SOA" ]
	soa_records = [ for r in var.records : r if r.type == "SOA" ]
}

resource "powerdns_zone" "zone" {
	for_each = toset(local.zones)
	name = each.value
	kind = "Native"
	nameservers = length(var.nameservers) == 0 ? local.nameservers_records_data : var.nameservers
	soa_edit_api = var.soa_edit_api

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
		for i, record in local.non_soa_records : join("-", compact([
			lower(record.type),
			try(lower(record.name), ""),
			])) => {
			type = record.type
			name = try(record.name, "")
			ttl  = try(record.ttl, null)
			idx  = i
		}
	}

	records_expanded_soa = {
		for i, record in local.soa_records : join("-", compact([
			lower(record.type),
			try(lower(record.name), ""),
			])) => {
			type = record.type
			name = try(record.name, "")
			ttl  = try(record.ttl, null)
			mname = try(record.mname, element(local.nameservers_records_data, 0)),
			rname = record.rname,
			serial = try(record.serial, 0),
			refresh = record.refresh,
			retry = record.retry,
			expire = record.expire,
			minimum = record.minimum,
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

	records_by_name_soa = {
		for product in setproduct(local.zones, keys(local.records_expanded_soa)) : "${product[1]}-${product[0]}" => {
			zone = powerdns_zone.zone[product[0]].name
			type = local.records_expanded_soa[product[1]].type
			name = local.records_expanded_soa[product[1]].name
			ttl  = local.records_expanded_soa[product[1]].ttl
			mname = local.records_expanded_soa[product[1]].mname,
			rname = local.records_expanded_soa[product[1]].rname,
			serial = local.records_expanded_soa[product[1]].serial,
			refresh = local.records_expanded_soa[product[1]].refresh,
			retry = local.records_expanded_soa[product[1]].retry,
			expire = local.records_expanded_soa[product[1]].expire,
			minimum = local.records_expanded_soa[product[1]].minimum,
			idx  = local.records_expanded_soa[product[1]].idx
		}
	}

	records = local.records_by_name
	records_soa = local.records_by_name_soa
}

resource "powerdns_record_soa" "record_soa" {
	for_each = local.records_soa
	name = each.value.name == "" ? each.value.zone : join(".", [each.value.name, each.value.zone])
	zone = each.value.zone
	type = each.value.type
	ttl  = each.value.ttl
	mname = each.value.mname
	rname = each.value.rname
	serial = each.value.serial
	refresh = each.value.refresh
	retry = each.value.retry
	expire = each.value.expire
	minimum = each.value.minimum

	lifecycle {
		ignore_changes = [
			serial,
		]
	}

}

resource "powerdns_record" "record" {
	for_each = local.records
	name = each.value.name == "" ? each.value.zone : join(".", [each.value.name, each.value.zone])
	zone = each.value.zone
	type = each.value.type
	ttl  = each.value.ttl
	records = can(local.non_soa_records[each.value.idx].records) ? [for r in local.non_soa_records[each.value.idx].records :
		each.value.type == "TXT" && length(regexall("(\\\"\\\")", r)) == 0 ?
		format("\"%s\"", r) : r
	] : null
}
