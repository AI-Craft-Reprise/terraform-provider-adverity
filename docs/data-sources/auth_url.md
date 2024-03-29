---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "adverity_auth_url Data Source - terraform-provider-adverity"
subcategory: ""
description: |-
  This datasource will generate an authentication url for a connection. This url, when followed, will authenticate the connection. The url will change everytime this datasource is run.
---

# adverity_auth_url (Data Source)

This datasource will generate an authentication url for a connection. This url, when followed, will authenticate the connection. The url will change everytime this datasource is run.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **connection_id** (String) The ID of the connection the auth url belongs to.
- **connection_type_id** (String) The connection type ID for the connection the auth url belongs to.

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **url** (String) The url to authorise the connection.


