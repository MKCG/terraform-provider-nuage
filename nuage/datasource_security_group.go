package nuage

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type securityGroupDataSourceType struct{}
type securityGroupDataSource struct {
    p provider
}

func (group securityGroupDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
                Computed: true,
            },
        },
    }, nil
}

func (source securityGroupDataSourceType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
    return securityGroupDataSource{
        p: *(p.(*provider)),
    }, nil
}

func (source securityGroupDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
    var search SecurityGroup

    diags := req.Config.Get(ctx, &search)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    groups, err := source.p.client.ListSecurityGroups()

    if err != nil {
        resp.Diagnostics.AddError("Read securityGroupDataSource", err.Error())
        return
    }

    for _, group := range *groups {
        if group.Name.Value == search.Name.Value {
            group.Id.Value = API_SECURITY_GROUPS + "/" + group.Id.Value
            diags = resp.State.Set(ctx, group)
            resp.Diagnostics.Append(diags...)
            return
        } 
    }

    resp.Diagnostics.AddError("Read securityGroupDataSource", "Not found")
}

func (client *Client) ListSecurityGroups() (*[]SecurityGroup, error) {
    content, err := client.Get(API_SECURITY_GROUPS)

    if err != nil {
        return nil, err
    }

    var groups []SecurityGroup

    for _, member := range content["hydra:member"].([]interface{}) {
        group := SecurityGroup{
            Id      : types.String{Value: member.(map[string]interface{})["id"].(string)},
            Name      : types.String{Value: member.(map[string]interface{})["name"].(string)},
        }

        groups = append(groups, group)
    }

    return &groups, nil
}
