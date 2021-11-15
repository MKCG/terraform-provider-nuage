package nuage

import (
	"context"
	"os"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var stderr = os.Stderr

const API_HOST = "https://api.nua.ge"

const API_AUTH 				= "/arya/auth"
const API_ORGANIZATIONS    	= "/arya/organizations"
const API_PROJECTS 			= "/arya/projects"
const API_KEYPAIRS 			= "/arya/keypairs"
const API_IMAGES 			= "/rockefeller/images"
const API_IPS 				= "/rockefeller/ips"
const API_SECURITY_GROUPS 	= "/rockefeller/security_groups"
const API_SECURITY_RULES 	= "/rockefeller/security_rules"
const API_SERVERS 			= "/rockefeller/servers"

type provider struct {
	configured bool
	client     *Client
}

type Client struct {
	organization 	string
	name 			string
	accessToken  	string
}

func New() tfsdk.Provider {
	return &provider{}
}

func (client *Client) CreateResource(api string, content map[string]interface{}) (string, error) {	
    payload, _ := json.Marshal(content)

    body, err := client.Post(api, payload)

    if err != nil {
    	return "", err
    }

    if val, ok := body["id"]; ok {
    	return val.(string), nil
    }

    return "", errors.New("Id is missing")
}

func (client* Client) Execute(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + client.accessToken)

	httpClient := &http.Client{}
	return httpClient.Do(req)
}

func (client *Client) Get(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", API_HOST + url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Execute(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var decoded map[string]interface{}

	if err = json.Unmarshal(body, &decoded) ; err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if val, ok := decoded["message"]; ok {
			return nil, errors.New(val.(string))
		}

		return nil, errors.New(string(body))
	}

	return decoded, nil
}

func (client *Client) Post(url string, payload []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", API_HOST + url, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	resp, err := client.Execute(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var decoded map[string]interface{}

	if err = json.Unmarshal(body, &decoded) ; err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		if val, ok := decoded["message"]; ok {
			return nil, errors.New(val.(string))
		}

		return nil, errors.New(string(body))
	}

	return decoded, nil
}

func (client *Client) Delete(url string, id string) error {
	req, err := http.NewRequest("DELETE", API_HOST + url + "/" + id, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer " + client.accessToken)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return errors.New("Could not delete the resource " + id)
	}

	return nil
}

func NewClient(organization *string, name *string, password *string, resp *tfsdk.ConfigureProviderResponse) (*Client, error) {
	payload, _ := json.Marshal(map[string]string{
		"organization": *organization,
		"name": *name,
		"password": *password,
	})

	auth, err := http.Post(API_HOST + API_AUTH, "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	defer auth.Body.Close()

	body, err := ioutil.ReadAll(auth.Body)

	if err != nil {
		return nil, err
	}

	var decoded map[string]interface{}

	if err = json.Unmarshal(body, &decoded) ; err != nil {
		return nil, err
	}

	if auth.StatusCode >= 400 {
		if val, ok := decoded["message"]; ok {
			return nil, errors.New(val.(string))
		}

		return nil, errors.New("Unable to authenticate")
	}

	val, found := decoded["token"]

	if !found {
		return nil, errors.New("JWT token missing")
	}

	return &Client{
		organization: *organization,
		name: *name,
		accessToken: val.(string),
	}, nil
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"organization": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Organization    types.String `tfsdk:"organization"`
	Name 			types.String `tfsdk:"name"`
	Password 		types.String `tfsdk:"password"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a user to the provider
	var name string
	if config.Name.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as name",
		)
		return
	}

	if config.Name.Null {
		name = os.Getenv("NUAGE_NAME")
	} else {
		name = config.Name.Value
	}

	if name == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find name",
			"Name cannot be an empty string",
		)
		return
	}

	// User must provide a password to the provider
	var password string
	if config.Password.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	if config.Password.Null {
		password = os.Getenv("NUAGE_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if password == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find password",
			"password cannot be an empty string",
		)
		return
	}

	// User must specify a organization
	var organization string
	if config.Organization.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as organization",
		)
		return
	}

	if config.Organization.Null {
		organization = os.Getenv("NUAGE_ORGANIZATION")
	} else {
		organization = config.Organization.Value
	}

	if organization == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find organization",
			"Organization cannot be an empty string",
		)
		return
	}

	// Create a new HashiCups client and set it to the provider client
	c, err := NewClient(&organization, &name, &password, resp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"nuage_keypair": resourceKeyPairType{},
		"nuage_project": resourceProjectType{},
		"nuage_security_group": resourceSecurityGroupType{},
		"nuage_security_rule": resourceSecurityRuleType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
