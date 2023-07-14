/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type CreatePoolParams struct {
	Name           string
	Desc           string
	IsFromVersion  bool
	SrcVersionName string
}
