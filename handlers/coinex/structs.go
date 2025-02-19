package coinex

type DefaultResponse struct {
}

type TickerUpdate struct {
	method string `json:"method"`
	data   struct {
		market         string `json:"market"`
		updated_at     uint64 `json:"updated_at"`
		best_bid_price string `json:"best_bid_price"`
		best_bid_size  string `json:"best_bid_size"`
		best_ask_price string `json:"best_ask_price"`
		best_ask_size  string `json:"best_ask_price"`
	} `json:"data"`
	id uint8 `json:"id"`
}

type TradeUpdate struct {
}
