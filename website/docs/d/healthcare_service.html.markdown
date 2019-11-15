---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_service"
sidebar_current: "docs-azurerm-datasource-healthcare-service-x"
description: |-
  Get information about an existing Healthcare Service
---

# Data Source: azurerm_healthcare_service

Use this data source to access information about an existing Healthcare Service

## Example Usage

```hcl
data "azurerm_healthcare_service" "example" {
  name                = "example-healthcare_service"
  resource_group_name = "example-resources"
  location            = "westus2"
}

output "healthcare_service_id" {
  name                = "uniquefhirname"
  resource_group_name = "sample-resource-group"
  kind                = "fhir-R4"
  cosmosdb_throughput = "2000"

  access_policy_object_ids = ["xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]

  tags = {
    "environment" = "testenv"
    "purpose"     = "AcceptanceTests"
  }
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Healthcare Service.

* `resource_group_name` - (Required) The name of the Resource Group in which the Healthcare Service exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the Service is located.

~> **Please Note**: Not all locations support this resource. Some are `West US 2`, `North Central US`, and `UK West`. 

* `kind` - The type of the service.
* `authentication_configuration` - An `authentication_configuration` block as defined below.
* `cosmosdb_offer_throughput` - The provisioned throughput for the backing database.
* `cors_configuration` - A `cors_configuration` block as defined below.
* `tags` - A mapping of tags to assign to the resource.

---
An `authentication_configuration` exports the following:

* `authority` - The Azure Active Directory (tenant) that serves as the authentication authority to access the service. 
* `audience` - The intended audience to receive authentication tokens for the service. 
* `smart_proxy_enabled` - Is the 'SMART on FHIR' option for mobile and web implementations enbled?

---
A `cors_configuration` block exports the following:

* `allowed_origins` - The set of origins to be allowed via CORS.
* `allowed_headers` - The set of headers to be allowed via CORS.
* `allowed_methods` - The methods to be allowed via CORS.
* `max_age_in_seconds` - The max age to be allowed via CORS.
* `allow_credentials` - Are credentials are allowed via CORS?

