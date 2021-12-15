package types

type NodeAddress struct {
	Address string `json:"address"`
	Score   uint64 `json:"score"`
}

type Node struct {
	Addresses  []*NodeAddress `json:"addresses"`
	IsOutbound bool           `json:"is_outbound"`
	NodeId     string         `json:"node_id"`
	Version    string         `json:"version"`
}

type BannedAddress struct {
	Address   string `json:"address"`
	BanReason string `json:"ban_reason"`
	BanUntil  uint64 `json:"ban_until"`
	CreatedAt uint64 `json:"created_at"`
}
