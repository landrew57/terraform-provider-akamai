package cps

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPSWarnings() *schema.Resource {
	return &schema.Resource{
		Description: "Returns a map of pre- and post-verification warnings alongside with identifiers to be used in acknowledging warnings lists",
		ReadContext: dataCPSWarningsRead,
		Schema: map[string]*schema.Schema{
			"warnings": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of pre- and post-verification warnings consisting of the warning id and description",
			},
		},
	}
}

func dataCPSWarningsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "dataCPSWarningsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	if err := d.Set("warnings", warningMap); err != nil {
		logger.WithError(err).Error("could not set cps warnings")
		return diag.FromErr(err)
	}

	d.SetId("akamai_cps_warnings")

	return nil
}
