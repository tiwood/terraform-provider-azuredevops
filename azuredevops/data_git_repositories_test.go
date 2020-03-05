package azuredevops

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestReposDataSource_ReadsCorrectRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectName := "a-project"
	repoName := "a-repo"

	projectID := uuid.New()
	resourceData := schema.TestResourceDataRaw(t, dataRepos().Schema, nil)
	resourceData.Set("project_id", projectID.String())
	resourceData.Set("name", repoName)

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	repoClient := azdosdkmocks.NewMockGitClient(ctrl)

	clients := &aggregatedClient{
		GitReposClient: repoClient,
		GraphClient:    graphClient,
		ctx:            context.Background()}

	expectedProjectDescriptorLookupArgs := graph.GetDescriptorArgs{StorageKey: &projectID}
	projectDescriptor := converter.String(projectName)
	projectDescriptorResponse := graph.GraphDescriptorResult{Value: projectDescriptor}
	graphClient.
		EXPECT().
		GetDescriptor(clients.ctx, expectedProjectDescriptorLookupArgs).
		Return(&projectDescriptorResponse, nil)

	expectedArgs := git.GetRepositoriesArgs{
		Project: converter.String(projectName),
	}

	gitReposResult := []git.GitRepository{{Name: converter.String(repoName)}}
	repoClient.
		EXPECT().
		GetRepositories(clients.ctx, expectedArgs).
		Return(&gitReposResult, nil)

	err := dataSourceReposRead(resourceData, clients)
	require.Nil(t, err)
	//require.Equal(t, repoName, *(*repos)[0].Name)

}

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

	resourceData := schema.TestResourceDataRaw(t, resourceGitRepository().Schema, nil)
	resourceData.Set("project_id", projectName)

	repos, err := resourceGitRepositoriesRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, repoName, *(*repos)[0].Name)
}
