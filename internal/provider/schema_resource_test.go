package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccSchemaResourceConfig("definition user {}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("spicedb_schema.test", "schema", "definition user {}"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccSchemaResourceConfig("definition organisation {}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("spicedb_schema.test", "schema", "definition organisation {}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSchemaResourceConfig(schema string) string {
	return fmt.Sprintf(`

resource "spicedb_schema" "test" {
  schema = %[1]q
}
`, schema)
}
