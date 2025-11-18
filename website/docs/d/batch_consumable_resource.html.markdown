---
subcategory: "Batch"
layout: "aws"
page_title: "AWS: aws_batch_consumable_resource"
description: |-
    Provides details about a Batch Consumable Resource
---

# Data Source: aws_batch_consumable_resource

The Batch Consumable Resource data source allows access to details of a specific Consumable Resource within AWS Batch.

## Example Usage

```terraform
data "aws_batch_consumable_resource" "example" {
  arn = "arn:aws:batch:us-east-1:012345678910:consumable-resource/example"
}
```

## Argument Reference

This data source supports the following arguments:

* `arn` - (Required) ARN of the consumable resource.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `name` - Name of the consumable resource.
* `resource_type` - Type of the consumable resource. Valid values are `REPLENISHABLE` and `NON_REPLENISHABLE`.
* `total_quantity` - Total amount of the consumable resource that is available.
* `tags` - Key-value map of resource tags.
