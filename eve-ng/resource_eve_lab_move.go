package eveng

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveLabMove() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabMoveCreate,
		ReadContext:   resourceEveLabMoveRead,
		UpdateContext: resourceEveLabMoveUpdate,
		DeleteContext: resourceEveLabMoveDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file":         {Type: schema.TypeString, Required: true, ForceNew: true},
			"source_path":      {Type: schema.TypeString, Required: true, ForceNew: true},
			"destination_path": {Type: schema.TypeString, Required: true},
			"new_name":         {Type: schema.TypeString, Optional: true},
		},
	}
}

func resourceEveLabMoveCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)
	destPath := d.Get("destination_path").(string)
	newName := d.Get("new_name").(string)

	// Normalize destination path
	if !strings.HasPrefix(destPath, "/") {
		destPath = "/" + destPath
	}
	if destPath != "/" && !strings.HasSuffix(destPath, "/") {
		destPath += "/"
	}

	payload := map[string]interface{}{
		"path": destPath,
	}
	if newName != "" {
		payload["name"] = newName
	}

	resp, err := c.Put("api/labs"+labFile+"/move", payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	// Update ID to reflect new location
	newLabFile := destPath
	if newName != "" {
		newLabFile += newName + ".unl"
	} else {
		// Extract name from original lab_file
		parts := strings.Split(labFile, "/")
		if len(parts) > 0 {
			newLabFile += parts[len(parts)-1]
		}
	}
	d.SetId(newLabFile + ":move")
	return resourceEveLabMoveRead(ctx, d, m)
}

func resourceEveLabMoveRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := strings.TrimSuffix(d.Id(), ":move")

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Filename string `json:"filename"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		d.SetId("")
		return nil
	}

	// Extract path and name from filename
	parts := strings.Split(strings.TrimSuffix(result.Data.Filename, ".unl"), "/")
	if len(parts) == 0 {
		d.SetId("")
		return nil
	}

	path := "/" + strings.Join(parts[:len(parts)-1], "/")
	if path == "//" {
		path = "/"
	}
	name := parts[len(parts)-1]

	if err := d.Set("lab_file", labFile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_path", d.Get("source_path")); err != nil { // Keep original
		return diag.FromErr(err)
	}
	if err := d.Set("destination_path", path); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("new_name", name); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEveLabMoveUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Move is a one-time operation, recreate if destination changes
	return resourceEveLabMoveCreate(ctx, d, m)
}

func resourceEveLabMoveDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Move operation doesn't need explicit deletion
	return nil
}
