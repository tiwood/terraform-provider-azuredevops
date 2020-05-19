---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_bitbucket"
description: |-
  Manages a Bitbucket service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_bitbucket
Manages a Bitbucket service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_bitbucket" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  username              = "xxxx"
  password              = "xxxx"
  service_endpoint_name = "test-bitbucket"
  description           = "test"
}
```

## Argument Reference

The following arguments are supported: 

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `username` - (Required) Bitbucket account username.
* `password` - (Required) Bitbucket account password.
* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)

## Import
Azure DevOps Service Endpoint Bitbucket can be imported using the **projectID/serviceEndpointID**, e.g.

```
 terraform import azuredevops_serviceendpoint_bitbucket.serviceendpoint xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```