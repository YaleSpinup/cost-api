package costexplorer

import (
	"github.com/YaleSpinup/cost-api/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/costexplorer"
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

			// ErrCodeDataUnavailableException for service response error code
			// "DataUnavailableException".
			//
			// The requested data is unavailable.
			costexplorer.ErrCodeDataUnavailableException:
			return apierror.New(apierror.ErrNotFound, msg, aerr)

		case
			// ErrCodeBillExpirationException for service response error code
			// "BillExpirationException".
			//
			// The requested report expired. Update the date interval and try again.
			costexplorer.ErrCodeBillExpirationException,

			// ErrCodeInvalidNextTokenException for service response error code
			// "InvalidNextTokenException".
			//
			// The pagination token is invalid. Try again without a pagination token.
			costexplorer.ErrCodeInvalidNextTokenException,

			// ErrCodeRequestChangedException for service response error code
			// "RequestChangedException".
			//
			// Your request parameters changed between pages. Try again with the old parameters
			// or without a pagination token.
			costexplorer.ErrCodeRequestChangedException,

			// ErrCodeUnresolvableUsageUnitException for service response error code
			// "UnresolvableUsageUnitException".
			//
			// Cost Explorer was unable to identify the usage unit. Provide UsageType/UsageTypeGroup
			// filter selections that contain matching units, for example: hours.
			costexplorer.ErrCodeUnresolvableUsageUnitException:

			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		case
			// Reduce your request rate.
			"SlowDown",

			// ErrCodeLimitExceededException for service response error code
			// "LimitExceededException".
			//
			// You made too many calls in a short period of time. Try again later.
			costexplorer.ErrCodeLimitExceededException:
			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		case
			// We encountered an internal error. Please try again.
			"InternalError",

			// A header you provided implies functionality that is not implemented.
			"NotImplemented",

			// Your socket connection to the server was not read from or written to within the timeout period.
			"RequestTimeout",

			// You are being redirected to the bucket while DNS updates.
			"TemporaryRedirect":
			return apierror.New(apierror.ErrServiceUnavailable, msg, aerr)
		default:
			m := msg + ": " + aerr.Message()
			return apierror.New(apierror.ErrBadRequest, m, aerr)
		}
	}

	log.Warnf("uncaught error: %s, returning Internal Server Error", err)
	return apierror.New(apierror.ErrInternalError, msg, err)
}
