package iam_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfiam "github.com/hashicorp/terraform-provider-aws/internal/service/iam"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccServiceSpecificCredential_basic(t *testing.T) {
	var cred iam.ServiceSpecificCredentialMetadata

	resourceName := "aws_iam_service_specific_credential.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckServiceSpecificCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccerviceSpecificCredentialBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceSpecificCredentialExists(resourceName, &cred),
					resource.TestCheckResourceAttrPair(resourceName, "user_name", "aws_iam_user.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "service_name", "codecommit.amazonaws.com"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "service_user_name"),
					resource.TestCheckResourceAttrSet(resourceName, "service_specific_credential_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_password"},
			},
		},
	})
}

func TestAccServiceSpecificCredential_multi(t *testing.T) {
	var cred iam.ServiceSpecificCredentialMetadata

	resourceName := "aws_iam_service_specific_credential.test"
	resourceName2 := "aws_iam_service_specific_credential.test2"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckServiceSpecificCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccerviceSpecificCredentialMultiConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceSpecificCredentialExists(resourceName, &cred),
					resource.TestCheckResourceAttrPair(resourceName, "user_name", "aws_iam_user.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "service_name", "codecommit.amazonaws.com"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "service_user_name"),
					resource.TestCheckResourceAttrSet(resourceName, "service_specific_credential_id"),
					resource.TestCheckResourceAttrPair(resourceName2, "user_name", "aws_iam_user.test", "name"),
					resource.TestCheckResourceAttr(resourceName2, "service_name", "codecommit.amazonaws.com"),
					resource.TestCheckResourceAttr(resourceName2, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName2, "service_user_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "service_specific_credential_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_password"},
			},
		},
	})
}

func TestAccServiceSpecificCredential_status(t *testing.T) {
	var cred iam.ServiceSpecificCredentialMetadata

	resourceName := "aws_iam_service_specific_credential.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckServiceSpecificCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceSpecificCredentialConfigStatus(rName, "Inactive"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceSpecificCredentialExists(resourceName, &cred),
					resource.TestCheckResourceAttr(resourceName, "status", "Inactive"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_password"},
			},
			{
				Config: testAccServiceSpecificCredentialConfigStatus(rName, "Active"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceSpecificCredentialExists(resourceName, &cred),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
				),
			},
			{
				Config: testAccServiceSpecificCredentialConfigStatus(rName, "Inactive"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceSpecificCredentialExists(resourceName, &cred),
					resource.TestCheckResourceAttr(resourceName, "status", "Inactive"),
				),
			},
		},
	})
}

func TestAccServiceSpecificCredential_disappears(t *testing.T) {
	var cred iam.ServiceSpecificCredentialMetadata
	resourceName := "aws_iam_service_specific_credential.test"

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckServiceSpecificCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccerviceSpecificCredentialBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceSpecificCredentialExists(resourceName, &cred),
					acctest.CheckResourceDisappears(acctest.Provider, tfiam.ResourceServiceSpecificCredential(), resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfiam.ResourceServiceSpecificCredential(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckServiceSpecificCredentialExists(n string, cred *iam.ServiceSpecificCredentialMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server Cert ID is set")
		}
		conn := acctest.Provider.Meta().(*conns.AWSClient).IAMConn

		serviceName, userName, credId, err := tfiam.DecodeServiceSpecificCredentialId(rs.Primary.ID)
		if err != nil {
			return err
		}

		output, err := tfiam.FindServiceSpecificCredential(conn, serviceName, userName, credId)
		if err != nil {
			return err
		}

		*cred = *output

		return nil
	}
}

func testAccCheckServiceSpecificCredentialDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).IAMConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_iam_service_specific_credential" {
			continue
		}

		serviceName, userName, credId, err := tfiam.DecodeServiceSpecificCredentialId(rs.Primary.ID)
		if err != nil {
			return err
		}

		output, err := tfiam.FindServiceSpecificCredential(conn, serviceName, userName, credId)

		if tfresource.NotFound(err) {
			continue
		}

		if output != nil {
			return fmt.Errorf("IAM Service Specific Credential (%s) still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccerviceSpecificCredentialBasicConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_user" "test" {
  name = %[1]q
}

resource "aws_iam_service_specific_credential" "test" {
  service_name = "codecommit.amazonaws.com"
  user_name    = aws_iam_user.test.name
}
`, rName)
}

func testAccerviceSpecificCredentialMultiConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_user" "test" {
  name = %[1]q
}

resource "aws_iam_service_specific_credential" "test" {
  service_name = "codecommit.amazonaws.com"
  user_name    = aws_iam_user.test.name
}

resource "aws_iam_service_specific_credential" "test2" {
  service_name = "codecommit.amazonaws.com"
  user_name    = aws_iam_user.test.name
}
`, rName)
}

func testAccServiceSpecificCredentialConfigStatus(rName, status string) string {
	return fmt.Sprintf(`
resource "aws_iam_user" "test" {
  name = %[1]q
}

resource "aws_iam_service_specific_credential" "test" {
  service_name = "codecommit.amazonaws.com"
  user_name    = aws_iam_user.test.name
  status       = %[2]q
}
`, rName, status)
}
