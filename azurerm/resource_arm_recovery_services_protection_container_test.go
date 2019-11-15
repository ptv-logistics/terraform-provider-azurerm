package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMRecoveryProtectionContainer_basic(t *testing.T) {
	resourceGroupName := "azurerm_resource_group.test"
	vaultName := "azurerm_recovery_services_vault.test"
	fabricName := "azurerm_recovery_services_fabric.test"
	resourceName := "azurerm_recovery_services_protection_container.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryProtectionContainer_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryProtectionContainerExists(resourceGroupName, vaultName, fabricName, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMRecoveryProtectionContainer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_fabric" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  name                = "acctest-fabric-%d"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_recovery_services_protection_container" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  recovery_vault_name  = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name = "${azurerm_recovery_services_fabric.test.name}"
  name                 = "acctest-protection-cont-%d"
}

`, rInt, location, rInt, rInt, rInt)
}

func testCheckAzureRMRecoveryProtectionContainerExists(resourceGroupStateName, vaultStateName string, resourceStateName string, protectionContainerStateName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		resourceGroupState, ok := s.RootModule().Resources[resourceGroupStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceGroupStateName)
		}
		vaultState, ok := s.RootModule().Resources[vaultStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", vaultStateName)
		}
		fabricState, ok := s.RootModule().Resources[resourceStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceStateName)
		}
		protectionContainerState, ok := s.RootModule().Resources[protectionContainerStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceStateName)
		}

		resourceGroupName := resourceGroupState.Primary.Attributes["name"]
		vaultName := vaultState.Primary.Attributes["name"]
		fabricName := fabricState.Primary.Attributes["name"]
		protectionContainerName := protectionContainerState.Primary.Attributes["name"]

		// Ensure fabric exists in API
		client := testAccProvider.Meta().(*ArmClient).RecoveryServices.ProtectionContainerClient(resourceGroupName, vaultName)
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, fabricName, protectionContainerName)
		if err != nil {
			return fmt.Errorf("Bad: Get on fabricClient: %+v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: fabric: %q does not exist", fabricName)
		}

		return nil
	}
}
