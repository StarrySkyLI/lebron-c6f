// Code generated by goctl. DO NOT EDIT.
package types

type SeckillReq struct {
	User_id    int64 `json:"userId"`
	Product_id int64 `json:"productId"`
}

type SeckillResponse struct {
}

type QuerReq struct {
	User_id    int64 `json:"userId"`
	Product_id int64 `json:"productId"`
}

type QuerResponse struct {
	Message string `json:"message"`
}

type SeckillResultReq struct {
	User_id      int64  `json:"userId"`
	Product_id   int64  `json:"productId"`
	SeckillToken string `json:"seckillToken"`
}

type SeckillResultResponse struct {
}
