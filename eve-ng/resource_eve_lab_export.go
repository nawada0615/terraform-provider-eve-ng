package eveng

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveLabExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabExportCreate,
		ReadContext:   resourceEveLabExportRead,
		DeleteContext: resourceEveLabExportDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file":        {Type: schema.TypeString, Required: true, ForceNew: true},
			"export_format":   {Type: schema.TypeString, Optional: true, Default: "unl", ForceNew: true},
			"include_configs": {Type: schema.TypeBool, Optional: true, Default: true, ForceNew: true},
			"export_data":     {Type: schema.TypeString, Computed: true},
			"export_filename": {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceEveLabExportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)
	exportFormat := d.Get("export_format").(string)
	includeConfigs := d.Get("include_configs").(bool)

	payload := map[string]interface{}{
		"format":  exportFormat,
		"configs": includeConfigs,
	}

	resp, err := c.Post("api/labs"+labFile+"/export", payload)
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ExportData string `json:"export_data"`
			Filename   string `json:"filename"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	exportID := labFile + ":export:" + exportFormat
	d.SetId(exportID)
	if err := d.Set("export_data", result.Data.ExportData); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("export_filename", result.Data.Filename); err != nil {
		return diag.FromErr(err)
	}
	return resourceEveLabExportRead(ctx, d, m)
}

func resourceEveLabExportRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Export is a one-time operation, just verify the lab exists
	c := m.(*client.Client)
	labFile := strings.Split(d.Id(), ":export:")[0]

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceEveLabExportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Export doesn't need explicit deletion
	return nil
}
