package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceAzureGitRepositoriesRead(d *schema.ResourceData, m interface{}) (*[]git.GitRepository, error) {
	projectID := d.Get("project_id").(string)
	clients := m.(*aggregatedClient)
	return clients.GitReposClient.GetRepositories(clients.ctx, git.GetRepositoriesArgs{
		Project: converter.String(projectID),
	})
}
