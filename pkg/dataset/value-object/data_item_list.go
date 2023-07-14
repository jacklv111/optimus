/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type DataItemList struct {
	RawDataType            string
	AnnotationTemplateId   string
	AnnotationTemplateType string
	RawData                []RawDataItem
	Annotations            map[string]AnnotationItem
}

func (dataItemList DataItemList) GetRawDataIdListFromAnnotations() (idList []string) {
	for _, anno := range dataItemList.Annotations {
		idList = append(idList, anno.RawDataId)
	}
	return
}

func (dataItemList DataItemList) GetRawDataIdListFromRawData() (idList []string) {
	for _, rd := range dataItemList.RawData {
		idList = append(idList, rd.RawDataId)
	}
	return
}

type RawDataItem struct {
	RawDataId  string
	RawDataUrl string
}

type AnnotationItem struct {
	RawDataId     string
	AnnotationId  string
	AnnotationUrl string
	LabelList     []string
}
