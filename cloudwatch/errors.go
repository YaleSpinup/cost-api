package cloudwatch

import (
	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ErrCode(msg string, err error) error {
	if aerr, ok := errors.Cause(err).(awserr.Error); ok {
		switch aerr.Code() {
		case
			// Access denied.
			"AccessDenied",

			// There is a problem with your AWS account that prevents the operation from completing successfully.
			"AccountProblem",

			// Access forbidden.
			"Forbidden",

			// The AWS access key ID you provided does not exist in our records.
			"InvalidAccessKeyId":

			return apierror.New(apierror.ErrForbidden, msg, aerr)
		case
			// The specified bucket does not exist.
			"NotFound",

			// ErrCodeResourceNotFound for service response error code
			// "ResourceNotFound".
			//
			// The named resource does not exist.
			cloudwatch.ErrCodeResourceNotFound,

			// ErrCodeResourceNotFoundException for service response error code
			// "ResourceNotFoundException".
			//
			// The named resource does not exist.
			cloudwatch.ErrCodeResourceNotFoundException:

			return apierror.New(apierror.ErrNotFound, msg, aerr)

		case
			// ErrCodeConcurrentModificationException for service response error code
			// "ConcurrentModificationException".
			//
			// More than one process tried to modify a resource at the same time.
			cloudwatch.ErrCodeConcurrentModificationException:

			return apierror.New(apierror.ErrConflict, msg, aerr)

		case
			// ErrCodeDashboardInvalidInputError for service response error code
			// "InvalidParameterInput".
			//
			// Some part of the dashboard data is invalid.
			cloudwatch.ErrCodeDashboardInvalidInputError,

			// ErrCodeInvalidFormatFault for service response error code
			// "InvalidFormat".
			//
			// Data was not syntactically valid JSON.
			cloudwatch.ErrCodeInvalidFormatFault,

			// ErrCodeInvalidNextToken for service response error code
			// "InvalidNextToken".
			//
			// The next token specified is invalid.
			cloudwatch.ErrCodeInvalidNextToken,

			// ErrCodeInvalidParameterCombinationException for service response error code
			// "InvalidParameterCombination".
			//
			// Parameters were used together that cannot be used together.
			cloudwatch.ErrCodeInvalidParameterCombinationException,

			// ErrCodeInvalidParameterValueException for service response error code
			// "InvalidParameterValue".
			//
			// The value of an input parameter is bad or out-of-range.
			cloudwatch.ErrCodeInvalidParameterValueException,

			// ErrCodeMissingRequiredParameterException for service response error code
			// "MissingParameter".
			//
			// An input parameter that is required is missing.
			cloudwatch.ErrCodeMissingRequiredParameterException:

			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		case
			// ErrCodeLimitExceededException for service response error code
			// "LimitExceededException".
			//
			// The operation exceeded one or more limits.
			cloudwatch.ErrCodeLimitExceededException,

			// ErrCodeLimitExceededFault for service response error code
			// "LimitExceeded".
			//
			// The quota for alarms for this customer has already been reached.
			cloudwatch.ErrCodeLimitExceededFault:

			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		case
			// ErrCodeInternalServiceFault for service response error code
			// "InternalServiceError".
			//
			// Request processing has failed due to some unknown error, exception, or failure.
			cloudwatch.ErrCodeInternalServiceFault:

			return apierror.New(apierror.ErrServiceUnavailable, msg, aerr)
		default:
			m := msg + ": " + aerr.Message()
			return apierror.New(apierror.ErrBadRequest, m, aerr)
		}
	}

	log.Warnf("uncaught error: %s, returning Internal Server Error", err)
	return apierror.New(apierror.ErrInternalError, msg, err)
}
