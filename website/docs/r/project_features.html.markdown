---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_features"
description: |-
  Manages features for Azure DevOps projects.
---

# azuredevops_group
Manages features for Azure DevOps projects

## Example Usage

```hcl
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "tf-project-test-001" {
  project_name = "Test Project"
}


resource "azuredevops_project_features" "my-project-features" {
       project_id = azuredevops_project.tf-project-test-001.id
       features = {
            Test = disabled
            Artifacts = enabled
            Board = enabled
       }
}

```

## Argument Reference

The following arguments are supported:

* `projectd_id` - (Required) 
* `features` - (Required) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


## Relevant Links


## Import
Azure DevOps feature settings can be imported using the project id, e.g.

```
terraform import azuredevops_project_features.project_id 2785562e-8f45-4534-a10e-b9ca1666b17e
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage
