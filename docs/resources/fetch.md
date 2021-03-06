---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "adverity_fetch Resource - terraform-provider-adverity"
subcategory: ""
description: |-
  Create a single data fetching job in Adverity for a particular datastream.
---

# adverity_fetch (Resource)

Create a single data fetching job in Adverity for a particular datastream.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **datastream_id** (String) The ID of the datastream this fetch belongs to.
- **days_to_fetch** (Number) The amount of days to go back for the fetch.
- **mode** (String) The mode of the fetching jobs specifies what time windows should be used. 'Days' will fetch all data from the amount of days specified until now. The 'current' options will fetch from the beginning of the current month/week. The 'previous' options will put the start date at the beginning of the week/month a specified number of days ago, and the enddate at the end of the previous week/month.

### Optional

- **disable** (Boolean) If set to true, the resource will be created, but the fetch will wait until this value is set to false before running. Useful if the configuration for the fetch is created before the connection for the datastream is authorised.
- **id** (String) The ID of this resource.
- **wait_until_completion** (Boolean) If set to true, Terraform will wait until the fetch has completed before reporting this resource as created.

### Read-Only

- **finished** (Boolean) Whether the job has finished.
- **is_waiting** (Boolean) Variable to check if the fetch job is disabled and is waiting to be enabled.
- **job_id** (Number) The ID in Adverity for this fetching job.
- **status** (String) The status of the job at the time this resource was last read.


