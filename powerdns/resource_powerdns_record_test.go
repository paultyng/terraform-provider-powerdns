package powerdns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPDNSRecord_Empty(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testPDNSRecordConfigRecordEmpty,
				ExpectError: regexp.MustCompile("'records' must not be empty"),
			},
		},
	})
}

func TestAccPDNSRecord_A(t *testing.T) {
	resourceName := "powerdns_record.test-a"
	resourceID := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfiga.sysa.xyz.:::A"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigA,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_WithPtr(t *testing.T) {
	resourceName := "powerdns_record.test-a-ptr"
	resourceID := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfigawithptr.sysa.xyz.:::A"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigAWithPtr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           resourceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"set_ptr"},
			},
		},
	})
}

func TestAccPDNSRecord_WithCount(t *testing.T) {
	resourceID0 := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighyphenedwithcount-0.sysa.xyz.:::A"}`
	resourceID1 := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighyphenedwithcount-1.sysa.xyz.:::A"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigHyphenedWithCount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists("powerdns_record.test-counted.0"),
					testAccCheckPDNSRecordExists("powerdns_record.test-counted.1"),
				),
			},
			{
				ResourceName:      "powerdns_record.test-counted[0]",
				ImportStateId:     resourceID0,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "powerdns_record.test-counted[1]",
				ImportStateId:     resourceID1,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_AAAA(t *testing.T) {
	resourceName := "powerdns_record.test-aaaa"
	resourceID := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfigaaaa.sysa.xyz.:::AAAA"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigAAAA,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_CNAME(t *testing.T) {
	resourceName := "powerdns_record.test-cname"
	resourceID := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfigcname.sysa.xyz.:::CNAME"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigCNAME,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_HINFO(t *testing.T) {
	resourceName := "powerdns_record.test-hinfo"
	resourceID := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighinfo.sysa.xyz.:::HINFO"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigHINFO,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_LOC(t *testing.T) {
	resourceName := "powerdns_record.test-loc"
	resourceID := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfigloc.sysa.xyz.:::LOC"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigLOC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_MX(t *testing.T) {
	resourceName := "powerdns_record.test-mx"
	resourceNameMulti := "powerdns_record.test-mx-multi"
	resourceID := `{"zone":"sysa.xyz.","id":"sysa.xyz.:::MX"}`
	resourceIDMulti := `{"zone":"sysa.xyz.","id":"multi.sysa.xyz.:::MX"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigMX,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testPDNSRecordConfigMXMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceNameMulti),
				),
			},
			{
				ResourceName:      resourceNameMulti,
				ImportStateId:     resourceIDMulti,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_NAPTR(t *testing.T) {
	resourceName := "powerdns_record.test-naptr"
	resourceID := `{"zone":"sysa.xyz.","id":"sysa.xyz.:::NAPTR"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigNAPTR,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_NS(t *testing.T) {
	resourceName := "powerdns_record.test-ns"
	resourceID := `{"zone":"sysa.xyz.","id":"lab.sysa.xyz.:::NS"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigNS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_SPF(t *testing.T) {
	resourceName := "powerdns_record.test-spf"
	resourceID := `{"zone":"sysa.xyz.","id":"sysa.xyz.:::SPF"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigSPF,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_SSHFP(t *testing.T) {
	resourceName := "powerdns_record.test-sshfp"
	resourceID := `{"zone":"sysa.xyz.","id":"ssh.sysa.xyz.:::SSHFP"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigSSHFP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_SRV(t *testing.T) {
	resourceName := "powerdns_record.test-srv"
	resourceID := `{"zone":"sysa.xyz.","id":"_redis._tcp.sysa.xyz.:::SRV"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigSRV,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_TXT(t *testing.T) {
	resourceName := "powerdns_record.test-txt"
	resourceID := `{"zone":"sysa.xyz.","id":"text.sysa.xyz.:::TXT"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigTXT,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_ALIAS(t *testing.T) {
	resourceName := "powerdns_record.test-alias"
	resourceID := `{"zone":"sysa.xyz.","id":"alias.sysa.xyz.:::ALIAS"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecord_SOA(t *testing.T) {
	resourceName := "powerdns_record.test-soa"
	resourceID := `{"zone":"test-soa-sysa.xyz.","id":"test-soa-sysa.xyz.:::SOA"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigSOA,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSRecordSOA_SOA(t *testing.T) {
	resourceName := "powerdns_record_soa.test-soa2"
	resourceID := `{"zone":"test-soa2-sysa.xyz.","id":"test-soa2-sysa.xyz.:::SOA"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigSOA2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPDNSRecordDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns_record" && rs.Type != "powerdns_record_soa" {
			continue
		}

		client := testAccProvider.Meta().(*Client)
		exists, err := client.RecordExistsByID(rs.Primary.Attributes["zone"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error checking if record still exists: %#v", rs.Primary.ID)
		}
		if exists {
			return fmt.Errorf("Record still exists: %#v", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckPDNSRecordExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		foundRecords, err := client.ListRecordsByID(rs.Primary.Attributes["zone"], rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(foundRecords) == 0 {
			return fmt.Errorf("Record does not exist")
		}
		for _, rec := range foundRecords {
			if rec.ID() == rs.Primary.ID {
				return nil
			}
		}
		return fmt.Errorf("Record does not exist: %#v", rs.Primary.ID)
	}
}

const testPDNSRecordConfigRecordEmpty = `
resource "powerdns_record" "test-a" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfigrecordempty.sysa.xyz."
	type = "A"
	ttl = 60
	records = [ ]
}`

const testPDNSRecordConfigA = `
resource "powerdns_record" "test-a" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfiga.sysa.xyz."
	type = "A"
	ttl = 60
	records = [ "1.1.1.1", "2.2.2.2" ]
}`

const testPDNSRecordConfigAWithPtr = `
resource "powerdns_record" "test-a-ptr" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfigawithptr.sysa.xyz."
	type = "A"
	ttl = 60
	set_ptr = true
	records = [ "1.1.1.1" ]
}`

const testPDNSRecordConfigHyphenedWithCount = `
resource "powerdns_record" "test-counted" {
	count = "2"
	zone = "sysa.xyz."
	name = "testpdnsrecordconfighyphenedwithcount-${count.index}.sysa.xyz."
	type = "A"
	ttl = 60
	records = [ "1.1.1.${count.index}" ]
}`

const testPDNSRecordConfigAAAA = `
resource "powerdns_record" "test-aaaa" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfigaaaa.sysa.xyz."
	type = "AAAA"
	ttl = 60
	records = [ "2001:db8:2000:bf0::1", "2001:db8:2000:bf1::1" ]
}`

const testPDNSRecordConfigCNAME = `
resource "powerdns_record" "test-cname" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfigcname.sysa.xyz."
	type = "CNAME"
	ttl = 60
	records = [ "redis.example.com." ]
}`

const testPDNSRecordConfigHINFO = `
resource "powerdns_record" "test-hinfo" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfighinfo.sysa.xyz."
	type = "HINFO"
	ttl = 60
	records = [ "\"PC-Intel-2.4ghz\" \"Linux\"" ]
}`

const testPDNSRecordConfigLOC = `
resource "powerdns_record" "test-loc" {
	zone = "sysa.xyz."
	name = "testpdnsrecordconfigloc.sysa.xyz."
	type = "LOC"
	ttl = 60
	records = [ "51 56 0.123 N 5 54 0.000 E 4.00m 1.00m 10000.00m 10.00m" ]
}`

const testPDNSRecordConfigMX = `
resource "powerdns_record" "test-mx" {
	zone = "sysa.xyz."
	name = "sysa.xyz."
	type = "MX"
	ttl = 60
	records = [ "10 mail.example.com." ]
}`

const testPDNSRecordConfigMXMulti = `
resource "powerdns_record" "test-mx-multi" {
	zone = "sysa.xyz."
	name = "multi.sysa.xyz."
	type = "MX"
	ttl = 60
	records = [ "10 mail1.example.com.", "20 mail2.example.com." ]
}`

const testPDNSRecordConfigNAPTR = `
resource "powerdns_record" "test-naptr" {
	zone = "sysa.xyz."
	name = "sysa.xyz."
	type = "NAPTR"
	ttl = 60
	records = [ "100 50 \"s\" \"z3950+I2L+I2C\" \"\" _z3950._tcp.gatech.edu'." ]
}`

const testPDNSRecordConfigNS = `
resource "powerdns_record" "test-ns" {
	zone = "sysa.xyz."
	name = "lab.sysa.xyz."
	type = "NS"
	ttl = 60
	records = [ "ns1.sysa.xyz.", "ns2.sysa.xyz." ]
}`

const testPDNSRecordConfigSPF = `
resource "powerdns_record" "test-spf" {
	zone = "sysa.xyz."
	name = "sysa.xyz."
	type = "SPF"
	ttl = 60
	records = [ "\"v=spf1 +all\"" ]
}`

const testPDNSRecordConfigSSHFP = `
resource "powerdns_record" "test-sshfp" {
	zone = "sysa.xyz."
	name = "ssh.sysa.xyz."
	type = "SSHFP"
	ttl = 60
	records = [ "1 1 123456789abcdef67890123456789abcdef67890" ]
}`

const testPDNSRecordConfigSRV = `
resource "powerdns_record" "test-srv" {
	zone = "sysa.xyz."
	name = "_redis._tcp.sysa.xyz."
	type = "SRV"
	ttl = 60
	records = [ "0 10 6379 redis1.sysa.xyz.", "0 10 6379 redis2.sysa.xyz.", "10 10 6379 redis-replica.sysa.xyz." ]
}`

const testPDNSRecordConfigTXT = `
resource "powerdns_record" "test-txt" {
	zone = "sysa.xyz."
	name = "text.sysa.xyz."
	type = "TXT"
	ttl = 60
	records = [ "\"text record payload\"" ]
}`

const testPDNSRecordConfigALIAS = `
resource "powerdns_record" "test-alias" {
	zone = "sysa.xyz."
	name = "alias.sysa.xyz."
	type = "ALIAS"
	ttl = 3600
	records = [ "www.some-alias.com." ]
}`

const testPDNSRecordConfigSOA = `
resource "powerdns_record" "test-soa" {
	zone = "test-soa-sysa.xyz."
	name = "test-soa-sysa.xyz."
	type = "SOA"
	ttl = 3600
	records = [ "something.something. hostmaster.sysa.xyz. 2019090301 10800 3600 604800 3600" ]
}`

const testPDNSRecordConfigSOA2 = `
resource "powerdns_record_soa" "test-soa2" {
	zone = "test-soa2-sysa.xyz."
	name = "test-soa2-sysa.xyz."
	type = "SOA"
	ttl = 3600
	mname = "ns1.sysa.xyz."
	rname = "hostmaster.sysa.xyz."
	serial = 0
	refresh = 7200
	retry = 600
	expire = 1209600
	minimum = 6400

	lifecycle {
		ignore_changes = [
			serial,
		]
	}
}`
