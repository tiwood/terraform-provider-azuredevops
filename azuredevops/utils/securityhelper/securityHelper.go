package securityhelper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func SetPrincipalPermissions(d *schema.ResourceData, sn *securityNamespace, aclToken *string, forcePermission *PermissionType, forceReplace bool) error {
	principal, ok := d.GetOk("principal")
	if !ok {
		return fmt.Errorf("Failed to get 'principal' from schema")
	}

	permissions, ok := d.GetOk("permissions")
	if !ok {
		return fmt.Errorf("Failed to get 'permissions' from schema")
	}

	bReplace := d.Get("replace")
	if forceReplace {
		bReplace = forceReplace
	}
	permissionMap := make(map[ActionName]PermissionType, len(permissions.(map[string]interface{})))
	for key, elem := range permissions.(map[string]interface{}) {
		if forcePermission != nil {
			permissionMap[ActionName(key)] = *forcePermission
		} else {
			permissionMap[ActionName(key)] = PermissionType(elem.(string))
		}
	}
	setPermissions := []SetPrincipalPermission{
		SetPrincipalPermission{
			Replace: bReplace.(bool),
			PrincipalPermission: PrincipalPermission{
				SubjectDescriptor: principal.(string),
				Permissions:       permissionMap,
			},
		}}

	return sn.SetPrincipalPermissions(&setPermissions, aclToken)
}

func GetPrincipalPermissions(d *schema.ResourceData, sn *securityNamespace, aclToken *string) (*PrincipalPermission, error) {
	principal, ok := d.GetOk("principal")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'principal' from schema")
	}

	permissions, ok := d.GetOk("permissions")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'permissions' from schema")
	}

	principalList := []string{*converter.StringFromInterface(principal)}
	principalPermissions, err := sn.GetPrincipalPermissions(aclToken, &principalList)
	if err != nil {
		return nil, err
	}
	if principalPermissions == nil || len(*principalPermissions) != 1 {
		return nil, fmt.Errorf("Failed to retrive current permissions for principal [%s]", principalList[0])
	}
	d.SetId(fmt.Sprintf("%s/%s", *aclToken, principal.(string)))
	for key := range ((*principalPermissions)[0]).Permissions {
		if _, ok := permissions.(map[string]interface{})[string(key)]; !ok {
			delete(((*principalPermissions)[0]).Permissions, key)
		}
	}
	return &(*principalPermissions)[0], nil
}
