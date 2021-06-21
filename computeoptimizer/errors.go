package computeoptimizer

import (
	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/computeoptimizer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ErrCode(msg string, err error) error {
	if aerr, ok := errors.Cause(err).(awserr.Error); ok {
		switch aerr.Code() {
		case
			// ErrCodeAccessDeniedException for service response error code
			// "AccessDeniedException".
			//
			// You do not have sufficient access to perform this action.
			computeoptimizer.ErrCodeAccessDeniedException,

			// ErrCodeMissingAuthenticationToken for service response error code
			// "MissingAuthenticationToken".
			//
			// The request must contain either a valid (registered) AWS access key ID or
			// X.509 certificate.
			computeoptimizer.ErrCodeMissingAuthenticationToken:

			return apierror.New(apierror.ErrForbidden, msg, aerr)
		case

			// ErrCodeResourceNotFoundException for service response error code
			// "ResourceNotFoundException".
			//
			// A resource that is required for the action doesn't exist.
			computeoptimizer.ErrCodeResourceNotFoundException:

			return apierror.New(apierror.ErrNotFound, msg, aerr)

		case

			"Conflict":

			return apierror.New(apierror.ErrConflict, msg, aerr)

		case

			// ErrCodeInvalidParameterValueException for service response error code
			// "InvalidParameterValueException".
			//
			// An invalid or out-of-range value was supplied for the input parameter.
			computeoptimizer.ErrCodeInvalidParameterValueException,

			// ErrCodeOptInRequiredException for service response error code
			// "OptInRequiredException".
			//
			// The account is not opted in to AWS Compute Optimizer.
			computeoptimizer.ErrCodeOptInRequiredException,

			// ErrCodeServiceUnavailableException for service response error code
			// "ServiceUnavailableException".
			//
			// The request has failed due to a temporary failure of the server.
			computeoptimizer.ErrCodeServiceUnavailableException:

			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		case

			// ErrCodeLimitExceededException for service response error code
			// "LimitExceededException".
			//
			// The request exceeds a limit of the service.
			computeoptimizer.ErrCodeLimitExceededException,

			// ErrCodeThrottlingException for service response error code
			// "ThrottlingException".
			//
			// The request was denied due to request throttling.
			computeoptimizer.ErrCodeThrottlingException:

			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		case

			// ErrCodeInternalServerException for service response error code
			// "InternalServerException".
			//
			// An internal error has occurred. Try your call again.
			computeoptimizer.ErrCodeInternalServerException:

			return apierror.New(apierror.ErrServiceUnavailable, msg, aerr)
		default:
			m := msg + ": " + aerr.Message()
			return apierror.New(apierror.ErrBadRequest, m, aerr)
		}
	}

	log.Warnf("uncaught error: %s, returning Internal Server Error", err)
	return apierror.New(apierror.ErrInternalError, msg, err)
}
