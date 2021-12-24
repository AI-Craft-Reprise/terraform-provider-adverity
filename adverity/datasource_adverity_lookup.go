package adverity

import (
	"strconv"
	"strings"
	"time"

	"github.com/fourcast/adverityclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceAdverityLookup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"argument": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"id_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"text": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"search_terms": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filtered_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"match_exact_term": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Read: lookupDataSource,
	}
}

func lookupDataSource(d *schema.ResourceData, m interface{}) error {
	parameters := d.Get("parameters").([]interface{})
	url := d.Get("url").(string)

	providerConfig := m.(*config)
	client := *providerConfig.Client

	params := []adverityclient.Query{}

	for _, parameter := range parameters {
		p := parameter.(map[string]interface{})
		param := adverityclient.Query{
			Key:   p["argument"].(string),
			Value: p["value"].(string),
		}
		params = append(params, param)
	}

	res, err := client.DoLookup(url, params)
	if err != nil {
		return err
	}
	if res.Error != "nil" {
		return errorString{"Error while doing lookup: " + res.Error}
	}

	idMappings := flattenLookup(&res.Results)
	if err := d.Set("id_mappings", idMappings); err != nil {
		return err
	}

	search_terms, exists := d.GetOk("search_terms")
	filtered_list := []int{}
	match_exact := d.Get("match_exact_term").(bool)

	if exists {
		for _, term := range search_terms.([]interface{}) {
			if term == nil {
				return errorString{"Failed doing lookup: empty string not permitted"}
			}
			found_match := false
			for _, mapping := range idMappings {
				mapping_cast := mapping.(map[string]interface{})
				id := mapping_cast["id"].(int)
				name := mapping_cast["text"].(string)
				if name == term.(string) || (strings.Contains(strings.ToLower(name), strings.ToLower(term.(string))) && !match_exact) {
					filtered_list = append(filtered_list, id)
					found_match = true
				}
			}
			if !found_match {
				return errorString{"Error while doing lookup: could not find a match for term \"" + term.(string) + "\""}
			}
		}
	}

	allFilters := make(map[int]bool)
	list := []int{}
	for _, item := range filtered_list {
		if _, value := allFilters[item]; !value {
			allFilters[item] = true
			list = append(list, item)
		}
	}

	d.Set("filtered_list", list)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func flattenLookup(idMappings *[]adverityclient.IDMapping) []interface{} {
	if idMappings != nil {
		mappings := make([]interface{}, len(*idMappings), len(*idMappings))

		for i, idMapping := range *idMappings {
			mapping := make(map[string]interface{})
			mapping["id"] = idMapping.ID
			mapping["text"] = idMapping.Name
			mappings[i] = mapping
		}
		return mappings
	}
	return make([]interface{}, 0)
}
