package eveng

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveLabClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabCloneCreate,
		ReadContext:   resourceEveLabCloneRead,
		DeleteContext: resourceEveLabCloneDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"source_lab_file":  {Type: schema.TypeString, Required: true, ForceNew: true},
			"destination_path": {Type: schema.TypeString, Required: true, ForceNew: true},
			"new_name":         {Type: schema.TypeString, Required: true, ForceNew: true},
			"cloned_lab_file":  {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceEveLabCloneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	sourceLabFile := d.Get("source_lab_file").(string)
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
		"name": newName,
	}

	resp, err := c.Post("api/labs"+sourceLabFile+"/clone", payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	clonedLabFile := destPath + newName + ".unl"
	d.SetId(clonedLabFile + ":clone")
	if err := d.Set("cloned_lab_file", clonedLabFile); err != nil {
		return diag.FromErr(err)
	}
	return resourceEveLabCloneRead(ctx, d, m)
}

func resourceEveLabCloneRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	clonedLabFile := strings.TrimSuffix(d.Id(), ":clone")

	resp, err := c.Get("api/labs" + clonedLabFile)
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

	if err := d.Set("cloned_lab_file", clonedLabFile); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEveLabCloneDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	clonedLabFile := strings.TrimSuffix(d.Id(), ":clone")

	// Delete the cloned lab
	resp, err := c.Delete("api/labs" + clonedLabFile)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
