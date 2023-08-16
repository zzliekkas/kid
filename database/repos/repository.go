package repos

import (
	"context"

	"github.com/leor-w/kid/database/repos/deleter"

	"github.com/leor-w/kid/database/repos/updater"

	"github.com/leor-w/kid/database/repos/creator"

	"github.com/leor-w/kid/database/repos/finder"
)

type (
	// IBasicRepository 基础查询与事务
	IBasicRepository interface {
		IRepository
		ITxRepository
	}
	// IRepository 基础查询
	IRepository interface {
		Exist(*finder.Finder) bool
		GetOne(*finder.Finder) error
		Find(*finder.Finder) error
		Create(*creator.Creator) error
		Update(*updater.Updater) error
		Delete(*deleter.Deleter) error
		Save(*creator.Creator) error
		Exec(context.Context, string, ...interface{}) error
		Raw(context.Context, string, ...interface{}) error
		GetUniqueID(finder *finder.Finder, min, max, ignoreStart, ignoreEnd int64) int64
		Count(*finder.Finder) error
		Sum(*finder.Sum) error
	}
	// ITxRepository 事务
	ITxRepository interface {
		Transaction(func(context.Context) error) error                             // 开启新的事物执行
		GetDb(context.Context) interface{}                                         // 获取数据库连接
		WhetherTx(context.Context) bool                                            // 检查是否在事务中
		ExecWithTx(hasTx context.Context, fn func(tx context.Context) error) error // 如果hasTx已经在事务中,则不开启新的事务, 否则开启新的事务执行 fn
	}
)

type TxKey struct{}
