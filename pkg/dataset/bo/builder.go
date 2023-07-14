/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/dataset/do"
)

func BuildFromCreateDatasetRequest(req openapi.CreateDatasetRequest, associationId string) DatasetBoInterface {
	return &DatasetBo{
		datasetDo: do.DatasetDo{
			ID:                     uuid.New().String(),
			Name:                   req.Name,
			Description:            req.Description,
			AssociationId:          associationId,
			RawDataType:            req.RawDataType,
			AnnotationTemplateType: req.AnnotationTemplateType,
			CoverImageUrl:          sql.NullString{String: req.CoverImageUrl, Valid: req.CoverImageUrl != ""},
			AnnotationTemplateId:   sql.NullString{String: req.AnnotationTemplateId, Valid: req.AnnotationTemplateId != ""},
		},
	}
}

func BuildWithId(id string) DatasetBoInterface {
	return &DatasetBo{
		datasetDo: do.DatasetDo{
			ID: id,
		},
	}
}
