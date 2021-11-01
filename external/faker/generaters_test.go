package faker

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	pb "github.com/dysnix/ai-scale-proto/external/proto/services"
)

func TestMetricsGenerator(t *testing.T) {
	a := pb.ReqSendMetrics{}

	err := MetricsGenerator()
	if err != nil {
		t.Error(err)
		return
	}

	err = faker.FakeData(&a)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", a)
}
