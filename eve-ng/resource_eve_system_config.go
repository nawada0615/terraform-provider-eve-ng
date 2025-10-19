package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveSystemConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveSystemConfigCreate,
		ReadContext:   resourceEveSystemConfigRead,
		UpdateContext: resourceEveSystemConfigUpdate,
		DeleteContext: resourceEveSystemConfigDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"cpu_limit":    {Type: schema.TypeInt, Optional: true},
			"ksm_enabled":  {Type: schema.TypeBool, Optional: true},
			"uksm_enabled": {Type: schema.TypeBool, Optional: true},
			"id":           {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceEveSystemConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Apply CPU limit if specified
	if cpuLimit, ok := d.GetOk("cpu_limit"); ok {
		payload := map[string]interface{}{"cpulimit": cpuLimit}
		resp, err := c.Post("api/cpulimit", payload)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := c.HandleResponse(resp, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	// Apply KSM setting if specified
	if ksmEnabled, ok := d.GetOk("ksm_enabled"); ok {
		payload := map[string]interface{}{"ksm": ksmEnabled}
		resp, err := c.Post("api/ksm", payload)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := c.HandleResponse(resp, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	// Apply UKSM setting if specified
	if uksmEnabled, ok := d.GetOk("uksm_enabled"); ok {
		payload := map[string]interface{}{"uksm": uksmEnabled}
		resp, err := c.Post("api/uksm", payload)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := c.HandleResponse(resp, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("system_config")
	return resourceEveSystemConfigRead(ctx, d, m)
}

func resourceEveSystemConfigRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("api/status")
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			CPULimit int `json:"cpu_limit"`
			KSM      struct {
				Enabled bool `json:"enabled"`
			} `json:"ksm"`
			UKSM struct {
				Enabled bool `json:"enabled"`
			} `json:"uksm"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("cpu_limit", result.Data.CPULimit); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ksm_enabled", result.Data.KSM.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("uksm_enabled", result.Data.UKSM.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("id", "system_config"); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEveSystemConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Update CPU limit if changed
	if d.HasChange("cpu_limit") {
		cpuLimit := d.Get("cpu_limit").(int)
		payload := map[string]interface{}{"cpulimit": cpuLimit}
		resp, err := c.Post("api/cpulimit", payload)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := c.HandleResponse(resp, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update KSM setting if changed
	if d.HasChange("ksm_enabled") {
		ksmEnabled := d.Get("ksm_enabled").(bool)
		payload := map[string]interface{}{"ksm": ksmEnabled}
		resp, err := c.Post("api/ksm", payload)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := c.HandleResponse(resp, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update UKSM setting if changed
	if d.HasChange("uksm_enabled") {
		uksmEnabled := d.Get("uksm_enabled").(bool)
		payload := map[string]interface{}{"uksm": uksmEnabled}
		resp, err := c.Post("api/uksm", payload)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := c.HandleResponse(resp, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEveSystemConfigRead(ctx, d, m)
}

func resourceEveSystemConfigDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// System config deletion doesn't reset settings, just remove from state
	return nil
}
