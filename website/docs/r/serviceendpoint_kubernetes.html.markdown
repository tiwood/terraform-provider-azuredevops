# azuredevops_serviceendpoint_kubernetes
Manages a Kubernetes service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
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
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `apiserver_url` - (Required) The Service Endpoint description.
* `authorization_type` - (Required) The authentication method used to authenticate on the Kubernetes cluster. The value should be one of AzureSubscription, Kubeconfig, ServiceAccount.
* `azure_subscription` - (Optional) The configuration for authorization_type="AzureSubscription".
  * `azure_environment` - (Optional) An Azure environment an independent deployment of Microsoft Azure, such as AzureCloud for global Azure. The default value is AzureCloud.
  * `cluster_name` - (Required) The name of the Kubernetes cluster.
  * `subscription_id` - (Required) The id of the Azure subscription.
  * `subscription_name` - (Required) The name of the Azure subscription.
  * `tenant_id` - (Required) The Azure Account tenant id.
  * `resourcegroup_id` - (Required) The resource group id, to which the Kubernetes cluster is deployed.
  * `namespace` - (Optional) The Kubernetes namespace. Default value is "default".
* `kubeconfig` - (Optional) The configuration for authorization_type="Kubeconfig".
  * `kube_config` - (Required) The content of the kubeconfig as yaml, for the user to be used to communicate with the API server of the Kubernetes cluster.
  * `accept_untrusted_certs` - (Optional) The username for Docker Hub account.
  * `cluster_context` - (Optional) The context of the Kubernetes cluster. Default value is the current-context in kubeconfig.
* `service_account` - (Optional) The configuration for authorization_type="ServiceAccount".
  * `token` - (Required) The token from a Kubernetes secret object.
  * `ca_cert` - (Required) The certificate from a Kubernetes secret object.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)
