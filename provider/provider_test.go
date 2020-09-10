package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"harbor": testAccProvider,
	}
}

func testProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {

	// apiClient := testAccProvider.Meta().(*client.Client)
	// _, _, err := apiClient.SendRequest("GET", "/health", nil, 200)
	// if err != nil {
	// 	t.Fatal("Harbor instamce is not healthy")
	// }

	if v := os.Getenv("HARBOR_URL"); v == "" {
		t.Fatal("HARBOR_URL must be set for acceptance tests")
	}
	if v := os.Getenv("HARBOR_USERNAME"); v == "" {
		t.Fatal("HARBOR_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("HARBOR_PASSWORD"); v == "" {
		t.Fatal("HARBOR_PASSWORD must be set for acceptance tests")
	}

}

func testAccCheckResourceExists(resource string) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID

		apiClient := testAccProvider.Meta().(*client.Client)
		_, _, err := apiClient.SendRequest("GET", name, nil, 200)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}
