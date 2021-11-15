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
                Type: types.BoolType,
                Optional: true,
            },
            "public_key": {
                Type: types.StringType,
                Required: true,
            },
            "user": {
                Type: types.StringType,
                Computed: true,
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
    var key KeyPair

    diags := req.Plan.Get(ctx, &key)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    id, err := r.p.client.CreateKeyPair(key)

    if err != nil {
        resp.Diagnostics.AddError("Create KeyPair", err.Error())
        return
    }

    // Generate resource state struct
    var result = KeyPair{
        Id:             types.String{Value: id},
        Description:    types.String{Value: key.Description.Value},
        IsDefault:      types.Bool{Value: key.IsDefault.Value},
        PublicKey:      types.String{Value: key.PublicKey.Value},
        User:           types.String{Value: r.p.client.name},
    }

    diags = resp.State.Set(ctx, result)
    resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r resourceKeyPair) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    var key KeyPair

    diags := req.State.Get(ctx, &key)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    id := key.Id.Value

    result, err := r.p.client.GetKeyPair(id)

    if err != nil {
        resp.Diagnostics.AddError("Read KeyPair", err.Error())
        return
    }

    diags = resp.State.Set(ctx, result)
    resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceKeyPair) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    resp.Diagnostics.AddError("Update KeyPair", "Not implemented")
}

// Delete resource
func (r resourceKeyPair) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
    resp.Diagnostics.AddError("Delete KeyPair", "Not implemented")
}

func (client *Client) CreateKeyPair(key KeyPair) (string, error) {
    payload := map[string]interface{}{
        "name": client.name,
        "description": key.Description.Value,
        "publicKey": key.PublicKey.Value,
        "isDefault": key.IsDefault.Value,
    }

    return client.CreateResource(API_KEYPAIRS, payload)
}

func (client *Client) GetKeyPair(id string) (*KeyPair, error) {
    content, err := client.Get(API_KEYPAIRS + "/" + id)

    if err != nil {
        return nil, err
    }

    key := KeyPair{
        Id              : types.String{Value: content["id"].(string)},
        Description     : types.String{Value: content["description"].(string)},
        IsDefault       : types.Bool{Value: content["isDefault"].(bool)},
        PublicKey       : types.String{Value: content["publicKey"].(string)},
        User            : types.String{Value: client.name},
    }

    return &key, nil
}
