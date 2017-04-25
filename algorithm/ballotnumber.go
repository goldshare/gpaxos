package algorithm

import (
	"github.com/lichuang/gpaxos/common"
)

type BallotNumber struct {
	proposalId uint64
	nodeId     common.NodeId
}

func (bn *BallotNumber) BE(other *BallotNumber) bool {
	if bn.proposalId == other.proposalId {
		return bn.nodeId >= other.nodeId
	}

	return bn.proposalId >= other.proposalId
}

func (bn *BallotNumber) NE(other *BallotNumber) bool {
	return bn.proposalId != other.proposalId ||
		bn.nodeId != other.nodeId
}

func (bn *BallotNumber) EQ(other *BallotNumber) bool {
	return !bn.NE(other)
}

func (bn *BallotNumber) BT(other *BallotNumber) bool {
	if bn.proposalId == other.proposalId {
		return bn.nodeId > other.nodeId
	}

	return bn.proposalId > other.proposalId
}

func (bn *BallotNumber) IsNull() bool {
	return bn.proposalId == 0
}

func (bn *BallotNumber) Reset() {
	bn.nodeId = 0
	bn.proposalId = 0
}
