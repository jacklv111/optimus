/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type GenVersionFromExistedViewParams struct {
	Name                  string
	Description           string
	TrainRawDataViewId    string
	TrainAnnotationViewId string
	TrainRawDataNum       int32
	TrainRawDataRatio     float32
	TrainTotalDataSize    int64

	ValRawDataViewId    string
	ValAnnotationViewId string
	ValRawDataNum       int32
	ValRawDataRatio     float32
	ValTotalDataSize    int64

	TestRawDataViewId    string
	TestAnnotationViewId string
	TestRawDataNum       int32
	TestRawDataRatio     float32
	TestTotalDataSize    int64
}
