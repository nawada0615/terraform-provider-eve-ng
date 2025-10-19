package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func dataSourceEveStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEveStatusRead,
		Schema: map[string]*schema.Schema{
			"cpu_usage":        {Type: schema.TypeFloat, Computed: true},
			"memory_usage":     {Type: schema.TypeFloat, Computed: true},
			"disk_usage":       {Type: schema.TypeFloat, Computed: true},
			"swap_usage":       {Type: schema.TypeFloat, Computed: true},
			"running_wrappers": {Type: schema.TypeInt, Computed: true},
			"ksm_enabled":      {Type: schema.TypeBool, Computed: true},
			"uksm_enabled":     {Type: schema.TypeBool, Computed: true},
			"cpu_limit":        {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceEveStatusRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			CPU struct {
				Usage float64 `json:"usage"`
			} `json:"cpu"`
			Memory struct {
				Usage float64 `json:"usage"`
			} `json:"memory"`
			Disk struct {
				Usage float64 `json:"usage"`
			} `json:"disk"`
			Swap struct {
				Usage float64 `json:"usage"`
			} `json:"swap"`
			RunningWrappers int `json:"running_wrappers"`
			KSM             struct {
				Enabled bool `json:"enabled"`
			} `json:"ksm"`
			UKSM struct {
				Enabled bool `json:"enabled"`
			} `json:"uksm"`
			CPULimit int `json:"cpu_limit"`
		} `json:"data"`
	}

	if err := c.HandleResponse(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("status")
	if err := d.Set("cpu_usage", result.Data.CPU.Usage); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_usage", result.Data.Memory.Usage); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_usage", result.Data.Disk.Usage); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("swap_usage", result.Data.Swap.Usage); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("running_wrappers", result.Data.RunningWrappers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ksm_enabled", result.Data.KSM.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("uksm_enabled", result.Data.UKSM.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_limit", result.Data.CPULimit); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
