/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type DataPoolStatistics struct {
	RawDataCount      int32
	RawDataTotalSize  int64 // in bytes
	LabelCount        int32
	LabelDistribution []LabelDistribution
}
type LabelDistribution struct {
	LabelId string
	Count   int32
	Ratio   float32
}
