package date

import (
	"strings"
	"time"

	"github.com/johannes-kuhfuss/services_utils/api_error"
)

const (
	ApiDateLayout = time.RFC3339
)

func IsValidTime(timeStr string) bool {
	_, err := time.Parse(ApiDateLayout, timeStr)
	return err == nil
}

func GetNowUtc() time.Time {
	return time.Now().UTC()
}

func GetNowLocal(location string) (*time.Time, api_error.ApiErr) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, api_error.NewBadRequestError("could not parse location")
	}
	localtime := time.Now().In(loc)
	return &localtime, nil
}

func GetNowUtcString() string {
	return GetNowUtc().Format(ApiDateLayout)
}

func GetNowLocalString(location string) (*string, api_error.ApiErr) {
	localtime, err := GetNowLocal(strings.TrimSpace(location))
	if err != nil {
		return nil, err
	} else {
		localtimeStr := localtime.Format(ApiDateLayout)
		return &localtimeStr, nil
	}
}
