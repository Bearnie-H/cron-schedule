// Package cronschedule implements a basic parser for purely numeric Cron timecodes.
// This package will parse a string-represented Cron timecode into a set of integer arrays, correponding to the times at which the schedule should be executed.
// The aim of this package is to provide a simple interface for adding Cron-like scheduling to other projects without the hassle of parsing timecodes more than once.
package cronschedule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Range definitions for the Time fields of a cron entry
const (
	MinuteMinimum int = 0
	MinuteMaximum int = 59

	HourMinimum int = 0
	HourMaximum int = 23

	DayOfMonthMinimum int = 1
	DayOfMonthMaximum int = 31

	MonthMinimum int = 1
	MonthMaximum int = 12

	DayOfWeekMinimum int = 0
	DayOfWeekMaximum int = 6
)

// Don't mind the EBNF.  This is in-place to define the allowable Cron codes which Dune will implement

/*
	Minute Timecode can be of the format:

	MinuteLiteral:		["0" ... "5"], "0"..."9"
	MinuteRange:		Literal, "-", Literal
	MinuteStepRange:	Range, "/", Literal
	MinuteTimeCode:		( Literal | Range | Step | StepRange ) { "," ( Literal | Range | Step | StepRange ) }
*/

/*
	Hour Timecode can be of the format:

	HourLiteral:		( ["0" ... "1"], "0"..."9" ) | ( "2", "0"..."3" )
	HourRange:			Literal, "-", Literal
	HourStepRange:		Range, "/", Literal
	HourTimeCode:		( Literal | Range | Step | StepRange ) { "," ( Literal | Range | Step | StepRange ) }
*/

/*
	DayOfMonth Timecode can be of the format:

	DayOfMonthLiteral:		( "0", "1"..."9" ) | ( ["1" ... "2"], "0"..."9" ) | ( "3", "0"..."1" )
	DayOfMonthRange:		Literal, "-", Literal
	DayOfMonthStepRange:	Range, "/", Literal
	DayOfMonthTimeCode:		( Literal | Range | Step | StepRange ) { "," ( Literal | Range | Step | StepRange ) }
*/

/*
	Month Timecode can be of the format:

	MonthLiteral:		( "0", "1"..."9" ) | ( "1", "0"..."2" )
	MonthRange:			Literal, "-", Literal
	MonthStepRange:		Range, "/", Literal
	MonthTimeCode:		( Literal | Range | Step | StepRange ) { "," ( Literal | Range | Step | StepRange ) }
*/

/*
	DayOfWeek Timecode can be of the format:

	DayOfWeekLiteral:		"0"..."6"
	DayOfWeekRange:			Literal, "-", Literal
	DayOfWeekStepRange:		Range, "/", Literal
	DayOfWeekTimeCode:		( Literal | Range | Step | StepRange ) { "," ( Literal | Range | Step | StepRange ) }
*/

/*
	Full Cron Timecode is fully defined as:

	CronTimeCode:	MinuteTimeCode, " ", HourTimeCode, " ", DayOfMonthTimeCode, " ", MonthTimeCode, " ", DayOfWeekTimeCode
*/

// ParseSchedule will convert a single Cron Timecode string into a set of integer arrays corresponding to:
//	a) Minutes
//	b) Hours
//	c) Days Of Month
//	d) Months
//  e) Days of Week
func ParseSchedule(Code string) (Schedule [5][]int, err error) {

	Fields := strings.Split(Code, " ")
	if len(Fields) != 5 {
		return [5][]int{}, errors.New("cron-schedule error - invalid timecode - Must be 5 whitespace-delimited fields")
	}

	if Schedule[0], err = ParseTimeCode(Fields[0], MinuteMinimum, MinuteMaximum); err != nil {
		return [5][]int{}, err
	}

	if Schedule[1], err = ParseTimeCode(Fields[1], HourMinimum, HourMaximum); err != nil {
		return [5][]int{}, err
	}

	if Schedule[2], err = ParseTimeCode(Fields[2], DayOfMonthMinimum, DayOfMonthMaximum); err != nil {
		return [5][]int{}, err
	}

	if Schedule[3], err = ParseTimeCode(Fields[3], MonthMinimum, MonthMaximum); err != nil {
		return [5][]int{}, err
	}

	if Schedule[4], err = ParseTimeCode(Fields[4], DayOfWeekMinimum, DayOfWeekMaximum); err != nil {
		return [5][]int{}, err
	}

	return Schedule, nil
}

// ParseTimeCode is the full parser for a single element of the timecode.  This will parse a single timecode into an array of corresponding matching times, as well as indicating if this is a valid timecode.
func ParseTimeCode(Code string, Min, Max int) (values []int, err error) {
	Code = strings.Replace(Code, "*", "0-60", -1)
	var tempValues []int
	tempValues, err = parseTimeCode(Code)
	if err != nil {
		return
	}

	for _, val := range tempValues {
		if val >= Min && val <= Max {
			values = append(values, val)
		}
	}

	if len(values) == 0 {
		return []int{}, fmt.Errorf("cron-schedule error - timecode parse error - Code %s corresponds to no matching times between %d and %d", Code, Min, Max)
	}

	return values, nil
}

func parseTimeCode(Code string) ([]int, error) {
	Values := []int{}
	SubFields := strings.Split(Code, ",")
	for _, field := range SubFields {
		if vals, valid := parseLiteral(field); valid {
			Values = append(Values, vals)
			continue
		}
		if vals, valid := parseRange(field); valid {
			Values = append(Values, vals...)
			continue
		}
		if vals, valid := parseStepRange(field); valid {
			Values = append(Values, vals...)
			continue
		}
		return nil, fmt.Errorf("cron-schedule error - timecode parse error - Unexpected token %s", field)
	}

	if len(Values) == 0 {
		return []int{}, fmt.Errorf("cron-schedule error - timecode parse error - Code %s corresponds to no matching times", Code)
	}

	return Values, nil
}

func parseLiteral(Code string) (int, bool) {
	val, err := strconv.Atoi(Code)
	if err != nil {
		return -1, false
	}
	return val, true
}

func parseRange(Code string) ([]int, bool) {
	r := strings.Split(Code, "-")
	if len(r) != 2 {
		return nil, false
	}

	Start, valid := parseLiteral(r[0])
	if !valid {
		return nil, valid
	}

	End, valid := parseLiteral(r[1])
	if !valid {
		return nil, valid
	}

	Values := []int{}
	if Start < End {
		for i := Start; i <= End; i++ {
			Values = append(Values, i)
		}
	} else {
		for i := End; i <= Start; i++ {
			Values = append(Values, i)
		}
	}

	return Values, true
}

func parseStepRange(Code string) ([]int, bool) {
	r := strings.Split(Code, "-")
	if len(r) != 2 {
		return nil, false
	}

	Start, valid := parseLiteral(r[0])
	if !valid {
		return nil, valid
	}

	step := strings.Split(r[1], "/")
	if len(step) != 2 {
		return nil, false
	}

	End, valid := parseLiteral(step[0])
	if !valid {
		return nil, valid
	}

	Step, valid := parseLiteral(step[1])
	if !valid {
		return nil, valid
	}

	Values := []int{}
	for i := Start; i <= End; i += Step {
		Values = append(Values, i)
	}
	return Values, true
}
