package seal

import "fmt"

type SealError struct{ Msg string }
func (e *SealError) Error() string { return e.Msg }

type UserError struct{ SealError }

type APIError struct {
	Code int
	Msg  string
}
func (e *APIError) Error() string { return fmt.Sprintf("seal api error %d: %s", e.Code, e.Msg) }

type InvalidPTBError struct{ APIError }
type InvalidPackageError struct{ APIError }
type InvalidParameterError struct{ APIError }
type InvalidUserSignatureError struct{ APIError }
type InvalidSessionKeySignatureError struct{ APIError }
type InvalidMVRNameError struct{ APIError }
type InvalidKeyServerObjectIDError struct{ APIError }
type UnsupportedPackageIDError struct{ APIError }
type InvalidSDKVersionError struct{ APIError }
type InvalidSDKTypeError struct{ APIError }
type DeprecatedSDKVersionError struct{ APIError }
type NoAccessError struct{ APIError }
type ExpiredSessionKeyError struct{ APIError }
type InternalError struct{ APIError }
type GeneralError struct{ APIError }

type InvalidPersonalMessageSignatureError struct{ UserError }
type InvalidGetObjectError struct{ UserError }
type UnsupportedFeatureError struct{ UserError }
type UnsupportedNetworkError struct{ UserError }
type InvalidKeyServerError struct{ UserError }
type InvalidKeyServerVersionError struct{ UserError }
type InvalidCiphertextError struct{ UserError }
type InvalidThresholdError struct{ UserError }
type InconsistentKeyServersError struct{ UserError }
type DecryptionError struct{ UserError }
type InvalidClientOptionsError struct{ UserError }
type TooManyFailedFetchKeyRequestsError struct{ UserError }

func ToMajorityError(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	counts := map[string]int{}
	byMsg := map[string]error{}
	for _, err := range errors {
		msg := err.Error()
		counts[msg]++
		byMsg[msg] = err
	}
	var (
		max int
		key string
	)
	for k, n := range counts {
		if n > max {
			max = n
			key = k
		}
	}
	return byMsg[key]
}
