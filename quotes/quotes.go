package quotes

import "time"

func GetQuote(date time.Time) string {
	switch date.Month() {
	case time.January:
		return january[date.Day()-1]
	case time.February:
		return february[date.Day()-1]
	case time.March:
		return march[date.Day()-1]
	case time.April:
		return april[date.Day()-1]
	case time.May:
		return may[date.Day()-1]
	case time.June:
		return june[date.Day()-1]
	case time.July:
		return july[date.Day()-1]
	case time.August:
		return august[date.Day()-1]
	case time.September:
		return september[date.Day()-1]
	case time.October:
		return october[date.Day()-1]
	case time.November:
		return november[date.Day()-1]
	case time.December:
		return december[date.Day()-1]
	}
	return ""
}
