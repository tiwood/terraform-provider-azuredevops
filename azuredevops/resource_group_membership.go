package azuredevops

import (
	"fmt"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Update: resourceGroupMembershipUpdate,
		Delete: resourceGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				MinItems: 1,
				Required: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	memberships := expandGroupMembers(d.Get("group").(string), d.Get("members").(*schema.Set))

	err := addMembers(clients, memberships)
	if err != nil {
		return fmt.Errorf("Error adding group memberships during create: %+v", err)
	}

	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}

func resourceGroupMembershipUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("members") {
		return nil
	}

	group := d.Get("group").(string)
	oldData, newData := d.GetChange("members")
	// members that need to be added will be missing from the old data, but present in the new data
	membersToAdd := newData.(*schema.Set).Difference(oldData.(*schema.Set))
	// members that need to be removed will be missing from the new data, but present in the old data
	membersToRemove := oldData.(*schema.Set).Difference(newData.(*schema.Set))
	return applyMembershipUpdate(m.(*config.AggregatedClient),
		expandGroupMembers(group, membersToAdd),
		expandGroupMembers(group, membersToRemove))
}

func applyMembershipUpdate(clients *config.AggregatedClient, toAdd *[]graph.GraphMembership, toRemove *[]graph.GraphMembership) error {
	err := removeMembers(clients, toRemove)
	if err != nil {
		return fmt.Errorf("Error removing group memberships during update: %+v", err)
	}

	err = addMembers(clients, toAdd)
	if err != nil {
		return fmt.Errorf("Error adding group memberships during update: %+v", err)
	}

	return nil
}

func resourceGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	memberships := expandGroupMembers(d.Get("group").(string), d.Get("members").(*schema.Set))

	err := removeMembers(clients, memberships)
	if err != nil {
		return fmt.Errorf("Error removing group memberships during delete: %+v", err)
	}

	// this marks the resource as deleted
	d.SetId("")
	return nil
}

// Add members to a group using the AzDO REST API. If any error is encountered, the function immediately returns.
func addMembers(clients *config.AggregatedClient, memberships *[]graph.GraphMembership) error {
	if memberships != nil {
		for _, membership := range *memberships {
			_, err := clients.GraphClient.AddMembership(clients.Ctx, graph.AddMembershipArgs{
				SubjectDescriptor:   membership.MemberDescriptor,
				ContainerDescriptor: membership.ContainerDescriptor,
			})

			if err != nil {
				return fmt.Errorf("Error adding member %s to group %s: %+v",
					converter.ToString(membership.MemberDescriptor, "nil"),
					converter.ToString(membership.ContainerDescriptor, "nil"),
					err)
			}
		}
	}
	return nil
}

// Remove members from a group using the AzDO REST API. If any error is encountered, the function immediately returns.
func removeMembers(clients *config.AggregatedClient, memberships *[]graph.GraphMembership) error {
	if memberships != nil {
		for _, membership := range *memberships {
			err := clients.GraphClient.RemoveMembership(clients.Ctx, graph.RemoveMembershipArgs{
				SubjectDescriptor:   membership.MemberDescriptor,
				ContainerDescriptor: membership.ContainerDescriptor,
			})

			if err != nil {
				return fmt.Errorf("Error removing member from group: %+v", err)
			}
		}
	}
	return nil
}

func expandGroupMembers(group string, memberSet *schema.Set) *[]graph.GraphMembership {
	if memberSet == nil || memberSet.Len() <= 0 {
		return &[]graph.GraphMembership{}
	}
	members := memberSet.List()
	memberships := make([]graph.GraphMembership, len(members))

	for i, member := range members {
		memberships[i] = *buildMembership(group, member.(string))
	}

	return &memberships
}

func buildMembership(group string, member string) *graph.GraphMembership {
	return &graph.GraphMembership{
		ContainerDescriptor: &group,
		MemberDescriptor:    &member,
	}
}

func resourceGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	group := d.Get("group").(string)

	actualMemberships, err := clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &group,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
	if err != nil {
		return fmt.Errorf("Error reading group memberships during read: %+v", err)
	}

	stateMembers := d.Get("members").(*schema.Set)
	members := make([]string, 0)
	for _, membership := range *actualMemberships {
		if stateMembers.Contains(*membership.MemberDescriptor) {
			members = append(members, *membership.MemberDescriptor)
		}
	}

	d.Set("members", members)
	return nil
}
