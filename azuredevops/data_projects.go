package azuredevops

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func dataProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectsRead,

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"state": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  core.ProjectStateValues.All,
				ValidateFunc: validation.StringInSlice([]string{
					string(core.ProjectStateValues.Deleting),
					string(core.ProjectStateValues.New),
					string(core.ProjectStateValues.WellFormed),
					string(core.ProjectStateValues.CreatePending),
					string(core.ProjectStateValues.All),
					string(core.ProjectStateValues.Unchanged),
					string(core.ProjectStateValues.Deleted),
				}, true),
			},

			"projects": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getProjectHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getProjectHash(v interface{}) int {
	return hashcode.String(v.(map[string]interface{})["project_id"].(string))
}

func dataSourceProjectsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	state := d.Get("state").(string)
	name := d.Get("project_name").(string)

	projects, err := getProjectsForStateAndName(clients, state, name)
	if err != nil {
		return fmt.Errorf("Error finding projects with state %s. Error: %v", state, err)
	}
	log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Read [%d] projects from current organization", len(projects))

	results, err := flattenProjectReferences(&projects)
	if err != nil {
		return fmt.Errorf("Error flattening projects. Error: %v", err)
	}

	h := sha1.New()
	projectNames, err := getAttributeValues(results, "name")
	if err != nil {
		return fmt.Errorf("Failed to get list of project names: %v", err)
	}
	if len(projectNames) <= 0 && name != "" {
		projectNames = append(projectNames, name)
	}
	if _, err := h.Write([]byte(state + strings.Join(projectNames, "-"))); err != nil {
		return fmt.Errorf("Unable to compute hash for project names: %v", err)
	}
	d.SetId("projects#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))
	err = d.Set("projects", results)
	if err != nil {
		return err
	}
	return nil
}

func getAttributeValues(items []interface{}, attributeName string) ([]string, error) {
	var result []string
	for _, element := range items {
		result = append(result, element.(map[string]interface{})[attributeName].(string))
	}
	return result, nil
}

func flattenProjectReferences(input *[]core.TeamProjectReference) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)

	for _, element := range *input {
		output := make(map[string]interface{})
		if element.Name != nil {
			output["name"] = *element.Name
		}

		if element.Id != nil {
			output["project_id"] = element.Id.String()
		}

		if element.Url != nil {
			output["project_url"] = *element.Url
		}

		if element.State != nil {
			output["state"] = string(*element.State)
		}

		results = append(results, output)
	}

	return results, nil
}

func getProjectsForStateAndName(clients *config.AggregatedClient, projectState string, projectName string) ([]core.TeamProjectReference, error) {
	var projects []core.TeamProjectReference
	var currentToken string

	for hasMore := true; hasMore; {
		newProjects, latestToken, err := getProjectsWithContinuationToken(clients, projectState, currentToken)
		currentToken = latestToken
		if err != nil {
			return nil, err
		}
		log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Received [%d] projects; Continuation token [%s]", len(newProjects), currentToken)

		if projectName != "" {
			log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Searching for project name [%s]", projectName)
			for _, project := range newProjects {
				if strings.EqualFold(*project.Name, projectName) {
					log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Found project [%s] in current project list", projectName)
					return []core.TeamProjectReference{project}, nil
				}
			}
		} else {
			projects = append(projects, newProjects...)
			log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Appended new projects to current project list (Length: %d)", len(projects))
		}
		hasMore = currentToken != ""
	}

	return projects, nil
}

func getProjectsWithContinuationToken(clients *config.AggregatedClient, projectState string, continuationToken string) ([]core.TeamProjectReference, string, error) {
	state := core.ProjectState(projectState)
	args := core.GetProjectsArgs{
		StateFilter: &state,
	}
	if continuationToken != "" {
		args.ContinuationToken = &continuationToken
	}

	response, err := clients.CoreClient.GetProjects(clients.Ctx, args)
	if err != nil {
		return nil, "", err
	}

	return response.Value, response.ContinuationToken, nil
}
