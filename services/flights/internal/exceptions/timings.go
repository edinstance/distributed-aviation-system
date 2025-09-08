package exceptions

import "fmt"

var ErrInvalidTimes = fmt.Errorf("arrival must be after departure")
