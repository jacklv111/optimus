/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type Permission struct {
	Domain string
	// Name         string
	ResourceType string
	ResourceId   string
	Action       []string
	Effect       string
}

type PermissionEnforce struct {
	Domain       string
	Name         string
	ResourceType string
	ResourceId   string
	Action       string
}
