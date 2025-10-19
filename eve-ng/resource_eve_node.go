package eveng

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveNodeCreate,
		ReadContext:   resourceEveNodeRead,
		UpdateContext: resourceEveNodeUpdate,
		DeleteContext: resourceEveNodeDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema:        nodeSchema(),
	}
}

func nodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"lab_file": {Type: schema.TypeString, Required: true, ForceNew: true},
		"id":       {Type: schema.TypeString, Computed: true},
		"name":     {Type: schema.TypeString, Required: true},
		"type":     {Type: schema.TypeString, Required: true},
		"template": {Type: schema.TypeString, Required: true},
		"image":    {Type: schema.TypeString, Optional: true, Default: ""},
		"icon":     {Type: schema.TypeString, Optional: true, Default: ""},
		"top":      {Type: schema.TypeInt, Optional: true, Default: 0},
		"left":     {Type: schema.TypeInt, Optional: true, Default: 0},
		"delay":    {Type: schema.TypeInt, Optional: true, Default: 0},
		"config":   {Type: schema.TypeString, Optional: true, Default: ""},
		"ethernet": {Type: schema.TypeInt, Optional: true, Default: 0},
		"serial":   {Type: schema.TypeInt, Optional: true, Default: 0},

		// lifecycle
		"desired_state":    {Type: schema.TypeString, Optional: true, Default: "stopped"},
		"reboot_on_change": {Type: schema.TypeBool, Optional: true, Default: false},
		"wipe_on_destroy":  {Type: schema.TypeBool, Optional: true, Default: false},

		// qemu-specific
		"cpu":                {Type: schema.TypeInt, Optional: true, Default: 0},
		"ram":                {Type: schema.TypeInt, Optional: true, Default: 0},
		"cpulimit":           {Type: schema.TypeBool, Optional: true, Default: false},
		"uuid":               {Type: schema.TypeString, Optional: true, Default: ""},
		"qemu_version":       {Type: schema.TypeString, Optional: true, Default: ""},
		"qemu_arch":          {Type: schema.TypeString, Optional: true, Default: ""},
		"qemu_nic":           {Type: schema.TypeString, Optional: true, Default: ""},
		"qemu_options":       {Type: schema.TypeString, Optional: true, Default: ""},
		"firstmac":           {Type: schema.TypeString, Optional: true, Default: ""},
		"timos_line":         {Type: schema.TypeString, Optional: true, Default: ""},
		"timos_license":      {Type: schema.TypeString, Optional: true, Default: ""},
		"management_address": {Type: schema.TypeString, Optional: true, Default: ""},
	}
}

func resourceEveNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)
	nodeName := d.Get("name").(string)
	nodeType := d.Get("type").(string)
	nodeTemplate := d.Get("template").(string)

	log.Printf("[DEBUG] Creating node '%s' of type '%s' with template '%s' in lab '%s'",
		nodeName, nodeType, nodeTemplate, labFile)

	payload := buildNodePayloadFromState(d, false)
	log.Printf("[DEBUG] Node payload: %+v", payload)

	resp, err := c.Post("api/labs"+labFile+"/nodes", payload)
	if err != nil {
		log.Printf("[ERROR] Failed to create node: %v", err)
		return diag.FromErr(fmt.Errorf("failed to create node: %w", err))
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ID interface{} `json:"id"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle node creation response: %v", err)
		return diag.FromErr(fmt.Errorf("failed to handle node creation response: %w", err))
	}

	if result.Code != 201 {
		log.Printf("[ERROR] Node creation failed with code %d: %s", result.Code, result.Message)
		return diag.Errorf("node creation failed: %s", result.Message)
	}

	// API may return number or array; handle both
	var nodeID int
	switch v := result.Data.ID.(type) {
	case float64:
		nodeID = int(v)
	case []interface{}:
		if len(v) > 0 {
			if f, ok := v[0].(float64); ok {
				nodeID = int(f)
			} else {
				log.Printf("[ERROR] Invalid node ID format in array: %v", v[0])
				return diag.Errorf("invalid node ID format")
			}
		} else {
			log.Printf("[ERROR] Empty node ID array")
			return diag.Errorf("empty node ID array")
		}
	default:
		log.Printf("[ERROR] Unexpected node ID type: %T", result.Data.ID)
		return diag.Errorf("unexpected node ID type")
	}

	setNodeID(d, nodeID, labFile)
	log.Printf("[DEBUG] Node created with ID: %d", nodeID)

	// converge desired_state
	if ds, _ := d.Get("desired_state").(string); ds == nodeStatusStarted {
		log.Printf("[DEBUG] Starting node %d", nodeID)
		if err := nodePower(ctx, c, labFile, nodeID, "start"); err != nil {
			log.Printf("[WARN] Failed to start node: %v", err)
		}
	}
	return resourceEveNodeRead(ctx, d, m)
}

func setNodeID(d *schema.ResourceData, id int, labFile string) {
	_ = d.Set("id", strconv.Itoa(id))
	d.SetId(labFile + ":node:" + strconv.Itoa(id))
}

func buildNodePayloadFromState(d *schema.ResourceData, includeID bool) map[string]interface{} {
	p := map[string]interface{}{
		"name":     d.Get("name"),
		"type":     d.Get("type"),
		"template": d.Get("template"),
	}
	if includeID {
		p["id"] = d.Get("id")
	}
	copyIf(d, p, "image", "icon", "top", "left", "delay", "config", "ethernet", "serial")
	copyIf(d, p, "cpu", "ram", "cpulimit", "uuid", "qemu_version", "qemu_arch", "qemu_nic", "qemu_options", "firstmac", "timos_line", "timos_license", "management_address")
	return p
}

func copyIf(d *schema.ResourceData, dst map[string]interface{}, keys ...string) {
	for _, k := range keys {
		if v, ok := d.GetOk(k); ok {
			dst[k] = v
		}
	}
}

func parseNodeID(id string) (labFile string, nodeID int, ok bool) {
	parts := strings.Split(id, ":node:")
	if len(parts) != 2 {
		return "", 0, false
	}
	lab := parts[0]
	nid, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, false
	}
	return lab, nid, true
}

func resourceEveNodeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, nodeID, ok := parseNodeID(d.Id())
	if !ok {
		log.Printf("[ERROR] Invalid node ID format: %s", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Reading node %d from lab '%s'", nodeID, labFile)

	resp, err := c.Get("api/labs" + labFile + "/nodes/" + strconv.Itoa(nodeID))
	if err != nil {
		log.Printf("[ERROR] Failed to get node: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get node: %w", err))
	}

	var result struct {
		Code    int                    `json:"code"`
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle node read response: %v", err)
		d.SetId("")
		return nil
	}

	if result.Code != 200 {
		log.Printf("[ERROR] Node read failed with code %d: %s", result.Code, result.Message)
		d.SetId("")
		return nil
	}

	// Set node data from response
	setNodeDataFromResponse(d, nodeID, result.Data)

	log.Printf("[DEBUG] Node read successfully: %s", result.Data["name"])
	return nil
}

func setNodeDataFromResponse(d *schema.ResourceData, nodeID int, data map[string]interface{}) {
	_ = d.Set("id", strconv.Itoa(nodeID))

	// Set basic fields
	setStringField(d, data, "name")
	setStringField(d, data, "template")
	setStringField(d, data, "type")
	setStringField(d, data, "image")
	setStringField(d, data, "icon")

	// Set numeric fields
	setIntField(d, data, "top")
	setIntField(d, data, "left")
	setIntField(d, data, "delay")
	setIntField(d, data, "cpu")
	setIntField(d, data, "ram")
	setIntField(d, data, "ethernet")
	setIntField(d, data, "serial")

	// Set additional string fields
	setStringField(d, data, "uuid")
	setStringField(d, data, "firstmac")
	setStringField(d, data, "qemu_version")
	setStringField(d, data, "qemu_arch")
	setStringField(d, data, "qemu_nic")
	setStringField(d, data, "qemu_options")
	setStringField(d, data, "timos_line")
	setStringField(d, data, "timos_license")
	setStringField(d, data, "management_address")
}

func setStringField(d *schema.ResourceData, data map[string]interface{}, field string) {
	if v, ok := data[field].(string); ok {
		_ = d.Set(field, v)
	}
}

func setIntField(d *schema.ResourceData, data map[string]interface{}, field string) {
	if v, ok := data[field].(float64); ok {
		_ = d.Set(field, int(v))
	}
}

func resourceEveNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, nodeID, ok := parseNodeID(d.Id())
	if !ok {
		return diag.Errorf("invalid ID format")
	}

	// manage power if reboot_on_change
	reboot := d.Get("reboot_on_change").(bool)
	if reboot {
		_ = nodePower(ctx, c, labFile, nodeID, "stop")
	}

	payload := buildNodePayloadFromState(d, true)
	resp, err := c.Put("api/labs"+labFile+"/nodes/"+strconv.Itoa(nodeID), payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	// converge desired_state
	if ds, _ := d.Get("desired_state").(string); ds == nodeStatusStarted {
		_ = nodePower(ctx, c, labFile, nodeID, "start")
	} else {
		_ = nodePower(ctx, c, labFile, nodeID, "stop")
	}
	return resourceEveNodeRead(ctx, d, m)
}

func resourceEveNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, nodeID, ok := parseNodeID(d.Id())
	if !ok {
		return diag.Errorf("invalid ID format")
	}

	if d.Get("wipe_on_destroy").(bool) {
		_ = nodePower(ctx, c, labFile, nodeID, "wipe")
	}
	resp, err := c.Delete("api/labs" + labFile + "/nodes/" + strconv.Itoa(nodeID))
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func nodePower(_ context.Context, c *client.Client, labFile string, id int, action string) error {
	log.Printf("[DEBUG] %s node %d in lab '%s'", action, id, labFile)

	resp, err := c.Get("api/labs" + labFile + "/nodes/" + strconv.Itoa(id) + "/" + action)
	if err != nil {
		log.Printf("[ERROR] Failed to %s node: %v", action, err)
		return fmt.Errorf("failed to %s node: %w", action, err)
	}

	if err := c.HandleResponse(resp, nil); err != nil {
		log.Printf("[ERROR] Failed to handle node %s response: %v", action, err)
		return fmt.Errorf("failed to handle node %s response: %w", action, err)
	}

	log.Printf("[DEBUG] Node %s successfully", action)
	return nil
}
