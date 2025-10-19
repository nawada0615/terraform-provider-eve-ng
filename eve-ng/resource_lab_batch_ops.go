package eveng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

// batchOperationType represents the type of batch operation
type batchOperationType string

const (
	batchStart batchOperationType = "start"
	batchStop  batchOperationType = "stop"
	batchWipe  batchOperationType = "wipe"
)

// createBatchOperation performs a batch operation on lab nodes
func createBatchOperation(ctx context.Context, d *schema.ResourceData, m interface{}, opType batchOperationType, readFunc func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)

	payload := map[string]interface{}{}
	if nodeIDs, ok := d.GetOk("node_ids"); ok {
		payload["nodes"] = nodeIDs
	}

	endpoint := "api/labs" + labFile + "/nodes/" + string(opType)
	resp, err := c.Post(endpoint, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(labFile + ":batch_" + string(opType))
	return readFunc(ctx, d, m)
}

func resourceEveLabBatchStart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabBatchStartCreate,
		ReadContext:   resourceEveLabBatchStartRead,
		DeleteContext: resourceEveLabBatchStartDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file": {Type: schema.TypeString, Required: true, ForceNew: true},
			"node_ids": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeInt}},
		},
	}
}

func resourceEveLabBatchStartCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return createBatchOperation(ctx, d, m, batchStart, resourceEveLabBatchStartRead)
}

func resourceEveLabBatchStartRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Batch operations are stateless, just verify lab exists
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceEveLabBatchStartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Batch start doesn't need explicit deletion
	return nil
}

func resourceEveLabBatchStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabBatchStopCreate,
		ReadContext:   resourceEveLabBatchStopRead,
		DeleteContext: resourceEveLabBatchStopDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file": {Type: schema.TypeString, Required: true, ForceNew: true},
			"node_ids": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeInt}},
		},
	}
}

func resourceEveLabBatchStopCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return createBatchOperation(ctx, d, m, batchStop, resourceEveLabBatchStopRead)
}

func resourceEveLabBatchStopRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceEveLabBatchStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEveLabBatchWipe() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveLabBatchWipeCreate,
		ReadContext:   resourceEveLabBatchWipeRead,
		DeleteContext: resourceEveLabBatchWipeDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file": {Type: schema.TypeString, Required: true, ForceNew: true},
			"node_ids": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeInt}},
		},
	}
}

func resourceEveLabBatchWipeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return createBatchOperation(ctx, d, m, batchWipe, resourceEveLabBatchWipeRead)
}

func resourceEveLabBatchWipeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)

	resp, err := c.Get("api/labs" + labFile)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceEveLabBatchWipeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
