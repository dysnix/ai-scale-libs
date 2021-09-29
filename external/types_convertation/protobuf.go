package types_convertation

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

var (
	TimeIsEmptyOrZero = errors.New("time parameter is empty or zero")
)

func AdaptTimeToPbTimestamp(currentTime *time.Time) (*tspb.Timestamp, error) {
	if currentTime != nil && !(*currentTime).IsZero() {
		protoTime := &tspb.Timestamp{}
		protoTime, err := ptypes.TimestampProto(TimePtrToTime(currentTime))
		if err != nil {
			return nil, err
		}
		return protoTime, nil
	}
	return nil, TimeIsEmptyOrZero
}

func AdaptPbTimestampToTime(protoTime *tspb.Timestamp) (*time.Time, error) {
	if protoTime == nil || (protoTime.GetNanos() == 0 || protoTime.GetSeconds() == 0) {
		return nil, fmt.Errorf("proto time parameter is empty or zero")
	}
	return TimeToTimePtr(time.Unix(protoTime.GetSeconds(), int64(protoTime.GetNanos()))), nil
}
