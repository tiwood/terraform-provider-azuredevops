package azuredevops

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

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
			//     display_name => GraphGroupVstsCreationContext
			// must be specified
			// ***

			"origin_id": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"mail", "display_name"},
			},

			"mail": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"origin_id", "display_name"},
			},

			"display_name": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
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
				Set:      schema.HashString,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: POST https://vssps.dev.azure.com/{organization}/_apis/graph/groups?api-version=5.1-preview.1
	cga := graph.CreateGroupArgs{}
	val, b := d.GetOk("scope")
	if b {
		uuid, _ := uuid.Parse(val.(string))
		desc, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{
			StorageKey: &uuid,
		})
		if err != nil {
			return err
		}
		cga.ScopeDescriptor = desc.Value
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
			val, b = d.GetOk("display_name")
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
	group, err := clients.GraphClient.CreateGroup(clients.Ctx, cga)
	if err != nil {
		return err
	}
	if group.Descriptor == nil {
		return fmt.Errorf("DevOps REST API returned group object without descriptor")
	}

	members := expandGroupMembers(*group.Descriptor, d.Get("members").(*schema.Set))
	if err := addMembers(clients, members); err != nil {
		return fmt.Errorf("Error adding group memberships during create: %+v", err)
	}
	return flattenGroup(d, group, members)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: GET https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}
	getGroupArgs := graph.GetGroupArgs{
		GroupDescriptor: converter.String(d.Id()),
	}
	group, err := clients.GraphClient.GetGroup(clients.Ctx, getGroupArgs)
	if err != nil {
		return err
	}
	if group.Descriptor == nil {
		return fmt.Errorf("DevOps REST API returned group object without descriptor; group %s", d.Id())
	}
	members, err := groupReadMembers(*group.Descriptor, clients)
	if err != nil {
		return err
	}
	return flattenGroup(d, group, members)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: PATCH https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}

	group, members, err := expandGroup(d)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		// FIXME: does not clear the description if empty! Error: value cannot be nil
		uptGroupArgs := graph.UpdateGroupArgs{
			GroupDescriptor: converter.String(d.Id()),
			PatchDocument: &[]webapi.JsonPatchOperation{
				{
					Op:    &webapi.OperationValues.Replace,
					From:  nil,
					Path:  converter.String("/description"),
					Value: group.Description,
				},
			},
		}
		group, err = clients.GraphClient.UpdateGroup(clients.Ctx, uptGroupArgs)
		if err != nil {
			return err
		}
	}
	if d.HasChange("members") {
		// FIXME: implement strategy to update memebrs like in azuredevops\resource_group_membership.go
		// Delete all members
		currentMembers, err := groupReadMembers(*group.Descriptor, clients)
		if err != nil {
			return err
		}
		if err := removeMembers(clients, currentMembers); err != nil {
			return err
		}
		// add all members from state
		if err := addMembers(clients, members); err != nil {
			return fmt.Errorf("Error adding group memberships during update: %+v", err)
		}
	}
	return flattenGroup(d, group, members)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: DELETE https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}
	delGroupArgs := graph.DeleteGroupArgs{
		GroupDescriptor: converter.String(d.Id()),
	}
	return clients.GraphClient.DeleteGroup(clients.Ctx, delGroupArgs)
}

// Convert internal Terraform data structure to an AzDO data structure
func expandGroup(d *schema.ResourceData) (*graph.GraphGroup, *[]graph.GraphMembership, error) {

	group := graph.GraphGroup{
		Descriptor:    converter.String(d.Id()),
		DisplayName:   converter.String(d.Get("display_name").(string)),
		Url:           converter.String(d.Get("url").(string)),
		Origin:        converter.String(d.Get("origin").(string)),
		OriginId:      converter.String(d.Get("origin_id").(string)),
		SubjectKind:   converter.String(d.Get("subject_kind").(string)),
		Domain:        converter.String(d.Get("domain").(string)),
		MailAddress:   converter.String(d.Get("mail").(string)),
		PrincipalName: converter.String(d.Get("principal_name").(string)),
		Description:   converter.String(d.Get("description").(string)),
	}

	dMembers := d.Get("members").(*schema.Set).List()

	members := make([]graph.GraphMembership, 0)
	for _, e := range dMembers {
		members = append(members, graph.GraphMembership{
			ContainerDescriptor: group.Descriptor,
			MemberDescriptor:    converter.String(e.(string)),
		})
	}
	return &group, &members, nil
}

func flattenGroup(d *schema.ResourceData, group *graph.GraphGroup, members *[]graph.GraphMembership) error {

	if group.Descriptor != nil {
		d.Set("descriptor", *group.Descriptor)
		d.SetId(*group.Descriptor)
	} else {
		return fmt.Errorf("Group Object does not contain a descriptor")
	}
	if group.DisplayName != nil {
		d.Set("display_name", *group.DisplayName)
	}
	if group.Url != nil {
		d.Set("url", *group.Url)
	}
	if group.Origin != nil {
		d.Set("origin", *group.Origin)
	}
	if group.OriginId != nil {
		d.Set("origin_id", *group.OriginId)
	}
	if group.SubjectKind != nil {
		d.Set("subject_kind", *group.SubjectKind)
	}
	if group.Domain != nil {
		d.Set("domain", *group.Domain)
	}
	if group.MailAddress != nil {
		d.Set("mail", *group.MailAddress)
	}
	if group.PrincipalName != nil {
		d.Set("principal_name", *group.PrincipalName)
	}
	if group.Description != nil {
		d.Set("description", *group.Description)
	}
	if members != nil {
		dMembers := make([]string, 0)
		for _, e := range *members {
			dMembers = append(dMembers, *e.MemberDescriptor)
		}
		d.Set("members", dMembers)
	}
	return nil
}

func groupReadMembers(groupDescriptor string, clients *config.AggregatedClient) (*[]graph.GraphMembership, error) {
	actualMembers, err := clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf("Error reading group memberships during read: %+v", err)
	}

	members := make([]graph.GraphMembership, len(*actualMembers))
	for i, membership := range *actualMembers {
		members[i] = *buildMembership(groupDescriptor, *membership.MemberDescriptor)
	}

	return &members, nil
}
