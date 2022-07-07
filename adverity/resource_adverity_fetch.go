package adverity

import (
	"context"
	"fmt"
	"time"

	"github.com/devoteamgcloud/adverityclient"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func fetch() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a single data fetching job in Adverity for a particular datastream.",
		CreateContext: fetchCreate,
		ReadContext:   fetchRead,
		UpdateContext: fetchUpdate,
		DeleteContext: fetchDelete,
		Schema: map[string]*schema.Schema{
			"datastream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the datastream this fetch belongs to.",
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"days", "previous_months", "current_month", "previous_weeks", "current_week", "custom"}, false),
				Description:  "The mode of the fetching jobs specifies what time windows should be used. 'Days' will fetch all data from the amount of days specified until now. The 'current' options will fetch from the beginning of the current month/week. The 'previous' options will put the start date at the beginning of the week/month a specified number of days ago, and the enddate at the end of the previous week/month.",
			},
			"days_to_fetch": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The amount of days to go back for the fetch.",
			},
			"wait_until_completion": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If set to true, Terraform will wait until the fetch has completed before reporting this resource as created.",
			},
			"disable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the resource will be created, but the fetch will wait until this value is set to false before running. Useful if the configuration for the fetch is created before the connection for the datastream is authorised.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID in Adverity for this fetching job.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the job at the time this resource was last read.",
			},
			"finished": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the job has finished.",
			},
			"is_waiting": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Variable to check if the fetch job is disabled and is waiting to be enabled.",
			},
		},
	}
}

func fetchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	disable := d.Get("disable").(bool)
	if disable {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "WARNING: The fetch is disabled, so it will not be executed.",
		})
		d.Set("is_waiting", true)
	} else {
		datastreamID := d.Get("datastream_id").(string)
		mode := d.Get("mode").(string)
		daysToFetch := d.Get("days_to_fetch").(int)
		wait := d.Get("wait_until_completion").(bool)
		providerConfig := m.(*config)
		client := *providerConfig.Client
		var err error
		var response *adverityclient.FetchResponse
		switch mode {
		case "days":
			response, err = client.FetchNumberOfDays(daysToFetch, datastreamID)
		case "previous_months":
			response, err = client.FetchPreviousMonths(daysToFetch, datastreamID)
		case "current_month":
			response, err = client.FetchCurrentMonth(datastreamID)
		case "previous_weeks":
			response, err = client.FetchPreviousWeeks(daysToFetch, datastreamID)
		case "current_week":
			response, err = client.FetchCurrentWeek(datastreamID)
		case "custom":
			err = errorString{"Custom mode not implemented yet."}
		default:
			err = errorString{fmt.Sprintf("%q is not implemented, should have been caught by schema validation.", mode)}
		}
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("job_id", response.Jobs[0].ID)
		if wait {
			complete := false
			for !complete {
				diagsFromRead := fetchRead(ctx, d, m)
				if diagsFromRead.HasError() {
					return append(diags, diagsFromRead...)
				}
				diags = append(diags, diagsFromRead...)
				complete = d.Get("finished").(bool)
				time.Sleep(10 * time.Second)
			}
		} else {
			diagsFromRead := fetchRead(ctx, d, m)
			if diagsFromRead.HasError() {
				return append(diags, diagsFromRead...)
			}
			diags = append(diags, diagsFromRead...)
		}
		d.Set("is_waiting", false)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(id)
	return diags
}

func fetchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	disable := d.Get("disable").(bool)
	if !disable {
		jobID := d.Get("job_id").(int)
		providerConfig := m.(*config)
		client := *providerConfig.Client
		res, err, code := client.ReadJob(jobID)
		if err != nil {
			if code == 404 {
				d.SetId("")
				return diags
			} else {
				return diag.FromErr(err)
			}
		}
		d.Set("status", res.StateLabel)
		d.Set("finished", res.JobEnd != "")
	}
	return diags
}

func fetchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if d.Get("is_waiting").(bool) && !d.Get("disable").(bool) {
		datastreamID := d.Get("datastream_id").(string)
		mode := d.Get("mode").(string)
		daysToFetch := d.Get("days_to_fetch").(int)
		wait := d.Get("wait_until_completion").(bool)
		providerConfig := m.(*config)
		client := *providerConfig.Client
		var err error
		var response *adverityclient.FetchResponse
		switch mode {
		case "days":
			response, err = client.FetchNumberOfDays(daysToFetch, datastreamID)
		case "previous_months":
			response, err = client.FetchPreviousMonths(daysToFetch, datastreamID)
		case "current_month":
			response, err = client.FetchCurrentMonth(datastreamID)
		case "previous_weeks":
			response, err = client.FetchPreviousWeeks(daysToFetch, datastreamID)
		case "current_week":
			response, err = client.FetchCurrentWeek(datastreamID)
		case "custom":
			err = errorString{"Custom mode not implemented yet."}
		default:
			err = errorString{fmt.Sprintf("%q is not implemented, should have been caught by schema validation.", mode)}
		}
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("job_id", response.Jobs[0].ID)
		if wait {
			complete := false
			for !complete {
				diagsFromRead := fetchRead(ctx, d, m)
				if diagsFromRead.HasError() {
					return append(diags, diagsFromRead...)
				}
				diags = append(diags, diagsFromRead...)
				complete = d.Get("finished").(bool)
				time.Sleep(10 * time.Second)
			}
			time.Sleep(10 * time.Second)
		} else {
			diags = append(diags, fetchRead(ctx, d, m)...)
		}
		d.Set("is_waiting", false)
	} else {
		return append(fetchRead(ctx, d, m), diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "WARNING: It is not possible to update a fetch job. If this happens, report this to the maintainers, since updating the fetch job should trigger a recreate.",
		})
	}
	return diags
}

func fetchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
