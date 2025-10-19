package eveng

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

const (
	nodeStatusStarted = "started"
)

func resourceEveLabMonitoring() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabMonitoringCreate,
		ReadContext:   resourceEveLabMonitoringRead,
		UpdateContext: resourceEveLabMonitoringUpdate,
		DeleteContext: resourceEveLabMonitoringDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file":               {Type: schema.TypeString, Required: true, ForceNew: true},
			"monitor_nodes":          {Type: schema.TypeBool, Optional: true, Default: true},
			"monitor_networks":       {Type: schema.TypeBool, Optional: true, Default: true},
			"alert_threshold_cpu":    {Type: schema.TypeFloat, Optional: true, Default: 80.0},
			"alert_threshold_memory": {Type: schema.TypeFloat, Optional: true, Default: 90.0},
			"monitoring_enabled":     {Type: schema.TypeBool, Computed: true},
			"node_count":             {Type: schema.TypeInt, Computed: true},
			"network_count":          {Type: schema.TypeInt, Computed: true},
			"running_nodes":          {Type: schema.TypeInt, Computed: true},
		},
	}
}

func resourceEveLabMonitoringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)

	// Get lab status for monitoring
	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}

	var labResult struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Nodes []struct {
				ID     int    `json:"id"`
				Status string `json:"status"`
			} `json:"nodes"`
			Networks []struct {
				ID int `json:"id"`
			} `json:"networks"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &labResult); err != nil {
		return diag.FromErr(err)
	}

	// Count running nodes
	runningCount := 0
	for _, node := range labResult.Data.Nodes {
		if node.Status == nodeStatusStarted {
			runningCount++
		}
	}

	monitoringID := labFile + ":monitoring"
	d.SetId(monitoringID)
	if err := d.Set("monitoring_enabled", true); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_count", len(labResult.Data.Nodes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_count", len(labResult.Data.Networks)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("running_nodes", runningCount); err != nil {
		return diag.FromErr(err)
	}
	return resourceEveLabMonitoringRead(ctx, d, m)
}

func resourceEveLabMonitoringRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := strings.TrimSuffix(d.Id(), ":monitoring")

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}

	var labResult struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Nodes []struct {
				ID     int    `json:"id"`
				Status string `json:"status"`
			} `json:"nodes"`
			Networks []struct {
				ID int `json:"id"`
			} `json:"networks"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &labResult); err != nil {
		d.SetId("")
		return nil
	}

	// Count running nodes
	runningCount := 0
	for _, node := range labResult.Data.Nodes {
		if node.Status == nodeStatusStarted {
			runningCount++
		}
	}

	if err := d.Set("monitoring_enabled", true); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_count", len(labResult.Data.Nodes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_count", len(labResult.Data.Networks)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("running_nodes", runningCount); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEveLabMonitoringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Monitoring settings are read-only for now
	return resourceEveLabMonitoringRead(ctx, d, m)
}

func resourceEveLabMonitoringDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Monitoring doesn't need explicit deletion
	return nil
}
