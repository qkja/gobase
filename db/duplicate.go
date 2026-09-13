package db

import (
	"errors"
	"strings"

	gerrors "github.com/qkja/gobase/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// MapDuplicateKey 将 Mongo 唯一键冲突（E11000）按索引名映射为 gobase 业务错误。
//
// 业务服务在仓储层「直接插入 → 捕获 E11000 → 按索引名判定冲突来源」，禁止「先查后插」查重
// （docs 租户管理需求说明 T6 / 平台账号与权限需求说明 P9）。
//
// indexToCode: 索引名 → 错误码常量（如 "idx_pfu_name_unique" → gerrors.CodePlatformUserNameExists）。
// 命中返回对应 BizError；非重复键错误返回 nil，调用方按原错误处理。
func MapDuplicateKey(err error, indexToCode map[string]string) error {
	if err == nil || !mongo.IsDuplicateKeyError(err) {
		return nil
	}
	var we mongo.WriteException
	if !errors.As(err, &we) {
		return nil
	}
	for _, we2 := range we.WriteErrors {
		if we2.Code != 11000 {
			continue
		}
		for idxName, code := range indexToCode {
			if strings.Contains(we2.Message, idxName) {
				return gerrors.New(code)
			}
		}
	}
	return nil
}
