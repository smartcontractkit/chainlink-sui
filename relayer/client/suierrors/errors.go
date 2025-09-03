/*
Package `errors` provides Sui error definitions and utilities for mapping raw error
messages returned by the Sui JSON-RPC API into well-defined Go errors.

The defined errors are based on the Sui error types found in the official Sui source:

	https://github.com/MystenLabs/sui/blob/main/crates/sui-types/src/error.rs

This package offers functions to parse error messages (e.g. ParseSuiErrorMessage) and
determine if an error is considered retryable (e.g. IsRetryable) according to Sui semantics.
*/

package suierrors

import (
	"errors"
	"strings"
)

// ErrorCategory represents an enumerated type for Sui error categories.
type ErrorCategory int

const (
	// Enumerated error categories.
	ObjectErrors ErrorCategory = iota
	MoveCallErrors
	GasErrors
	SignatureAndTransactionErrors
	CheckpointAndConsensusErrors
	PublishingErrors
	SoftBundleErrors
)

func (c ErrorCategory) String() string {
	switch c {
	case ObjectErrors:
		return "Object Errors"
	case MoveCallErrors:
		return "Move Call Errors"
	case GasErrors:
		return "Gas Errors"
	case SignatureAndTransactionErrors:
		return "Signature & Transaction Errors"
	case CheckpointAndConsensusErrors:
		return "Checkpoint & Consensus Errors"
	case PublishingErrors:
		return "Publishing Errors"
	case SoftBundleErrors:
		return "Soft Bundle Errors"
	default:
		return "Unknown Error Category"
	}
}

// SuiError is a custom error type that carries a category and a message.
type SuiError struct {
	Category ErrorCategory
	Message  string
}

func (e *SuiError) Error() string {
	return e.Message
}

// NewSuiError creates a new SuiError with the given category and message.
func NewSuiError(category ErrorCategory, message string) *SuiError {
	return &SuiError{
		Category: category,
		Message:  message,
	}
}

// ====================
// Error Definitions
// ====================

// Object Errors
var ErrMutableObjectUsedMoreThanOnce = NewSuiError(ObjectErrors, "MutableObjectUsedMoreThanOnce")
var ErrObjectInputArityViolation = NewSuiError(ObjectErrors, "ObjectInputArityViolation")
var ErrObjectNotFound = NewSuiError(ObjectErrors, "ObjectNotFound")
var ErrObjectVersionUnavailableForConsumption = NewSuiError(ObjectErrors, "ObjectVersionUnavailableForConsumption")
var ErrInvalidChildObjectArgument = NewSuiError(ObjectErrors, "InvalidChildObjectArgument")
var ErrInvalidObjectDigest = NewSuiError(ObjectErrors, "InvalidObjectDigest")
var ErrInvalidSequenceNumber = NewSuiError(ObjectErrors, "InvalidSequenceNumber")
var ErrObjectSequenceNumberTooHigh = NewSuiError(ObjectErrors, "ObjectSequenceNumberTooHigh")
var ErrObjectDeleted = NewSuiError(ObjectErrors, "ObjectDeleted")
var ErrInaccessibleSystemObject = NewSuiError(ObjectErrors, "InaccessibleSystemObject")
var ErrNotSharedObjectError = NewSuiError(ObjectErrors, "NotSharedObjectError")
var ErrNotOwnedObjectError = NewSuiError(ObjectErrors, "NotOwnedObjectError")
var ErrDuplicateObjectRefInput = NewSuiError(ObjectErrors, "DuplicateObjectRefInput")

// Move Call Errors
var ErrMovePackageAsObject = NewSuiError(MoveCallErrors, "MovePackageAsObject")
var ErrMoveObjectAsPackage = NewSuiError(MoveCallErrors, "MoveObjectAsPackage")
var ErrBlockedMoveFunction = NewSuiError(MoveCallErrors, "BlockedMoveFunction")
var ErrTransferObjectWithoutPublicTransferError = NewSuiError(MoveCallErrors, "TransferObjectWithoutPublicTransferError")
var ErrUnsupported = NewSuiError(MoveCallErrors, "Unsupported")
var ErrMoveFunctionInputError = NewSuiError(MoveCallErrors, "MoveFunctionInputError")
var ErrPostRandomCommandRestrictions = NewSuiError(MoveCallErrors, "PostRandomCommandRestrictions")

// Gas Errors
var ErrMissingGasPayment = NewSuiError(GasErrors, "MissingGasPayment")
var ErrGasObjectNotOwnedObject = NewSuiError(GasErrors, "GasObjectNotOwnedObject")
var ErrGasBudgetTooHigh = NewSuiError(GasErrors, "GasBudgetTooHigh")
var ErrGasBudgetTooLow = NewSuiError(GasErrors, "GasBudgetTooLow")
var ErrGasBalanceTooLow = NewSuiError(GasErrors, "GasBalanceTooLow")
var ErrGasPriceUnderRGP = NewSuiError(GasErrors, "GasPriceUnderRGP")
var ErrGasPriceTooHigh = NewSuiError(GasErrors, "GasPriceTooHigh")
var ErrInvalidGasObject = NewSuiError(GasErrors, "InvalidGasObject")
var ErrInsufficientBalanceToCoverMinimalGas = NewSuiError(GasErrors, "InsufficientBalanceToCoverMinimalGas")
var ErrUnexpectedGasPaymentObject = NewSuiError(GasErrors, "UnexpectedGasPaymentObject")
var ErrGasPriceMismatchError = NewSuiError(GasErrors, "GasPriceMismatchError")
var ErrInsufficientGas = NewSuiError(GasErrors, "InsufficientGas")

// Signature & Transaction Errors
var ErrIncorrectUserSignature = NewSuiError(SignatureAndTransactionErrors, "IncorrectUserSignature")
var ErrInvalidBatchTransaction = NewSuiError(SignatureAndTransactionErrors, "InvalidBatchTransaction")
var ErrTransactionDenied = NewSuiError(SignatureAndTransactionErrors, "TransactionDenied")
var ErrUnsupportedSponsoredTransactionKind = NewSuiError(SignatureAndTransactionErrors, "UnsupportedSponsoredTransactionKind")
var ErrEmptyCommandInput = NewSuiError(SignatureAndTransactionErrors, "EmptyCommandInput")
var ErrEmptyInputCoins = NewSuiError(SignatureAndTransactionErrors, "EmptyInputCoins")
var ErrInvalidIdentifier = NewSuiError(SignatureAndTransactionErrors, "InvalidIdentifier")

// Checkpoint & Consensus Errors
var ErrPackageVerificationTimeout = NewSuiError(CheckpointAndConsensusErrors, "PackageVerificationTimeout")
var ErrVerifiedCheckpointNotFound = NewSuiError(CheckpointAndConsensusErrors, "VerifiedCheckpointNotFound")
var ErrVerifiedCheckpointDigestNotFound = NewSuiError(CheckpointAndConsensusErrors, "VerifiedCheckpointDigestNotFound")
var ErrLatestCheckpointSequenceNumberNotFound = NewSuiError(CheckpointAndConsensusErrors, "LatestCheckpointSequenceNumberNotFound")
var ErrCheckpointContentsNotFound = NewSuiError(CheckpointAndConsensusErrors, "CheckpointContentsNotFound")
var ErrGenesisTransactionNotFound = NewSuiError(CheckpointAndConsensusErrors, "GenesisTransactionNotFound")
var ErrTransactionCursorNotFound = NewSuiError(CheckpointAndConsensusErrors, "TransactionCursorNotFound")

// Publishing Errors
var ErrDependentPackageNotFound = NewSuiError(PublishingErrors, "DependentPackageNotFound")
var ErrMaxPublishCountExceeded = NewSuiError(PublishingErrors, "MaxPublishCountExceeded")
var ErrMutableParameterExpected = NewSuiError(PublishingErrors, "MutableParameterExpected")
var ErrImmutableParameterExpectedError = NewSuiError(PublishingErrors, "ImmutableParameterExpectedError")
var ErrSizeLimitExceeded = NewSuiError(PublishingErrors, "SizeLimitExceeded")
var ErrAddressDeniedForCoin = NewSuiError(PublishingErrors, "AddressDeniedForCoin")
var ErrCoinTypeGlobalPause = NewSuiError(PublishingErrors, "CoinTypeGlobalPause")

// Soft Bundle Errors
var ErrTooManyTransactionsInSoftBundle = NewSuiError(SoftBundleErrors, "TooManyTransactionsInSoftBundle")
var ErrSoftBundleTooLarge = NewSuiError(SoftBundleErrors, "SoftBundleTooLarge")
var ErrNoSharedObjectError = NewSuiError(SoftBundleErrors, "NoSharedObjectError")
var ErrAlreadyExecutedError = NewSuiError(SoftBundleErrors, "AlreadyExecutedError")
var ErrCertificateAlreadyProcessed = NewSuiError(SoftBundleErrors, "CertificateAlreadyProcessed")

// ========================================
// Error Mapping and Retry Functions
// ========================================

// suiErrorMappings maps raw error substrings to structured Sui errors,
// using the SuiError.Error() value to ensure the actual error string is used.
var suiErrorMappings = []struct {
	substring string
	err       *SuiError
}{
	// Object Errors & Move Call Errors
	{ErrInvalidChildObjectArgument.Error(), ErrInvalidChildObjectArgument},
	{ErrInvalidObjectDigest.Error(), ErrInvalidObjectDigest},
	{ErrInvalidSequenceNumber.Error(), ErrInvalidSequenceNumber},
	{ErrObjectSequenceNumberTooHigh.Error(), ErrObjectSequenceNumberTooHigh},
	{ErrObjectDeleted.Error(), ErrObjectDeleted},
	{ErrInaccessibleSystemObject.Error(), ErrInaccessibleSystemObject},
	{ErrNotSharedObjectError.Error(), ErrNotSharedObjectError},
	{ErrNotOwnedObjectError.Error(), ErrNotOwnedObjectError},
	{ErrDuplicateObjectRefInput.Error(), ErrDuplicateObjectRefInput},
	{ErrMovePackageAsObject.Error(), ErrMovePackageAsObject},
	{ErrMoveObjectAsPackage.Error(), ErrMoveObjectAsPackage},
	{ErrBlockedMoveFunction.Error(), ErrBlockedMoveFunction},
	{ErrTransferObjectWithoutPublicTransferError.Error(), ErrTransferObjectWithoutPublicTransferError},
	{ErrUnsupported.Error(), ErrUnsupported},
	{ErrMoveFunctionInputError.Error(), ErrMoveFunctionInputError},
	{ErrPostRandomCommandRestrictions.Error(), ErrPostRandomCommandRestrictions},

	// Gas Errors
	{ErrMissingGasPayment.Error(), ErrMissingGasPayment},
	{ErrGasObjectNotOwnedObject.Error(), ErrGasObjectNotOwnedObject},
	{ErrGasBudgetTooHigh.Error(), ErrGasBudgetTooHigh},
	{ErrGasBudgetTooLow.Error(), ErrGasBudgetTooLow},
	{ErrGasBalanceTooLow.Error(), ErrGasBalanceTooLow},
	{ErrGasPriceUnderRGP.Error(), ErrGasPriceUnderRGP},
	{ErrGasPriceTooHigh.Error(), ErrGasPriceTooHigh},
	{ErrInvalidGasObject.Error(), ErrInvalidGasObject},
	{ErrInsufficientBalanceToCoverMinimalGas.Error(), ErrInsufficientBalanceToCoverMinimalGas},
	{ErrUnexpectedGasPaymentObject.Error(), ErrUnexpectedGasPaymentObject},
	{ErrGasPriceMismatchError.Error(), ErrGasPriceMismatchError},
	{ErrInsufficientGas.Error(), ErrInsufficientGas},

	// Signature & Transaction Errors
	{ErrIncorrectUserSignature.Error(), ErrIncorrectUserSignature},
	{ErrInvalidBatchTransaction.Error(), ErrInvalidBatchTransaction},
	{ErrTransactionDenied.Error(), ErrTransactionDenied},
	{ErrUnsupportedSponsoredTransactionKind.Error(), ErrUnsupportedSponsoredTransactionKind},
	{ErrEmptyCommandInput.Error(), ErrEmptyCommandInput},
	{ErrEmptyInputCoins.Error(), ErrEmptyInputCoins},
	{ErrInvalidIdentifier.Error(), ErrInvalidIdentifier},

	// Checkpoint & Consensus Errors
	{ErrPackageVerificationTimeout.Error(), ErrPackageVerificationTimeout},
	{ErrVerifiedCheckpointNotFound.Error(), ErrVerifiedCheckpointNotFound},
	{ErrVerifiedCheckpointDigestNotFound.Error(), ErrVerifiedCheckpointDigestNotFound},
	{ErrLatestCheckpointSequenceNumberNotFound.Error(), ErrLatestCheckpointSequenceNumberNotFound},
	{ErrCheckpointContentsNotFound.Error(), ErrCheckpointContentsNotFound},
	{ErrGenesisTransactionNotFound.Error(), ErrGenesisTransactionNotFound},
	{ErrTransactionCursorNotFound.Error(), ErrTransactionCursorNotFound},

	// Publishing Errors
	{ErrDependentPackageNotFound.Error(), ErrDependentPackageNotFound},
	{ErrMaxPublishCountExceeded.Error(), ErrMaxPublishCountExceeded},
	{ErrMutableParameterExpected.Error(), ErrMutableParameterExpected},
	{ErrImmutableParameterExpectedError.Error(), ErrImmutableParameterExpectedError},
	{ErrSizeLimitExceeded.Error(), ErrSizeLimitExceeded},
	{ErrAddressDeniedForCoin.Error(), ErrAddressDeniedForCoin},
	{ErrCoinTypeGlobalPause.Error(), ErrCoinTypeGlobalPause},

	// Soft Bundle Errors
	{ErrTooManyTransactionsInSoftBundle.Error(), ErrTooManyTransactionsInSoftBundle},
	{ErrSoftBundleTooLarge.Error(), ErrSoftBundleTooLarge},
	{ErrNoSharedObjectError.Error(), ErrNoSharedObjectError},
	{ErrAlreadyExecutedError.Error(), ErrAlreadyExecutedError},
	{ErrCertificateAlreadyProcessed.Error(), ErrCertificateAlreadyProcessed},
}

// ParseSuiErrorMessage maps a raw RPC error message to a structured error.
// It iterates over the known substrings in suiErrorMappings. If a substring is found,
// the corresponding error is returned. Otherwise, it returns an "Unknown error".
func ParseSuiErrorMessage(msg string) *SuiError {
	for _, mapping := range suiErrorMappings {
		if strings.Contains(msg, mapping.substring) {
			return mapping.err
		}
	}

	return nil
}

// retryableErrors is a list of errors that are considered transient and retriable.
// Adjust this list as needed according to Sui semantics.
var retryableErrors = []error{
	ErrPackageVerificationTimeout,
	ErrVerifiedCheckpointNotFound,
	ErrVerifiedCheckpointDigestNotFound,
	ErrLatestCheckpointSequenceNumberNotFound,
	ErrCheckpointContentsNotFound,
	ErrGenesisTransactionNotFound,
	ErrTransactionCursorNotFound,
	ErrGasBudgetTooLow,
	ErrGasBudgetTooHigh,
	ErrGasBalanceTooLow,
	ErrGasPriceUnderRGP,
	ErrGasPriceTooHigh,
	ErrInsufficientGas,
}

// IsRetryable determines if a Sui error is retryable (transient).
// It uses errors.Is to correctly recognize wrapped errors.
func IsRetryable(err error) bool {
	for _, retryErr := range retryableErrors {
		if errors.Is(err, retryErr) {
			return true
		}
	}

	return false
}
