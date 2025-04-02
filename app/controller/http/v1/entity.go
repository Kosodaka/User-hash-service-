package v1

import "mainHashService/app/entity"

type CheckerRequest struct {
	Hash   []entity.Hash `json:"hash"`
	Domain int64         `json:"domain"`
}

type QueryRequest struct {
	Fields     []string           `json:"fields"`
	Statements []entity.QueryStmt `json:"statements"`
}

type GetHashRequest struct {
	Query string `json:"query"`
}

type UnhashFromFileRequest struct {
	Bucket  string `json:"bucket_name"`
	ObjName string `json:"object_name"`
}
