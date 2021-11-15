package nuage

import (
    "context"
    // "math/big"
    // "strconv"
    // "time"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceSecurityGroupType struct{}

// Order Resource schema
func (r resourceSecurityGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{
            "id": {
                Type: types.StringType,
                Computed: true,
            },
            "name": {
                Type: types.StringType,
                Optional: true,
            },
            "description": {
                Type: types.StringType,
                Optional: true,
            },
            "rules": {
                Type: types.ListType{
                    ElemType: types.StringType,
                },
                Required: true,
            },
        },
    }, nil
}

// New resource instance
func (r resourceSecurityGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return resourceSecurityGroup{
        p: *(p.(*provider)),
    }, nil
}

type resourceSecurityGroup struct {
    p provider
}

// Create a new resource
func (r resourceSecurityGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
    resp.Diagnostics.AddError("Create SecurityGroup", "Not implemented")
}

// Read resource information
func (r resourceSecurityGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    resp.Diagnostics.AddError("Read SecurityGroup", "Not implemented")
}

// Update resource
func (r resourceSecurityGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    resp.Diagnostics.AddError("Update SecurityGroup", "Not implemented")
}

// Delete resource
func (r resourceSecurityGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
    resp.Diagnostics.AddError("Delete SecurityGroup", "Not implemented")
}
