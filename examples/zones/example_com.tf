module "example_com" {
  source = "../../terraform-provider-powerdns"
  zones = [
    "example.com.",
  ]

  soa_edit_api = "INCREASE"

  # preferably declare NS records instead of this
  # https://github.com/pan-net/terraform-provider-powerdns/issues/63
  nameservers = [
    "ns1.example.com.",
    "ns2.example.com.",
  ]

  records = [
    {
      type    = "SOA"
      ttl     = 43200
      rname   = "admin.opensuse.org."
      refresh = 7200
      retry   = 600
      expire  = 1209600
      minimum = 6400
      # this can be used to set an initial serial number for new zones
      # serial number changes to existing zones will be ignored, the user is expected to use SOA-EDIT-API
      serial = 1
    },
    {
      type = "SOA",
      ttl  = 300,
      records = [
        "ns1.example.com. hostmaster.example.com. 0 10800 3600 604800 3600"
      ]
    },
    {
      name = "www",
      type = "AAAA",
      ttl  = 300,
      records = [
        "::1",
      ]
    }
  ]
}

