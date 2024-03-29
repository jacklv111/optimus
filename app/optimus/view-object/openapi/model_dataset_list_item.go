/*
 * optimus
 *
 * optimus api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DatasetListItem struct {

	// The dataset id
	Id string `json:"id"`

	// The dataset name
	Name string `json:"name"`

	// The dataset description
	Description string `json:"description"`

	// The dataset created time
	CreatedAt int64 `json:"createdAt"`

	// The dataset updated time
	UpdatedAt int64 `json:"updatedAt,omitempty"`
}
