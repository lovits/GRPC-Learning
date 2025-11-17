package server

import (
	v1 "bubble/api/bubble/v1"
	"bubble/internal/conf"
	"bubble/internal/service"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"

	// "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	// jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Middleware1 自定义中间件
// type middleware func（Handler) Handler
// type Handler func(ctx context.Context, req any) (any, error)
func Middleware1() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			//执行之前做
			fmt.Println("Middleware:执行handeler之前")
			//做token校验 postman里在header添加token参数
			if tr, ok := transport.FromServerContext(ctx); ok {
				token := tr.RequestHeader().Get("token")
				fmt.Printf("token:%v\n", token)
			}
			defer func() {
				//执行之后做
				fmt.Println("Middleware:执行handeler之后")
			}()
			return handler(ctx, req)
		}
	}
}

// 自定义中间件2，相比于中间件1，失去了灵活性
func Middleware2(middleware.Handler) middleware.Handler {
	return nil
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, todo *service.TodoService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			Middleware1(),
			// Middleware1("a"),
			// Middleware1("b"),
			// Middleware1("c"),
			// Middleware2,

			// jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
			// 	return []byte("123"), nil
			// }),

			// recovery.Recovery(), //全局中间件
			// selector.Server( //特定path才执行的中间件

			// 	Middleware1(),
			// ).
			// 	Path("/api.bubble.v1.Todo/CreateTodo").
			// 	Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	//48 替换默认的Http响应编码器
	opts = append(opts, http.ResponseEncoder(responseEncoder))
	//48 替换默认的错误响应编码器
	opts = append(opts, http.ErrorEncoder(errorEncoder))

	srv := http.NewServer(opts...)
	v1.RegisterTodoHTTPServer(srv, todo)
	return srv
}
