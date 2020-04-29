package azuredevops

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/featuremanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
)

// ProjectFeatureType Project feature in Azure DevOps
type ProjectFeatureType string

type projectFeatureTypeValuesType struct {
	Boards       ProjectFeatureType
	Repositories ProjectFeatureType
	Pipelines    ProjectFeatureType
	TestPlans    ProjectFeatureType
	Artifacts    ProjectFeatureType
}

// ProjectFeatureTypeValues valid projects features in Azure DevOps
var ProjectFeatureTypeValues = projectFeatureTypeValuesType{
	Boards:       "boards",
	Repositories: "repositories",
	Pipelines:    "pipelines",
	TestPlans:    "testplans",
	Artifacts:    "artifacts",
}

var projectFeatureNameMap = map[ProjectFeatureType]string{
	ProjectFeatureTypeValues.Boards:       "ms.vss-work.agile",
	ProjectFeatureTypeValues.Repositories: "ms.vss-code.version-control",
	ProjectFeatureTypeValues.Pipelines:    "ms.vss-build.pipelines",
	ProjectFeatureTypeValues.TestPlans:    "ms.vss-test-web.test",
	ProjectFeatureTypeValues.Artifacts:    "ms.feed.feed",
}

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
			"project_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"features": {
				Type:     schema.TypeMap,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ProjectFeatureTypeValues.Boards),
					string(ProjectFeatureTypeValues.Repositories),
					string(ProjectFeatureTypeValues.Pipelines),
					string(ProjectFeatureTypeValues.TestPlans),
					string(ProjectFeatureTypeValues.Artifacts),
				}, false),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled),
						string(featuremanagement.ContributedFeatureEnabledValueValues.Disabled),
					}, false),
				},
			},
		},
	}
}

func resourceProjectFeaturesCreate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectFeaturesRead(d, m)
}

func resourceProjectFeaturesRead(d *schema.ResourceData, m interface{}) error {
	//clients := m.(*config.AggregatedClient)

	return errors.New("Not implemented")
}

func resourceProjectFeaturesUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectFeaturesRead(d, m)
}

func resourceProjectFeaturesDelete(d *schema.ResourceData, m interface{}) error {
	return errors.New("Not implemented")
}
