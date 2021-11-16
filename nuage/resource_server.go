package nuage

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceServerType struct{}

type resourceServer struct {
	p provider
}

func (r resourceServerType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"flavor": {
				Type:     types.StringType,
				Required: true,
			},
			"image": {
				Type:     types.StringType,
				Required: true,
			},
			"keypair": {
				Type:     types.StringType,
				Optional: true,
			},
			"project": {
				Type:     types.StringType,
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

func (r resourceServerType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceServer{
		p: *(p.(*provider)),
	}, nil
}

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

	var result = Server{
		Id:             types.String{Value: id},
		Name:           types.String{Value: server.Name.Value},
		Description:    types.String{Value: server.Description.Value},
		Project:        types.String{Value: server.Project.Value},
		Flavor:         types.String{Value: server.Flavor.Value},
		Image:          types.String{Value: server.Image.Value},
		KeyPair:        types.String{Value: server.KeyPair.Value},
		SecurityGroups: server.SecurityGroups,
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceServer) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var server Server

	diags := req.State.Get(ctx, &server)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := server.Id.Value

	result, err := r.p.client.GetServer(id)

	if err != nil {
		resp.Diagnostics.AddError("Read Server", err.Error())
		return
	}

	/**
	 * This is very ugly : the KeyPair can no longer be updated
	 * However there is no other way around since the KeyPair id is not returned by the Nuage API
	 */
	result.KeyPair.Value = server.KeyPair.Value

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceServer) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	resp.Diagnostics.AddError("Update Server", "Not implemented")
}

func (r resourceServer) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var server Server

	diags := req.State.Get(ctx, &server)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := server.Id.Value

	if err := r.p.client.DeleteServer(id); err != nil {
		resp.Diagnostics.AddError("Delete Server", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (client *Client) CreateServer(server Server) (string, error) {
	groups := []string{}

	for _, group := range server.SecurityGroups {
		groups = append(groups, group.Value)
	}

	payload := map[string]interface{}{
		"name":           server.Name.Value,
		"description":    server.Description.Value,
		"project":        server.Project.Value,
		"flavor":         server.Flavor.Value,
		"image":          server.Image.Value,
		"keypair":        server.KeyPair.Value,
		"securityGroups": groups,
	}

	return client.CreateResource(API_SERVERS, payload)
}

func (client *Client) GetServer(id string) (*Server, error) {
	content, err := client.Get(API_SERVERS + "/" + id)

	if err != nil {
		return nil, err
	}

	project := strings.Replace(content["project"].(string), API_PROJECTS+"/", "", -1)

	var groups []types.String

	for _, group := range content["securityGroups"].([]interface{}) {
		groupId := group.(map[string]interface{})["@id"].(string)

		groups = append(groups, types.String{Value: groupId})
	}

	server := Server{
		Id:             types.String{Value: content["id"].(string)},
		Name:           types.String{Value: content["name"].(string)},
		Description:    types.String{Value: content["description"].(string)},
		Project:        types.String{Value: project},
		Flavor:         types.String{Value: content["flavor"].(string)},
		Image:          types.String{Value: content["image"].(string)},
		SecurityGroups: groups,
		// KeyPair         : types.String{Value: content["keypair"].(string)},
	}

	return &server, nil
}

func (client *Client) DeleteServer(id string) error {
	return client.Delete(API_SERVERS, id)
}
