package nuage

import (
    "context"
    "errors"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceProjectType struct{}

// Order Resource schema
func (r resourceProjectType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
            "name": {
                Type: types.StringType,
                Required: true,
            },
            "organization": {
                Type: types.StringType,
                Computed: true,
            },
        },
    }, nil
}

// New resource instance
func (r resourceProjectType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return resourceProject{
        p: *(p.(*provider)),
    }, nil
}

type resourceProject struct {
    p provider
}

// Create a new resource
func (r resourceProject) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
    var project Project

    diags := req.Plan.Get(ctx, &project)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    organizationId, err := r.p.client.GetOrganizationId()

    if err != nil {
        resp.Diagnostics.AddError("Create project", err.Error())
        return
    }

    id, err := r.p.client.CreateProject(project)

    if err != nil {
        resp.Diagnostics.AddError("Create project", err.Error())
        return
    }

    // Generate resource state struct
    var result = Project{
        Id:             types.String{Value: id},
        Name:           types.String{Value: project.Name.Value},
        Description:    types.String{Value: project.Description.Value},
        Organization:   types.String{Value: organizationId},
    }

    diags = resp.State.Set(ctx, result)
    resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r resourceProject) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    var project Project

    diags := req.State.Get(ctx, &project)
    resp.Diagnostics.Append(diags...)

    if resp.Diagnostics.HasError() {
        return
    }

    id := project.Id.Value

    result, err := r.p.client.GetProject(id)

    if err != nil {
        resp.Diagnostics.AddError("Read project", err.Error())
        return
    }

    diags = resp.State.Set(ctx, result)
    resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceProject) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    resp.Diagnostics.AddError("Update project", "Not implemented")
}

// Delete resource
func (r resourceProject) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
    resp.Diagnostics.AddError("Delete project", "Not implemented")
}

func (client *Client) CreateProject(project Project) (string, error) {
    organizationId, err := client.GetOrganizationId()

    if err != nil {
        return "", err
    }

    payload := map[string]interface{}{
        "name": project.Name.Value,
        "description": project.Description.Value,
        "organization": organizationId,
    }

    return client.CreateResource(API_PROJECTS, payload)
}

func (client *Client) GetProject(id string) (*Project, error) {
    content, err := client.Get(API_PROJECTS + "/" + id)

    if err != nil {
        return nil, err
    }

    organizationId, err := client.GetOrganizationId()

    if err != nil {
        return nil, err
    }

    project := Project{
        Id              : types.String{Value: content["id"].(string)},
        Name            : types.String{Value: content["name"].(string)},
        Description     : types.String{Value: content["description"].(string)},
        Organization    : types.String{Value: organizationId},
    }

    return &project, nil
}

func (client *Client) DeleteProject(id string) error {
    return client.Delete(API_PROJECTS, id)
}

func (client *Client) ListOrganizations() (*[]Organization, error) {
    content, err := client.Get(API_ORGANIZATIONS)

    if err != nil {
        return nil, err
    }

    var organizations []Organization

    for _, member := range content["hydra:member"].([]interface{}) {
        organization := Organization{
            Id          : types.String{Value: member.(map[string]interface{})["id"].(string)},
            Name        : types.String{Value: member.(map[string]interface{})["name"].(string)},
            Description : types.String{Value: member.(map[string]interface{})["description"].(string)},
            State       : types.String{Value: member.(map[string]interface{})["state"].(string)},
        }

        organizations = append(organizations, organization)
    }

    return &organizations, nil
}

func (client *Client) GetOrganizationId() (string, error) {
    organizations, err := client.ListOrganizations()

    if err != nil {
        return "", err
    }

    if len(*organizations) == 0 {
        return "", errors.New("No organization defined")
    }

    organizationId := API_ORGANIZATIONS + "/" + (*organizations)[0].Id.Value

    return organizationId, nil
}
