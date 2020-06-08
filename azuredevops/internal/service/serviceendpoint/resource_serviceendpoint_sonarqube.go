package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointSonarqube schema and implementation for Sonarqube service endpoint resource
func ResourceServiceEndpointSonarqube() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointSonarqube, expandServiceEndpointSonarqube)
	makeUnprotectedSchema(r, "url", "AZDO_SONARQUBE_SERVICE_CONNECTION_URL", "The URL of the Sonarqube instance.")
	makeProtectedSchema(r, "token", "AZDO_SONARQUBE_SERVICE_CONNECTION_TOKEN", "The Sonarqube user token which should be used.")
	return r
}

func expandServiceEndpointSonarqube(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("token").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("sonarqube")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointSonarqube(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	tfhelper.HelpFlattenSecret(d, "token")
	d.Set("token", (*serviceEndpoint.Authorization.Parameters)["username"])
	d.Set("url", (*serviceEndpoint.Url))
}
