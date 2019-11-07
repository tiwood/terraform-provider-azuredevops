package azuredevops

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestReposDataSource_ListRepos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoName := "a-repo"
	projectName := "a-project"

	reposClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &aggregatedClient{
		GitReposClient: reposClient,
		ctx:            context.Background()}

	expectedArgs := git.GetRepositoriesArgs{
		Project: converter.String(projectName),
	}

	gitReposResult := []git.GitRepository{{Name: converter.String(repoName)}}
	reposClient.
		EXPECT().
		GetRepositories(clients.ctx, expectedArgs).
		Return(&gitReposResult, nil)

	resourceData := schema.TestResourceDataRaw(t, resourceAzureGitRepository().Schema, nil)
	resourceData.Set("project_id", projectName)

	repos, err := resourceAzureGitRepositoriesRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, repoName, *(*repos)[0].Name)
}
