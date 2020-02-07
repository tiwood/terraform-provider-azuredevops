# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
provider "azuredevops" {
  version = ">= 0.0.1"
}


// This section creates a project
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}


// This section assigns users from AAD into a pre-existing group in AzDO
data "azuredevops_group" "group" {
  project_id = azuredevops_project.project.id
  name       = "Build Administrators"
}

resource "azuredevops_user_entitlement" "users" {
  for_each             = toset(var.aad_users)
  principal_name       = "${each.value}"
  account_license_type = "stakeholder"
}

resource "azuredevops_group_membership" "membership" {
  group   = data.azuredevops_group.group.descriptor
  members = values(azuredevops_user_entitlement.users)[*].descriptor
}



// This section configures variable groups and a build definition
resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"
  path       = "\\ExampleFolder"

  repository {
    repo_type   = "TfsGit"
    repo_name   = azuredevops_azure_git_repository.repository.name
    branch_name = azuredevops_azure_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }

  # https://github.com/microsoft/terraform-provider-azuredevops/issues/171
  # variables_groups = [azuredevops_variable_group.vg.id]
}

// This section configures an Azure DevOps Variable Group
# https://github.com/microsoft/terraform-provider-azuredevops/issues/170
resource "azuredevops_variable_group" "vg" {
  project_id   = azuredevops_project.project.id
  name         = "Sample VG 1"
  description  = "A sample variable group."
  allow_access = true

  variable {
    name      = "key1"
    value     = "value1"
    is_secret = true
  }

  variable {
    name      = "key2"
    value     = "value2"
  }

  variable {
    name      = "key3"
  }
}

// This section configures an Azure DevOps Git Repository with branch policies
resource "azuredevops_azure_git_repository" "repository" {
  project_id = azuredevops_project.project.id
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

// Configuration of AzureRm service end point
resource "azuredevops_serviceendpoint_azurerm" "endpoint1" {
  project_id                = azuredevops_project.project.id
  service_endpoint_name     = "TestServiceAzureRM"
  azurerm_spn_clientid      = "ee7f75a0-8553-4e6a-xxxx-xxxxxxxx"
  azurerm_spn_clientsecret  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  azurerm_spn_tenantid      = "2e3a33f9-66b1-4xxx-xxxx-xxxxxxxxx"
  azurerm_subscription_id   = "8a7aace5-xxxx-xxxx-xxxx-xxxxxxxxxx"
  azurerm_subscription_name = "Microsoft Azure DEMO"
  azurerm_scope             = "/subscriptions/1da42ac9-xxxx-xxxxx-xxxx-xxxxxxxxxxx"
}

resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = "Sample Kubernetes"
  apiserver_url          = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type = "AzureSubscription"
  
  azure_subscription {
    subscription_id = "8a7aace5-xxxx-xxxx-xxxx-xxxxxxxxxx"
    subscription_name = "Microsoft Azure DEMO"
    tenant_id = "2e3a33f9-66b1-4xxx-xxxx-xxxxxxxxx"
    resourcegroup_id = "sample-rg"
    namespace = "default"
    cluster_name = "sample-aks"
  }
}

resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = "Sample Kubernetes"
  apiserver_url          = "https://sample-aks.hcp.westeurope.azmk8s.io"
  authorization_type = "Kubeconfig"
  
  kubeconfig {
    kube_config = <<EOT
                apiVersion: v1
                clusters:
                - cluster:
                    certificate-authority: fake-ca-file
                    server: https://1.2.3.4
                  name: development
                contexts:
                - context:
                    cluster: development
                    namespace: frontend
                    user: developer
                  name: dev-frontend
                current-context: dev-frontend
                kind: Config
                preferences: {}
                users:
                - name: developer
                  user:
                    client-certificate: fake-cert-file
                    client-key: fake-key-file
                EOT
    accept_untrusted_certs = true
    cluster_context = "dev-frontend"
  } 
}

resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = "Sample Kubernetes"
  apiserver_url          = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type = "ServiceAccount"
  
  service_account {
    token = "bXktYXBw[...]K8bPxc2uQ=="
    ca_cert = "Mzk1MjgkdmRnN0pi[...]mHHRUH14gw4Q=="
  }
}

resource "azuredevops_resource_authorization" "auth" {
  project_id = azuredevops_project.project.id
  resource_id = azuredevops_serviceendpoint_kubernetes.serviceendpoint.id
  authorized = true
}

#
# https://github.com/microsoft/terraform-provider-azuredevops/issues/83
# resource "azuredevops_policy_build" "p1" {
#   scope {
#     repository_id  = azuredevops_azure_git_repository.repository.id
#     repository_ref = azuredevops_azure_git_repository.repository.default_branch
#     match_type     = "Exact"
#   }
#   settings {
#     build_definition_id    = azuredevops_build_definition.build.id
#     queue_on_source_update = true
#   }
# }
# resource "azuredevops_policy_min_reviewers" "p1" {
#   scope {
#     repository_id  = azuredevops_azure_git_repository.repository.id
#     repository_ref = azuredevops_azure_git_repository.repository.default_branch
#     match_type     = "Exact"
#   }
#   settings {
#     reviewer_count     = 2
#     submitter_can_vote = false
#   }
# }


// This section configures service connections to Azure and ACR
#
# https://github.com/microsoft/terraform-provider-azuredevops/issues/3
# resource "azuredevops_serviceendpoint_azurerm" "arm" {
#   project_id            = azuredevops_project.project.id
#   service_endpoint_name = "Sample ARM Service Connection"

#   configuration = {
#     service_principal_username = "..."
#     service_principal_password = "..."
#     subscription_id            = "..."
#     tenant_id                  = "..."
#   }
# }
# resource "azuredevops_serviceendpoint_acr" "acr" {
#   project_id            = azuredevops_project.project.id
#   service_endpoint_name = "Sample ACR Service Connection"

#   configuration = {
#     ...
#   }
# }
