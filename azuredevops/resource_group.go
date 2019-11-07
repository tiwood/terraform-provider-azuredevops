package azuredevops

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceAzureGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureGroupCreate,
		Read:   resourceAzureGroupRead,
		Update: resourceAzureGroupUpdate,
		Delete: resourceAzureGroupDelete,

		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional:     true,
				ForceNew:     true,
			},

			// ***
			// One of
			//     origin_id => GraphGroupOriginIdCreationContext
			//     mail => GraphGroupMailAddressCreationContext
			//     displayName => GraphGroupVstsCreationContext
			// must be specified
			// ***

			"origin_id": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{ "mail", "displayName" },
			},

			"mail": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{ "origin_id", "displayName" },
			},

			"displayName": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{ "origin_id", "mail" },
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subject_kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"principal_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
		},
	}
}

func resourceAzureGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	return resourceAzureGroupRead(d, m)
}

func resourceAzureGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	return nil
}

func resourceAzureGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	return resourceAzureGroupRead(d, m)
}

func resourceAzureGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	return nil
}

// Convert internal Terraform data structure to an AzDO data structure
func expandGroup(clients *aggregatedClient, d *schema.ResourceData, forCreate bool) (*graph.GraphGroup, error) {

	group := &graph.GraphGroup{
		Descriptor:    converter.String(d.Get("descriptor").(string)),
		DisplayName:   converter.String(d.Get("displayName").(string)),
		Url:           converter.String(d.Get("url").(string)),
		Origin:        converter.String(d.Get("origin").(string)),
		OriginId:      converter.String(d.Get("origin_id").(string)),
		SubjectKind:   converter.String(d.Get("subject_kind").(string)),
 		Domain:        converter.String(d.Get("domain").(string)),
		MailAddress:   converter.String(d.Get("mail").(string)),
		PrincipalName: converter.String(d.Get("principal_name").(string)),
		Description:   converter.String(d.Get("description").(string)),
	}

	return group, nil
}

func flattenGroup(clients *aggregatedClient, d *schema.ResourceData, group *graph.GraphGroup) error {

	d.SetId(*group.Descriptor)
	d.Set("descriptor", *group.Descriptor)
	d.Set("displayName", *group.DisplayName)
	d.Set("url", *group.Url)
	d.Set("origin", *group.Origin)
	d.Set("origin_id", *group.OriginId)
	d.Set("subject_kind", *group.SubjectKind)
	d.Set("domain", *group.Domain)
	d.Set("mail", *group.MailAddress)
	d.Set("principal_name", *group.PrincipalName
	d.Set("description", *group.Description)

	return nil
}
