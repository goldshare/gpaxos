package algorithm

import (
  "github.com/lichuang/gpaxos/config"
  "github.com/lichuang/gpaxos/logstorage"
  "github.com/lichuang/gpaxos/common"
)

type Instance struct {
}

func NewInstance(config *config.Config, logStorage logstorage.LogStorage,
  transport common.MsgTransport, ) *Instance {
  instance := new(Instance)

  return instance
}

func (self *Instance) GetLastChecksum() uint32 {
  return 0
}

func (self *Instance) GetNowInstanceID() uint64 {
  return 1
}

func (self *Instance) CheckNewValue() {

}

func (self *Instance) OnReceivePaxosMsg(msg *common.PaxosMsg) error {
  return nil
}

func (self *Instance) OnReceiveMsg(msg string) error {
  return nil
}

func (self *Instance) OnTimeout(timerId uint32, timerType int) {

}
