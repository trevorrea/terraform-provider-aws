---
subcategory: "SESv2 (Simple Email V2)"
layout: "aws"
page_title: "AWS: aws_sesv2_account_suppression_attributes"
description: |-
  Manages AWS SESv2 (Simple Email V2) account-level suppression attributes.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_sesv2_account_suppression_attributes

Manages AWS SESv2 (Simple Email V2) account-level suppression attributes.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { Sesv2AccountSuppressionAttributes } from "./.gen/providers/aws/sesv2-account-suppression-attributes";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new Sesv2AccountSuppressionAttributes(this, "example", {
      suppressedReasons: ["COMPLAINT"],
    });
  }
}

```

## Argument Reference

The following arguments are required:

* `suppressedReasons` - (Required) A list that contains the reasons that email addresses will be automatically added to the suppression list for your account. Valid values: `COMPLAINT`, `BOUNCE`.

## Attribute Reference

This resource exports no additional attributes.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import account-level suppression attributes using the account ID. For example:

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { Sesv2AccountSuppressionAttributes } from "./.gen/providers/aws/sesv2-account-suppression-attributes";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    Sesv2AccountSuppressionAttributes.generateConfigForImport(
      this,
      "example",
      "123456789012"
    );
  }
}

```

Using `terraform import`, import account-level suppression attributes using the account ID. For example:

```console
% terraform import aws_sesv2_account_suppression_attributes.example 123456789012
```

<!-- cache-key: cdktf-0.20.8 input-e4fe45ccad6391ce28893fb1f786a670740820da345f9e02df98116ee33cd05d -->