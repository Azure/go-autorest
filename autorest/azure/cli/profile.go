package cli

import "github.com/mitchellh/go-homedir"

// AzureCLIProfile represents a Profile from the Azure CLI
type AzureCLIProfile struct {
	InstallationID string                 `json:"installationId"`
	Subscriptions  []AzureCLISubscription `json:"subscriptions"`
}

// AzureCLISubscription represents a Subscription from the Azure CLI
type AzureCLISubscription struct {
	EnvironmentName string `json:"environmentName"`
	ID              string `json:"id"`
	IsDefault       bool   `json:"isDefault"`
	Name            string `json:"name"`
	State           string `json:"state"`
	TenantID        string `json:"tenantId"`
}

// AzureCLIProfilePath returns the path where the Azure Profile is stored from the Azure CLI
func AzureCLIProfilePath() (string, error) {
	return homedir.Expand("~/.azure/azureProfile.json")
}
