// +build all resource_serviceendpoint_github
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointSonarqube_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_sonarqube"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t,
				&[]string{"AZDO_SONARQUBE_SERVICE_CONNECTION_URL", "AZDO_SONARQUBE_SERVICE_CONNECTION_TOKEN"})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointSonarqubeResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "url"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "token_hash"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				Config: testutils.HclServiceEndpointSonarqubeResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "url"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "token_hash"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token_hash"},
			},
		},
	})
}
