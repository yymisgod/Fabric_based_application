package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

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

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 120):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
