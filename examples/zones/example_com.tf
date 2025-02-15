module "example_com" {
  source = "../../terraform-provider-powerdns"
  zones = [
    "example.com.",
  ]

  nameservers = [
    "ns1.example.com.",
    "ns2.example.com.",
  ]

  records = [
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

