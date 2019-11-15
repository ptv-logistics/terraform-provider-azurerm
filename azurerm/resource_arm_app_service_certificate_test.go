package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceCertificate_Pfx(t *testing.T) {
	resourceName := "azurerm_app_service_certificate.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	config := testAccAzureRMAppServiceCertificatePfx(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "password", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "thumbprint", "7B985BF42467791F23E52B364A3E8DEBAB9C606E"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pfx_blob", "password"},
			},
		},
	})
}

func TestAccAzureRMAppServiceCertificate_PfxNoPassword(t *testing.T) {
	resourceName := "azurerm_app_service_certificate.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	config := testAccAzureRMAppServiceCertificatePfxNoPassword(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "thumbprint", "7B985BF42467791F23E52B364A3E8DEBAB9C606E"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pfx_blob"},
			},
		},
	})
}

func TestAccAzureRMAppServiceCertificate_KeyVault(t *testing.T) {
	resourceName := "azurerm_app_service_certificate.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	config := testAccAzureRMAppServiceCertificateKeyVault(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "thumbprint", "7B985BF42467791F23E52B364A3E8DEBAB9C606E"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_vault_secret_id"},
			},
		},
	})
}

func testAccAzureRMAppServiceCertificatePfx(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestwebcert%d"
  location = "%s"
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  pfx_blob            = filebase64("testdata/app_service_certificate.pfx")
  password            = "terraform"
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServiceCertificatePfxNoPassword(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestwebcert%d"
  location = "%s"
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  pfx_blob            = filebase64("testdata/app_service_certificate_nopassword.pfx")
}
`, rInt, location, rInt)
}

func testAccAzureRMAppServiceCertificateKeyVault(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {}

data "azuread_service_principal" "test" {
  display_name = "Microsoft Azure App Service"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestwebcert%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tenant_id = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.service_principal_object_id
    secret_permissions      = ["delete", "get", "set"]
    certificate_permissions = ["create", "delete", "get", "import"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azuread_service_principal.test.object_id
    secret_permissions      = ["get"]
    certificate_permissions = ["get"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest%d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_secret_id = azurerm_key_vault_certificate.test.id
}
`, rInt, location, rInt, rInt, rInt)
}

func testCheckAzureRMAppServiceCertificateDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).Web.CertificatesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
