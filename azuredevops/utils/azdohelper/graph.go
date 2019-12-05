package azdohelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
)

// AzDOGraphCreateGroupArgs Arguments for the AzDOGraphCreateGroup function
type AzDOGraphCreateGroupArgs struct {
	// (required) The subset of the full graph group used to uniquely find the graph subject in an external provider.
	CreationContext interface{}
	// (optional) A descriptor referencing the scope (collection, project) in which the group should be created. If omitted, will be created in the scope of the enclosing account or organization. Valid only for VSTS groups.
	ScopeDescriptor *string
	// (optional) A comma separated list of descriptors referencing groups you want the graph group to join
	GroupDescriptors *[]string
}

// AzDOGraphCreateGroup This function is temporary fix, as long as the Azure DevOps GO API can't handle different group creation args properly
func AzDOGraphCreateGroup(ctx context.Context, client graph.Client, args AzDOGraphCreateGroupArgs) (*graph.GraphGroup, error) {
	if args.CreationContext == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CreationContext"}
	}
	queryParams := url.Values{}
	if args.ScopeDescriptor != nil {
		queryParams.Add("scopeDescriptor", *args.ScopeDescriptor)
	}
	if args.GroupDescriptors != nil {
		listAsString := strings.Join((*args.GroupDescriptors)[:], ",")
		queryParams.Add("groupDescriptors", listAsString)
	}

	t := reflect.TypeOf(args.CreationContext)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t != reflect.TypeOf((*graph.GraphUserOriginIdCreationContext)(nil)).Elem() &&
		t != reflect.TypeOf((*graph.GraphUserMailAddressCreationContext)(nil)).Elem() &&
		t != reflect.TypeOf((*graph.GraphUserMailAddressCreationContext)(nil)).Elem() {
		return nil, fmt.Errorf("Unsupported user creation context: %T", t)
	}

	body, marshalErr := json.Marshal(args.CreationContext)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationID, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	clientImpl := client.(*graph.ClientImpl)
	resp, err := clientImpl.Client.Send(ctx, http.MethodPost, locationID, "5.1-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue graph.GraphGroup
	err = clientImpl.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
