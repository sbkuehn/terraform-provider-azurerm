package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/automation"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationRunbook_PSWorkflow(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAutomationRunbook_PSWorkflow(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExistsAndType("azurerm_automation_runbook.test", automation.PowerShellWorkflow),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMAutomationRunbook_PSWorkflowWithHash(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAutomationRunbook_PSWorkflowWithHash(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationRunbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationRunbookExistsAndType("azurerm_automation_runbook.test", automation.PowerShellWorkflow),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMAutomationRunbookDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationRunbookClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_runbook" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Runbook still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationRunbookExistsAndType(name string, runbookType automation.RunbookTypeEnum) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Runbook: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automationRunbookClient

		resp, err := conn.Get(resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Runbook '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationRunbookClient: %s", err)
		}

		if resp.RunbookType != runbookType {
			return fmt.Errorf("Current runbook type %s is not consistent with the checked value %s", resp.RunbookType, runbookType)
		}

		return nil
	}
}

func testAccAzureRMAutomationRunbook_PSWorkflow(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
 name = "acctestRG-%d"
 location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
        name = "Free"
  }
}

resource "azurerm_automation_runbook" "test" {
  name	 	      = "Get-AzureVMTutorial"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
 
  account_name        = "${azurerm_automation_account.test.name}"
  log_verbose	      = "true"
  log_progress	      = "true"
  description	      = "This is a test runbook for terraform acceptance test"
  runbook_type	      = "PowerShellWorkflow"
  publish_content_link {
	uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAutomationRunbook_PSWorkflowWithHash(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
 name = "acctestRG-%d"
 location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
        name = "Free"
  }
}

resource "azurerm_automation_runbook" "test" {
  name                = "Get-AzureVMTutorial"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  account_name        = "${azurerm_automation_account.test.name}"
  log_verbose         = "true"
  log_progress        = "true"
  description         = "This is a test runbook for terraform acceptance test"
  runbook_type        = "PowerShellWorkflow"
  publish_content_link {
        uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
	version = "1.0.0.0"
	hash {
		algorithm = "SHA256"
		value = "115775B8FF2BE672D8A946BD0B489918C724DDE15A440373CA54461D53010A80"
	}
  }
}
`, rInt, location, rInt)
}