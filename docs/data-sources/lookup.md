---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "adverity_lookup Data Source - terraform-provider-adverity"
subcategory: ""
description: |-
  
---

# adverity_lookup (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **expect_string** (Boolean)
- **url** (String)

### Optional

- **disable_lookup** (Boolean)
- **id** (String) The ID of this resource.
- **match_exact_term** (Boolean)
- **parameters** (Block List) (see [below for nested schema](#nestedblock--parameters))
- **search_terms** (List of String)

### Read-Only

- **filtered_list** (List of String)
- **id_mappings** (List of Object) (see [below for nested schema](#nestedatt--id_mappings))

<a id="nestedblock--parameters"></a>
### Nested Schema for `parameters`

Required:

- **argument** (String)
- **value** (String)


<a id="nestedatt--id_mappings"></a>
### Nested Schema for `id_mappings`

Read-Only:

- **id** (String)
- **text** (String)


