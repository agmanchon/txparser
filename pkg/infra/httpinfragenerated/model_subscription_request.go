// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Blockchain Transaction Parser API
 *
 * API for parsing and tracking blockchain transactions with subscription capabilities.
 *
 * API version: 1.0.0
 */

package httpinfragenerated

type SubscriptionRequest struct {

	// Valid blockchain address to subscribe to
	Address string `json:"address" validate:"regexp=^0x[a-fA-F0-9]{40}$"`
}

// AssertSubscriptionRequestRequired checks if the required fields are not zero-ed
func AssertSubscriptionRequestRequired(obj SubscriptionRequest) error {
	elements := map[string]interface{}{
		"address": obj.Address,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertSubscriptionRequestConstraints checks if the values respects the defined constraints
func AssertSubscriptionRequestConstraints(obj SubscriptionRequest) error {
	return nil
}
