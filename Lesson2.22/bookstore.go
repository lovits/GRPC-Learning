package main

import (
	"context"
	"lesson222/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	proto.UnimplementedBookstoreServer            // 必须嵌入
	bs                                 *bookstore // 书店数据库操作
}

// 返回所有书架
func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*proto.ListShelvesResponse, error) {
	sl, err := s.bs.ListShelves(ctx)
	if err == gorm.ErrEmptySlice {
		return &proto.ListShelvesResponse{}, nil
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "ListShelves fail")
	}
	vl := make([]*proto.Shelf, 0, len(sl))
	for _, v := range sl {
		vl = append(vl, &proto.Shelf{
			Id:    v.ID,
			Theme: v.Theme,
			Size:  v.Size,
		})
	}
	return &proto.ListShelvesResponse{Shelves: vl}, nil
}

// 创建书架
func (s *server) CreateShelf(ctx context.Context, in *proto.CreateShelfRequest) (*proto.Shelf, error) {
	if len(in.GetShelf().GetTheme()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid theme")
	}
	data := Shelf{Theme: in.Shelf.GetTheme(), Size: in.GetShelf().GetSize()}
	ns, err := s.bs.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &proto.Shelf{Id: ns.ID, Theme: ns.Theme, Size: ns.Size}, nil
}

// 查询书架
func (s *server) GetShelf(ctx context.Context, in *proto.GetShelfRequest) (*proto.Shelf, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	shelf, err := s.bs.GetShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	return &proto.Shelf{Id: shelf.ID, Theme: shelf.Theme, Size: shelf.Size}, nil
}

// 删除书架
func (s *server) DeleteShelf(ctx context.Context, in *proto.DeleteShelfRequest) (*emptypb.Empty, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	err := s.bs.DeleteShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "delete failed")
	}
	return &emptypb.Empty{}, nil
}

// func (s *server) ListBooks(ctx context.Context, in *proto.ListBooksRequest) (*proto.ListBooksResponse, error) {
// 	if in.GetShelf() <= 0 {
// 		return nil, status.Error(codes.InvalidArgument, "invalid vaild argument")
// 	}
// 	if pageToken := in.GetPageToken(); pageToken == "" {

// 	} else {
// 		pageInfo := Token(in.GetPageToken()).Decode()
// 	}

// 	if pageInfo.NextID==""||pageInfo.NextTimeAtUTC==0||pageInfo.NextTimeAtUTC>time.Now.Unix()||pageInfo.PageSize<=0 {

// 	}

// }
