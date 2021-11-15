package nuage

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceKeyPairType struct{}

// Order Resource schema
func (r resourceKeyPairType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{
            "id": {
                Type: types.StringType,
                Computed: true,
            },
            "description": {
                Type: types.StringType,
                Optional: true,
            },
            "is_default": {
                Type: types.StringType,
                Optional: true,
            },
            "public_key": {
                Type: types.StringType,
                Required: true,
            },
            "user": {
                Type: types.StringType,
                Required: true,
            },
        },
    }, nil
}

// New resource instance
func (r resourceKeyPairType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return resourceKeyPair{
        p: *(p.(*provider)),
    }, nil
}

type resourceKeyPair struct {
    p provider
}

// Create a new resource
func (r resourceKeyPair) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
}

// Read resource information
func (r resourceKeyPair) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
}

// Update resource
func (r resourceKeyPair) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// Delete resource
func (r resourceKeyPair) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}
