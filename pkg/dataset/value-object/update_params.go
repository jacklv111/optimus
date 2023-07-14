/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type DatasetUpdateParams struct {
	Name                   string
	Description            string
	AnnotationTemplateId   string
	AnnotationTemplateType string
	RawDataType            string
	CoverImageUrl          string
}
