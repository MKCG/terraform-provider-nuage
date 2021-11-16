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
            // "security_groups": {
            //     Type: types.ListType{
            //         ElemType: types.StringType,
            //     },
            //     Optional: true,
            // },
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
    var server Server

    diags := req.Plan.Get(ctx, &server)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    id, err := r.p.client.CreateServer(server)

    if err != nil {
        resp.Diagnostics.AddError("Create Server", err.Error())
        return
    }

    // Generate resource state struct
    var result = Server{
        Id:             types.String{Value: id},
        Name:           types.String{Value: server.Name.Value},
        Description:    types.String{Value: server.Description.Value},
        Project:        types.String{Value: server.Project.Value},
        Flavor:         types.String{Value: server.Flavor.Value},
        Image:          types.String{Value: server.Image.Value},
        KeyPair:        types.String{Value: server.KeyPair.Value},
    }

    diags = resp.State.Set(ctx, result)
    resp.Diagnostics.Append(diags...)
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

func (client *Client) CreateServer(server Server) (string, error) {
    payload := map[string]interface{}{
        "name":         server.Name.Value,
        "description":  server.Description.Value,
        "project":      server.Project.Value,
        "flavor":       server.Flavor.Value,
        "image":        server.Image.Value,
        "keypair":      server.KeyPair.Value,
    }

    return client.CreateResource(API_SERVERS, payload)
}
