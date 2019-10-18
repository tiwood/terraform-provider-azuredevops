package azuredevops

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,

		//https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"principal_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ad", "aad", "vsts", "msa", "imp"}, true),
			},
			"origin_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
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
	clients := m.(*aggregatedClient)
	users := make([]map[string]interface{}, 0)

	principalName := d.Get("principal_name").(string)
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)

	subjectTypes := strings.Split(origin, ",")
	var continuationToken string

	// continously query the API until no ContinuationToken is in the response
	for {
		listArgs := graph.ListUsersArgs{
			SubjectTypes:      &subjectTypes,
			ContinuationToken: &continuationToken,
		}

		resp, err := clients.GraphClient.ListUsers(clients.ctx, listArgs)
		if err != nil {
			return fmt.Errorf("Error listing users: %q", err)
		}

		for _, user := range *resp.GraphUsers {

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

			if principalName != "" {
				if !strings.EqualFold(principalName, s["principal_name"].(string)) {
					continue
				}
			}

			if originID != "" {
				if !strings.EqualFold(originID, s["origin_id"].(string)) {
					continue
				}
			}

			users = append(users, s)
		}

		if resp.ContinuationToken == nil || (*resp.ContinuationToken)[0] == "" {
			// done listing users, break out of the for loop
			break
		}

		continuationToken = (*resp.ContinuationToken)[0]
	}

	d.SetId("users") //TODO: Maybe add a guid to the ID..
	if err := d.Set("users", users); err != nil {
		return fmt.Errorf("Error setting `users`: %+v", err)
	}

	return nil
}
