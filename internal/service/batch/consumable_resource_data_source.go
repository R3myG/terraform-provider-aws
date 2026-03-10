// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_batch_consumable_resource", name="Consumable Resource")
// @Tags
// @Testing(tagsIdentifierAttribute="arn")
func dataSourceConsumableResource() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceConsumableResourceRead,

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			names.AttrCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"in_use_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrResourceType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			names.AttrTags: tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceConsumableResourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BatchClient(ctx)

	consumableResource, err := findConsumableResourceByARN(ctx, conn, d.Get(names.AttrARN).(string))

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("Batch Consumable Resource", err))
	}

	d.SetId(aws.ToString(consumableResource.ConsumableResourceArn))
	d.Set("available_quantity", aws.ToInt64(consumableResource.AvailableQuantity))
	if consumableResource.CreatedAt != nil {
		d.Set(names.AttrCreatedAt, consumableResource.CreatedAt.Format(time.RFC3339))
	}
	d.Set("in_use_quantity", aws.ToInt64(consumableResource.InUseQuantity))
	d.Set(names.AttrName, consumableResource.ConsumableResourceName)
	d.Set(names.AttrResourceType, consumableResource.ResourceType)
	d.Set("total_quantity", aws.ToInt64(consumableResource.TotalQuantity))

	setTagsOut(ctx, consumableResource.Tags)

	return diags
}
