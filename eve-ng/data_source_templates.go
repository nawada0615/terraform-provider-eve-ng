package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func dataSourceEveTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEveTemplatesRead,
		Schema: map[string]*schema.Schema{
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":        {Type: schema.TypeString, Computed: true},
						"type":        {Type: schema.TypeString, Computed: true},
						"description": {Type: schema.TypeString, Computed: true},
						"icon":        {Type: schema.TypeString, Computed: true},
						"category":    {Type: schema.TypeString, Computed: true},
						"defaults": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEveTemplatesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("api/list/templates/")
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

	templates := make([]map[string]interface{}, 0, len(result.Data))
	for name, templateData := range result.Data {
		templateMap, ok := templateData.(map[string]interface{})
		if !ok {
			continue
		}

		template := map[string]interface{}{
			"name": name,
		}

		if t, ok := templateMap["type"].(string); ok {
			template["type"] = t
		}
		if desc, ok := templateMap["description"].(string); ok {
			template["description"] = desc
		}
		if icon, ok := templateMap["icon"].(string); ok {
			template["icon"] = icon
		}
		if cat, ok := templateMap["category"].(string); ok {
			template["category"] = cat
		}

		// Extract defaults
		if defaults, ok := templateMap["defaults"].(map[string]interface{}); ok {
			template["defaults"] = defaults
		}

		templates = append(templates, template)
	}

	d.SetId("templates")
	if err := d.Set("templates", templates); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
