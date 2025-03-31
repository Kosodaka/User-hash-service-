package v1

import "mainHashService/app/entity"

type CheckerRequest struct {
	Hash   []entity.Hash `json:"hash"`
	Domain int64         `json:"domain"`
}

type UnhasherRequest struct {
	Query string `json:"query"`
}

type UnhashFromFileRequest struct {
	Bucket  string `json:"bucket_name"`
	ObjName string `json:"object_name"`
}
