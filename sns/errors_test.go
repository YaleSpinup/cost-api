package sns

import (
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
)

func TestErrCode(t *testing.T) {
	apiErrorTestCases := map[string]string{
		"":                                     apierror.ErrBadRequest,
		sns.ErrCodeAuthorizationErrorException: apierror.ErrForbidden,
		sns.ErrCodeInvalidSecurityException:    apierror.ErrForbidden,
		sns.ErrCodeKMSAccessDeniedException:    apierror.ErrForbidden,

		sns.ErrCodeNotFoundException:         apierror.ErrNotFound,
		sns.ErrCodeResourceNotFoundException: apierror.ErrNotFound,

		sns.ErrCodeConcurrentAccessException: apierror.ErrConflict,

		sns.ErrCodeInvalidParameterException:            apierror.ErrBadRequest,
		sns.ErrCodeInvalidParameterValueException:       apierror.ErrBadRequest,
		sns.ErrCodeKMSDisabledException:                 apierror.ErrBadRequest,
		sns.ErrCodeKMSInvalidStateException:             apierror.ErrBadRequest,
		sns.ErrCodeKMSNotFoundException:                 apierror.ErrBadRequest,
		sns.ErrCodeKMSOptInRequired:                     apierror.ErrBadRequest,
		sns.ErrCodeOptedOutException:                    apierror.ErrBadRequest,
		sns.ErrCodePlatformApplicationDisabledException: apierror.ErrBadRequest,
		sns.ErrCodeStaleTagException:                    apierror.ErrBadRequest,
		sns.ErrCodeEndpointDisabledException:            apierror.ErrBadRequest,
		sns.ErrCodeTagPolicyException:                   apierror.ErrBadRequest,
		sns.ErrCodeUserErrorException:                   apierror.ErrBadRequest,
		sns.ErrCodeValidationException:                  apierror.ErrBadRequest,
		sns.ErrCodeVerificationException:                apierror.ErrBadRequest,

		sns.ErrCodeFilterPolicyLimitExceededException: apierror.ErrLimitExceeded,
		sns.ErrCodeKMSThrottlingException:             apierror.ErrLimitExceeded,
		sns.ErrCodeSubscriptionLimitExceededException: apierror.ErrLimitExceeded,
		sns.ErrCodeTagLimitExceededException:          apierror.ErrLimitExceeded,
		sns.ErrCodeThrottledException:                 apierror.ErrLimitExceeded,
		sns.ErrCodeTopicLimitExceededException:        apierror.ErrLimitExceeded,

		sns.ErrCodeInternalErrorException: apierror.ErrServiceUnavailable,
	}

	for awsErr, apiErr := range apiErrorTestCases {
		err := ErrCode("test error", awserr.New(awsErr, awsErr, nil))
		if aerr, ok := errors.Cause(err).(apierror.Error); ok {
			t.Logf("got apierror '%s'", aerr)
		} else {
			t.Errorf("expected cloudwatch error %s to be an apierror.Error %s, got %s", awsErr, apiErr, err)
		}
	}

	err := ErrCode("test error", errors.New("Unknown"))
	if aerr, ok := errors.Cause(err).(apierror.Error); ok {
		t.Logf("got apierror '%s'", aerr)
	} else {
		t.Errorf("expected unknown error to be an apierror.ErrInternalError, got %s", err)
	}
}
