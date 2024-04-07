// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ce

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	awstypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ce_cost_allocation_tag")
func ResourceCostAllocationTag() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCostAllocationTagUpdate,
		ReadWithoutTimeout:   resourceCostAllocationTagRead,
		UpdateWithoutTimeout: resourceCostAllocationTagUpdate,
		DeleteWithoutTimeout: resourceCostAllocationTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: enum.Validate[awstypes.CostAllocationTagStatus](),
			},
			"tag_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCostAllocationTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).CEClient(ctx)

	costAllocTag, err := FindCostAllocationTagByKey(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		create.LogNotFoundRemoveState(names.CE, create.ErrActionReading, ResNameCostAllocationTag, d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return create.AppendDiagError(diags, names.CE, create.ErrActionReading, ResNameCostAllocationTag, d.Id(), err)
	}

	d.Set("tag_key", costAllocTag.TagKey)
	d.Set("status", costAllocTag.Status)
	d.Set("type", costAllocTag.Type)

	return diags
}

func resourceCostAllocationTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	key := d.Get("tag_key").(string)

	updateTagStatus(ctx, d, meta, false)

	d.SetId(key)

	return append(diags, resourceCostAllocationTagRead(ctx, d, meta)...)
}

func resourceCostAllocationTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateTagStatus(ctx, d, meta, true)
}

func updateTagStatus(ctx context.Context, d *schema.ResourceData, meta interface{}, delete bool) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).CEClient(ctx)

	key := d.Get("tag_key").(string)
	tagStatus := awstypes.CostAllocationTagStatusEntry{
		TagKey: aws.String(key),
		Status: awstypes.CostAllocationTagStatus(d.Get("status").(string)),
	}

	if delete {
		tagStatus.Status = awstypes.CostAllocationTagStatusInactive
	}

	input := &costexplorer.UpdateCostAllocationTagsStatusInput{
		CostAllocationTagsStatus: []awstypes.CostAllocationTagStatusEntry{tagStatus},
	}

	_, err := conn.UpdateCostAllocationTagsStatus(ctx, input)

	if err != nil {
		return create.AppendDiagError(diags, names.CE, create.ErrActionUpdating, ResNameCostAllocationTag, d.Id(), err)
	}

	return diags
}
