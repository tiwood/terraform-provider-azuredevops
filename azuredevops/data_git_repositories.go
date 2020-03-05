package azuredevops

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func dataGitRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validate.UUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"repository_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"include_hidden": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"repositories": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getGitRepositoryHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssh_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getGitRepositoryHash(v interface{}) int {
	return hashcode.String(v.(map[string]interface{})["id"].(string))
}

func dataSourceGitRepositoriesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	repoName, projectID := d.Get("name").(string), d.Get("project_id").(string)
	includeHidden := d.Get("include_hidden").(bool)

	projectRepos, err := getGitRepositoriesByNameAndProject(clients, repoName, projectID, includeHidden)
	if err != nil {
		return fmt.Errorf("Error finding repos for project with ID %s. Error: %v", projectID, err)
	}
	log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Read [%d] Git repositories", len(*projectRepos))

	results, err := flattenGitRepositories(projectRepos)
	if err != nil {
		return fmt.Errorf("Error flattening projects. Error: %v", err)
	}

	h := sha1.New()
	repoNames, err := getAttributeValues(results, "name")
	if err != nil {
		return fmt.Errorf("Failed to get list of repository names: %v", err)
	}
	if len(repoNames) <= 0 && repoName != "" {
		repoNames = append(repoNames, repoName)
	}
	if _, err := h.Write([]byte(strings.Join(repoNames, "-"))); err != nil {
		return fmt.Errorf("Unable to compute hash for Git repository names: %v", err)
	}
	d.SetId("gitRepos#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))
	err = d.Set("repositories", results)
	if err != nil {
		return err
	}
	return nil
}

func flattenGitRepositories(repos *[]git.GitRepository) ([]interface{}, error) {
	if repos == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)

	for _, element := range *repos {
		output := make(map[string]interface{})
		if element.Name != nil {
			output["name"] = *element.Name
		}

		if element.Id != nil {
			output["id"] = element.Id.String()
		}

		if element.Url != nil {
			output["url"] = *element.Url
		}

		if element.RemoteUrl != nil {
			output["remote_url"] = *element.RemoteUrl
		}

		if element.SshUrl != nil {
			output["ssh_url"] = *element.SshUrl
		}

		if element.WebUrl != nil {
			output["web_url"] = *element.WebUrl
		}

		if element.Project != nil && element.Project.Id != nil {
			output["project_id"] = *element.Project.Id
		}

		if element.Size != nil {
			output["project_id"] = *element.Size
		}

		results = append(results, output)
	}

	return results, nil
}

func getGitRepositoriesByNameAndProject(clients *config.AggregatedClient, name string, projectID string, includeHidden bool) (*[]git.GitRepository, error) {
	var repos *[]git.GitRepository
	var err error
	if name != "" && projectID != "" {
		repo, err := gitRepositoryRead(clients, "", name, projectID)
		if err != nil {
			return nil, err
		}
		repos = &[]git.GitRepository{*repo}
	} else {
		repos, err = clients.GitReposClient.GetRepositories(clients.Ctx, git.GetRepositoriesArgs{
			Project:       converter.String(projectID),
			IncludeHidden: converter.Bool(includeHidden),
		})
		if err != nil {
			return nil, err
		}
		if name != "" {
			for _, repo := range *repos {
				if strings.EqualFold(*repo.Name, name) {
					repos = &[]git.GitRepository{repo}
					break
				}
			}
		}
	}
	return repos, nil
}
