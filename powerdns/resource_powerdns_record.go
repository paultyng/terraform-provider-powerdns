package powerdns

import (
	"encoding/json"
	"fmt"
	"log"
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

			"comment": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"account": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				Optional:    true,
				ForceNew:    true,
				Description: "A comment about an RRSet.",
			},
		},
	}
}

func resourcePDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	rrSet := ResourceRecordSet{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		TTL:  d.Get("ttl").(int),
	}

	comments := d.Get("comment").(*schema.Set).List()
	if len(comments) > 0 {
		commentObjs := make([]Comment, 0, len(comments))
		for _, comment := range comments {
			commentMap := comment.(map[string]interface{})
			commentObjs = append(
				commentObjs,
				Comment{
					Content: commentMap["content"].(string),
					Account: commentMap["account"].(string),
				})
		}
		rrSet.Comments = commentObjs
	}

	zone := d.Get("zone").(string)
	ttl := d.Get("ttl").(int)
	recs := d.Get("records").(*schema.Set).List()
	setPtr := false

	if v, ok := d.GetOk("set_ptr"); ok {
		setPtr = v.(bool)
	}

	// begin: ValidateFunc
	// https://www.terraform.io/docs/extend/schemas/schema-behaviors.html
	// "ValidateFunc is not yet supported on lists or sets"
	// when terraform will support ValidateFunc for non-primitives
	// we can move this block there
	for _, recs := range recs {
		if len(strings.Trim(recs.(string), " ")) == 0 {
			log.Printf("[WARN] One or more values in 'records' contain empty '' value(s)")
		}
	}
	if !(len(recs) > 0) {
		return fmt.Errorf("'records' must not be empty")
	}
	// end: ValidateFunc

	if len(recs) > 0 {
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

		log.Printf("[DEBUG] Creating PowerDNS Record: %#v", rrSet)

		recID, err := client.ReplaceRecordSet(zone, rrSet)
		if err != nil {
			return fmt.Errorf("Failed to create PowerDNS Record: %s", err)
		}

		d.SetId(recID)
		log.Printf("[INFO] Created PowerDNS Record with ID: %s", d.Id())

	}

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
	for _, r := range records {
		recs = append(recs, r.Content)
	}
	d.Set("records", recs)

	if len(records) > 0 {
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
	rrset, err := client.GetRRSetOfRecord(zoneName, recordID)

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

	commentsSet := make([]map[string]interface{}, len(rrset.Comments))
	for i, comment := range rrset.Comments {
		commentsSet[i] = map[string]interface{}{
			"content": comment.Content,
			"account": comment.Account,
		}
	}

	d.Set("zone", zoneName)
	d.Set("name", records[0].Name)
	d.Set("ttl", records[0].TTL)
	d.Set("type", records[0].Type)
	d.Set("records", recs)
	d.Set("comment", commentsSet)
	d.SetId(recordID)

	return []*schema.ResourceData{d}, nil
}
