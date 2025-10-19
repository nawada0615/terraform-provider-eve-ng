package eveng

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveLab() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabCreate,
		ReadContext:   resourceEveLabRead,
		DeleteContext: resourceEveLabDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"path":          {Type: schema.TypeString, Required: true, ForceNew: true},
			"name":          {Type: schema.TypeString, Required: true, ForceNew: true},
			"author":        {Type: schema.TypeString, Optional: true, Default: "", ForceNew: true},
			"description":   {Type: schema.TypeString, Optional: true, Default: "", ForceNew: true},
			"body":          {Type: schema.TypeString, Optional: true, Default: "", ForceNew: true},
			"version":       {Type: schema.TypeString, Optional: true, Default: "1", ForceNew: true},
			"scripttimeout": {Type: schema.TypeInt, Optional: true, Default: 300, ForceNew: true},
			"lock":          {Type: schema.TypeBool, Optional: true, Default: false, ForceNew: true},
			"file":          {Type: schema.TypeString, Computed: true},
		},
	}
}

func normalizePath(p string) string {
	if p != "/" && !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

func resourceEveLabCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	path := normalizePath(d.Get("path").(string))
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Creating lab '%s' in path '%s'", name, path)

	payload := map[string]interface{}{
		"path": path,
		"name": name,
	}
	if v, ok := d.GetOk("author"); ok {
		payload["author"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		payload["description"] = v
	}
	if v, ok := d.GetOk("body"); ok {
		payload["body"] = v
	}
	if v, ok := d.GetOk("version"); ok {
		payload["version"] = v
	}
	if v, ok := d.GetOk("scripttimeout"); ok {
		payload["scripttimeout"] = v
	}

	log.Printf("[DEBUG] Lab payload: %+v", payload)

	resp, err := c.Post("api/labs", payload)
	if err != nil {
		log.Printf("[ERROR] Failed to create lab: %v", err)
		return diag.FromErr(fmt.Errorf("failed to create lab: %w", err))
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle lab creation response: %v", err)
		return diag.FromErr(fmt.Errorf("failed to handle lab creation response: %w", err))
	}

	if result.Code != 200 {
		log.Printf("[ERROR] Lab creation failed with code %d: %s", result.Code, result.Message)
		return diag.Errorf("lab creation failed: %s", result.Message)
	}

	labFile := path
	if labFile != "/" && !strings.HasSuffix(labFile, "/") {
		labFile += "/"
	}
	labFile += name + ".unl"
	d.SetId(labFile)
	if err := d.Set("file", labFile); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Lab created with file: %s", labFile)

	return resourceEveLabRead(ctx, d, m)
}

func resourceEveLabRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	labFile := d.Id()
	log.Printf("[DEBUG] Reading lab: %s", labFile)

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		log.Printf("[ERROR] Failed to get lab: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get lab: %w", err))
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Author        string      `json:"author"`
			Description   string      `json:"description"`
			Body          string      `json:"body"`
			Filename      string      `json:"filename"`
			ID            string      `json:"id"`
			Name          string      `json:"name"`
			Version       interface{} `json:"version"` // Can be string or int
			ScriptTimeout int         `json:"scripttimeout"`
			Lock          interface{} `json:"lock"` // Can be bool or int
		} `json:"data"`
	}

	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle lab read response: %v", err)
		// NotFound -> clear state
		d.SetId("")
		return nil
	}

	if result.Code != 200 {
		log.Printf("[ERROR] Lab read failed with code %d: %s", result.Code, result.Message)
		d.SetId("")
		return nil
	}

	// Handle version and lock fields
	version := handleVersionField(result.Data.Version)
	lock := handleLockField(result.Data.Lock)

	// Set lab data from response
	if err := setLabDataFromResponse(d, &result.Data, labFile, version, lock); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Lab read successfully: %s", result.Data.Name)
	return nil
}

func handleVersionField(version interface{}) string {
	switch v := version.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return "1" // default value
	}
}

func handleLockField(lock interface{}) bool {
	switch v := lock.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

func setLabDataFromResponse(d *schema.ResourceData, data *struct {
	Author        string      `json:"author"`
	Description   string      `json:"description"`
	Body          string      `json:"body"`
	Filename      string      `json:"filename"`
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Version       interface{} `json:"version"` // Can be string or int
	ScriptTimeout int         `json:"scripttimeout"`
	Lock          interface{} `json:"lock"` // Can be bool or int
}, labFile, version string, lock bool) error {
	if err := d.Set("author", data.Author); err != nil {
		return err
	}
	if err := d.Set("description", data.Description); err != nil {
		return err
	}
	if err := d.Set("body", data.Body); err != nil {
		return err
	}
	if err := d.Set("version", version); err != nil {
		return err
	}
	if err := d.Set("scripttimeout", data.ScriptTimeout); err != nil {
		return err
	}
	if err := d.Set("lock", lock); err != nil {
		return err
	}
	if err := d.Set("file", labFile); err != nil {
		return err
	}
	return nil
}

func resourceEveLabDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	labFile := d.Id()
	log.Printf("[DEBUG] Deleting lab: %s", labFile)

	resp, err := c.Delete("api/labs" + labFile)
	if err != nil {
		log.Printf("[ERROR] Failed to delete lab: %v", err)
		return diag.FromErr(fmt.Errorf("failed to delete lab: %w", err))
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle lab delete response: %v", err)
		return diag.FromErr(fmt.Errorf("failed to handle lab delete response: %w", err))
	}

	if result.Code != 200 {
		log.Printf("[ERROR] Lab deletion failed with code %d: %s", result.Code, result.Message)
		return diag.Errorf("lab deletion failed: %s", result.Message)
	}

	log.Printf("[DEBUG] Lab deleted successfully")
	return nil
}
