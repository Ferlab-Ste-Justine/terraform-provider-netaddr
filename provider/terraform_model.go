package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetRangeIdsFromResource(d *schema.ResourceData) []string {
	keyPrefixes := d.Get("range_ids")
	rangeIds := []string{}
	for _, val := range (keyPrefixes.(*schema.Set)).List() {
		rangeIds = append(rangeIds, val.(string))
	}

	return rangeIds
}