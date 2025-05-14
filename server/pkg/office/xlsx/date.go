package xlsx

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// nanoSecondsPerDay defines a constant value of the number of nanoseconds in a day.
const nanoSecondsPerDay = 24 * 60 * 60 * 1000 * 1000 * 1000

// excelEpoch specifies the epoch of all excel dates. Dates are internally represented
// as relative to this value.
var excelEpoch = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)

// convertExcelDateToDateString takes an excel numeric representation of a date, and
// converts it to a human-readable RFC3339 or ISO formatted string.
//
// Excel dates are stored within excel as a signed floating point number.
// The integer portion determines the number of days ahead of 30/12/1899 the date is.
// The portion after the decimal point represents the proportion through the day.
// For example, 6am would be 1/4 of the way through a 24hr day, so it is stored as 0.25.
func convertExcelDateToDateString(value string) (string, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", fmt.Errorf("unable to parse date float value: %w", err)
	}

	numberOfDays := math.Trunc(floatValue)
	numberOfNanoSeconds := (floatValue - numberOfDays) * nanoSecondsPerDay

	actualTime := excelEpoch.AddDate(0, 0, int(numberOfDays)).Add(time.Duration(numberOfNanoSeconds))

	formatString := time.RFC3339
	if floatValue == numberOfDays {
		// We are dealing with a date, and not a datetime
		formatString = "2006-01-02"
	}

	return actualTime.Format(formatString), nil
}
