package nuage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type flavorDataSourceType struct{}
type flavorDataSource struct {
	p provider
}

func (flavor *Flavor) Equals(other *Flavor) bool {
	if flavor.Core != other.Core {
		return false
	}

	if flavor.Disk != other.Disk {
		return false
	}

	if flavor.Ram != other.Ram {
		return false
	}

	return true
}

func (flavor flavorDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"core": {
				Type:     types.NumberType,
				Required: true,
			},
			"ram": {
				Type:     types.NumberType,
				Required: true,
			},
			"disk": {
				Type:     types.NumberType,
				Required: true,
			},
		},
	}, nil
}

func (source flavorDataSourceType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return flavorDataSource{
		p: *(p.(*provider)),
	}, nil
}

func (source flavorDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var search Flavor

	diags := req.Config.Get(ctx, &search)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	flavors, err := source.p.client.ListFlavors()

	if err != nil {
		resp.Diagnostics.AddError("Read flavorDataSource", err.Error())
		return
	}

	for _, flavor := range *flavors {
		if flavor.Equals(&search) {
			flavor.Id.Value = API_FLAVORS + "/" + flavor.Id.Value
			diags = resp.State.Set(ctx, flavor)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resp.Diagnostics.AddError("Read flavorDataSource", "Not found")
}

func (client *Client) ListFlavors() (*[]Flavor, error) {
	content, err := client.Get(API_FLAVORS)

	if err != nil {
		return nil, err
	}

	var flavors []Flavor

	for _, member := range content["hydra:member"].([]interface{}) {
		flavor := Flavor{
			Id:   types.String{Value: member.(map[string]interface{})["id"].(string)},
			Ram:  member.(map[string]interface{})["ram"].(float64),
			Core: member.(map[string]interface{})["core"].(float64),
			Disk: member.(map[string]interface{})["disk"].(float64),
		}

		flavors = append(flavors, flavor)
	}

	return &flavors, nil
}
