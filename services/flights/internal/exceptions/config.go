package exceptions

import "errors"

var ErrDownstreamClientDown = errors.New("an external service is down")
