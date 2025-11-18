// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"context"
	"fmt"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfbatch "github.com/hashicorp/terraform-provider-aws/internal/service/batch"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccBatchConsumableResource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var consumableResource1 awstypes.ConsumableResourceDetail
	resourceName := "aws_batch_consumable_resource.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConsumableResourceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConsumableResourceConfig_basic(rName, "REPLENISHABLE", 100),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource1),
					acctest.CheckResourceAttrRegionalARNFormat(ctx, resourceName, names.AttrARN, "batch", "consumable-resource/{name}"),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, names.AttrResourceType, "REPLENISHABLE"),
					resource.TestCheckResourceAttr(resourceName, "total_quantity", "100"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "1"),
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

func TestAccBatchConsumableResource_update(t *testing.T) {
	ctx := acctest.Context(t)
	var consumableResource1, consumableResource2 awstypes.ConsumableResourceDetail
	resourceName := "aws_batch_consumable_resource.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConsumableResourceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConsumableResourceConfig_basic(rName, "REPLENISHABLE", 100),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource1),
					resource.TestCheckResourceAttr(resourceName, "total_quantity", "100"),
				),
			},
			{
				Config: testAccConsumableResourceConfig_basic(rName, "REPLENISHABLE", 200),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource2),
					resource.TestCheckResourceAttr(resourceName, "total_quantity", "200"),
				),
			},
		},
	})
}

func TestAccBatchConsumableResource_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var consumableResource1 awstypes.ConsumableResourceDetail
	resourceName := "aws_batch_consumable_resource.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConsumableResourceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConsumableResourceConfig_basic(rName, "REPLENISHABLE", 100),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource1),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfbatch.ResourceConsumableResource(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBatchConsumableResource_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var consumableResource1, consumableResource2, consumableResource3 awstypes.ConsumableResourceDetail
	resourceName := "aws_batch_consumable_resource.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BatchServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConsumableResourceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConsumableResourceConfig_tags1(rName, acctest.CtKey1, acctest.CtValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource1),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "1"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConsumableResourceConfig_tags2(rName, acctest.CtKey1, acctest.CtValue1Updated, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource2),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "2"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1Updated),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
			{
				Config: testAccConsumableResourceConfig_tags1(rName, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumableResourceExists(ctx, resourceName, &consumableResource3),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "1"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
		},
	})
}

func testAccCheckConsumableResourceExists(ctx context.Context, n string, v *awstypes.ConsumableResourceDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).BatchClient(ctx)

		output, err := tfbatch.FindConsumableResourceByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckConsumableResourceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_batch_consumable_resource" {
				continue
			}
			conn := acctest.Provider.Meta().(*conns.AWSClient).BatchClient(ctx)

			_, err := tfbatch.FindConsumableResourceByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Batch Consumable Resource %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccConsumableResourceConfig_basic(rName, resourceType string, totalQuantity int) string {
	return fmt.Sprintf(`
resource "aws_batch_consumable_resource" "test" {
  name           = %[1]q
  resource_type  = %[2]q
  total_quantity = %[3]d

  tags = {
    Name = "Test Batch Consumable Resource"
  }
}
`, rName, resourceType, totalQuantity)
}

func testAccConsumableResourceConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_batch_consumable_resource" "test" {
  name           = %[1]q
  resource_type  = "REPLENISHABLE"
  total_quantity = 100

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccConsumableResourceConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_batch_consumable_resource" "test" {
  name           = %[1]q
  resource_type  = "REPLENISHABLE"
  total_quantity = 100

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
