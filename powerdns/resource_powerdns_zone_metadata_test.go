package powerdns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPDNSZoneMetadata_Empty(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testPDNSZoneMetadata_Empty,
				ExpectError: regexp.MustCompile("'metadata' must not be empty"),
			},
		},
	})
}

func TestAccPDNSZoneMetadata_AxfrFrom(t *testing.T) {
	resourceName := "powerdns_zone_metadata.test-axfr-from"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSZoneMetadataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSZoneMetadata_AxfrFrom,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "kind", "ALLOW-AXFR-FROM"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSZoneMetadata_AxfrSource(t *testing.T) {
	resourceName := "powerdns_zone_metadata.test-axfr-source"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSZoneMetadataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSZoneMetadata_AxfrSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "kind", "AXFR-SOURCE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPDNSZoneMetadata_XTest(t *testing.T) {
	resourceName := "powerdns_zone_metadata.test-x-test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSZoneMetadataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSZoneMetadata_XTest,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "kind", "X-TEST"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPDNSZoneMetadataDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "powerdns_zone_metadata" {
			continue
		}

		client := testAccProvider.Meta().(*Client)
		exists, err := client.ZoneMetadataExists(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error checking if zone metadata still exists: %#v", resource.Primary.ID)
		}
		if exists {
			return fmt.Errorf("Zone still exists: %#v", resource.Primary.ID)
		}

	}
	return nil
}

const testPDNSZoneMetadata_Empty = `
resource "powerdns_zone_metadata" "empty" {
	zone = "sysa.xyz."
	kind = "ALLOW-AXFR-FROM"
	metadata = [ ]
}`

const testPDNSZoneMetadata_AxfrFrom = `
resource "powerdns_zone_metadata" "test-axfr-from" {
	zone = "sysa.xyz."
	kind = "ALLOW-AXFR-FROM"
	metadata = ["AUTO-NS"]
}`

const testPDNSZoneMetadata_AxfrSource = `
resource "powerdns_zone_metadata" "test-axfr-source" {
	zone = "sysa.xyz."
	kind = "AXFR-SOURCE"
	metadata = ["10.0.0.1"]
}`

const testPDNSZoneMetadata_XTest = `
resource "powerdns_zone_metadata" "test-x-test" {
	zone = "sysa.xyz."
	kind = "X-TEST"
	metadata = ["test1", "test2"]
}`
