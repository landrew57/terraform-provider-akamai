package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAkamaiBotCategoryAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAkamaiBotCategoryActionCreate,
		ReadContext:   resourceAkamaiBotCategoryActionRead,
		UpdateContext: resourceAkamaiBotCategoryActionUpdate,
		DeleteContext: resourceAkamaiBotCategoryActionDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
			verifyCategoryIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"akamai_bot_category_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceAkamaiBotCategoryActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionCreate")
	logger.Debugf("in resourceAkamaiBotCategoryActionCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "akamaiBotCategoryAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID, err := tools.GetStringValue("category_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayload, err := getJSONPayload(d, "akamai_bot_category_action", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateAkamaiBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateAkamaiBotCategoryAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateAkamaiBotCategoryAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, categoryID))

	return resourceAkamaiBotCategoryActionRead(ctx, d, m)
}

func resourceAkamaiBotCategoryActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionRead")
	logger.Debugf("in resourceAkamaiBotCategoryActionRead")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:categoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	categoryID := iDParts[2]

	request := botman.GetAkamaiBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
	}

	response, err := client.GetAkamaiBotCategoryAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetAkamaiBotCategoryAction': %s", err.Error())
		return diag.FromErr(err)
	}

	// Removing categoryId from response to suppress diff
	delete(response, "categoryId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                  configID,
		"security_policy_id":         securityPolicyID,
		"category_id":                categoryID,
		"akamai_bot_category_action": string(jsonBody),
	}
	if err = tools.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAkamaiBotCategoryActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionUpdate")
	logger.Debugf("in resourceAkamaiBotCategoryActionUpdate")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:customBotCategoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "akamaiBotCategoryAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	categoryID := iDParts[2]

	jsonPayload, err := getJSONPayload(d, "akamai_bot_category_action", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateAkamaiBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateAkamaiBotCategoryAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateAkamaiBotCategoryAction': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAkamaiBotCategoryActionRead(ctx, d, m)
}

func resourceAkamaiBotCategoryActionDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionDelete")
	logger.Debugf("in resourceAkamaiBotCategoryActionDelete")
	logger.Info("Botman API does not support akamai bot category action deletion - resource will only be removed from state")

	return nil
}
