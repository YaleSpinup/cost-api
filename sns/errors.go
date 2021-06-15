package sns

import (
	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ErrCode(msg string, err error) error {
	if aerr, ok := errors.Cause(err).(awserr.Error); ok {
		switch aerr.Code() {
		case

			// ErrCodeAuthorizationErrorException for service response error code
			// "AuthorizationError".
			//
			// Indicates that the user has been denied access to the requested resource.
			sns.ErrCodeAuthorizationErrorException,

			// ErrCodeInvalidSecurityException for service response error code
			// "InvalidSecurity".
			//
			// The credential signature isn't valid. You must use an HTTPS endpoint and
			// sign your request using Signature Version 4.
			sns.ErrCodeInvalidSecurityException,

			// ErrCodeKMSAccessDeniedException for service response error code
			// "KMSAccessDenied".
			//
			// The ciphertext references a key that doesn't exist or that you don't have
			// access to.
			sns.ErrCodeKMSAccessDeniedException:

			return apierror.New(apierror.ErrForbidden, msg, aerr)
		case

			// ErrCodeNotFoundException for service response error code
			// "NotFound".
			//
			// Indicates that the requested resource does not exist.
			sns.ErrCodeNotFoundException,

			// ErrCodeResourceNotFoundException for service response error code
			// "ResourceNotFound".
			//
			// Can’t perform the action on the specified resource. Make sure that the
			// resource exists.
			sns.ErrCodeResourceNotFoundException:

			return apierror.New(apierror.ErrNotFound, msg, aerr)

		case

			// ErrCodeConcurrentAccessException for service response error code
			// "ConcurrentAccess".
			//
			// Can't perform multiple operations on a tag simultaneously. Perform the operations
			// sequentially.
			sns.ErrCodeConcurrentAccessException:

			return apierror.New(apierror.ErrConflict, msg, aerr)

		case

			// ErrCodeInvalidParameterException for service response error code
			// "InvalidParameter".
			//
			// Indicates that a request parameter does not comply with the associated constraints.
			sns.ErrCodeInvalidParameterException,

			// ErrCodeInvalidParameterValueException for service response error code
			// "ParameterValueInvalid".
			//
			// Indicates that a request parameter does not comply with the associated constraints.
			sns.ErrCodeInvalidParameterValueException,

			// ErrCodeKMSDisabledException for service response error code
			// "KMSDisabled".
			//
			// The request was rejected because the specified customer master key (CMK)
			// isn't enabled.
			sns.ErrCodeKMSDisabledException,

			// ErrCodeKMSInvalidStateException for service response error code
			// "KMSInvalidState".
			//
			// The request was rejected because the state of the specified resource isn't
			// valid for this request. For more information, see How Key State Affects Use
			// of a Customer Master Key (https://docs.aws.amazon.com/kms/latest/developerguide/key-state.html)
			// in the AWS Key Management Service Developer Guide.
			sns.ErrCodeKMSInvalidStateException,

			// ErrCodeKMSNotFoundException for service response error code
			// "KMSNotFound".
			//
			// The request was rejected because the specified entity or resource can't be
			// found.
			sns.ErrCodeKMSNotFoundException,

			// ErrCodeKMSOptInRequired for service response error code
			// "KMSOptInRequired".
			//
			// The AWS access key ID needs a subscription for the service.
			sns.ErrCodeKMSOptInRequired,

			// ErrCodeOptedOutException for service response error code
			// "OptedOut".
			//
			// Indicates that the specified phone number opted out of receiving SMS messages
			// from your AWS account. You can't send SMS messages to phone numbers that
			// opt out.
			sns.ErrCodeOptedOutException,

			// ErrCodePlatformApplicationDisabledException for service response error code
			// "PlatformApplicationDisabled".
			//
			// Exception error indicating platform application disabled.
			sns.ErrCodePlatformApplicationDisabledException,

			// ErrCodeStaleTagException for service response error code
			// "StaleTag".
			//
			// A tag has been added to a resource with the same ARN as a deleted resource.
			// Wait a short while and then retry the operation.
			sns.ErrCodeStaleTagException,

			// ErrCodeEndpointDisabledException for service response error code
			// "EndpointDisabled".
			//
			// Exception error indicating endpoint disabled.
			sns.ErrCodeEndpointDisabledException,

			// ErrCodeTagPolicyException for service response error code
			// "TagPolicy".
			//
			// The request doesn't comply with the IAM tag policy. Correct your request
			// and then retry it.
			sns.ErrCodeTagPolicyException,

			// ErrCodeUserErrorException for service response error code
			// "UserError".
			//
			// Indicates that a request parameter does not comply with the associated constraints.
			sns.ErrCodeUserErrorException,

			// ErrCodeValidationException for service response error code
			// "ValidationException".
			//
			// Indicates that a parameter in the request is invalid.
			sns.ErrCodeValidationException,

			// ErrCodeVerificationException for service response error code
			// "VerificationException".
			//
			// Indicates that the one-time password (OTP) used for verification is invalid.
			sns.ErrCodeVerificationException:

			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		case

			// ErrCodeFilterPolicyLimitExceededException for service response error code
			// "FilterPolicyLimitExceeded".
			//
			// Indicates that the number of filter polices in your AWS account exceeds the
			// limit. To add more filter polices, submit an sns.ErrCodeLimit Increase case in the
			// AWS Support Center.
			sns.ErrCodeFilterPolicyLimitExceededException,

			// ErrCodeKMSThrottlingException for service response error code
			// "KMSThrottling".
			//
			// The request was denied due to request throttling. For more information about
			// throttling, see Limits (https://docs.aws.amazon.com/kms/latest/developerguide/limits.html#requests-per-second)
			// in the AWS Key Management Service Developer Guide.
			sns.ErrCodeKMSThrottlingException,

			// ErrCodeSubscriptionLimitExceededException for service response error code
			// "SubscriptionLimitExceeded".
			//
			// Indicates that the customer already owns the maximum allowed number of subscriptions.
			sns.ErrCodeSubscriptionLimitExceededException,

			// ErrCodeTagLimitExceededException for service response error code
			// "TagLimitExceeded".
			//
			// Can't add more than 50 tags to a topic.
			sns.ErrCodeTagLimitExceededException,

			// ErrCodeThrottledException for service response error code
			// "Throttled".
			//
			// Indicates that the rate at which requests have been submitted for this action
			// exceeds the limit for your account.
			sns.ErrCodeThrottledException,

			// ErrCodeTopicLimitExceededException for service response error code
			// "TopicLimitExceeded".
			//
			// Indicates that the customer already owns the maximum allowed number of topics.
			sns.ErrCodeTopicLimitExceededException:

			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		case

			// ErrCodeInternalErrorException for service response error code
			// "InternalError".
			//
			// Indicates an internal service error.
			sns.ErrCodeInternalErrorException:

			return apierror.New(apierror.ErrServiceUnavailable, msg, aerr)
		default:
			m := msg + ": " + aerr.Message()
			return apierror.New(apierror.ErrBadRequest, m, aerr)
		}
	}

	log.Warnf("uncaught error: %s, returning Internal Server Error", err)
	return apierror.New(apierror.ErrInternalError, msg, err)
}
