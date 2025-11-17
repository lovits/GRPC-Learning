package server

import (
	"net/http"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/grpc/status"
)

// Http Encoder
// 自定义HTTP响应编码器：生成自定义的响应格式
type httpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	// Data interface{} `json:"data"`
	Data any `json:"data"`
}

// 自定义编码方法
func responseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {

	if v == nil {
		return nil
	}
	//判断是不是重定向
	if rd, ok := v.(kratoshttp.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	//构造自定义响应结构体
	resp := &httpResponse{
		Code: http.StatusOK,
		Msg:  "success",
		Data: v,
	}
	codec, _ := kratoshttp.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(resp) //json序列化
	if err != nil {
		return err
	}
	//设置响应头  Content-Type:application/json
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, err = w.Write(data)
	return err

}

// 默认http编码器
// func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
// 	if v == nil {
// 		return nil
// 	}
// 	if rd, ok := v.(Redirector); ok {
// 		url, code := rd.Redirect()
// 		http.Redirect(w, r, url, code)
// 		return nil
// 	}
// 	codec, _ := CodecForRequest(r, "Accept")
// 	data, err := codec.Marshal(v)
// 	if err != nil {
// 		return err
// 	}
// 	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
// 	_, err = w.Write(data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// 自定义的错误响应编码器器
func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	resp := new(httpResponse)
	// 能从err⾥面解析出错误码的
	if gs, ok := status.FromError(err); ok {
		resp = &httpResponse{
			Code: kratosstatus.FromGRPCCode(gs.Code()),
			Msg:  gs.Message(),
			Data: nil,
		}
	} else {
		resp = &httpResponse{
			Code: http.StatusInternalServerError, // 500
			Msg:  "内部错误",
		}
	}
	codec, _ := kratoshttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(resp.Code)
	_, _ = w.Write(body)
}

// DefaultErrorEncoder encodes the error to the HTTP response.
// func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
// 	se := errors.FromError(err)
// 	codec, _ := CodecForRequest(r, "Accept")
// 	body, err := codec.Marshal(se)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
// 	w.WriteHeader(int(se.Code))
// 	_, _ = w.Write(body)
// }
