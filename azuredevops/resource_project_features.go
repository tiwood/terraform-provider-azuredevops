package azuredevops

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/featuremanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
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

var projectFeatureNameMap = map[string]ProjectFeatureType{
	"ms.vss-work.agile":           ProjectFeatureTypeValues.Boards,
	"ms.vss-code.version-control": ProjectFeatureTypeValues.Repositories,
	"ms.vss-build.pipelines":      ProjectFeatureTypeValues.Pipelines,
	"ms.vss-test-web.test":        ProjectFeatureTypeValues.TestPlans,
	"ms.feed.feed":                ProjectFeatureTypeValues.Artifacts,
}

var projectFeatureNameMapReverse = map[ProjectFeatureType]string{}

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"features": {
				Type:         schema.TypeMap,
				Required:     true,
				ValidateFunc: validateFeatures,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func validateFeatures(i interface{}, k string) ([]string, []error) {
	var errors []error
	var warnings []string

	m := i.(map[string]interface{})

	for feature, state := range m {
		lfeature := strings.ToLower(feature)
		if _, ok := getProjectFeatureID(ProjectFeatureType(lfeature)); !ok {
			errors = append(errors, fmt.Errorf("unknown feature: %s", feature))
		}

		if state != string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled) && state != string(featuremanagement.ContributedFeatureEnabledValueValues.Disabled) {
			errors = append(errors, fmt.Errorf("invalid state: %s", state))
		}
	}

	return warnings, errors
}

func getProjectFeatureIDs() *[]string {
	keys := make([]string, len(projectFeatureNameMap))
	for k := range projectFeatureNameMap {
		keys = append(keys, k)
	}
	return &keys
}

func getProjectFeatureID(fname ProjectFeatureType) (string, bool) {
	if len(projectFeatureNameMapReverse) <= 0 {
		for k, v := range projectFeatureNameMap {
			projectFeatureNameMapReverse[v] = k
		}
	}
	f, ok := projectFeatureNameMapReverse[fname]
	return f, ok
}

func resourceProjectFeaturesCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	projectID := d.Get("project_id").(string)
	featureStates := d.Get("features").(map[string]interface{})

	err := setProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID, featureStates)
	if err != nil {
		return err
	}

	d.SetId(projectID)
	return resourceProjectFeaturesRead(d, m)
}

func resourceProjectFeaturesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	projectID := d.Get("project_id").(string)

	currentFeatureStates, err := readProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID)
	if err != nil {
		return err
	}

	featureStates := d.Get("features").(map[string]interface{})
	for k := range *currentFeatureStates {
		if _, ok := featureStates[k]; !ok {
			delete(*currentFeatureStates, k)
		}
	}

	d.Set("features", *currentFeatureStates)
	return nil
}

func setProjectFeatureStates(ctx context.Context, fc featuremanagement.Client, projectID string, featureStates map[string]interface{}) error {
	for k, v := range featureStates {
		enabledValue := featuremanagement.ContributedFeatureEnabledValue(v.(string))
		f, ok := getProjectFeatureID(ProjectFeatureType(k))
		if !ok {
			return fmt.Errorf("unknown feature: %s", k)
		}
		_, err := fc.SetFeatureStateForScope(ctx, featuremanagement.SetFeatureStateForScopeArgs{
			Feature: &featuremanagement.ContributedFeatureState{
				FeatureId: converter.String(f),
				State:     &enabledValue,
				Scope: &featuremanagement.ContributedFeatureSettingScope{
					SettingScope: converter.String("project"),
					UserScoped:   converter.Bool(false),
				},
			},
			FeatureId:  converter.String(f),
			UserScope:  converter.String("host"),
			ScopeName:  converter.String("project"),
			ScopeValue: &projectID,
		})
		if nil != err {
			return err
		}
	}

	return nil
}

func readProjectFeatureStates(ctx context.Context, fc featuremanagement.Client, projectID string) (*map[string]string, error) {
	states, err := fc.QueryFeatureStates(ctx, featuremanagement.QueryFeatureStatesArgs{
		Query: &featuremanagement.ContributedFeatureStateQuery{
			FeatureIds: getProjectFeatureIDs(),
			ScopeValues: &map[string]string{
				"project": projectID,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	featureStates := make(map[string]string)
	for k := range projectFeatureNameMap {
		state, ok := (*states.FeatureStates)[k]
		if ok {
			featureStates[k] = string(*state.State)
		}
	}
	return &featureStates, nil
}

func resourceProjectFeaturesUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectFeaturesCreate(d, m)
}

func resourceProjectFeaturesDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	projectID := d.Get("project_id").(string)
	featureStates := d.Get("features").(map[string]interface{})
	for k := range featureStates {
		featureStates[k] = string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled)
	}
	err := setProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID, featureStates)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
