## 1.5.0 (Unreleased)

FEATURES:
  * **Added option to cache PowerDNS API response** ([#81](https://github.com/pan-net/terraform-provider-powerdns/pull/81), @menai34)
  * **Added comments attribute for powerdns zone resource** ([#87](https://github.com/pan-net/terraform-provider-powerdns/pull/87), @PetrusHahol)

## 1.4.1 (January 21, 2021)

FEATURES:
  * **Added PowerDNS Zone Account support**  ([#71](https://github.com/pan-net/terraform-provider-powerdns/issues/71), @jbe-dw)

FIXES:
  * **Added support for port along with IP in the masters attribute** ([#64](https://github.com/pan-net/terraform-provider-powerdns/issues/64), @mbag)

ENHANCEMENTS:

  * **Add note in documentation about usage of SQLite3** ([#75](https://github.com/pan-net/terraform-provider-powerdns/issues/75), @dkowis)
  * **Improve _Using_ section in README** ([#67](https://github.com/pan-net/terraform-provider-powerdns/pull/67), @Nowaker)

## 1.4.0 (April 27, 2020)

FEATURES:
  * **Added ServerVersion attribute to client** ([#52](https://github.com/pan-net/terraform-provider-powerdns/issues/52))
  * **Added masters zone attribute for Slave zone kind** ([#59](https://github.com/pan-net/terraform-provider-powerdns/issues/59))

FIXES:
  * **Updated client tests to test sanitizeURL directly** ([#51](https://github.com/pan-net/terraform-provider-powerdns/issues/51))
  * **Fixed case sensitivity of kind zone attribute** ([#58](https://github.com/pan-net/terraform-provider-powerdns/issues/58))

ENHANCEMENTS:
  * **Updated documentation with examples based on user feedback** ([#57](https://github.com/pan-net/terraform-provider-powerdns/issues/57))

## 1.3.0 (December 20, 2019)

FEATURES:
  * **Move to using ParallelTest** - making tests faster ([#38](https://github.com/pan-net/terraform-provider-powerdns/issues/38))
  * **Added soa_edit_api option** ([#40](https://github.com/pan-net/terraform-provider-powerdns/issues/40))

FIXES:
  * **Fixed formatting in docs regarding import function** ([#31](https://github.com/pan-net/terraform-provider-powerdns/issues/31))

ENHANCEMENTS:
  * **Added tests for ALIAS type** ([#42](https://github.com/pan-net/terraform-provider-powerdns/issues/42))
  * **Migrated to terraform plugin SDK** ([#47](https://github.com/pan-net/terraform-provider-powerdns/issues/47))
  * **Updated vedor dependencies** ([#48](https://github.com/pan-net/terraform-provider-powerdns/issues/48))

## 1.2.0 (October 11, 2019)

FEATURES:
  * **Added support for terraform resource import** ([#31](https://github.com/pan-net/terraform-provider-powerdns/issues/31))

FIXES:
  * **Validate value of records** - record with empty records deleted the record from the PowerDNS remote but not from state file ([#33](https://github.com/pan-net/terraform-provider-powerdns/issues/33))

## 1.1.0 (August 13, 2019)

FEATURES:
  * **HTTPS Custom CA**: added option for custom Root CA for HTTPS Certificate validation (option `ca_certificate`) ([#22](https://github.com/pan-net/terraform-provider-powerdns/issues/22))
  * **HTTPS**: added option to skip HTTPS certificate validation - insecure HTTPS (option `insecure_https`) ([#22](https://github.com/pan-net/terraform-provider-powerdns/issues/22))

ENHANCEMENTS:
  * The provider doesn't attempt to connect to the PowerDNS endpoint if there is nothing to be done ([#24](https://github.com/pan-net/terraform-provider-powerdns/issues/24))
  * `server_url` (`PDNS_SERVER_URL`) can now be declared with/without scheme, port, trailing slashes or path ([#28](https://github.com/pan-net/terraform-provider-powerdns/issues/28))

## 1.0.0 (August 06, 2019)

NOTES:
 * provider: This release includes only a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and this update should have no significant changes in behavior for the provider. Please report any unexpected behavior in new GitHub issues (Terraform core: https://github.com/hashicorp/terraform/issues or Terraform PowerDNS Provider: https://github.com/pan-net/terraform-provider-powerdns/issues) ([#16](https://github.com/pan-net/terraform-provider-powerdns/issues/16))

ENHANCEMENTS:
  * Switch to go modules and Terraform v0.12 SDK ([#16](https://github.com/pan-net/terraform-provider-powerdns/issues/16))

## 0.2.0 (July 31, 2019)

FEATURES:
  * **New resource**: `powerdns_zone` ([#8](https://github.com/pan-net/terraform-provider-powerdns/issues/8))

ENHANCEMENTS:
  * resource/powerdns_record: Add support for set-ptr option ([#4](https://github.com/pan-net/terraform-provider-powerdns/issues/4))
  * build: Added docker-compose tests ([#9](https://github.com/pan-net/terraform-provider-powerdns/issues/9))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
