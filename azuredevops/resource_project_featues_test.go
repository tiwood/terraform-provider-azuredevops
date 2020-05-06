// +build all core resource_project_features

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"testing"

	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"

	"github.com/golang/mock/gomock"
)

func init() {
	/* add code for test setup here */
}

/**
 * Begin unit tests
 */

func TestProjectFeatures_Read_TestDontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &config.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	/* start writing test here */
}
