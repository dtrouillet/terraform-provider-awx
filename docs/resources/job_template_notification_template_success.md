---
layout: "awx"
page_title: "AWX: awx_job_template_notification_template_success"
sidebar_current: "docs-awx-resource-job_template_notification_template_success"
description: |-
  *TBD*
---

# awx_job_template_notification_template_success

*TBD*

## Example Usage

```hcl
resource "awx_job_template_notification_template_success" "baseconfig" {
  job_template_id            = awx_job_template.baseconfig.id
  notification_template_id   = awx_notification_template.default.id
}
```

## Argument Reference

The following arguments are supported:

* `job_template_id` - (Required, ForceNew) 
* `notification_template_id` - (Required, ForceNew) 

