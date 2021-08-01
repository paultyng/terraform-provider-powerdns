package powerdns

import (
	"fmt"
	"hash/crc32"
	"regexp"
	"strconv"
	"strings"
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
				Config:      testPDNSRecordConfigRecordEmpty().ResourceDeclaration(),
				ExpectError: regexp.MustCompile("'records' must not be empty"),
			},
		},
	})
}

func TestAccPDNSRecord_A(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigA)
}

func TestAccPDNSRecord_WithPtr(t *testing.T) {
	recordConfig := testPDNSRecordConfigAWithPtr()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: recordConfig.ResourceDeclaration(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordContents(recordConfig),
				),
			},
			{
				ResourceName:            recordConfig.ResourceName(),
				ImportStateId:           recordConfig.ResourceID(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"set_ptr"}, // Variance from common function
			},
		},
	})
}

// Use a basic existance check on the count resources, to avoid having to resolve interpolations in the names.
func TestAccPDNSRecord_WithCount(t *testing.T) {
	recordConfig := testPDNSRecordConfigHyphenedWithCount()
	resourceID0 := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighyphenedwithcount-0.sysa.xyz.:::A"}`
	resourceID1 := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighyphenedwithcount-1.sysa.xyz.:::A"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: recordConfig.ResourceDeclaration(),
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
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigAAAA)
}

func TestAccPDNSRecord_CNAME(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigCNAME)
}

func TestAccPDNSRecord_HINFO(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigHINFO)
}

func TestAccPDNSRecord_LOC(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigLOC)
}

func TestAccPDNSRecord_MX(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigMX)
}

func TestAccPDNSRecord_MXMulti(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigMXMulti)
}

func TestAccPDNSRecord_NAPTR(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigNAPTR)
}

func TestAccPDNSRecord_NS(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigNS)
}

func TestAccPDNSRecord_SPF(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSPF)
}

func TestAccPDNSRecord_SSHFP(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSSHFP)
}

func TestAccPDNSRecord_SRV(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSRV)
}

func TestAccPDNSRecord_TXT(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigTXT)
}

func TestAccPDNSRecord_ALIAS(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigALIAS)
}

func TestAccPDNSRecord_SOA(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSOA)
}

//
// Test Helper Functions
//

// Common Test Core: This function builds a create / update test for the majority of test cases
// Takes a function variable to avoid deep copy issues on updates.
func testPDNSRecordCommonTestCore(t *testing.T, recordConfigGenerator func() *PowerDNSRecordResource) {
	// Update test resources.
	recordConfig := recordConfigGenerator()

	updateRecordConfig := recordConfigGenerator()
	updateRecordConfig.Arguments.TTL += 100
	updateRecordConfig.Arguments.Records = updateRecordConfig.Arguments.UpdateRecords

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: recordConfig.ResourceDeclaration(),
				Check:  recordConfig.ResourceChecks(),
			},
			{
				Config: updateRecordConfig.ResourceDeclaration(),
				Check:  updateRecordConfig.ResourceChecks(),
			},
			{
				ResourceName:      recordConfig.ResourceName(),
				ImportStateId:     recordConfig.ResourceID(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPDNSRecordDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns_record" {
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

func testAccCheckPDNSRecordContents(recordConfig *PowerDNSRecordResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[recordConfig.ResourceName()]
		if !ok {
			return fmt.Errorf("Not found: %s", recordConfig.ResourceName())
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

		var matchingRecords []Record
		for _, rec := range foundRecords {
			// ListRecordsByID returns a list of records in v0 record structure format, which is a flat array of one entry per record content.
			if rec.ID() == rs.Primary.ID {
				matchingRecords = append(matchingRecords, rec)
			}
		}

		if len(matchingRecords) == 0 {
			return fmt.Errorf("Record does not exist: %#v", rs.Primary.ID)
		}

		// Assumption: Order will match between foundRecords and recordConfig.Arguments.Records
		for idx, desiredRecordContents := range recordConfig.Arguments.Records {
			errorPrefix := ("[#" + strconv.Itoa(idx) + "/" + desiredRecordContents + "] ")

			if idx >= len(matchingRecords) {
				return fmt.Errorf(errorPrefix + "Record not found")
			}

			rec := matchingRecords[idx]

			if rec.Name != recordConfig.Arguments.Name {
				return fmt.Errorf(errorPrefix+"Record name field does not match: %#v : %#v", rec.Name, recordConfig.Arguments.Name)
			}

			if rec.Type != recordConfig.Arguments.Type {
				return fmt.Errorf(errorPrefix+"Record type field does not match: %#v : %#v", rec.Type, recordConfig.Arguments.Type)
			}

			if rec.Content != desiredRecordContents {
				return fmt.Errorf(errorPrefix+"Record content field does not match: %#v : %#v", rec.Content, desiredRecordContents)
			}

			if rec.TTL != recordConfig.Arguments.TTL {
				return fmt.Errorf(errorPrefix+"Record TTL field does not match: %#v : %#v", rec.TTL, recordConfig.Arguments.TTL)
			}

			// Skipping check of SetPtr: this setting has been deprecated since PowerDNS 4.3.0 so the check will fail.

		}

		return nil
	}
}

//
// Resource Declaration types and methods
// These types & methods define a object layout declare resources during test, to allow for easy update tests and code deduplication
//
type PowerDNSRecordResourceArguments struct {
	Count   int
	Zone    string
	Name    string
	Type    string
	TTL     int
	Records []string
	// UpdateRecords are recordsets used for testing update behavior.
	UpdateRecords []string
	SetPtr        bool
}

type PowerDNSRecordResource struct {
	Name      string
	Arguments *PowerDNSRecordResourceArguments
}

// This function returns the record attribute ID in "records.[hash]" format for the given record
// Record hash lookup per https://github.com/pan-net/terraform-provider-powerdns/pull/78#issuecomment-793288653
func RecordAttributeIDForRecord(record string) string {
	crc := int(crc32.ChecksumIEEE([]byte(record)))
	if -crc >= 0 {
		crc = -crc
	}

	return "records." + strconv.Itoa(crc)
}

// This function builds out a suite of checks for the resource, suitable for passing to a TestStep as Check
func (resourceConfig *PowerDNSRecordResource) ResourceChecks() resource.TestCheckFunc {
	var checks []resource.TestCheckFunc

	checks = append(checks, testAccCheckPDNSRecordContents(resourceConfig))
	checks = append(checks, resource.TestCheckResourceAttr(resourceConfig.ResourceName(), "zone", resourceConfig.Arguments.Zone))
	checks = append(checks, resource.TestCheckResourceAttr(resourceConfig.ResourceName(), "name", resourceConfig.Arguments.Name))
	checks = append(checks, resource.TestCheckResourceAttr(resourceConfig.ResourceName(), "type", resourceConfig.Arguments.Type))
	checks = append(checks, resource.TestCheckResourceAttr(resourceConfig.ResourceName(), "ttl", strconv.Itoa(resourceConfig.Arguments.TTL)))

	for _, record := range resourceConfig.Arguments.Records {
		checks = append(checks, resource.TestCheckResourceAttr(resourceConfig.ResourceName(), RecordAttributeIDForRecord(record), record))
	}

	return resource.ComposeTestCheckFunc(checks...)
}

// This function builds out the Terraform DSL for the resource, suitable for passing to a TestStep as Config
func (resourceConfig *PowerDNSRecordResource) ResourceDeclaration() string {
	var encapsulatedRecords []string
	for _, record := range resourceConfig.Arguments.Records {
		encapsulatedRecords = append(encapsulatedRecords, (`"` + strings.Replace(record, `"`, `\"`, -1) + `"`))
	}

	resourceDeclaration := `resource "powerdns_record" "` + resourceConfig.Name + "\" {\n"
	if resourceConfig.Arguments.Count > 0 {
		resourceDeclaration += "  count = " + strconv.Itoa(resourceConfig.Arguments.Count) + "\n"
	}

	// zone, name, type, ttl, and records are mandatory
	resourceDeclaration += `  zone = "` + resourceConfig.Arguments.Zone + "\"\n"
	resourceDeclaration += `  name = "` + resourceConfig.Arguments.Name + "\"\n"
	resourceDeclaration += `  type = "` + resourceConfig.Arguments.Type + "\"\n"
	resourceDeclaration += "  ttl = " + strconv.Itoa(resourceConfig.Arguments.TTL) + "\n"
	resourceDeclaration += "  records = [ " + strings.Join(encapsulatedRecords, ", ") + " ]\n"

	if resourceConfig.Arguments.SetPtr {
		resourceDeclaration += "  set_ptr = true\n"
	}

	resourceDeclaration += "}"

	return resourceDeclaration
}

// This function builds out the Terraform resource ID for the resource
func (resourceConfig *PowerDNSRecordResource) ResourceID() string {
	return `{"zone":"` + resourceConfig.Arguments.Zone + `","id":"` + resourceConfig.Arguments.Name + ":::" + resourceConfig.Arguments.Type + `"}`
}

// This function is a trivial helper to return the Terraform resource name
func (resourceConfig *PowerDNSRecordResource) ResourceName() string {
	return "powerdns_record." + resourceConfig.Name
}

func NewPowerDNSRecordResource() *PowerDNSRecordResource {
	record := &PowerDNSRecordResource{}
	record.Arguments = &PowerDNSRecordResourceArguments{}

	// The zone argument is common across all the tests
	record.Arguments.Zone = "sysa.xyz."
	// TTL is set to 60 in the majority of the tests, default to 60 do deduplicate code.
	record.Arguments.TTL = 60
	record.Arguments.Records = make([]string, 0)
	return record
}

//
// Test resource declaration functions
//
// Pattern: testPDNSRecordConfigXXX() returns a PowerDNSRecordResource struct
// The PowerDNSRecordResource struct can be used to query test config, update attributes for update tests,
// and can have ResourceDeclaration() called against it to generate the Terraform DSL resource block string.
//
func testPDNSRecordConfigRecordEmpty() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-a"
	record.Arguments.Name = "testpdnsrecordconfigrecordempty.sysa.xyz."
	record.Arguments.Type = "A"
	return record
}

func testPDNSRecordConfigA() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-a"
	record.Arguments.Name = "testpdnsrecordconfigrecorda.sysa.xyz."
	record.Arguments.Type = "A"
	record.Arguments.Records = append(record.Arguments.Records, "1.1.1.1")
	record.Arguments.Records = append(record.Arguments.Records, "2.2.2.2")

	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "2.2.2.2")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "3.3.3.3")
	return record
}

func testPDNSRecordConfigAWithPtr() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-a"
	record.Arguments.Name = "testpdnsrecordconfigrecordawithptr.sysa.xyz."
	record.Arguments.Type = "A"
	record.Arguments.Records = append(record.Arguments.Records, "1.1.1.1")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "2.2.2.2")
	record.Arguments.SetPtr = true
	return record
}

func testPDNSRecordConfigHyphenedWithCount() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-counted"
	record.Arguments.Count = 2
	record.Arguments.Name = "testpdnsrecordconfighyphenedwithcount-${count.index}.sysa.xyz."
	record.Arguments.Type = "A"
	record.Arguments.Records = append(record.Arguments.Records, "1.1.1.${count.index}")
	return record
}

func testPDNSRecordConfigAAAA() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-aaaa"
	record.Arguments.Name = "testpdnsrecordconfigaaaa.sysa.xyz."
	record.Arguments.Type = "AAAA"
	record.Arguments.Records = append(record.Arguments.Records, "2001:db8:2000:bf0::1")
	record.Arguments.Records = append(record.Arguments.Records, "2001:db8:2000:bf1::1")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "2001:db8:2000:bf3::1")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "2001:db8:2000:bf4::1")
	return record
}

func testPDNSRecordConfigCNAME() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-cname"
	record.Arguments.Name = "testpdnsrecordconfigcname.sysa.xyz."
	record.Arguments.Type = "CNAME"
	record.Arguments.Records = append(record.Arguments.Records, "redis.example.com.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "redis.example.net.")
	return record
}

func testPDNSRecordConfigHINFO() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-hinfo"
	record.Arguments.Name = "testpdnsrecordconfighinfo.sysa.xyz."
	record.Arguments.Type = "HINFO"
	record.Arguments.Records = append(record.Arguments.Records, `"PC-Intel-2.4ghz" "Linux"`)
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, `"PC-Intel-3.2ghz" "Linux"`)
	return record
}

func testPDNSRecordConfigLOC() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-loc"
	record.Arguments.Name = "testpdnsrecordconfigloc.sysa.xyz."
	record.Arguments.Type = "LOC"
	record.Arguments.Records = append(record.Arguments.Records, "51 56 0.123 N 5 54 0.000 E 4.00m 1.00m 10000.00m 10.00m")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "51 10 43.900 N 1 49 34.300 E 4.00m 1.00m 10000.00m 10.00m")
	return record
}

func testPDNSRecordConfigMX() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-mx"
	record.Arguments.Name = "sysa.xyz."
	record.Arguments.Type = "MX"
	record.Arguments.Records = append(record.Arguments.Records, "10 mail.example.com.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "10 mail2.example.net.")
	return record
}

func testPDNSRecordConfigMXMulti() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-mx-multi"
	record.Arguments.Name = "multi.sysa.xyz."
	record.Arguments.Type = "MX"
	record.Arguments.Records = append(record.Arguments.Records, "10 mail.example.com.")
	record.Arguments.Records = append(record.Arguments.Records, "20 mail2.example.com.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "10 mail3.example.com.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "10 mail4.example.com.")
	return record
}

func testPDNSRecordConfigNAPTR() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-naptr"
	record.Arguments.Name = "sysa.xyz."
	record.Arguments.Type = "NAPTR"
	record.Arguments.Records = append(record.Arguments.Records, `100 50 "s" "z3950+I2L+I2C" "" _z3950._tcp.gatech.edu'.`)
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, `100 70 "s" "z3950+I2L+I2C" "" _z3950._tcp.gatech.edu'.`)
	return record
}

func testPDNSRecordConfigNS() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-ns"
	record.Arguments.Name = "lab.sysa.xyz."
	record.Arguments.Type = "NS"
	record.Arguments.Records = append(record.Arguments.Records, "ns1.sysa.xyz.")
	record.Arguments.Records = append(record.Arguments.Records, "ns2.sysa.xyz.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "ns3.sysa.xyz.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "ns4.sysa.xyz.")
	return record
}

func testPDNSRecordConfigSPF() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-spf"
	record.Arguments.Name = "sysa.xyz."
	record.Arguments.Type = "SPF"
	record.Arguments.Records = append(record.Arguments.Records, `"v=spf1 +all"`)
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, `"v=spf1 -all"`)
	return record
}

func testPDNSRecordConfigSSHFP() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-sshfp"
	record.Arguments.Name = "ssh.sysa.xyz."
	record.Arguments.Type = "SSHFP"
	record.Arguments.Records = append(record.Arguments.Records, "1 1 123456789abcdef67890123456789abcdef67890")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "1 1 fedcba9876543210fedcba9876543210fedcba98")
	return record
}

func testPDNSRecordConfigSRV() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-srv"
	record.Arguments.Name = "_redis._tcp.sysa.xyz."
	record.Arguments.Type = "SRV"
	record.Arguments.Records = append(record.Arguments.Records, "0 10 6379 redis1.sysa.xyz.")
	record.Arguments.Records = append(record.Arguments.Records, "0 10 6379 redis2.sysa.xyz.")
	record.Arguments.Records = append(record.Arguments.Records, "10 10 6379 redis-replica.sysa.xyz.")

	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "0 10 6379 redis1.sysa.xyz.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "0 10 6379 redis2.sysa.xyz.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "0 10 6379 redis3.sysa.xyz.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "10 10 6379 redis-replica.sysa.xyz.")
	return record
}

func testPDNSRecordConfigTXT() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-txt"
	record.Arguments.Name = "text.sysa.xyz."
	record.Arguments.Type = "TXT"
	record.Arguments.Records = append(record.Arguments.Records, `"text record payload"`)
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, `"updated text record payload"`)
	return record
}

func testPDNSRecordConfigALIAS() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-alias"
	record.Arguments.Name = "alias.sysa.xyz."
	record.Arguments.Type = "ALIAS"
	record.Arguments.TTL = 3600
	record.Arguments.Records = append(record.Arguments.Records, "www.some-alias.com.")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "www.some-other-alias.com.")
	return record
}

func testPDNSRecordConfigSOA() *PowerDNSRecordResource {
	record := NewPowerDNSRecordResource()
	record.Name = "test-soa"
	record.Arguments.Zone = "test-soa-sysa.xyz."
	record.Arguments.Name = "test-soa-sysa.xyz."
	record.Arguments.Type = "SOA"
	record.Arguments.TTL = 3600
	record.Arguments.Records = append(record.Arguments.Records, "something.something. hostmaster.sysa.xyz. 2019090301 10800 3600 604800 3600")
	record.Arguments.UpdateRecords = append(record.Arguments.UpdateRecords, "something.something. hostmaster.sysa.xyz. 2021021801 10800 3600 604800 3600")
	return record
}
