package types_convertation

import (
	"strconv"
	"time"
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

	sec := ts / 1000
	msec := ts % 1000

	return time.Unix(sec, msec*int64(time.Millisecond)), nil
}
