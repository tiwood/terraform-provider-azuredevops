package azuredevops

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceAzureGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureGroupCreate,
		Read:   resourceAzureGroupRead,
		Update: resourceAzureGroupUpdate,
		Delete: resourceAzureGroupDelete,

		Schema: map[string]*schema.Schema{
			// add properties here
		},
	}
}

func resourceAzureGroupCreate(d *schema.ResourceData, m interface{}) error {

	return resourceAzureGroupRead(d, m)
}

func resourceAzureGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	return nil
}

func resourceAzureGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAzureGroupRead(d, m)
}

func resourceAzureGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
