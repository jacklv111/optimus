/*
 * iam
 *
 * iam api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

type TokenInfo struct {

	Username string `json:"username,omitempty"`

	ExpiredAt time.Time `json:"expiredAt,omitempty"`
}
