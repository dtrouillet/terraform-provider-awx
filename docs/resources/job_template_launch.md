---
layout: "awx"
page_title: "AWX: awx_job_template_launch"
sidebar_current: "docs-awx-resource-job_template_launch"
description: |-
  *TBD*
---

# awx_job_template_launch

*TBD*

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `job_template_id` - (Required, ForceNew) Job template ID
* `extra_vars` - (Optional, ForceNew) Host or group to operate
* `limit` - (Optional, ForceNew) Host or group to operate

