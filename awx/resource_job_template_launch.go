/*
*TBD*

Example Usage

```hcl
data "awx_inventory" "default" {
  name            = "private_services"
  organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "baseconfig" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = data.awx_inventory.default.id
  project_id     = awx_project.base_service_config.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}

resource "awx_job_template_launch" "now" {
  job_template_id = awx_job_template.baseconfig.id
}
```

*/

package awx

import (
	"context"
	"fmt"
	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"time"
)

func resourceJobTemplateLaunch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateLaunchCreate,
		ReadContext:   resourceJobRead,
		DeleteContext: resourceJobDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Job template ID",
				ForceNew:    true,
			}, "limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host or group to operate",
				ForceNew:    true,
			}, "extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host or group to operate",
				ForceNew:    true,
			},
		},
	}
}

func resourceJobTemplateLaunchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	jobTemplateID := d.Get("job_template_id").(int)
	_, err := awxService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", jobTemplateID, err)
	}

	res, err := awxService.Launch(jobTemplateID, map[string]interface{}{
		"extra_vars": d.Get("extra_vars").(string),
		"limit":      d.Get("limit").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Failed to create Template Launch %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d, failed to create %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceJobRead(ctx, d, m)
	//return diags
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	jobService := client.JobService
	jobId, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("Failed to get JobID %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to get JobID",
			Detail:   fmt.Sprintf("Unable to get JobID %s : %s", d.Id(), err.Error()),
		})
		return diags
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate)-time.Minute, func() *retry.RetryError {
		job, err := jobService.GetJob(jobId, make(map[string]string))
		if err != nil {
			log.Printf("Failed to get Job %v", err)

			return retry.NonRetryableError(fmt.Errorf("Unable to get Job %s : %s", d.Id(), err.Error()))
		}

		if job.Finished.IsZero() {
			return retry.RetryableError(fmt.Errorf("Job is not finished"))
		} else if job.Status == "failed" || job.Status == "error" || job.Status == "canceled" {
			id := d.Id()
			d.SetId("")
			return retry.NonRetryableError(fmt.Errorf("Job %s is in error state : %s", id, job.Status))
		}

		return nil
	})

	if err != nil {
		diags = diag.FromErr(err)
	}

	return diags
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobService
	jobID, diags := convertStateIDToNummeric("Delete Job", d)
	_, err := awxService.GetJob(jobID, map[string]string{})
	if err != nil {
		return buildDiagNotFoundFail("job", jobID, err)
	}

	d.SetId("")
	return diags
}
