package main

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

var cloud = azure.PublicCloud

func main() {
	tenantID := os.Getenv("TENANT_ID")
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")
	resourceGroupName := os.Getenv("RESOURCE_GROUP")

	oauthConfig, err := adal.NewOAuthConfig(cloud.ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		panic(err)
	}

	spt, err := adal.NewServicePrincipalTokenFromMSI(
		*oauthConfig,
		cloud.ResourceManagerEndpoint)
	if err != nil {
		panic(err)
	}

	client := resources.NewGroupsClientWithBaseURI(cloud.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(spt)

	resources, err := client.ListResources(resourceGroupName, "", "", nil)
	if err != nil {
		panic(err)
	}

	result := "resources: "
	for _, r := range *(resources.Value) {
		result += *r.Name + ", "
	}

	log.Println(result)

	// update RG, add a tag
	rg, err := client.Get(resourceGroupName)
	if err != nil {
		panic(err)
	}

	value := "success"
	var tags map[string]*string
	if rg.Tags != nil {
		tags = *rg.Tags
	} else {
		tags = make(map[string]*string)
	}
	tags["msie2e"] = &value

	rg.Tags = &tags

	// this is a silly necessity
	rg.Properties.ProvisioningState = nil
	rg.ID = nil

	_, err = client.CreateOrUpdate(
		resourceGroupName,
		rg)
	if err != nil {
		panic(err)
	}
}
