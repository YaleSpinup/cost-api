package api

import (
	"context"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	log "github.com/sirupsen/logrus"
)

func (o *inventoryOrchestrator) GetResourceInventory(ctx context.Context, account, spaceid string) ([]*InventoryResponse, error) {
	if spaceid == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "spaceid is required", nil)
	}

	list := []*InventoryResponse{}
	input := resourcegroupstaggingapi.GetResourcesInput{
		TagFilters: []*resourcegroupstaggingapi.TagFilter{
			{
				Key:    aws.String("spinup:spaceid"),
				Values: []*string{aws.String(spaceid)},
			},
		},
		ResourcesPerPage: aws.Int64(100),
	}

	for {
		out, err := o.client.ListResourcesWithTags(ctx, &input)
		if err != nil {
			return nil, err
		}

		for _, res := range out.ResourceTagMappingList {
			list = append(list, toInventoryResponse(res))
		}

		log.Debugf("%+v:%s", out.PaginationToken, aws.StringValue(out.PaginationToken))

		if aws.StringValue(out.PaginationToken) == "" {
			break
		}

		input.PaginationToken = out.PaginationToken

		log.Debugf("%d resources in list", len(list))
	}

	return list, nil
}
