package nuage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type imageDataSourceType struct{}
type imageDataSource struct {
	p provider
}

func (image *Image) Equals(other *Image) bool {
    if image.OsName != other.OsName {
        return false
    }

    if image.OsVersion != other.OsVersion {
        return false
    }

    return true
}

func (image imageDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{
            "id": {
                Type: types.StringType,
                Computed: true,
            },
            "name": {
                Type: types.StringType,
                Computed: true,
            },
            "description": {
                Type: types.StringType,
                Computed: true,
            },
            "is_default": {
                Type: types.BoolType,
                Computed: true,
            },
            "is_public": {
                Type: types.BoolType,
                Computed: true,
            },
            "os_admin_user": {
                Type: types.StringType,
                Computed: true,
            },
            "os_name": {
                Type: types.StringType,
                Required: true,
            },
            "os_version": {
                Type: types.StringType,
                Required: true,
            },
        },
    }, nil
}

func (image imageDataSourceType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
    return imageDataSource{
    	p: *(p.(*provider)),
    }, nil
}

func (source imageDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
    var search Image

    diags := req.Config.Get(ctx, &search)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    images, err := source.p.client.ListImages()

    if err != nil {
        resp.Diagnostics.AddError("Read imageDataSource", err.Error())
        return
    }

    for _, image := range *images {
        if image.Equals(&search) {
            diags = resp.State.Set(ctx, image)
            resp.Diagnostics.Append(diags...)
            return
        } 
    }

    resp.Diagnostics.AddError("Read imageDataSource", "Not found")
}

func (client *Client) ListImages() (*[]Image, error) {
    content, err := client.Get(API_IMAGES)

    if err != nil {
        return nil, err
    }

    var images []Image

    for _, member := range content["hydra:member"].([]interface{}) {
        image := Image{
            Id          : types.String{Value: member.(map[string]interface{})["id"].(string)},
            Name        : types.String{Value: member.(map[string]interface{})["name"].(string)},
            Description : types.String{Value: member.(map[string]interface{})["description"].(string)},
            OsName      : types.String{Value: member.(map[string]interface{})["osName"].(string)},
            OsVersion   : types.String{Value: member.(map[string]interface{})["osVersion"].(string)},
        }

        images = append(images, image)
    }

    return &images, nil
}
