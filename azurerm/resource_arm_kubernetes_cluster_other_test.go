package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMKubernetesCluster_basicAvailabilitySet(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_basicAvailabilitySet(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config_raw", ""),
					resource.TestCheckResourceAttr(resourceName, "network_profile.0.load_balancer_sku", "Basic"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_basicVMSS(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_basicVMSS(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config_raw", ""),
					resource.TestCheckResourceAttr(resourceName, "network_profile.0.load_balancer_sku", "Basic"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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
				Config: testAccAzureRMKubernetesCluster_basicVMSS(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMKubernetesCluster_requiresImport(ri, clientId, clientSecret, location),
				ExpectError: testRequiresImportError("azurerm_kubernetes_cluster"),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_linuxProfile(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_linuxProfile(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttrSet(resourceName, "linux_profile.0.admin_username"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_nodeTaints(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_nodeTaints(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.node_taints.0", "key=value:PreferNoSchedule"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_nodeResourceGroup(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_nodeResourceGroup(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeConfig(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgrade(ri, location, clientId, clientSecret, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_version", olderKubernetesVersion),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_upgrade(ri, location, clientId, clientSecret, currentKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_version", currentKubernetesVersion),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_tags(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_tags(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
			{
				Config: testAccAzureRMKubernetesCluster_tagsUpdated(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_windowsProfile(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_windowsProfile(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttrSet(resourceName, "default_node_pool.0.max_pods"),
					resource.TestCheckResourceAttrSet(resourceName, "linux_profile.0.admin_username"),
					resource.TestCheckResourceAttrSet(resourceName, "windows_profile.0.admin_username"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"windows_profile.0.admin_password",
					"service_principal.0.client_secret",
				},
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_basicAvailabilitySet(rInt int, clientId string, clientSecret string, location string) string {
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

  default_node_pool {
    name       = "default"
    node_count = 1
  	type       = "AvailabilitySet"
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_basicVMSS(rInt int, clientId string, clientSecret string, location string) string {
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

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_requiresImport(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesCluster_basicVMSS(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster" "import" {
  name                = azurerm_kubernetes_cluster.test.name
  location            = azurerm_kubernetes_cluster.test.location
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
  dns_prefix          = azurerm_kubernetes_cluster.test.dns_prefix

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, template, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_linuxProfile(rInt int, clientId string, clientSecret string, location string) string {
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

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_nodeTaints(rInt int, clientId string, clientSecret string, location string) string {
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

  default_node_pool {
    name        = "default"
    node_count  = 1
    vm_size     = "Standard_DS2_v2"
    node_taints = [
      "key=value:PreferNoSchedule"
    ]
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_nodeResourceGroup(rInt int, clientId string, clientSecret string, location string) string {
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
  node_resource_group = "acctestRGAKS-%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_tags(rInt int, clientId string, clientSecret string, location string) string {
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

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  tags = {
    dimension = "C-137"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_tagsUpdated(rInt int, clientId string, clientSecret string, location string) string {
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

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  tags = {
    dimension = "D-99"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_upgrade(rInt int, location, clientId, clientSecret, version string) string {
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
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, version, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_windowsProfile(rInt int, clientId, clientSecret, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!"
  }

  # the default node pool /has/ to be Linux agents - Windows agents can be added via the node pools resource
  default_node_pool {
    name       = "np"
    node_count = 3
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}
