package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

// NOTE: all of the tests in this file are for functionality which will be removed in 2.0

func TestAccAzureRMKubernetesCluster_legacyAgentPoolProfileAvailabilitySet(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_legacyAgentPoolProfileAvailabilitySet(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.type", "AvailabilitySet"),
				),
				// since users are prompted to move to `default_node_pool`
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_legacyAgentPoolProfileVMSS(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_legacyAgentPoolProfileVMSS(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
				),
				// since users are prompted to move to `default_node_pool`
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_legacyAgentPoolProfileAvailabilitySet(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  agent_pool_profile {
    name    = "default"
    count   = 1
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_legacyAgentPoolProfileVMSS(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  agent_pool_profile {
    name    = "default"
    count   = 1
    type    = "VirtualMachineScaleSets"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}
