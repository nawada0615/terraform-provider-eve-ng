package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveUserCreate,
		ReadContext:   resourceEveUserRead,
		UpdateContext: resourceEveUserUpdate,
		DeleteContext: resourceEveUserDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"username": {Type: schema.TypeString, Required: true, ForceNew: true},
			"password": {Type: schema.TypeString, Required: true, Sensitive: true},
			"email":    {Type: schema.TypeString, Optional: true},
			"name":     {Type: schema.TypeString, Optional: true},
			"role":     {Type: schema.TypeString, Optional: true, Default: "user"},
			"enabled":  {Type: schema.TypeBool, Optional: true, Default: true},
			"expires":  {Type: schema.TypeString, Optional: true},
			"id":       {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceEveUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	payload := map[string]interface{}{
		"username": d.Get("username").(string),
		"password": d.Get("password").(string),
	}
	if v, ok := d.GetOk("email"); ok {
		payload["email"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		payload["name"] = v
	}
	if v, ok := d.GetOk("role"); ok {
		payload["role"] = v
	}
	if v, ok := d.GetOk("enabled"); ok {
		payload["enabled"] = v
	}
	if v, ok := d.GetOk("expires"); ok {
		payload["expires"] = v
	}

	resp, err := c.Post("api/users", payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	username := d.Get("username").(string)
	d.SetId(username)
	return resourceEveUserRead(ctx, d, m)
}

func resourceEveUserRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	username := d.Id()

	resp, err := c.Get("api/users/" + username)
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Role     string `json:"role"`
			Enabled  bool   `json:"enabled"`
			Expires  string `json:"expires"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("username", result.Data.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", result.Data.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", result.Data.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role", result.Data.Role); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", result.Data.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expires", result.Data.Expires); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("id", username); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEveUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	username := d.Id()

	payload := map[string]interface{}{}
	if d.HasChange("password") {
		payload["password"] = d.Get("password").(string)
	}
	if d.HasChange("email") {
		payload["email"] = d.Get("email").(string)
	}
	if d.HasChange("name") {
		payload["name"] = d.Get("name").(string)
	}
	if d.HasChange("role") {
		payload["role"] = d.Get("role").(string)
	}
	if d.HasChange("enabled") {
		payload["enabled"] = d.Get("enabled").(bool)
	}
	if d.HasChange("expires") {
		payload["expires"] = d.Get("expires").(string)
	}

	resp, err := c.Put("api/users/"+username, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}
	return resourceEveUserRead(ctx, d, m)
}

func resourceEveUserDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	username := d.Id()

	resp, err := c.Delete("api/users/" + username)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
