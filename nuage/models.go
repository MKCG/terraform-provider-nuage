package nuage

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KeyPair struct {
	Id				types.String 	`tfsdk:"id"`
	Description 	types.String 	`tfsdk:"description"`
	IsDefault 		types.Bool 		`tfsdk:"is_default"`
	PublicKey 		types.String 	`tfsdk:"public_key"`
	User  			types.String 	`tfsdk:"user"`
}

type Organization struct {
	Id				types.String 	`tfsdk:"id"`
	Name 			types.String 	`tfsdk:"name"`
	Description 	types.String 	`tfsdk:"description"`
	State  			types.String 	`tfsdk:"state"`
}

type Project struct {
	Id				types.String 	`tfsdk:"id"`
	Name 			types.String 	`tfsdk:"name"`
	Description 	types.String 	`tfsdk:"description"`
	Organization 	types.String 	`tfsdk:"organization"`
}

type SecurityGroup struct {
	Id				types.String 	`tfsdk:"id"`
	Name 			types.String 	`tfsdk:"name"`
	Description 	types.String 	`tfsdk:"description"`
	Rules			[]types.String 	`tfsdk:"rules"`
}

type SecurityRule struct {
	Id			types.String 	`tfsdk:"id"`
	// Group 		types.String 	`tfsdk:"group"`
	// ReadOnly 	types.Bool 		`tfsdk:"read_only"`
	Direction	types.String 	`tfsdk:"direction"`
	Protocol	types.String 	`tfsdk:"protocol"`
	// Ethertype	types.String 	`tfsdk:"ether_type"`
	PortMin		uint 			`tfsdk:"port_min"`
	PortMax		uint 			`tfsdk:"port_max"`
	Remote		types.String 	`tfsdk:"remote"`
}
