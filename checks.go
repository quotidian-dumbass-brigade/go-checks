package checks

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	Red        = "\033[31m"
	BoldRed    = "\033[1;31m"
	BoldYellow = "\033[1;33m"
	Blue       = "\033[34m"
	BoldCyan   = "\033[1;36m"
	Reset      = "\033[0m"
	// Black       = "\033[30m"
	// BoldBlack   = "\033[1;30m"
	// Green       = "\033[32m"
	// BoldGreen   = "\033[1;32m"
	// Yellow      = "\033[33m"
	// BoldBlue    = "\033[1;34m"
	// Magenta     = "\033[35m"
	// BoldMagenta = "\033[1;35m"
	// Cyan        = "\033[36m"
	// White     = "\033[37m"
	// BoldWhite = "\033[1;37m"
)

type Check struct {
	Failed   bool
	FailedBy []*Check
	Value    interface{}
	Rule     string
	Result   string
}

// Create a new check
func New(rule string) *Check {
	return &Check{Rule: rule}
}

// Mark this check as failed
func (check *Check) Fail(value interface{}) *Check {
	check.Failed = true
	check.Value = value

	caller := getCallerName()

	check.Result = fmt.Sprintf(
		"%sfailed check%s in %s%s%s: %s`%s`%s, %sgot:%s %s`%v`%s",
		BoldRed, Reset, BoldCyan, caller, Reset, Blue, check.Rule, Reset, BoldRed, Reset, Red, check.Value, Reset,
	)
	return check
}

// Propagate multiple failures correctly
func (check *Check) FailBy(others ...*Check) *Check {
	caller := getCallerName()

	for _, other := range others {
		if other.Failed {
			check.Failed = true
			check.FailedBy = append(check.FailedBy, other)
		}
	}
	check.Result = fmt.Sprintf("%sfailed check%s in %s%s%s: %s`%s`%s",
		BoldRed, Reset, BoldCyan, caller, Reset, Blue, check.Rule, Reset)
	return check
}

// Get the name of the calling function where a check failed
func getCallerName() string {
	pc, _, _, ok := runtime.Caller(2) // 2 = Up 2 levels
	callerName := "unknown"
	if ok {
		callerName = runtime.FuncForPC(pc).Name()
	}
	return callerName
}

// Print a structured failure report
func (check *Check) Blame(message string) {
	fmt.Printf("%s%s%s\n%s\n\n", BoldYellow, message, Reset, check.formatResult(0))
}

// Proper recursive formatting with indentation
func (check *Check) formatResult(indent int) string {
	indentation := strings.Repeat("    ", indent) // 4 spaces per level
	result := fmt.Sprintf("%s↳ %s", indentation, check.Result)

	for _, failedCheck := range check.FailedBy {
		result += "\n" + failedCheck.formatResult(indent+1) // Increase indentation
	}
	return result
}

// Create an empty passing check
func Pass() *Check {
	return &Check{}
}
