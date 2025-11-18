---
subcategory: "Batch"
layout: "aws"
page_title: "AWS: aws_batch_consumable_resource"
description: |-
  Provides a Batch Consumable Resource resource.
---

# Resource: aws_batch_consumable_resource

Provides a Batch Consumable Resource resource.

AWS Batch now supports job scheduling that takes into account consumable resources such as third-party license tokens, database access bandwidth, budgetary limits, and more. This feature enables resource-aware scheduling to help reduce job failures and wasted compute time.

## Example Usage

```terraform
resource "aws_batch_consumable_resource" "example" {
  name           = "example"
  resource_type  = "REPLENISHABLE"
  total_quantity = 100

  tags = {
    Name = "Example Batch Consumable Resource"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `name` - (Required) Specifies the name of the consumable resource. Must be unique.
* `resource_type` - (Required) Indicates whether the resource is reusable after job completion. Valid values are `REPLENISHABLE` and `NON_REPLENISHABLE`.
* `total_quantity` - (Required) The total amount of the consumable resource that is available. Must be non-negative.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - The Amazon Resource Name of the consumable resource.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Batch Consumable Resource using the `arn`. For example:

```terraform
import {
  to = aws_batch_consumable_resource.example
  id = "arn:aws:batch:us-east-1:123456789012:consumable-resource/example"
}
```

Using `terraform import`, import Batch Consumable Resource using the `arn`. For example:

```console
% terraform import aws_batch_consumable_resource.example arn:aws:batch:us-east-1:123456789012:consumable-resource/example
```
