package azuredevops

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

func dataUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,

		//https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"principal_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"origin", "origin_id"},
			},
			"subject_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"origin": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"principal_name"},
			},
			"origin_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"principal_name"},
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getUserHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"descriptor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mail_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataUserRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	users := make([]interface{}, 0)

	subjectTypes := convertToStringSlice(d.Get("subject_types").(*schema.Set).List())
	principalName := d.Get("principal_name").(string)
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)
	comp := []AttributeComparison{}
	if principalName != "" {
		comp = append(comp, AttributeComparison{
			Name:       "PrincipalName",
			Value:      principalName,
			IgnoreCase: true,
			AllowNil:   false,
		})
	}
	if origin != "" {
		comp = append(comp, AttributeComparison{
			Name:       "Origin",
			Value:      origin,
			IgnoreCase: true,
			AllowNil:   false,
		})
	}
	if originID != "" {
		comp = append(comp, AttributeComparison{
			Name:       "OriginId",
			Value:      originID,
			IgnoreCase: true,
			AllowNil:   false,
		})
	}

	var currentToken string
	for hasMore := true; hasMore; {
		newUsers, latestToken, err := getUsersWithContinuationToken(clients, subjectTypes, currentToken)
		currentToken = latestToken
		hasMore = currentToken != ""
		if err != nil {
			return err
		}

		newUsers, err = filterUsersByAttributeValues(&newUsers, &comp)
		if err != nil {
			return err
		}
		fusers, err := flattenUsers(&newUsers)
		if err != nil {
			return err
		}
		users = append(users, fusers...)
	}

	descriptors, err := getValueSlice(&users, "descriptor")
	if err != nil {
		return err
	}
	h := sha1.New()
	if _, err := h.Write([]byte(strings.Join(descriptors, "-"))); err != nil {
		return fmt.Errorf("Unable to compute hash for user descriptors: %v", err)
	}
	d.SetId("users#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))
	if err := d.Set("users", users); err != nil {
		return fmt.Errorf("Error setting `users`: %+v", err)
	}

	return nil
}

func convertToStringSlice(input []interface{}) []string {
	result := make([]string, len(input))
	for i, k := range input {
		result[i] = k.(string)
	}

	return result
}

func getValueSlice(input *[]interface{}, attributeName string) ([]string, error) {
	if input == nil {
		return []string{}, nil
	}

	output := make([]string, len(*input))
	for i, user := range *input {
		s := reflect.ValueOf(user)
		if s.Kind() == reflect.Ptr {
			s = s.Elem()
		}
		if s.Kind() == reflect.Struct {
			f := s.FieldByName(attributeName)
			v := reflect.Indirect(f).Interface().(string)
			output = append(output, v)
		} else if s.Kind() == reflect.Map {
			ifc := s.Interface()
			imap := ifc.(map[string]interface{})
			output[i] = imap[attributeName].(string)
		} else {
			panic("Unsupported type")
		}
	}
	return output, nil
}

// AttributeComparison defines a comparison on an (struct) attribute
type AttributeComparison struct {
	Name       string
	Value      string
	IgnoreCase bool
	AllowNil   bool
}

func filterUsersByAttributeValues(input *[]graph.GraphUser, comparison *[]AttributeComparison) ([]graph.GraphUser, error) {
	if input == nil {
		return []graph.GraphUser{}, nil
	}
	if comparison == nil || len(*comparison) <= 0 {
		return *input, nil
	}

	output := []graph.GraphUser{}
	for _, user := range *input {
		b := true
		s := reflect.ValueOf(user)
		for _, comp := range *comparison {
			f := s.FieldByName(comp.Name)
			if f.Kind() == reflect.Ptr && f.IsNil() {
				if comp.AllowNil {
					continue
				} else {
					b = false
					break
				}
			}
			v := reflect.Indirect(f).Interface().(string)
			if comp.IgnoreCase {
				if !strings.EqualFold(comp.Value, v) {
					b = false
					break
				}
			} else {
				if comp.Value != v {
					b = false
					break
				}
			}
		}
		if b {
			output = append(output, user)
		}
	}
	return output, nil
}

func getUserHash(v interface{}) int {
	return hashcode.String(v.(map[string]interface{})["descriptor"].(string))
}

func flattenUsers(input *[]graph.GraphUser) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	results := make([]interface{}, len(*input))
	for i, element := range *input {
		output, err := flattenUser(&element)
		if err != nil {
			return nil, err
		}
		results[i] = output
	}
	return results, nil
}

func flattenUser(user *graph.GraphUser) (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if v := user.Descriptor; v != nil {
		s["descriptor"] = *v
	}
	if v := user.PrincipalName; v != nil {
		s["principal_name"] = *v
	}
	if v := user.Origin; v != nil {
		s["origin"] = *v
	}
	if v := user.OriginId; v != nil {
		s["origin_id"] = *v
	}
	if v := user.DisplayName; v != nil {
		s["display_name"] = *v
	}
	if v := user.MailAddress; v != nil {
		s["mail_address"] = *v
	}

	return s, nil
}

func getUsersWithContinuationToken(clients *config.AggregatedClient, subjectTypes []string, continuationToken string) ([]graph.GraphUser, string, error) {
	args := graph.ListUsersArgs{
		SubjectTypes: &subjectTypes,
	}
	if continuationToken != "" {
		args.ContinuationToken = &continuationToken
	}
	response, err := clients.GraphClient.ListUsers(clients.Ctx, args)
	if err != nil {
		return nil, "", fmt.Errorf("Error listing users: %q", err)
	}

	continuationToken = ""
	if response.ContinuationToken != nil && (*response.ContinuationToken)[0] != "" {
		continuationToken = (*response.ContinuationToken)[0]
	}

	return *response.GraphUsers, continuationToken, nil
}
