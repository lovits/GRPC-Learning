package main

import (
	"fmt"
	"protobuf_demo/api"
)

func oneofdemo() {
	req1 := api.NoticeReaderRequest{
		Msg: "李文周的博客更新啦~",
		NoticeWay: &api.NoticeReaderRequest_Email{
			Email: "123@xx.com",
		},
	}
	// 使用短信通知的请求消息
	// req2 := api.NoticeReaderRequest{
	// 	Msg: "李文周的博客更新啦~",
	// 	NoticeWay: &api.NoticeReaderRequest_Phone{
	// 		Phone: "123456789",
	// 	},
	// }
	req := req1
	switch v := req.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		noticeWithEmail(v)
	case *api.NoticeReaderRequest_Phone:
		noticeWithPhone(v)
	}

}
func noticeWithEmail(in *api.NoticeReaderRequest_Email) {
	fmt.Printf("notice reader by email:%v\n", in.Email)
}

func noticeWithPhone(in *api.NoticeReaderRequest_Phone) {
	fmt.Printf("notice reader by phone:%v\n", in.Phone)
}
func main() {
	oneofdemo()
}
