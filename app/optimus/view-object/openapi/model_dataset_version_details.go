/*
 * optimus
 *
 * optimus api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DatasetVersionDetails struct {

	// The dataset version name
	Name string `json:"name,omitempty"`

	// The dataset version number
	Number int32 `json:"number,omitempty"`

	// The dataset version description
	Description string `json:"description,omitempty"`

	// The dataset version created time
	CreatedAt int64 `json:"createdAt,omitempty"`

	// The dataset version updated time
	UpdatedAt int64 `json:"updatedAt,omitempty"`

	// The dataset version train raw data number
	TrainRawDataNum int32 `json:"trainRawDataNum,omitempty"`

	// The dataset version test raw data number
	TestRawDataNum int32 `json:"testRawDataNum,omitempty"`

	// The dataset version validation raw data number
	ValidationRawDataNum int32 `json:"validationRawDataNum,omitempty"`
}
