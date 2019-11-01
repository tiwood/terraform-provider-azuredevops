package azuredevops

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

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
			// scope (optional)
			"scope": {
				Type:     schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional: true,
				ForceNew: true,
			},

			// ***
			// One of 
			//     origin_id
			//     mail
			//     displayName
			// must be specified
			// **

			// origin_id
			"origin_id": {
				Type:     schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional: true,
				ForceNew: true,
			},
			// mail
			"mail": {
				Type:     schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional: true,
				ForceNew: true,
			},			
			// displayName (required if )
			"displayName": {
				Type:     schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional: true,
				ForceNew: true,
			},

			// cross_project (optional, default: false)
			"cross_project" : {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,	
			},

			// restricted_visibility (optional, default: false)
			"restricted_visibility" : {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,	
			},

			// description (optional)
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			
			// members (optional, list)
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
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
