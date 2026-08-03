package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// KeyVal 索引键（保序）：Field 字段名，Value 为 1/-1（排序方向）或 "text" 等字符串
type KeyVal struct {
	Field string `json:"f"`
	Value any    `json:"v"`
}

// IndexSpec 单个索引定义（对应 JSON 中一个索引）
type IndexSpec struct {
	Keys                   []KeyVal `json:"keys"`                   // [{"f":"tenant_id","v":1}, {"f":"domain","v":1}]（有序）
	Unique                 bool     `json:"unique,omitempty"`       // 是否唯一索引
	Name                   string   `json:"name,omitempty"`         // 索引名，缺省由 mongo 自动命名
	PartialFilterExpression bson.M  `json:"partialFilterExpression,omitempty"` // 部分索引过滤条件（如 {"is_deleted": false}，配合假删除）
}

// CollectionIndexes 一个集合的索引清单
type CollectionIndexes struct {
	Collection string      `json:"collection"`
	Indexes    []IndexSpec `json:"indexes"`
}

// EnsureIndexesFromFile 按 JSON 索引清单同步集合索引：
//   - 创建清单中缺失的索引；
//   - 删除集合中不在清单里的索引（清单是唯一事实来源），但保留系统索引 _id_。
//
// 清单格式（数组）：
//
//	[
//	  { "collection": "Directory", "indexes": [
//	    { "name": "idx_x", "keys": [{"f":"tenant_id","v":1}, {"f":"domain","v":1}], "unique": true }
//	  ]}
//	]
func EnsureIndexesFromFile(db *mongo.Database, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read index file %s: %w", filePath, err)
	}
	var specs []CollectionIndexes
	if err := json.Unmarshal(data, &specs); err != nil {
		return fmt.Errorf("parse index file %s: %w", filePath, err)
	}

	for _, cs := range specs {
		coll := db.Collection(cs.Collection)

		// 1、创建/补齐清单索引
		models := make([]mongo.IndexModel, 0, len(cs.Indexes))
		for _, idx := range cs.Indexes {
			keys := make(bson.D, 0, len(idx.Keys))
			for _, kv := range idx.Keys {
				keys = append(keys, bson.E{Key: kv.Field, Value: kv.Value})
			}
			opts := options.Index()
			if idx.Unique {
				opts = opts.SetUnique(true)
			}
			if idx.Name != "" {
				opts = opts.SetName(idx.Name)
			}
			if len(idx.PartialFilterExpression) > 0 {
				opts = opts.SetPartialFilterExpression(idx.PartialFilterExpression)
			}
			models = append(models, mongo.IndexModel{Keys: keys, Options: opts})
		}
		if _, err := coll.Indexes().CreateMany(context.Background(), models); err != nil {
			return fmt.Errorf("create indexes for %s: %w", cs.Collection, err)
		}

		// 2、删除不在清单里的索引（保留 _id_）
		keep := make(map[string]bool, len(cs.Indexes))
		for _, idx := range cs.Indexes {
			if idx.Name != "" {
				keep[idx.Name] = true
			}
		}
		cur, err := coll.Indexes().List(context.Background())
		if err != nil {
			return fmt.Errorf("list indexes for %s: %w", cs.Collection, err)
		}
		var existing []string
		for cur.Next(context.Background()) {
			var doc struct {
				Name string `bson:"name"`
			}
			if err := cur.Decode(&doc); err != nil {
				continue
			}
			existing = append(existing, doc.Name)
		}
		cur.Close(context.Background())

		for _, name := range existing {
			if name == "_id_" || keep[name] {
				continue
			}
			if _, err := coll.Indexes().DropOne(context.Background(), name); err != nil {
				return fmt.Errorf("drop index %s.%s: %w", cs.Collection, name, err)
			}
		}
	}
	return nil
}
