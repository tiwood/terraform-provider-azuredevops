package azuredevops

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func resourceServiceEndpointPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceEndpointPermissionsCreate,
		Read:   resourceServiceEndpointPermissionsRead,
		Update: resourceServiceEndpointPermissionsUpdate,
		Delete: resourceServiceEndpointPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServiceEndpointPermissionsImporter,
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Required:     true,
				ForceNew:     true,
			},
			"service_endpoint_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Optional:     true,
				ForceNew:     true,
			},
		}),
	}
}

func createServiceEndpointToken(clients *config.AggregatedClient, d *schema.ResourceData) (*string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'project_id' from schema")
	}

	/*
	 * Token format
	 * ACL for ALL Service Endpoints in a project:	endpoints/#ProjectID#
	 * ACL for a Service Endpoint in a project:			endpoints/#ProjectID#/#ServiceEndpointID#
	 */
	aclToken := "endpoints/" + projectID.(string)
	serviceEndpointID, serviceEndpointOk := d.GetOkExists("service_endpoint_id")
	if serviceEndpointOk {
		aclToken += "/" + serviceEndpointID.(string)
	}
	return &aclToken, nil
}

func resourceServiceEndpointPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.ServiceEndpoints,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createServiceEndpointToken(clients, d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false)
	if err != nil {
		return err
	}

	return resourceServiceEndpointPermissionsRead(d, m)
}

func resourceServiceEndpointPermissionsRead(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createServiceEndpointToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.ServiceEndpoints,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn, aclToken)
	if err != nil {
		return err
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceServiceEndpointPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	return resourceServiceEndpointPermissionsCreate(d, m)
}

func resourceServiceEndpointPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createServiceEndpointToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.ServiceEndpoints,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, &securityhelper.PermissionTypeValues.NotSet, true)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceServiceEndpointPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	debugWait()

	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchName#/#SubjectDescriptor#
	return nil, errors.New("resourceServiceEndpointPermissionsImporter: Not implemented")
}
