package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// WithTransaction 在 MongoDB 事务中执行 fn。
//
// 用法：跨多个集合/仓储的批量写（如"删除目录域连带删用户/组织"）需要
// 要么全成功要么全回滚时，用此函数包裹；fn 收到的 ctx 为事务会话上下文，
// 各仓储方法传入它即可自动进入同一事务（前提是各仓储复用同一个 client）。
//
// 注意：多文档事务要求 MongoDB 副本集（replica set），standalone 会报错。
func WithTransaction(client *mongo.Client, ctx context.Context, fn func(ctx context.Context) error) error {
	return client.UseSession(ctx, func(sctx mongo.SessionContext) error {
		if err := sctx.StartTransaction(); err != nil {
			return err
		}
		if err := fn(sctx); err != nil {
			_ = sctx.AbortTransaction(sctx)
			return err
		}
		return sctx.CommitTransaction(sctx)
	})
}
