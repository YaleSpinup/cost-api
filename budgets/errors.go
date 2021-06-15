package budgets

import (
	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/budgets"
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
			// You are not authorized to use this operation with the given parameters.
			budgets.ErrCodeAccessDeniedException:

			return apierror.New(apierror.ErrForbidden, msg, aerr)
		case

			// ErrCodeNotFoundException for service response error code
			// "NotFoundException".
			//
			// We can’t locate the resource that you specified.
			budgets.ErrCodeNotFoundException:

			return apierror.New(apierror.ErrNotFound, msg, aerr)

		case

			// ErrCodeCreationLimitExceededException for service response error code
			// "CreationLimitExceededException".
			//
			// You've exceeded the notification or subscriber limit.
			budgets.ErrCodeCreationLimitExceededException,

			// ErrCodeDuplicateRecordException for service response error code
			// "DuplicateRecordException".
			//
			// The budget name already exists. Budget names must be unique within an account.
			budgets.ErrCodeDuplicateRecordException:

			return apierror.New(apierror.ErrConflict, msg, aerr)

		case

			// ErrCodeExpiredNextTokenException for service response error code
			// "ExpiredNextTokenException".
			//
			// The pagination token expired.
			budgets.ErrCodeExpiredNextTokenException,

			// ErrCodeInvalidNextTokenException for service response error code
			// "InvalidNextTokenException".
			//
			// The pagination token is invalid.
			budgets.ErrCodeInvalidNextTokenException,

			// ErrCodeInvalidParameterException for service response error code
			// "InvalidParameterException".
			//
			// An error on the client occurred. Typically, the cause is an invalid input
			// value.
			budgets.ErrCodeInvalidParameterException,

			// ErrCodeResourceLockedException for service response error code
			// "ResourceLockedException".
			//
			// The request was received and recognized by the server, but the server rejected
			// that particular method for the requested resource.
			budgets.ErrCodeResourceLockedException:

			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		case
			"LimitExceeded":

			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		case

			// ErrCodeInternalErrorException for service response error code
			// "InternalErrorException".
			//
			// An error on the server occurred during the processing of your request. Try
			// again later.
			budgets.ErrCodeInternalErrorException:

			return apierror.New(apierror.ErrServiceUnavailable, msg, aerr)
		default:
			m := msg + ": " + aerr.Message()
			return apierror.New(apierror.ErrBadRequest, m, aerr)
		}
	}

	log.Warnf("uncaught error: %s, returning Internal Server Error", err)
	return apierror.New(apierror.ErrInternalError, msg, err)
}
