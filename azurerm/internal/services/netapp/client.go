package netapp

import (
	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient *netapp.AccountsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	accountClient := netapp.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient: &accountClient,
	}
}
