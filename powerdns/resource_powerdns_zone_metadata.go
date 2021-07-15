package powerdns

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePDNSZoneMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourcePDNSZoneMetadataCreate,
		Read:   resourcePDNSZoneMetadataRead,
		Delete: resourcePDNSZoneMetadataDelete,
		Exists: resourcePDNSZoneMetadataExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"metadata": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePDNSZoneMetadataCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	zone := d.Get("zone").(string)
	mtdata := d.Get("metadata").(*schema.Set).List()

	for _, mt := range mtdata {
		if len(strings.TrimSpace(mt.(string))) == 0 {
			log.Printf("[WARN] One or more values in 'metadata' contain empty value(s)")
			break
		}
	}
	if !(len(mtdata) > 0) {
		return fmt.Errorf("'metadata' must not be empty")
	}

	metadata := make([]string, 0, len(mtdata))
	for _, mt := range mtdata {
		metadata = append(metadata, mt.(string))
	}

	zoneMetadata := ResourceZoneMetadata{
		Kind:     d.Get("kind").(string),
		Metadata: metadata,
	}

	log.Printf("[DEBUG] Creating PowerDNS Zone Metadata: %#v", zoneMetadata)

	metaid, err := client.UpdateZoneMetadata(zone, zoneMetadata)
	if err != nil {
		return fmt.Errorf("Failed to create PowerDNS Zone Metadata: %s", err)
	}

	d.SetId(metaid)
	log.Printf("[INFO] Created PowerDNS Zone Metadata with ID: %s", d.Id())

	return nil
}

func resourcePDNSZoneMetadataRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[DEBUG] Reading PowerDNS Zone Metadata: %s", d.Id())
	record, err := client.GetZoneMetadata(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't fetch PowerDNS Zone Metadata: %s", err)
	}

	zone, _, err := parseID(d.Id())
	if err != nil {
		return err
	}

	d.SetId(d.Id())
	d.Set("kind", record.Kind)
	d.Set("metadata", record.Metadata)
	d.Set("zone", zone)

	return nil
}

func resourcePDNSZoneMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Deleting PowerDNS Zone Metadata: %s", d.Id())
	err := client.DeleteZoneMetadata(d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting PowerDNS Zone Metadata: %s", err)
	}

	return nil
}

func resourcePDNSZoneMetadataExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("[INFO] Checking existence of PowerDNS Zone Metadata: %s", d.Id())

	client := meta.(*Client)
	exists, err := client.ZoneMetadataExists(d.Id())

	if err != nil {
		return false, fmt.Errorf("Error checking PowerDNS Zone Metadata: %s", err)
	}
	return exists, nil
}
