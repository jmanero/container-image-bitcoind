package client

// BlockchainInfo decodes getblockchaininfo RPC messages
type BlockchainInfo struct {
	Chain         string  `json:"chain"`
	Blocks        uint64  `json:"blocks"`
	Headers       uint64  `json:"headers"`
	BestBlockHash string  `json:"bestblockhash"`
	Difficulty    float64 `json:"difficulty"`
	Time          int64   `json:"time"`
	MedianTime    int64   `json:"mediantime"`
	Progress      float64 `json:"verificationprogress"`
	Initializing  bool    `json:"initialblockdownload"`
	ChainWork     string  `json:"chainwork"`
	SizeOnDisk    uint64  `json:"size_on_disk"`
	Pruned        bool    `json:"pruned"`
	Warnings      string  `json:"warnings"`
}
