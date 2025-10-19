package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func dataSourceEveNetworkTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEveNetworkTypesRead,
		Schema: map[string]*schema.Schema{
			"network_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":        {Type: schema.TypeString, Computed: true},
						"type":        {Type: schema.TypeString, Computed: true},
						"description": {Type: schema.TypeString, Computed: true},
						"icon":        {Type: schema.TypeString, Computed: true},
						"category":    {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceEveNetworkTypesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("api/list/networks")
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int                    `json:"code"`
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := c.HandleResponse(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	networkTypes := make([]map[string]interface{}, 0, len(result.Data))
	for name, networkData := range result.Data {
		networkMap, ok := networkData.(map[string]interface{})
		if !ok {
			continue
		}

		networkType := map[string]interface{}{
			"name": name,
		}

		if t, ok := networkMap["type"].(string); ok {
			networkType["type"] = t
		}
		if desc, ok := networkMap["description"].(string); ok {
			networkType["description"] = desc
		}
		if icon, ok := networkMap["icon"].(string); ok {
			networkType["icon"] = icon
		}
		if cat, ok := networkMap["category"].(string); ok {
			networkType["category"] = cat
		}

		networkTypes = append(networkTypes, networkType)
	}

	d.SetId("network_types")
	if err := d.Set("network_types", networkTypes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
