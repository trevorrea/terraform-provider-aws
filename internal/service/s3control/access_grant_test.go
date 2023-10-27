// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package s3control_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfs3control "github.com/hashicorp/terraform-provider-aws/internal/service/s3control"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccAccessGrant_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_access_grant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3ControlEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAccessGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessGrantConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAccessGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_grant_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "access_grant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "access_grants_location_id"),
					resource.TestCheckResourceAttr(resourceName, "access_grants_location_configuration.#", "0"),
					acctest.CheckResourceAttrAccountID(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "grant_scope"),
					resource.TestCheckResourceAttr(resourceName, "permission", "READ"),
					resource.TestCheckNoResourceAttr(resourceName, "s3_prefix_type"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func testAccAccessGrant_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3control_access_grant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3ControlEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAccessGrantsLocationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessGrantConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessGrantExists(ctx, resourceName),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfs3control.ResourceAccessGrantsLocation, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAccessGrantDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3ControlClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_s3control_access_grant" {
				continue
			}

			_, err := tfs3control.FindAccessGrantByTwoPartKey(ctx, conn, rs.Primary.Attributes["account_id"], rs.Primary.Attributes["access_grant_id"])

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("S3 Access Grant %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAccessGrantExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3ControlClient(ctx)

		_, err := tfs3control.FindAccessGrantByTwoPartKey(ctx, conn, rs.Primary.Attributes["account_id"], rs.Primary.Attributes["access_grant_id"])

		return err
	}
}

func testAccAccessGrantConfig_baseCustomLocation(rName string) string {
	return acctest.ConfigCompose(testAccAccessGrantsLocationConfig_baseCustomLocation(rName), `
data "aws_iam_user" "test" {
  user_name = "teamcity"
}
`)
}

func testAccAccessGrantConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccAccessGrantConfig_baseCustomLocation(rName), fmt.Sprintf(`
resource "aws_s3control_access_grants_location" "test" {
  depends_on = [aws_iam_role_policy.test, aws_s3control_access_grants_instance.test]

  iam_role_arn   = aws_iam_role.test.arn
  location_scope = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.test.key}*"
}

resource "aws_s3control_access_grant" "test" {
  access_grants_location_id = aws_s3control_access_grants_location.test.access_grants_location_id
  permission                = "READ"

  grantee {
    grantee_type       = "IAM"
    grantee_identifier = data.aws_iam_user.test.arn
  }
}
`, rName))
}
