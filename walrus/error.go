package walrus

type ClientError struct{ Msg string }
func (e *ClientError) Error() string { return e.Msg }

type RetryableClientError struct{ Msg string }
func (e *RetryableClientError) Error() string { return e.Msg }

type NoBlobStatusReceivedError struct{ ClientError }
type NoVerifiedBlobStatusReceivedError struct{ ClientError }
type NoBlobMetadataReceivedError struct{ RetryableClientError }
type NotEnoughSliversReceivedError struct{ RetryableClientError }
type NotEnoughBlobConfirmationsError struct{ RetryableClientError }
type BehindCurrentEpochError struct{ RetryableClientError }
type BlobNotCertifiedError struct{ RetryableClientError }
type InconsistentBlobError struct{ ClientError }
type BlobBlockedError struct{ Msg string }
func (e *BlobBlockedError) Error() string { return e.Msg }
