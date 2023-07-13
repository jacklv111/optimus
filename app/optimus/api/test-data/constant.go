/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package testdata

import "net/http"

func GetApplicationJsonHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return header
}
