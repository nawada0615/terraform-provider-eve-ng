package eveng

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

const (
	pathSeparator = "//"
)

func resourceEveFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveFolderCreate,
		ReadContext:   resourceEveFolderRead,
		DeleteContext: resourceEveFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Parent folder path",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Folder name",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full folder path",
			},
		},
	}
}

func resourceEveFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	path := d.Get("path").(string)
	name := d.Get("name").(string)

	// Normalize path
	if path != "/" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if path != "/" && !strings.HasSuffix(path, "/") {
		path += "/"
	}

	createData := map[string]interface{}{
		"path": path,
		"name": name,
	}

	resp, err := c.Post("api/folders", createData)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	// Set full path as ID
	fullPath := path + name
	if fullPath == pathSeparator {
		fullPath = "/" + name
	}
	d.SetId(fullPath)

	return resourceEveFolderRead(ctx, d, m)
}

func resourceEveFolderRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	fullPath := d.Id()

	// Normalize path for API call
	apiPath := fullPath
	if apiPath == "/" {
		apiPath = ""
	}

	resp, err := c.Get("api/folders" + apiPath)
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Folders []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"folders"`
		} `json:"data"`
	}

	if err := c.HandleResponse(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	// Find the folder in the response
	found := false
	for _, folder := range result.Data.Folders {
		if folder.Path == fullPath {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	// Extract path and name from full path
	if err := setFolderPathAndName(d, fullPath); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("full_path", fullPath); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setFolderPathAndName(d *schema.ResourceData, fullPath string) error {
	parts := strings.Split(strings.Trim(fullPath, "/"), "/")
	if len(parts) == 0 {
		if err := d.Set("path", "/"); err != nil {
			return err
		}
		if err := d.Set("name", ""); err != nil {
			return err
		}
	} else if len(parts) == 1 {
		if err := d.Set("path", "/"); err != nil {
			return err
		}
		if err := d.Set("name", parts[0]); err != nil {
			return err
		}
	} else {
		path := "/" + strings.Join(parts[:len(parts)-1], "/")
		if err := d.Set("path", path); err != nil {
			return err
		}
		if err := d.Set("name", parts[len(parts)-1]); err != nil {
			return err
		}
	}
	return nil
}

func resourceEveFolderDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	fullPath := d.Id()

	// Normalize path for API call
	apiPath := fullPath
	if apiPath == "/" {
		return diag.Errorf("cannot delete root folder")
	}

	resp, err := c.Delete("api/folders" + apiPath)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
