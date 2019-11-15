package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccResourceArmDashboard_basic(t *testing.T) {
	resourceName := "azurerm_dashboard.test"
	ri := tf.AccRandTimeInt()
	config := testResourceArmDashboard_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDashboardExists(resourceName),
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

func testCheckAzureRMDashboardExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dashboardName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Dashboard: %s", dashboardName)
		}

		client := testAccProvider.Meta().(*ArmClient).Portal.DashboardsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, dashboardName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dashboardsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Dashboard %q (resource group: %q) does not exist", dashboardName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Portal.DashboardsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dashboard" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Dashboard still exists:\n%+v", resp)
		}
	}

	return nil
}

func testResourceArmDashboard_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test-group" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dashboard" "test" {
  name                 = "my-test-dashboard"
  resource_group_name  = azurerm_resource_group.test-group.name
  location             = azurerm_resource_group.test-group.location
  dashboard_properties = <<DASH
{
   "lenses": {
        "0": {
            "order": 0,
            "parts": {
                "0": {
                    "position": {
                        "x": 0,
                        "y": 0,
                        "rowSpan": 2,
                        "colSpan": 3
                    },
                    "metadata": {
                        "inputs": [],
                        "type": "Extension/HubsExtension/PartType/MarkdownPart",
                        "settings": {
                            "content": {
                                "settings": {
                                    "content": "## This is only a test :)",
                                    "subtitle": "",
                                    "title": "Test MD Tile"
                                }
                            }
                        }
                    }
				}
			}
		}
	}
}
DASH
}
`, rInt, location)
}
