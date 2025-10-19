// Package eveng provides a Terraform provider for EVE-NG network emulation platform.
// It supports managing labs, networks, nodes, and interface attachments.
package eveng

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

// Provider returns the EVE-NG provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The EVE-NG server endpoint URL",
				DefaultFunc: schema.EnvDefaultFunc("EVE_NG_ENDPOINT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username for EVE-NG authentication",
				DefaultFunc: schema.EnvDefaultFunc("EVE_NG_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for EVE-NG authentication",
				DefaultFunc: schema.EnvDefaultFunc("EVE_NG_PASSWORD", nil),
			},
			"insecure_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip TLS certificate verification",
				DefaultFunc: schema.EnvDefaultFunc("EVE_NG_INSECURE", false),
			},
			"timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "30s",
				Description: "Timeout for API requests",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"eve_lab_lock":             resourceEveLabLock(),
			"eve_lab_move":             resourceEveLabMove(),
			"eve_lab_clone":            resourceEveLabClone(),
			"eve_lab_batch_start":      resourceEveLabBatchStart(),
			"eve_lab_batch_stop":       resourceEveLabBatchStop(),
			"eve_lab_batch_wipe":       resourceEveLabBatchWipe(),
			"eve_folder":               resourceEveFolder(),
			"eve_lab":                  resourceEveLab(),
			"eve_network":              resourceEveNetwork(),
			"eve_node":                 resourceEveNode(),
			"eve_interface_attachment": resourceEveInterfaceAttachment(),
			"eve_user":                 resourceEveUser(),
			"eve_system_config":        resourceEveSystemConfig(),
			"eve_lab_export":           resourceEveLabExport(),
			"eve_lab_monitoring":       resourceEveLabMonitoring(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"eve_templates":     dataSourceEveTemplates(),
			"eve_network_types": dataSourceEveNetworkTypes(),
			"eve_icons":         dataSourceEveIcons(),
			"eve_status":        dataSourceEveStatus(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	timeoutStr := d.Get("timeout").(string)
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("invalid timeout format: %w", err))
	}

	config := &client.Config{
		Endpoint:           d.Get("endpoint").(string),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		InsecureSkipVerify: d.Get("insecure_skip_verify").(bool),
		Timeout:            timeout,
	}

	client, err := client.NewClient(config)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to create EVE-NG client: %w", err))
	}

	return client, nil
}
