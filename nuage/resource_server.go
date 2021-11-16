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

type resourceServerType struct{}

// Order Resource schema
func (r resourceServerType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{
            "id": {
                Type: types.StringType,
                Computed: true,
            },
            "name": {
                Type: types.StringType,
                Required: true,
            },
            "description": {
                Type: types.StringType,
                Optional: true,
            },
            "flavor": {
                Type: types.StringType,
                Required: true,
            },
            "image": {
                Type: types.StringType,
                Required: true,
            },
            "keypair": {
                Type: types.StringType,
                Optional: true,
            },
            "project": {
                Type: types.StringType,
                Required: true,
            },
            "security_groups": {
                Type: types.ListType{
                    ElemType: types.StringType,
                },
                Optional: true,
            },
        },
    }, nil
}

// New resource instance
func (r resourceServerType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return resourceServer{
        p: *(p.(*provider)),
    }, nil
}

type resourceServer struct {
    p provider
}

// Create a new resource
func (r resourceServer) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
    resp.Diagnostics.AddError("Create Server", "Not implemented")
}

// Read resource information
func (r resourceServer) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    resp.Diagnostics.AddError("Read Server", "Not implemented")
}

// Update resource
func (r resourceServer) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    resp.Diagnostics.AddError("Update Server", "Not implemented")
}

// Delete resource
func (r resourceServer) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
    resp.Diagnostics.AddError("Delete Server", "Not implemented")
}
