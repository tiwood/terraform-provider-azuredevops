package azuredevops

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func resourceProjectFeatures() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectFeaturesCreate,
		Read:   resourceProjectFeaturesRead,
		Update: resourceProjectFeaturesUpdate,
		Delete: resourceProjectFeaturesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// add properties here
		},
	}
}

func resourceProjectFeaturesCreate(d *schema.ResourceData, m interface{}) error {

	return resourceProjectFeaturesRead(d, m)
}

func resourceProjectFeaturesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	return nil
}

func resourceProjectFeaturesUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectFeaturesRead(d, m)
}

func resourceProjectFeaturesDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
