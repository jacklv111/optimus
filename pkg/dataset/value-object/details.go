/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type DatasetDetails struct {
	DatasetId              string
	Name                   string
	Description            string
	CreatedAt              int64
	AssociationId          string
	Versions               []VersionDetails
	Pools                  []PoolDetails
	RawDataType            string
	AnnotationTemplateType string
	AnnotationTemplateId   string
	CoverImageUrl          string
}

type VersionDetails struct {
	VersionName     string
	Description     string
	CreatedAt       int64
	UpdatedAt       int64
	TrainRawDataNum int
	ValRawDataNum   int
	TestRawDataNum  int
}

type PoolDetails struct {
	PoolName    string
	Description string
	CreatedAt   int64
	UpdatedAt   int64
}
