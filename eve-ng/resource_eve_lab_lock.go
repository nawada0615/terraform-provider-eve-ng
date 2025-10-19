package eveng

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveLabLock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabLockCreate,
		ReadContext:   resourceEveLabLockRead,
		DeleteContext: resourceEveLabLockDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file": {Type: schema.TypeString, Required: true, ForceNew: true},
			"locked":   {Type: schema.TypeBool, Computed: true},
		},
	}
}

func resourceEveLabLockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)

	resp, err := c.Put("api/labs"+labFile+"/Lock", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(labFile + ":lock")
	return resourceEveLabLockRead(ctx, d, m)
}

func resourceEveLabLockRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := strings.TrimSuffix(d.Id(), ":lock")

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Lock bool `json:"lock"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("lab_file", labFile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("locked", result.Data.Lock); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEveLabLockDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := strings.TrimSuffix(d.Id(), ":lock")

	resp, err := c.Put("api/labs"+labFile+"/Unlock", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
