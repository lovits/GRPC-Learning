package data

import (
	"context"

	"bubble/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// 实现了biz层定义的repo接口
type todoRepo struct {
	data *Data
	log  *log.Helper
}

// NewTodoeRepo .
func NewTodoRepo(data *Data, logger log.Logger) biz.TodoRepo {
	return &todoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *todoRepo) Save(ctx context.Context, t *biz.Todo) (*biz.Todo, error) {
	//实现数据库的操作
	err := r.data.db.Create(t).Error
	return t, err
}

func (r *todoRepo) Update(ctx context.Context, t *biz.Todo) error {
	return nil
}

func (r *todoRepo) FindByID(ctx context.Context, id int64) (*biz.Todo, error) {
	t := biz.Todo{ID: id}
	err := r.data.db.First(&t).Error
	return &t, err
}

func (r *todoRepo) Delete(context.Context, int64) error {
	return nil
}

func (r *todoRepo) ListAll(context.Context) ([]*biz.Todo, error) {
	return nil, nil
}
