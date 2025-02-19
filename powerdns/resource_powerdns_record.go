package powerdns

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourcePDNSRecordCreate,
		Read:   resourcePDNSRecordRead,
		Delete: resourcePDNSRecordDelete,
		Exists: resourcePDNSRecordExists,
		Importer: &schema.ResourceImporter{
			State: resourcePDNSRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
				Set:      schema.HashString,
			},

			"set_ptr": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "For A and AAAA records, if true, create corresponding PTR.",
			},
		},
	}
}

func resourcePDNSRecordSOA() *schema.Resource {
	return &schema.Resource{
		Create: resourcePDNSRecordCreateSOA,
		Read:   resourcePDNSRecordRead,
		Delete: resourcePDNSRecordDelete,
		Exists: resourcePDNSRecordExists,
		Importer: &schema.ResourceImporter{
			State: resourcePDNSRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"mname": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"rname": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"serial": {
				Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"refresh": {
				Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"retry": {
				Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"expire": {
				Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"minimum": {
				Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePDNSRecordCreatePrepare(d *schema.ResourceData, meta interface{}) (ResourceRecordSet, string, int) {
	rrSet := ResourceRecordSet{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		TTL:  d.Get("ttl").(int),
	}

	zone := d.Get("zone").(string)
	ttl := d.Get("ttl").(int)

	return rrSet, zone, ttl
}

func resourcePDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	rrSet, zone, ttl := resourcePDNSRecordCreatePrepare(d, meta)
	recs := d.Get("records").(*schema.Set).List()
	setPtr := false

	// begin: ValidateFunc
	// https://www.terraform.io/docs/extend/schemas/schema-behaviors.html
	// "ValidateFunc is not yet supported on lists or sets"
	// when terraform will support ValidateFunc for non-primitives
	// we can move this block there
	if len(recs) == 0 {
		return fmt.Errorf("'records' must not be empty")
	}

	for _, recs := range recs {
		if len(strings.Trim(recs.(string), " ")) == 0 {
			log.Printf("[WARN] One or more values in 'records' contain empty '' value(s)")
		}
	}
	// end: ValidateFunc

	if v, ok := d.GetOk("set_ptr"); ok {
		setPtr = v.(bool)
	}

	records := make([]Record, 0, len(recs))
	for _, recContent := range recs {
		records = append(records,
			Record{Name: rrSet.Name,
				Type:    rrSet.Type,
				TTL:     ttl,
				Content: recContent.(string),
				SetPtr:  setPtr})
	}

	rrSet.Records = records

	return (resourcePDNSRecordCreateFinish(d, meta, zone, rrSet))
}

func resourcePDNSRecordCreateSOA(d *schema.ResourceData, meta interface{}) error {
	rrSet, zone, ttl := resourcePDNSRecordCreatePrepare(d, meta)
	client := meta.(*Client)

	log.Printf("[DEBUG] Searching existing SOA record at %s => %s", zone, d.Get("name").(string))
	soa_records, err := client.ListRecordsInRRSet(zone, d.Get("name").(string), "SOA")
	log.Printf("[DEBUG] Found existing SOA records %v", soa_records)
	if err != nil {
		return fmt.Errorf("Failed to fetch old SOA record: %s", err)
	}
	var serial int
	if len(soa_records) > 0 {
		serial, err = strconv.Atoi(strings.Fields(soa_records[0].Content)[2])
		if err != nil {
			return fmt.Errorf("Failed to parse old serial value in SOA record: %s", err)
		}
	} else {
		serial = d.Get("serial").(int)
	}
	log.Printf("[DEBUG] Set serial number to %d", serial)

	records := make([]Record, 0, 1)
	records = append(records,
		Record{Name: rrSet.Name,
			Type:    rrSet.Type,
			TTL:     ttl,
			Content: fmt.Sprintf("%s %s %d %d %d %d %d", d.Get("mname"), d.Get("rname"), serial, d.Get("refresh"), d.Get("retry"), d.Get("expire"), d.Get("minimum")),
			SetPtr:  false})

	rrSet.Records = records

	return (resourcePDNSRecordCreateFinish(d, meta, zone, rrSet))
}

func resourcePDNSRecordCreateFinish(d *schema.ResourceData, meta interface{}, zone string, rrSet ResourceRecordSet) error {
	client := meta.(*Client)

	log.Printf("[DEBUG] Creating PowerDNS Record: %#v", rrSet)

	recID, err := client.ReplaceRecordSet(zone, rrSet)
	if err != nil {
		return fmt.Errorf("Failed to create PowerDNS Record: %s", err)
	}

	d.SetId(recID)
	log.Printf("[INFO] Created PowerDNS Record with ID: %s", d.Id())

	return resourcePDNSRecordRead(d, meta)
}

func resourcePDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[DEBUG] Reading PowerDNS Record: %s", d.Id())
	records, err := client.ListRecordsByID(d.Get("zone").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't fetch PowerDNS Record: %s", err)
	}

	recs := make([]string, 0, len(records))
	if d.Get("type") == "SOA" {
		rsplit := strings.Fields(records[0].Content)
		mname := rsplit[0]
		d.Set("mname", mname)
		rname := rsplit[1]
		d.Set("rname", rname)

		serial, err := strconv.Atoi(rsplit[2])
		if err != nil {
			return fmt.Errorf("Failed to parse serial value in SOA record: %s", err)
		}
		d.Set("serial", serial)

		refresh, err := strconv.Atoi(rsplit[3])
		if err != nil {
			return fmt.Errorf("Failed to parse refresh value in SOA record: %s", err)
		}
		d.Set("refresh", refresh)

		retry, err := strconv.Atoi(rsplit[4])
		if err != nil {
			return fmt.Errorf("Failed to parse retry value in SOA record: %s", err)
		}
		d.Set("retry", retry)

		expire, err := strconv.Atoi(rsplit[5])
		if err != nil {
			return fmt.Errorf("Failed to parse expire value in SOA record: %s", err)
		}
		d.Set("expire", expire)

		minimum, err := strconv.Atoi(rsplit[6])
		if err != nil {
			return fmt.Errorf("Failed to parse minimum value in SOA record: %s", err)
		}
		d.Set("minimum", minimum)

		log.Printf("[DEBUG] Parsed PowerDNS SOA Record contents: mname %s rname %s serial %d refresh %d expire %d minimum %d", mname, rname, serial, refresh, expire, minimum)
	} else {
		for _, r := range records {
			recs = append(recs, r.Content)
		}
		d.Set("records", recs)
	}

	if len(records) > 0 || d.Get("Type") == "SOA" {
		d.Set("ttl", records[0].TTL)
	}

	return nil
}

func resourcePDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[INFO] Deleting PowerDNS Record: %s", d.Id())
	err := client.DeleteRecordSetByID(d.Get("zone").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting PowerDNS Record: %s", err)
	}

	return nil
}

func resourcePDNSRecordExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	tpe := d.Get("type").(string)

	log.Printf("[INFO] Checking existence of PowerDNS Record: %s, %s", name, tpe)

	client := meta.(*Client)
	exists, err := client.RecordExists(zone, name, tpe)

	if err != nil {
		return false, fmt.Errorf("Error checking PowerDNS Record: %s", err)
	}
	return exists, nil
}

func resourcePDNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	client := meta.(*Client)

	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}

	zoneName, ok := data["zone"]
	if !ok {
		return nil, fmt.Errorf("missing zone name in input data")
	}

	recordID, ok := data["id"]
	if !ok {
		return nil, fmt.Errorf("missing record id in input data")
	}

	log.Printf("[INFO] importing PowerDNS Record %s in Zone: %s", recordID, zoneName)

	records, err := client.ListRecordsByID(zoneName, recordID)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch PowerDNS Record: %s", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("rrset has no records to import")
	}

	recs := make([]string, 0, len(records))
	for _, r := range records {
		recs = append(recs, r.Content)
	}

	d.Set("zone", zoneName)
	d.Set("name", records[0].Name)
	d.Set("ttl", records[0].TTL)
	d.Set("type", records[0].Type)
	d.Set("records", recs)
	d.SetId(recordID)

	return []*schema.ResourceData{d}, nil
}
