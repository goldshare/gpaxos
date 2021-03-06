package checkpoint

import (
  "github.com/lichuang/gpaxos/config"
  "github.com/lichuang/gpaxos/log"
  "github.com/lichuang/gpaxos/sm_base"
  "github.com/lichuang/gpaxos/logstorage"
  "github.com/lichuang/gpaxos/util"
  "github.com/lichuang/saul/common/errors"
)

type CheckpointManager struct {
  Config                   *config.Config
  LogStorage               logstorage.LogStorage
  Factory                  *sm_base.StateMachineFactory
  replayer                 *Replayer
  cleaner                  *Cleaner
  MinChosenInstanceID      uint64
  MaxChosenInstanceID      uint64
  inAskforCheckpointMode   bool
  NeedAskNodes             map[uint64]bool
  LastAskforCheckpointTime uint64
  UseCheckpointReplayer    bool
}

func NewCheckpointManager(config *config.Config, factory *sm_base.StateMachineFactory,
  storage logstorage.LogStorage, useReplayer bool) *CheckpointManager {
  mng := &CheckpointManager{
    Config:                config,
    LogStorage:            storage,
    Factory:               factory,
    UseCheckpointReplayer: useReplayer,
  }
  mng.replayer = NewReplayer(config, factory, storage, mng)
  mng.cleaner = NewCleaner(config, factory, storage, mng)

  return mng
}

func (self *CheckpointManager) Init() error {
  err := self.LogStorage.GetMinChosenInstanceId(self.Config.GetMyGroupIdx(), &self.MinChosenInstanceID)
  if err != nil {
    return err
  }

  err = self.cleaner.FixMinChosenInstanceID(self.MinChosenInstanceID)
  return err
}

func (self *CheckpointManager) Start() {
  if self.UseCheckpointReplayer {
    self.replayer.Start()
  }

  self.cleaner.Start()
}

func (self *CheckpointManager) Stop() {
  if self.UseCheckpointReplayer {
    self.replayer.Stop()
  }

  self.cleaner.Stop()
}

func (self *CheckpointManager) PrepareForAskforCheckpoint(sendNodeId uint64) error {
  _, exist := self.NeedAskNodes[sendNodeId]
  if !exist {
    self.NeedAskNodes[sendNodeId] = true
  }

  if self.LastAskforCheckpointTime == 0 {
    self.LastAskforCheckpointTime = util.NowTimeMs()
  }

  nowTime := util.NowTimeMs()
  if nowTime > self.LastAskforCheckpointTime + 60000 {
    log.Info("no majority reply, just ask for checkpoint")
  } else {
    if len(self.NeedAskNodes) < self.Config.GetMajorityCount() {
      log.Info("Need more other tell us need to ask for checkpoint")

      return errors.New("Need more other tell us need to ask for checkpoint")
    }
  }

  self.LastAskforCheckpointTime = 0
  self.inAskforCheckpointMode = true

  return nil
}

func (self *CheckpointManager) InAskforcheckpointMode() bool {
  return self.inAskforCheckpointMode
}

func (self *CheckpointManager) ExitCheckpointMode() {
  self.inAskforCheckpointMode = false
}

func (self *CheckpointManager) GetCheckpointInstanceID() uint64 {
  return self.Factory.GetCheckpointInstanceID(self.Config.GetMyGroupIdx())
}

func (self *CheckpointManager) SetMinChosenInstanceID(minChosenInstanceId uint64) error {
  options := logstorage.WriteOptions{
    Sync:true,
  }

  err := self.LogStorage.SetMinChosenInstanceId(options, self.Config.GetMyGroupIdx(), minChosenInstanceId)
  if err != nil {
    return err
  }

  self.MinChosenInstanceID = minChosenInstanceId
  return nil
}

func (self *CheckpointManager) SetMaxChosenInstanceID(maxChosenInstanceId uint64) error {
  self.MaxChosenInstanceID = maxChosenInstanceId
}

func (self *CheckpointManager) SetMinChosenInstanceIDCache(minChosenInstanceID uint64) {
  self.MinChosenInstanceID = minChosenInstanceID
}

func (self *CheckpointManager) GetMinChosenInstanceID() uint64 {
  return self.MinChosenInstanceID
}

func (self *CheckpointManager) GetMaxChosenInstanceID() uint64 {
  return self.MaxChosenInstanceID
}

func (self *CheckpointManager) GetCleaner() *Cleaner {
  return self.cleaner
}

func (self *CheckpointManager) GetReplayer() *Replayer {
  return self.replayer
}

