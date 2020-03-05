package azuredevops

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func dataRepos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceReposRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"weburl": {
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

func dataSourceReposRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	repoName, projectID := d.Get("name").(string), d.Get("project_id")

	projectDescriptor, err := getProjectDescriptor(clients, projectID.(string))
	if err != nil {
		return fmt.Errorf("Error finding descriptor for project with ID %s. Error: %v", projectID, err)
	}

	projectRepos, err := getReposForDescriptor(clients, projectDescriptor)
	if err != nil {
		return fmt.Errorf("Error finding repos for project with ID %s. Error: %v", projectID, err)
	}

	targetRepo := selectRepo(projectRepos, repoName)
	if targetRepo == nil {
		return fmt.Errorf("Could not find repo with name %s in project with ID %s", repoName, projectID)
	}

	d.SetId(*targetRepo.Name)
	d.Set("descriptor", *targetRepo.Name)
	return nil
}

func getReposForDescriptor(clients *aggregatedClient, projectDescriptor string) (*[]git.GitRepository, error) {
	var repos *[]git.GitRepository
	projectID := projectDescriptor
	repos, err := clients.GitReposClient.GetRepositories(clients.ctx, git.GetRepositoriesArgs{
		Project: converter.String(projectID),
	})

	return repos, err
}

func selectRepo(repos *[]git.GitRepository, repoName string) *git.GitRepository {

	return &(*repos)[0]
}

func resourceAzureGitRepositoriesRead(d *schema.ResourceData, m interface{}) (*[]git.GitRepository, error) {
	projectID := d.Get("project_id").(string)
	clients := m.(*aggregatedClient)
	return clients.GitReposClient.GetRepositories(clients.ctx, git.GetRepositoriesArgs{
		Project: converter.String(projectID),
	})
}
