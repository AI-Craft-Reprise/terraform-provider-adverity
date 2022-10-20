package adverity

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/devoteamgcloud/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func connection() *schema.Resource {
	return &schema.Resource{
		CreateContext: connectionCreate,
		ReadContext:   connectionRead,
		UpdateContext: connectionUpdate,
		DeleteContext: connectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: connectionImportHelper,
		},

		Schema: map[string]*schema.Schema{
			NAME: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connection.",
			},
			STACK: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the workspace to create this connection in.",
			},
			CONNECTION_TYPE_ID: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The type ID of the connection.",
			},
			CONNECTION_PARAMETERS: {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A map of extra parameters needed for connection creation. For example 'app' or 'service_account_data'.",
			},
			IS_AUTHORIZED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the connection has been authorised.",
			},
		},
		Description: "This resource will create a connection (a.k.a. authorization) of the given type in the given workspace.",
	}
}

func connectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get(NAME).(string)
	stack := d.Get(STACK).(int)
	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)

	connectionParameters, exists := d.GetOk(CONNECTION_PARAMETERS)

	parameters := []*adverityclient.Parameters{}

	if exists {
		for n, v := range connectionParameters.(map[string]interface{}) {
			parameter := new(adverityclient.Parameters)
			parameter.Value = v.(string)
			parameter.Name = n
			parameters = append(parameters, parameter)
		}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.ConnectionConfig{
		Name:       name,
		Stack:      stack,
		Parameters: parameters,
	}

	res, err := client.CreateConnection(conf, connectionTypeId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(res.ID))

	return connectionRead(ctx, d, m)
}

func connectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)

	providerConfig := m.(*config)

	client := *providerConfig.Client

	res, err, code := client.ReadConnection(d.Id(), connectionTypeId)
	if err != nil {
		if code == 404 {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set(NAME, res.Name)
	d.Set(STACK, res.Stack)
	d.Set(IS_AUTHORIZED, res.IsAuthorized)

	return diags
}

func connectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get(NAME).(string)
	stack := d.Get(STACK).(int)
	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)

	connectionParameters, exists := d.GetOk(CONNECTION_PARAMETERS)

	parameters := []*adverityclient.Parameters{}
	if exists {
		for n, v := range connectionParameters.(map[string]interface{}) {
			parameter := new(adverityclient.Parameters)
			parameter.Value = v.(string)
			parameter.Name = n
			parameters = append(parameters, parameter)
		}
	}

	providerConfig := m.(*config)

	client := *providerConfig.Client

	conf := adverityclient.ConnectionConfig{
		Name:       name,
		Stack:      stack,
		Parameters: parameters,
	}

	_, err := client.UpdateConnection(conf, d.Id(), connectionTypeId)

	if err != nil {
		return diag.FromErr(err)
	}
	return connectionRead(ctx, d, m)
}

func connectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	connectionTypeId := d.Get(CONNECTION_TYPE_ID).(int)
	providerConfig := m.(*config)

	client := *providerConfig.Client

	_, err := client.DeleteConnection(d.Id(), connectionTypeId)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func connectionImportHelper(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected connection_type:connection_id", d.Id())
	}
	connection_type_id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("could not convert connection_type (%s) to an integer", parts[0])
	}
	d.Set(CONNECTION_TYPE_ID, connection_type_id)
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
