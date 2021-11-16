/**
 * @author Kévin Masseix <masseix.kevin@gmail.com>
 */

package nuage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceSecurityRuleType struct{}

type resourceSecurityRule struct {
	p provider
}

func (r resourceSecurityRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"direction": {
				Type:     types.StringType,
				Required: true,
			},
			"protocol": {
				Type:     types.StringType,
				Required: true,
			},
			"port_min": {
				Type:     types.NumberType,
				Required: true,
			},
			"port_max": {
				Type:     types.NumberType,
				Required: true,
			},
			"remote": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

func (r resourceSecurityRuleType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceSecurityRule{
		p: *(p.(*provider)),
	}, nil
}

func (r resourceSecurityRule) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var rule SecurityRule

	diags := req.Plan.Get(ctx, &rule)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.p.client.CreateSecurityRule(rule)

	if err != nil {
		resp.Diagnostics.AddError("Create SecurityRule", err.Error())
		return
	}

	var result = SecurityRule{
		Id:        types.String{Value: id},
		Direction: types.String{Value: rule.Direction.Value},
		Protocol:  types.String{Value: rule.Protocol.Value},
		PortMin:   rule.PortMin,
		PortMax:   rule.PortMax,
		Remote:    types.String{Value: rule.Remote.Value},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceSecurityRule) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	resp.Diagnostics.AddError("Read SecurityRule", "Not implemented")
}

func (r resourceSecurityRule) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	resp.Diagnostics.AddError("Update SecurityRule", "Not implemented")
}

func (r resourceSecurityRule) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	resp.Diagnostics.AddError("Delete SecurityRule", "Not implemented")
}

func (client *Client) CreateSecurityRule(rule SecurityRule) (string, error) {
	payload := map[string]interface{}{
		"direction": rule.Direction.Value,
		"protocol":  rule.Protocol.Value,
		"portMin":   rule.PortMin,
		"portMax":   rule.PortMax,
		"remote":    rule.Remote.Value,
	}

	return client.CreateResource(API_SECURITY_RULES, payload)
}
