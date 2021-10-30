package main

type ModelWorker struct {
	Worker_id   int    `json:"workerid"`
	Worker_name string `json:"name"`
	Remark      string `json:"remark"`
	Del_flag    string `json:"delflag"`
	Create_by   string `json:"createby"`
	Create_time string `json:"createtime"`
	Update_by   string `json:"updateby"`
	Update_time string `json:"updatetime"`
	Gid         string `json:"gid"`

	Historys []HistoryItem
}

type HistoryItem struct {
	TxId        string
	ModelWorker ModelWorker
}
