package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func dataSourceEveIcons() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEveIconsRead,
		Schema: map[string]*schema.Schema{
			"icons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":        {Type: schema.TypeString, Computed: true},
						"filename":    {Type: schema.TypeString, Computed: true},
						"description": {Type: schema.TypeString, Computed: true},
						"category":    {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceEveIconsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("api/list/icons")
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

	icons := make([]map[string]interface{}, 0, len(result.Data))
	for name, iconData := range result.Data {
		iconMap, ok := iconData.(map[string]interface{})
		if !ok {
			continue
		}

		icon := map[string]interface{}{
			"name": name,
		}

		if filename, ok := iconMap["filename"].(string); ok {
			icon["filename"] = filename
		}
		if desc, ok := iconMap["description"].(string); ok {
			icon["description"] = desc
		}
		if cat, ok := iconMap["category"].(string); ok {
			icon["category"] = cat
		}

		icons = append(icons, icon)
	}

	d.SetId("icons")
	if err := d.Set("icons", icons); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
