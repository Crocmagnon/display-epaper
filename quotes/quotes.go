package quotes

import "time"

func GetQuote(date time.Time) string {
	switch date.Month() {
	case time.January:
		return january[date.Day()]
	case time.February:
		return february[date.Day()]
	case time.March:
		return march[date.Day()]
	case time.April:
		return april[date.Day()]
	case time.May:
		return may[date.Day()]
	case time.June:
		return june[date.Day()]
	case time.July:
		return july[date.Day()]
	case time.August:
		return august[date.Day()]
	case time.September:
		return september[date.Day()]
	case time.October:
		return october[date.Day()]
	case time.November:
		return november[date.Day()]
	case time.December:
		return december[date.Day()]
	}
	return ""
}
