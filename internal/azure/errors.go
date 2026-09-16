package azure

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// IsSkippableResponseError preserves the reference scanner behavior for Azure scopes that
// cannot be assessed because a provider/operation is unavailable rather than because the
// overall assessment is invalid.
func IsSkippableResponseError(err error) (string, bool) {
	var responseError *azcore.ResponseError
	if !errors.As(err, &responseError) {
		return "", false
	}

	switch responseError.ErrorCode {
	case "MissingRegistrationForResourceProvider", "MissingSubscriptionRegistration", "DisallowedOperation", "NotFound":
		return responseError.ErrorCode, true
	default:
		return responseError.ErrorCode, false
	}
}
