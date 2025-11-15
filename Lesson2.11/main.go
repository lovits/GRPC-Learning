package main

import (
	"fmt"

	"lesson2.11/api"
)

type Book struct {
	Price int64 //无法区分Price传没传值
	Title string
	// Price sql.NullInt64
	// price *int64
}

func foo() {
	// var book Book
	// if book.Price == nil{
	// 	//没有赋值
	// }else{}
	// book = Book{
	// 	Price:0
	// }
}
func wrapValueDemo() {
	//client
	// book := api.Book{
	// 	Title: "《跟七米学Go语言》",
	// 	Price: &wrapperspb.Int64Value{Value: 9900},
	// 	Memo
	// }

	// //server
	// if book.GetPrice() == nil { // price没赋值
	// 	fmt.Println("book with no price")
	// } else {
	// 	fmt.Printf("book with price:%v\n", book.GetPrice().GetValue())
	// }

}
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
