package midtrans

import "errors"

var ErrMarshaller = errors.New("error marshalling request json")
var ErrUnMarshaller = errors.New("error un marshalling request json")
var ErrApiCall = errors.New("error calling api")
var ErrRateLimitExceeded = errors.New("rate limit exceeded")
var ErrInternalServerPaymentGatewayError = errors.New("internal server payment gateway error")
