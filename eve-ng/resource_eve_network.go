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

func resourceEveNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveNetworkCreate,
		ReadContext:   resourceEveNetworkRead,
		UpdateContext: resourceEveNetworkUpdate,
		DeleteContext: resourceEveNetworkDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file":   {Type: schema.TypeString, Required: true, ForceNew: true},
			"name":       {Type: schema.TypeString, Required: true},
			"type":       {Type: schema.TypeString, Required: true},
			"top":        {Type: schema.TypeInt, Optional: true, Default: 0},
			"left":       {Type: schema.TypeInt, Optional: true, Default: 0},
			"icon":       {Type: schema.TypeString, Optional: true, Default: ""},
			"visibility": {Type: schema.TypeString, Optional: true, Default: "1"},
			"id":         {Type: schema.TypeString, Computed: true},
			"node_count": {Type: schema.TypeInt, Computed: true},
		},
	}
}

func resourceEveNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	labFile := d.Get("lab_file").(string)
	networkName := d.Get("name").(string)
	networkType := d.Get("type").(string)

	log.Printf("[DEBUG] Creating network '%s' of type '%s' in lab '%s'", networkName, networkType, labFile)

	payload := map[string]interface{}{
		"name": networkName,
		"type": networkType,
	}
	if v, ok := d.GetOk("top"); ok {
		payload["top"] = v
	}
	if v, ok := d.GetOk("left"); ok {
		payload["left"] = v
	}
	if v, ok := d.GetOk("icon"); ok {
		payload["icon"] = v
	}
	if v, ok := d.GetOk("visibility"); ok {
		payload["visibility"] = v
	}

	log.Printf("[DEBUG] Network payload: %+v", payload)

	resp, err := c.Post("api/labs"+labFile+"/networks", payload)
	if err != nil {
		log.Printf("[ERROR] Failed to create network: %v", err)
		return diag.FromErr(fmt.Errorf("failed to create network: %w", err))
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle network creation response: %v", err)
		return diag.FromErr(fmt.Errorf("failed to handle network creation response: %w", err))
	}

	if result.Code != 201 {
		log.Printf("[ERROR] Network creation failed with code %d: %s", result.Code, result.Message)
		return diag.Errorf("network creation failed: %s", result.Message)
	}

	// ID format: <lab_file>:network:<id>
	id := result.Data.ID
	networkID := labFile + ":network:" + strconv.Itoa(id)
	d.SetId(networkID)

	log.Printf("[DEBUG] Network created with ID: %s", networkID)

	return resourceEveNetworkRead(ctx, d, m)
}

func parseNetworkID(id string) (labFile string, netID int, ok bool) {
	// format: <lab_file>:network:<id>
	parts := strings.Split(id, ":network:")
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

func resourceEveNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, netID, ok := parseNetworkID(d.Id())
	if !ok {
		log.Printf("[ERROR] Invalid network ID format: %s", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Reading network %d from lab '%s'", netID, labFile)

	// Try individual network read first
	resp, err := c.Get("api/labs" + labFile + "/networks/" + strconv.Itoa(netID))
	if err != nil {
		log.Printf("[ERROR] Failed to get network: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get network: %w", err))
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Count      int         `json:"count"`
			Left       int         `json:"left"`
			Name       string      `json:"name"`
			Top        int         `json:"top"`
			Type       string      `json:"type"`
			Visibility interface{} `json:"visibility"` // Can be string or int
			Icon       string      `json:"icon"`
		} `json:"data"`
	}

	// Handle individual network read response
	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[WARN] Individual network read failed: %v", err)
		// Fallback to network list read
		return resourceEveNetworkReadFromList(ctx, d, m, labFile, netID)
	}

	if result.Code != 200 {
		log.Printf("[WARN] Individual network read failed with code %d: %s, trying network list", result.Code, result.Message)
		// Fallback to network list read
		return resourceEveNetworkReadFromList(ctx, d, m, labFile, netID)
	}

	// Set network data from individual read
	return setNetworkData(d, labFile, netID, &result.Data, "individual")
}

// setNetworkData sets network data from either individual read or list read
func setNetworkData(d *schema.ResourceData, labFile string, netID int, data *struct {
	Count      int         `json:"count"`
	Left       int         `json:"left"`
	Name       string      `json:"name"`
	Top        int         `json:"top"`
	Type       string      `json:"type"`
	Visibility interface{} `json:"visibility"`
	Icon       string      `json:"icon"`
}, source string) diag.Diagnostics {
	// Handle visibility field which can be string or int
	visibility := convertVisibilityToString(data.Visibility)

	if err := d.Set("lab_file", labFile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("id", strconv.Itoa(netID)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", data.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", data.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("top", data.Top); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("left", data.Left); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("icon", data.Icon); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("visibility", visibility); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_count", data.Count); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Network read successfully from %s: %s", source, data.Name)
	return nil
}

func convertVisibilityToString(visibility interface{}) string {
	switch v := visibility.(type) {
	case string:
		return v
	case float64:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	default:
		return "1" // default value
	}
}

// Fallback function to read network from network list
func resourceEveNetworkReadFromList(_ context.Context, d *schema.ResourceData, m interface{}, labFile string, netID int) diag.Diagnostics {
	c := m.(*client.Client)

	log.Printf("[DEBUG] Reading network %d from network list in lab '%s'", netID, labFile)

	resp, err := c.Get("api/labs" + labFile + "/networks")
	if err != nil {
		log.Printf("[ERROR] Failed to get network list: %v", err)
		d.SetId("")
		return nil
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    map[string]struct {
			Count      int         `json:"count"`
			Left       int         `json:"left"`
			Name       string      `json:"name"`
			Top        int         `json:"top"`
			Type       string      `json:"type"`
			Visibility interface{} `json:"visibility"` // Can be string or int
			Icon       string      `json:"icon"`
		} `json:"data"`
	}

	if err := c.HandleResponse(resp, &result); err != nil {
		log.Printf("[ERROR] Failed to handle network list response: %v", err)
		d.SetId("")
		return nil
	}

	if result.Code != 200 {
		log.Printf("[ERROR] Network list read failed with code %d: %s", result.Code, result.Message)
		d.SetId("")
		return nil
	}

	// Look for the specific network ID in the list
	netIDStr := strconv.Itoa(netID)
	networkData, exists := result.Data[netIDStr]
	if !exists {
		log.Printf("[ERROR] Network %d not found in network list", netID)
		d.SetId("")
		return nil
	}

	// Set network data from list read
	return setNetworkData(d, labFile, netID, &networkData, "list")
}

func resourceEveNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, netID, ok := parseNetworkID(d.Id())
	if !ok {
		log.Printf("[ERROR] Invalid network ID format: %s", d.Id())
		return diag.Errorf("invalid ID format")
	}

	log.Printf("[DEBUG] Updating network %d in lab '%s'", netID, labFile)

	payload := map[string]interface{}{
		"id": netID,
	}
	if d.HasChange("name") {
		payload["name"] = d.Get("name").(string)
	}
	if d.HasChange("type") {
		payload["type"] = d.Get("type").(string)
	}
	if d.HasChange("top") {
		payload["top"] = d.Get("top").(int)
	}
	if d.HasChange("left") {
		payload["left"] = d.Get("left").(int)
	}
	if d.HasChange("icon") {
		payload["icon"] = d.Get("icon").(string)
	}
	if d.HasChange("visibility") {
		payload["visibility"] = d.Get("visibility").(string)
	}

	log.Printf("[DEBUG] Network update payload: %+v", payload)

	resp, err := c.Put("api/labs"+labFile+"/networks/"+strconv.Itoa(netID), payload)
	if err != nil {
		log.Printf("[ERROR] Failed to update network: %v", err)
		return diag.FromErr(fmt.Errorf("failed to update network: %w", err))
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		log.Printf("[ERROR] Failed to handle network update response: %v", err)
		return diag.FromErr(fmt.Errorf("failed to handle network update response: %w", err))
	}

	log.Printf("[DEBUG] Network updated successfully")
	return resourceEveNetworkRead(ctx, d, m)
}

func resourceEveNetworkDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, netID, ok := parseNetworkID(d.Id())
	if !ok {
		log.Printf("[ERROR] Invalid network ID format: %s", d.Id())
		return diag.Errorf("invalid ID format")
	}

	log.Printf("[DEBUG] Deleting network %d from lab '%s'", netID, labFile)

	resp, err := c.Delete("api/labs" + labFile + "/networks/" + strconv.Itoa(netID))
	if err != nil {
		log.Printf("[ERROR] Failed to delete network: %v", err)
		return diag.FromErr(fmt.Errorf("failed to delete network: %w", err))
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		log.Printf("[ERROR] Failed to handle network delete response: %v", err)
		return diag.FromErr(fmt.Errorf("failed to handle network delete response: %w", err))
	}

	log.Printf("[DEBUG] Network deleted successfully")
	return nil
}
