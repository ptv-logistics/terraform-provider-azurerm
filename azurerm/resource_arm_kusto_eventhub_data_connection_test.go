package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKustoEventHubDataConnection_basic(t *testing.T) {
	resourceName := "azurerm_kusto_eventhub_data_connection.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(6)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKustoEventHubDataConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoEventHubDataConnection_basic(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoEventHubDataConnectionExists(resourceName),
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

func testAccAzureRMKustoEventHubDataConnection_basic(rInt int, rs string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_name        = "${azurerm_kusto_cluster.test.name}"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_kusto_eventhub_data_connection" "test" {
  name                = "acctestkedc-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_name        = "${azurerm_kusto_cluster.test.name}"
  database_name       = "${azurerm_kusto_database.test.name}"

  eventhub_id    = "${azurerm_eventhub.test.id}"
  consumer_group = "${azurerm_eventhub_consumer_group.test.name}"
}
`, rInt, location, rs, rInt, rInt, rInt, rInt, rInt)
}

func testCheckAzureRMKustoEventHubDataConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Kusto.DataConnectionsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_eventhub_data_connection" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		name := rs.Primary.Attributes["name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMKustoEventHubDataConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Kusto EventHub Data Connection: %s", name)
		}

		clusterName, hasClusterName := rs.Primary.Attributes["cluster_name"]
		if !hasClusterName {
			return fmt.Errorf("Bad: no resource group found in state for Kusto EventHub Data Connection: %s", name)
		}

		databaseName, hasDatabaseName := rs.Primary.Attributes["database_name"]
		if !hasDatabaseName {
			return fmt.Errorf("Bad: no resource group found in state for Kusto EventHub Data Connection: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).Kusto.DataConnectionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Kusto EventHub Data Connection %q (resource group: %q, cluster: %q, database: %q) does not exist", name, resourceGroup, clusterName, databaseName)
			}

			return fmt.Errorf("Bad: Get on DataConnectionsClient: %+v", err)
		}

		return nil
	}
}
