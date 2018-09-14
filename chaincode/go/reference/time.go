package main

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	timeBasicArgumentsNumber = 2
)

type TimePeriod struct {
	From int64   `json:"timestamp"`
	To int64     `json:"timestamp"`
}

func (time *TimePeriod) FillFromArguments(args []string) error {
	//0				  1
	//timePeriodFrom  timePeriodTo
	if len(args) < timeBasicArgumentsNumber {
		return errors.New(fmt.Sprintf("arguments array must contain at least %d items", timeBasicArgumentsNumber))
	}

	timePeriodFrom, err := strconv.ParseInt(args[0],10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse the timePeriodFrom: %s", err.Error()))
	}

	if timePeriodFrom < 0 {
		return errors.New("timePeriodFrom must be larger than zero")
	}

	timePeriodTo, err := strconv.ParseInt(args[1],10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse the timePeriodTo: %s", err.Error()))
	}

	if timePeriodTo < timePeriodFrom {
		return errors.New("timePeriodTo must be larger than timePeriodFrom")
	}

	time.From = int64(timePeriodFrom)
	time.To   = int64(timePeriodTo)

	return nil
}
