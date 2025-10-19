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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceEveIconsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("api/list/networks")
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int               `json:"code"`
		Status  string            `json:"status"`
		Message string            `json:"message"`
		Icons   map[string]string `json:"icons"`
	}

	if err := c.HandleResponse(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	// アイコン名のリストを返す（シンプルな文字列配列）
	iconNames := make([]string, 0, len(result.Icons))
	for iconName := range result.Icons {
		iconNames = append(iconNames, iconName)
	}

	d.SetId("icons")
	if err := d.Set("icons", iconNames); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
