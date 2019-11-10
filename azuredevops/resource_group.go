package azuredevops

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
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
				ConflictsWith: []string{"mail", "displayName"},
			},

			"mail": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"origin_id", "displayName"},
			},

			"displayName": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"origin_id", "mail"},
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

	// using: POST https://vssps.dev.azure.com/{organization}/_apis/graph/groups?api-version=5.1-preview.1
	cga := graph.CreateGroupArgs{}
	val, b := d.GetOk("scope")
	if b {
		cga.ScopeDescriptor = converter.String(val.(string))
	}
	val, b = d.GetOk("origin_id")
	if b {
		cga.CreationContext = &graph.GraphGroupOriginIdCreationContext{
			OriginId: converter.String(val.(string)),
		}
	} else {
		val, b = d.GetOk("mail")
		if b {
			cga.CreationContext = &graph.GraphGroupMailAddressCreationContext{
				MailAddress: converter.String(val.(string)),
			}
		} else {
			val, b = d.GetOk("displayName")
			if b {
				cga.CreationContext = &graph.GraphGroupVstsCreationContext{
					DisplayName: converter.String(val.(string)),
					Description: converter.String(d.Get("description").(string)),
				}
			} else {
				return fmt.Errorf("INTERNAL ERROR: Unable to determine strategy to create group")
			}
		}
	}
	group, err := clients.GraphClient.CreateGroup(clients.ctx, cga)
	if err != nil {
		return err
	}
	return flattenGroup(d, group)
}

func resourceAzureGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	// using: GET https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}
	getGroupArgs := graph.GetGroupArgs{
		GroupDescriptor: converter.String(d.Get("descriptor").(string)),
	}
	group, err := clients.GraphClient.GetGroup(clients.ctx, getGroupArgs)
	if err != nil {
		return err
	}
	return flattenGroup(d, group)
}

func resourceAzureGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	// using: PATCH https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}

	group, err := expandGroup(d)
	if err != nil {
		return err
	}
	// as specified in the schema, the only updatable attribute is the description

	uptGroupArgs := graph.UpdateGroupArgs{
		GroupDescriptor: converter.String(d.Get("descriptor").(string)),
		PatchDocument: &[]webapi.JsonPatchOperation{
			{
				Op:    &webapi.OperationValues.Replace,
				From:  nil,
				Path:  converter.String("/description"),
				Value: group.Description,
			},
		},
	}
	group, err = clients.GraphClient.UpdateGroup(clients.ctx, uptGroupArgs)
	if err != nil {
		return err
	}
	return flattenGroup(d, group)
}

func resourceAzureGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	// using: DELETE https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}
	delGroupArgs := graph.DeleteGroupArgs{
		GroupDescriptor: converter.String(d.Get("descriptor").(string)),
	}
	return clients.GraphClient.DeleteGroup(clients.ctx, delGroupArgs)
}

// Convert internal Terraform data structure to an AzDO data structure
func expandGroup(d *schema.ResourceData) (*graph.GraphGroup, error) {

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

func flattenGroup(d *schema.ResourceData, group *graph.GraphGroup) error {

	d.SetId(*group.Descriptor)
	d.Set("descriptor", *group.Descriptor)
	d.Set("displayName", *group.DisplayName)
	d.Set("url", *group.Url)
	d.Set("origin", *group.Origin)
	d.Set("origin_id", *group.OriginId)
	d.Set("subject_kind", *group.SubjectKind)
	d.Set("domain", *group.Domain)
	d.Set("mail", *group.MailAddress)
	d.Set("principal_name", *group.PrincipalName)
	d.Set("description", *group.Description)

	return nil
}