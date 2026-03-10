// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_batch_consumable_resource", name="Consumable Resource")
// @Tags(identifierAttribute="arn")
// @Testing(existsType="github.com/aws/aws-sdk-go-v2/service/batch;batch.DescribeConsumableResourceOutput")
func resourceConsumableResource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceConsumableResourceCreate,
		ReadWithoutTimeout:   resourceConsumableResourceRead,
		UpdateWithoutTimeout: resourceConsumableResourceUpdate,
		DeleteWithoutTimeout: resourceConsumableResourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validName,
			},
			names.AttrResourceType: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"REPLENISHABLE", "NON_REPLENISHABLE"}, false)),
			},
			"total_quantity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceConsumableResourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BatchClient(ctx)

	name := d.Get(names.AttrName).(string)
	input := &batch.CreateConsumableResourceInput{
		ConsumableResourceName: aws.String(name),
		ResourceType:           aws.String(d.Get(names.AttrResourceType).(string)),
		TotalQuantity:          aws.Int64(int64(d.Get("total_quantity").(int))),
		Tags:                   getTagsIn(ctx),
	}

	output, err := conn.CreateConsumableResource(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Batch Consumable Resource (%s): %s", name, err)
	}

	d.SetId(aws.ToString(output.ConsumableResourceArn))

	return append(diags, resourceConsumableResourceRead(ctx, d, meta)...)
}

func resourceConsumableResourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BatchClient(ctx)

	cr, err := findConsumableResourceByARN(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Batch Consumable Resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Batch Consumable Resource (%s): %s", d.Id(), err)
	}

	d.Set(names.AttrARN, cr.ConsumableResourceArn)
	d.Set(names.AttrName, cr.ConsumableResourceName)
	d.Set(names.AttrResourceType, cr.ResourceType)
	d.Set("total_quantity", aws.ToInt64(cr.TotalQuantity))

	setTagsOut(ctx, cr.Tags)

	return diags
}

func resourceConsumableResourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BatchClient(ctx)

	if d.HasChange("total_quantity") {
		input := &batch.UpdateConsumableResourceInput{
			ConsumableResource: aws.String(d.Id()),
			Operation:          aws.String("SET"),
			Quantity:           aws.Int64(int64(d.Get("total_quantity").(int))),
		}

		_, err := conn.UpdateConsumableResource(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Batch Consumable Resource (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceConsumableResourceRead(ctx, d, meta)...)
}

func resourceConsumableResourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BatchClient(ctx)

	log.Printf("[DEBUG] Deleting Batch Consumable Resource: %s", d.Id())
	_, err := conn.DeleteConsumableResource(ctx, &batch.DeleteConsumableResourceInput{
		ConsumableResource: aws.String(d.Id()),
	})

	if errs.IsAErrorMessageContains[*awstypes.ClientException](err, "does not exist") {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Batch Consumable Resource (%s): %s", d.Id(), err)
	}

	return diags
}

func findConsumableResourceByARN(ctx context.Context, conn *batch.Client, arn string) (*batch.DescribeConsumableResourceOutput, error) {
	input := &batch.DescribeConsumableResourceInput{
		ConsumableResource: aws.String(arn),
	}

	output, err := conn.DescribeConsumableResource(ctx, input)

	if errs.IsAErrorMessageContains[*awstypes.ClientException](err, "does not exist") {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
