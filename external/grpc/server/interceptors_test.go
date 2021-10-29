package server

import (
	pb "github.com/dysnix/ai-scale-proto/external/proto/services"
	"reflect"
	"testing"
)

func TestInjectClientMetadataInterceptor(t *testing.T) {
	var req interface{} = &pb.ReqSendMetrics{}
	st := reflect.TypeOf(req)
	_, ok := st.MethodByName("GetHeader")
	if ok {
		var b interface{} = pb.Header{ClusterId: "bsc-1"}
		field := reflect.New(reflect.TypeOf(b))
		field.Elem().Set(reflect.ValueOf(b))
		reflect.ValueOf(req).Elem().FieldByName("Header").Set(field)
	}

	t.Log(req.(*pb.ReqSendMetrics).Header.ClusterId)
}
