package types_convertation

import (
	"errors"
	"strconv"
	"time"
)

var Epoch = time.Unix(0, 0)

const (
	ZeroDurationErr = "zero time duration"
)

func ParseMillisecondUnixTimestamp(s interface{}) (res time.Time, err error) {
	var ts int64
	switch tmp := s.(type) {
	case string:
		ts, err = strconv.ParseInt(tmp, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
	case int8:
		ts = int64(tmp)
	case int16:
		ts = int64(tmp)
	case int32:
		ts = int64(tmp)
	case int:
		ts = int64(tmp)
	case int64:
		ts = tmp
	case uint:
		ts = int64(tmp)
	case uint8:
		ts = int64(tmp)
	case uint16:
		ts = int64(tmp)
	case uint32:
		ts = int64(tmp)
	case uint64:
		ts = int64(tmp)
	}

	if ts > 0 {
		return Epoch.Add(time.Duration(ts) * time.Millisecond), nil
	}

	return res, errors.New(ZeroDurationErr)
}
