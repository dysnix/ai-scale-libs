package faker

import (
	"reflect"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/bxcodec/faker/v3"

	tc "github.com/dysnix/ai-scale-libs/external/types_convertation"
	pb "github.com/dysnix/ai-scale-proto/external/proto/commonproto"
	"github.com/dysnix/ai-scale-proto/external/proto/enums"
)

// MetricsGenerator function for generate random metrics response from Provider
func MetricsGenerator(diffDuration time.Duration) (err error) {
	if err = faker.AddProvider("unixTime", func(v reflect.Value) (interface{}, error) {
		start := time.Now()
		t := gofakeit.DateRange(start.Add(-diffDuration), start)
		return uint64(t.Unix()), nil
	}); err != nil {
		return err
	}

	if err = faker.AddProvider("metricsSlice", func(v reflect.Value) (interface{}, error) {
		var result []*pb.MetricValue

		for i := 0; i < gofakeit.Number(1, 10); i++ {
			tmpMetricType := gofakeit.Number(0, 6)
			tmpPrometheusResponseType := gofakeit.Number(0, 3)
			result = append(result, &pb.MetricValue{
				MetricType:             enums.MetricsType(tmpMetricType),
				PrometheusResponseType: enums.ValueType(tmpPrometheusResponseType),
				Values:                 generateRandItemsSlice(1, 10, diffDuration),
			})
		}

		return result, nil
	}); err != nil {
		return err
	}

	return nil
}

func generateRandItemsSlice(min, max int, diffDuration time.Duration) (result []*pb.Item) {
	start := time.Now()

	for i := 0; i < gofakeit.Number(min, max); i++ {
		tmpTime, _ := tc.AdaptTimeToPbTimestamp(tc.TimeToTimePtr(gofakeit.DateRange(start.Add(-diffDuration), start)))

		result = append(result, &pb.Item{
			Timestamp: tmpTime,
			Value:     gofakeit.Float64Range(0.1, 1000),
		})
	}

	return result
}
