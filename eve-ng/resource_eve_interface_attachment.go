package eveng

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nawada0615/terraform-provider-eve-ng/internal/client"
)

func resourceEveInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEveIfAttachApply,
		ReadContext:   resourceEveIfAttachRead,
		UpdateContext: resourceEveIfAttachApply,
		DeleteContext: resourceEveIfAttachDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"lab_file":        {Type: schema.TypeString, Required: true, ForceNew: true},
			"node_id":         {Type: schema.TypeInt, Required: true, ForceNew: true},
			"interface_index": {Type: schema.TypeInt, Required: true, ForceNew: true},
			"target":          {Type: schema.TypeString, Required: true}, // network:<id> or node:<remote_node_id>[:<remote_if>]
		},
	}
}

func makeIfAttachID(labFile string, nodeID, ifIndex int) string {
	return fmt.Sprintf("%s:ifattach:%d:%d", labFile, nodeID, ifIndex)
}

func parseIfAttachID(id string) (labFile string, nodeID, ifIndex int, ok bool) {
	// <lab_file>:ifattach:<node_id>:<if_index>
	parts := strings.Split(id, ":ifattach:")
	if len(parts) != 2 {
		return "", 0, 0, false
	}
	lab := parts[0]
	rest := strings.Split(parts[1], ":")
	if len(rest) != 2 {
		return "", 0, 0, false
	}
	nid, err := strconv.Atoi(rest[0])
	if err != nil {
		return "", 0, 0, false
	}
	idx, err := strconv.Atoi(rest[1])
	if err != nil {
		return "", 0, 0, false
	}
	return lab, nid, idx, true
}

func resourceEveIfAttachApply(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile := d.Get("lab_file").(string)
	nodeID := d.Get("node_id").(int)
	ifIndex := d.Get("interface_index").(int)
	target := d.Get("target").(string)

	payload := map[string]interface{}{}
	// Decide value form based on target
	if strings.HasPrefix(target, "network:") {
		idStr := strings.TrimPrefix(target, "network:")
		nid, _ := strconv.Atoi(idStr)
		payload[strconv.Itoa(ifIndex)] = nid
	} else if strings.HasPrefix(target, "node:") {
		// pass-through remote mapping string (e.g., "<remote_id>" or "<remote_id>:<remote_if>")
		payload[strconv.Itoa(ifIndex)] = strings.TrimPrefix(target, "node:")
	} else {
		return diag.Errorf("invalid target format: %s", target)
	}

	resp, err := c.Put("api/labs"+labFile+"/nodes/"+strconv.Itoa(nodeID)+"/interfaces", payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(makeIfAttachID(labFile, nodeID, ifIndex))
	return resourceEveIfAttachRead(ctx, d, m)
}

func resourceEveIfAttachRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, nodeID, _, ok := parseIfAttachID(d.Id())
	if !ok {
		d.SetId("")
		return nil
	}

	resp, err := c.Get("api/labs" + labFile + "/nodes/" + strconv.Itoa(nodeID) + "/interfaces")
	if err != nil {
		return diag.FromErr(err)
	}

	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Ethernet []struct {
				Name      string `json:"name"`
				NetworkID int    `json:"network_id"`
			} `json:"ethernet"`
			Serial []struct {
				Name     string `json:"name"`
				RemoteID *int   `json:"remote_id"`
				RemoteIf *int   `json:"remote_if"`
			} `json:"serial"`
			ID   int    `json:"id"`
			Sort string `json:"sort"`
		} `json:"data"`
	}
	if err := c.HandleResponse(resp, &result); err != nil {
		d.SetId("")
		return nil
	}

	// Best-effort: we don't strictly verify mapping equality here
	_ = result
	return nil
}

func resourceEveIfAttachDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	labFile, nodeID, ifIndex, ok := parseIfAttachID(d.Id())
	if !ok {
		return diag.Errorf("invalid ID format")
	}

	payload := map[string]interface{}{strconv.Itoa(ifIndex): 0}
	resp, err := c.Put("api/labs"+labFile+"/nodes/"+strconv.Itoa(nodeID)+"/interfaces", payload)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.HandleResponse(resp, nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
