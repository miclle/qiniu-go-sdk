package clientv2

import (
	"crypto/tls"
	"fmt"
	clientV1 "github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/internal/log"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
)

var (
	printRequestTrace         = false
	printRequest        *bool = nil
	printRequestDetail  *bool = nil
	printResponse       *bool = nil
	printResponseDetail *bool = nil
)

func PrintRequestTrace(isPrint bool) {
	printRequestTrace = isPrint
}

func IsPrintRequestTrace() bool {
	return printRequestTrace
}

func PrintRequest(isPrint bool) {
	printRequest = &isPrint
}

func IsPrintRequest() bool {
	if printRequest != nil {
		return *printRequest
	}
	return clientV1.DebugMode
}

func PrintRequestDetail(isPrint bool) {
	printRequestDetail = &isPrint
}

func IsPrintRequestDetail() bool {
	if printRequestDetail != nil {
		return *printRequestDetail
	}
	return clientV1.DeepDebugInfo
}

func PrintResponse(isPrint bool) {
	printResponse = &isPrint
}

func IsPrintResponse() bool {
	if printResponse != nil {
		return *printResponse
	}
	return clientV1.DebugMode
}

func PrintResponseDetail(isPrint bool) {
	printResponseDetail = &isPrint
}

func IsPrintResponseDetail() bool {
	if printResponseDetail != nil {
		return *printResponseDetail
	}
	return clientV1.DeepDebugInfo
}

type debugInterceptor struct {
}

func newDebugInterceptor() Interceptor {
	return &debugInterceptor{}
}

func (interceptor *debugInterceptor) Priority() InterceptorPriority {
	return InterceptorPriorityDebug
}

func (interceptor *debugInterceptor) Intercept(req *http.Request, handler Handler) (*http.Response, error) {
	if interceptor == nil {
		return handler(req)
	}

	label := interceptor.requestLabel(req)

	if e := interceptor.printRequest(label, req); e != nil {
		return nil, e
	}

	req = interceptor.printRequestTrace(label, req)

	resp, err := handler(req)

	if e := interceptor.printResponse(label, resp); e != nil {
		return nil, e
	}

	return resp, err
}

func (interceptor *debugInterceptor) requestLabel(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	return fmt.Sprintf("Url:%s", req.URL.String())
}

func (interceptor *debugInterceptor) printRequest(label string, req *http.Request) error {
	if req == nil {
		return nil
	}

	printReq := IsPrintRequest()
	printReqDetail := IsPrintRequestDetail()
	if !printReq && !printReqDetail {
		return nil
	}

	info := label + " request:\n"
	d, dErr := httputil.DumpRequest(req, printReqDetail)
	if dErr != nil {
		return dErr
	}
	info += string(d) + "\n"

	log.Debug(info)
	return nil
}

func (interceptor *debugInterceptor) printRequestTrace(label string, req *http.Request) *http.Request {
	if !IsPrintRequestTrace() || req == nil {
		return req
	}

	label += "\n"
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Debug(label + fmt.Sprintf("GetConn, %s \n", hostPort))
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr := connInfo.Conn.RemoteAddr()
			log.Debug(label + fmt.Sprintf("GotConn, Network:%s RemoteAddr:%s \n", remoteAddr.Network(), remoteAddr.String()))
		},
		PutIdleConn: func(err error) {
			log.Debug(label + fmt.Sprintf("PutIdleConn, err:%v \n", err))
		},
		GotFirstResponseByte: func() {
			log.Debug(label + fmt.Sprint("GotFirstResponseByte \n"))
		},
		Got100Continue: func() {
			log.Debug(label + fmt.Sprint("Got100Continue \n"))
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Debug(label + fmt.Sprintf("DNSStart, host:%s \n", info.Host))
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.Debug(label + fmt.Sprintf("DNSDone, addr:%+v \n", info.Addrs))
		},
		ConnectStart: func(network, addr string) {
			log.Debug(label + fmt.Sprintf("ConnectStart, network:%+v ip:%s \n", network, addr))
		},
		ConnectDone: func(network, addr string, err error) {
			log.Debug(label + fmt.Sprintf("ConnectDone, network:%s ip:%s err:%v \n", network, addr, err))
		},
		TLSHandshakeStart: func() {
			log.Debug(label + fmt.Sprint("TLSHandshakeStart \n"))
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			log.Debug(label + fmt.Sprintf("TLSHandshakeDone, state:%+v err:%s \n", state, err))
		},
		// go1.10 不支持
		//WroteHeaderField: func(key string, value []string) {
		//	log.Debug(label + fmt.Sprintf("WroteHeaderField, key:%s value:%s \n", key, value))
		//},
		WroteHeaders: func() {
			log.Debug(label + fmt.Sprint("WroteHeaders \n"))
		},
		Wait100Continue: func() {
			log.Debug(label + fmt.Sprint("Wait100Continue \n"))
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			log.Debug(label + fmt.Sprintf("WroteRequest, err:%v \n", info.Err))
		},
	}
	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}

func (interceptor *debugInterceptor) printResponse(label string, resp *http.Response) error {
	if resp == nil {
		return nil
	}

	printResp := IsPrintResponse()
	printRespDetail := IsPrintResponseDetail()
	if !printResp && !printRespDetail {
		return nil
	}

	info := label + " response:\n"
	d, dErr := httputil.DumpResponse(resp, printRespDetail)
	if dErr != nil {
		return dErr
	}
	info += string(d) + "\n"

	log.Debug(info)
	return nil
}
