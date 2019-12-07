package azuredevops

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

func dataGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"project_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Performs a lookup of a project group. This involves the following actions:
//	(1) Identify AzDO graph descriptor for the project in which the group exists
//	(2) Query for all AzDO groups that exist within the project. This leverages the AzDO graph descriptor for the project.
//		This involves querying a paginated API, so multiple API calls may be needed for this step.
//	(3) Select group that has the name identified by the schema
func dataSourceGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	groupName, projectID := d.Get("name").(string), d.Get("project_id").(string)

	projectDescriptor, err := getProjectDescriptor(clients, projectID)
	if err != nil {
		return fmt.Errorf("Error finding descriptor for project with ID %s. Error: %v", projectID, err)
	}

	projectGroups, err := getGroupsForDescriptor(clients, projectDescriptor)
	if err != nil {
		return fmt.Errorf("Error finding groups for project with ID %s. Error: %v", projectID, err)
	}

	targetGroup := selectGroup(projectGroups, groupName)
	if targetGroup == nil {
		return fmt.Errorf("Could not find group with name %s in project with ID %s", groupName, projectID)
	}

	d.SetId(*targetGroup.Descriptor)
	d.Set("descriptor", *targetGroup.Descriptor)
	return nil
}

func getProjectDescriptor(clients *config.AggregatedClient, projectID string) (string, error) {
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return "", err
	}

	descriptor, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{StorageKey: &projectUUID})
	if err != nil {
		return "", err
	}

	return *descriptor.Value, nil
}

func getGroupsForDescriptor(clients *config.AggregatedClient, projectDescriptor string) (*[]graph.GraphGroup, error) {
	var groups []graph.GraphGroup
	var currentToken string

	for hasMore := true; hasMore; {
		newGroups, latestToken, err := getGroupsWithContinuationToken(clients, projectDescriptor, currentToken)
		currentToken = latestToken
		if err != nil {
			return nil, err
		}

		groups = append(groups, *newGroups...)
		hasMore = currentToken != ""
	}

	return &groups, nil
}

func getGroupsWithContinuationToken(clients *config.AggregatedClient, projectDescriptor string, continuationToken string) (*[]graph.GraphGroup, string, error) {
	args := graph.ListGroupsArgs{ScopeDescriptor: &projectDescriptor}
	if continuationToken != "" {
		args.ContinuationToken = &continuationToken
	}

	response, err := clients.GraphClient.ListGroups(clients.Ctx, args)
	if err != nil {
		return nil, "", err
	}

	if response.ContinuationToken != nil && len(*response.ContinuationToken) > 1 {
		return nil, "", fmt.Errorf("Expected at most 1 continuation token, but found %d", len(*response.ContinuationToken))
	}

	var newToken string
	if response.ContinuationToken != nil && len(*response.ContinuationToken) > 0 {
		newToken = (*response.ContinuationToken)[0]
	}

	return response.GraphGroups, newToken, nil
}

func selectGroup(groups *[]graph.GraphGroup, groupName string) *graph.GraphGroup {
	for _, group := range *groups {
		if strings.EqualFold(*group.DisplayName, groupName) {
			return &group
		}
	}
	return nil
}
