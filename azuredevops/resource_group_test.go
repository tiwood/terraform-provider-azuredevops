package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

var descriptor = "vssgp.Uy0xLTktMTU1MTM3NDI0NS01OTMwNjE4OTktMTUzMjM2ODQ0OC0yNjEwNDc0OTEzLTIwMTI3MjY3MjgtMS00MTA1Mjg5ODQ0LTUxNzgwOTc0My0yNDc0MDIwNDA4LTI5NDAwMzQ4NTk"
var origin = "TEST_ORIGIN"
var originID = "5d466068-fe00-47c8-80d7-bb268165820c"
var displayName = "TEST_GROUP"
var description = "TEST_DESCRIPTION"
var url = "https://dev.azure.com/_test_organization"
var email = "test_group@test.local"
var subjectKind = "group"
var domain = "test.domain.local"
var principalName = "test@domain.local"

func init() {
	/* add code for test setup here */
}

func TestGroupResource_Create_TestHandleErrorVstsContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{
		CreationContext: &graph.GraphGroupVstsCreationContext{
			DisplayName: &displayName,
			Description: &description,
		},
	}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(nil, errors.New("CreateGroup() Failed")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("display_name", displayName)
	resourceData.Set("description", description)

	err := resourceGroupCreate(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CreateGroup() Failed")
}

func TestGroupResource_Create_TestHandleErrorMailContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{
		CreationContext: &graph.GraphGroupMailAddressCreationContext{
			MailAddress: &email,
		},
	}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(nil, errors.New("CreateGroup() Failed")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("mail", email)

	err := resourceGroupCreate(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CreateGroup() Failed")
}

func TestGroupResource_Create_TestHandleErrorOriginIdContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{
		CreationContext: &graph.GraphGroupOriginIdCreationContext{
			OriginId: &originID,
		},
	}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(nil, errors.New("CreateGroup() Failed")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("origin_id", originID)

	err := resourceGroupCreate(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CreateGroup() Failed")
}

func TestGroupResource_Create_TestVstsContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{
		CreationContext: &graph.GraphGroupVstsCreationContext{
			DisplayName: &displayName,
			Description: &description,
		},
	}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(&graph.GraphGroup{
			Descriptor:    &descriptor,
			DisplayName:   &displayName,
			Description:   &description,
			Origin:        &origin,
			OriginId:      &originID,
			MailAddress:   &email,
			Url:           &url,
			SubjectKind:   &subjectKind,
			Domain:        &domain,
			PrincipalName: &principalName,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("display_name", displayName)
	resourceData.Set("description", description)

	err := resourceGroupCreate(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, descriptor, resourceData.Id())
	require.Equal(t, descriptor, resourceData.Get("descriptor"))
	require.Equal(t, displayName, resourceData.Get("display_name"))
	require.Equal(t, description, resourceData.Get("description"))
	require.Equal(t, origin, resourceData.Get("origin"))
	require.Equal(t, originID, resourceData.Get("origin_id"))
	require.Equal(t, email, resourceData.Get("mail"))
	require.Equal(t, url, resourceData.Get("url"))
	require.Equal(t, subjectKind, resourceData.Get("subject_kind"))
	require.Equal(t, domain, resourceData.Get("domain"))
	require.Equal(t, principalName, resourceData.Get("principal_name"))
}

func TestGroupResource_Create_TestMailContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{
		CreationContext: &graph.GraphGroupMailAddressCreationContext{
			MailAddress: &email,
		},
	}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(&graph.GraphGroup{
			Descriptor:    &descriptor,
			DisplayName:   &displayName,
			Description:   &description,
			Origin:        &origin,
			OriginId:      &originID,
			MailAddress:   &email,
			Url:           &url,
			SubjectKind:   &subjectKind,
			Domain:        &domain,
			PrincipalName: &principalName,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("mail", email)

	err := resourceGroupCreate(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, descriptor, resourceData.Id())
	require.Equal(t, descriptor, resourceData.Get("descriptor"))
	require.Equal(t, displayName, resourceData.Get("display_name"))
	require.Equal(t, description, resourceData.Get("description"))
	require.Equal(t, origin, resourceData.Get("origin"))
	require.Equal(t, originID, resourceData.Get("origin_id"))
	require.Equal(t, email, resourceData.Get("mail"))
	require.Equal(t, url, resourceData.Get("url"))
	require.Equal(t, subjectKind, resourceData.Get("subject_kind"))
	require.Equal(t, domain, resourceData.Get("domain"))
	require.Equal(t, principalName, resourceData.Get("principal_name"))
}

func TestGroupResource_Create_TestOriginIdContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{
		CreationContext: &graph.GraphGroupOriginIdCreationContext{
			OriginId: &originID,
		},
	}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(&graph.GraphGroup{
			Descriptor:    &descriptor,
			DisplayName:   &displayName,
			Description:   &description,
			Origin:        &origin,
			OriginId:      &originID,
			MailAddress:   &email,
			Url:           &url,
			SubjectKind:   &subjectKind,
			Domain:        &domain,
			PrincipalName: &principalName,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("origin_id", originID)

	err := resourceGroupCreate(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, descriptor, resourceData.Id())
	require.Equal(t, descriptor, resourceData.Get("descriptor"))
	require.Equal(t, displayName, resourceData.Get("display_name"))
	require.Equal(t, description, resourceData.Get("description"))
	require.Equal(t, origin, resourceData.Get("origin"))
	require.Equal(t, originID, resourceData.Get("origin_id"))
	require.Equal(t, email, resourceData.Get("mail"))
	require.Equal(t, url, resourceData.Get("url"))
	require.Equal(t, subjectKind, resourceData.Get("subject_kind"))
	require.Equal(t, domain, resourceData.Get("domain"))
	require.Equal(t, principalName, resourceData.Get("principal_name"))
}

func TestGroupResource_Create_TestParameterCollisions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedCreateGroupArgs := graph.CreateGroupArgs{}

	graphClient.
		EXPECT().
		CreateGroup(clients.Ctx, expectedCreateGroupArgs).
		Return(nil, errors.New("CreateGroup() INVALID CALL")).
		Times(0)

	var resourceData *schema.ResourceData
	var err error

	resourceData = schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("mail", email)
	resourceData.Set("origin_id", originID)

	err = resourceGroupCreate(resourceData, clients)
	require.NotNil(t, err)

	resourceData = schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("display_name", displayName)
	resourceData.Set("origin_id", originID)

	err = resourceGroupCreate(resourceData, clients)
	require.NotNil(t, err)

	resourceData = schema.TestResourceDataRaw(t, resourceGroup().Schema, nil)
	resourceData.Set("display_name", displayName)
	resourceData.Set("mail", originID)

	err = resourceGroupCreate(resourceData, clients)
	require.NotNil(t, err)
}
