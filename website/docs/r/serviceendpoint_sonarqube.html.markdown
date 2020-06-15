---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_sonarqube"
description: |-
  Manages a Sonarqube service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_sonarqube

Manages a Sonarqube service endpoint within Azure DevOps.

-> **NOTE:** The `Sonarqube` Azure DevOps extension is required for this resource.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_sonarqube" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  url                   = "https://my-sonarqube-instance.com"
  token                 = "my-sonarqube-token"
  service_endpoint_name = "test-sonarqube"
  description           = "test"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `url` - (Required) Sonarqube URL.
* `token` - (Required) Sonarqube user token.
* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

* [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)
* [Sonarqube Azure DevOps Extension](https://marketplace.visualstudio.com/items?itemName=SonarSource.sonarqube)

## Import

Azure DevOps Service Endpoint Sonarqube can be imported using the **projectID/serviceEndpointID**, e.g.

```shell
 terraform import azuredevops_serviceendpoint_sonarqube.serviceendpoint xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
