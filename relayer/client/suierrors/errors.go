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

// Object Errors
var ErrMutableObjectUsedMoreThanOnce = errors.New("MutableObjectUsedMoreThanOnce")
var ErrObjectInputArityViolation = errors.New("ObjectInputArityViolation")
var ErrObjectNotFound = errors.New("ObjectNotFound")
var ErrObjectVersionUnavailableForConsumption = errors.New("ObjectVersionUnavailableForConsumption")
var ErrInvalidChildObjectArgument = errors.New("InvalidChildObjectArgument")
var ErrInvalidObjectDigest = errors.New("InvalidObjectDigest")
var ErrInvalidSequenceNumber = errors.New("InvalidSequenceNumber")
var ErrObjectSequenceNumberTooHigh = errors.New("ObjectSequenceNumberTooHigh")
var ErrObjectDeleted = errors.New("ObjectDeleted")
var ErrInaccessibleSystemObject = errors.New("InaccessibleSystemObject")
var ErrNotSharedObjectError = errors.New("NotSharedObjectError")
var ErrNotOwnedObjectError = errors.New("NotOwnedObjectError")
var ErrDuplicateObjectRefInput = errors.New("DuplicateObjectRefInput")

// Move Call Errors
var ErrMovePackageAsObject = errors.New("MovePackageAsObject")
var ErrMoveObjectAsPackage = errors.New("MoveObjectAsPackage")
var ErrBlockedMoveFunction = errors.New("BlockedMoveFunction")
var ErrTransferObjectWithoutPublicTransferError = errors.New("TransferObjectWithoutPublicTransferError")
var ErrUnsupported = errors.New("Unsupported")
var ErrMoveFunctionInputError = errors.New("MoveFunctionInputError")
var ErrPostRandomCommandRestrictions = errors.New("PostRandomCommandRestrictions")

// Gas Errors
var ErrMissingGasPayment = errors.New("MissingGasPayment")
var ErrGasObjectNotOwnedObject = errors.New("GasObjectNotOwnedObject")
var ErrGasBudgetTooHigh = errors.New("GasBudgetTooHigh")
var ErrGasBudgetTooLow = errors.New("GasBudgetTooLow")
var ErrGasBalanceTooLow = errors.New("GasBalanceTooLow")
var ErrGasPriceUnderRGP = errors.New("GasPriceUnderRGP")
var ErrGasPriceTooHigh = errors.New("GasPriceTooHigh")
var ErrInvalidGasObject = errors.New("InvalidGasObject")
var ErrInsufficientBalanceToCoverMinimalGas = errors.New("InsufficientBalanceToCoverMinimalGas")
var ErrUnexpectedGasPaymentObject = errors.New("UnexpectedGasPaymentObject")
var ErrGasPriceMismatchError = errors.New("GasPriceMismatchError")

// Signature & Transaction Errors
var ErrIncorrectUserSignature = errors.New("IncorrectUserSignature")
var ErrInvalidBatchTransaction = errors.New("InvalidBatchTransaction")
var ErrTransactionDenied = errors.New("TransactionDenied")
var ErrUnsupportedSponsoredTransactionKind = errors.New("UnsupportedSponsoredTransactionKind")
var ErrEmptyCommandInput = errors.New("EmptyCommandInput")
var ErrEmptyInputCoins = errors.New("EmptyInputCoins")
var ErrInvalidIdentifier = errors.New("InvalidIdentifier")

// Checkpoint & Consensus Errors
var ErrPackageVerificationTimeout = errors.New("PackageVerificationTimeout")
var ErrVerifiedCheckpointNotFound = errors.New("VerifiedCheckpointNotFound")
var ErrVerifiedCheckpointDigestNotFound = errors.New("VerifiedCheckpointDigestNotFound")
var ErrLatestCheckpointSequenceNumberNotFound = errors.New("LatestCheckpointSequenceNumberNotFound")
var ErrCheckpointContentsNotFound = errors.New("CheckpointContentsNotFound")
var ErrGenesisTransactionNotFound = errors.New("GenesisTransactionNotFound")
var ErrTransactionCursorNotFound = errors.New("TransactionCursorNotFound")

// Publishing Errors
var ErrDependentPackageNotFound = errors.New("DependentPackageNotFound")
var ErrMaxPublishCountExceeded = errors.New("MaxPublishCountExceeded")
var ErrMutableParameterExpected = errors.New("MutableParameterExpected")
var ErrImmutableParameterExpectedError = errors.New("ImmutableParameterExpectedError")
var ErrSizeLimitExceeded = errors.New("SizeLimitExceeded")
var ErrAddressDeniedForCoin = errors.New("AddressDeniedForCoin")
var ErrCoinTypeGlobalPause = errors.New("CoinTypeGlobalPause")

// Soft Bundle Errors
var ErrTooManyTransactionsInSoftBundle = errors.New("TooManyTransactionsInSoftBundle")
var ErrSoftBundleTooLarge = errors.New("SoftBundleTooLarge")
var ErrNoSharedObjectError = errors.New("NoSharedObjectError")
var ErrAlreadyExecutedError = errors.New("AlreadyExecutedError")
var ErrCertificateAlreadyProcessed = errors.New("CertificateAlreadyProcessed")

// Error mappings from Sui raw error substrings to structured Go errors.
// The underlying errors are based on the definitions found in:
//
//	https://github.com/MystenLabs/sui/blob/main/crates/sui-types/src/error.rs
var suiErrorMappings = []struct {
	substring string
	err       error
}{
	// Object Errors & Move Call Errors (add others as needed)
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
// It iterates over known substrings from Sui and returns the corresponding error.
// If no mapping is found, it returns an "Unknown error" with the original message.
func ParseSuiErrorMessage(msg string) error {
	for _, mapping := range suiErrorMappings {
		if strings.Contains(msg, mapping.substring) {
			return mapping.err
		}
	}

	return errors.New("Unknown error: " + msg)
}

// List of errors that are considered retryable.
// Adjust this list as needed based on Sui semantics.
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
}

// IsRetryable determines if a Sui error is transient or can be retried.
// It now uses errors.Is to correctly recognize wrapped errors.
func IsRetryable(err error) bool {
	for _, retryErr := range retryableErrors {
		if errors.Is(err, retryErr) {
			return true
		}
	}

	return false
}
