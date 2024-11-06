package tracing

import "fmt"

const (
	// caller labels
	OrionMetricLabelKeyCallerService = "xo-caller-service"
	OrionMetricLabelKeyCallerIP      = "xo-caller-ip"

	// callee labels
	OrionMetricLabelKeyCalleeService = "xo-callee-service"
	OrionMetricLabelKeyCalleeIp      = "xo-callee-ip"

	// codes
	OrionMetricLabelKeyGrpcCode = "xo-grpc-code"
	OrionMetricLabelKeyHttpCode = "xo-http-code"
)

type OrionMetricLabel map[string]string

func (o OrionMetricLabel) SetCaller(serviceName, ip string) {
	o[OrionMetricLabelKeyCallerService] = serviceName
	o[OrionMetricLabelKeyCallerIP] = ip
}

func (o OrionMetricLabel) SetCallee(serviceName, ip string) {
	o[OrionMetricLabelKeyCalleeService] = serviceName
	o[OrionMetricLabelKeyCalleeIp] = ip
}

func (o OrionMetricLabel) SetCodes(gRPCCode, httpCode int) {
	o[OrionMetricLabelKeyGrpcCode] = fmt.Sprintf("%d", gRPCCode)
	o[OrionMetricLabelKeyHttpCode] = fmt.Sprintf("%d", httpCode)
}
