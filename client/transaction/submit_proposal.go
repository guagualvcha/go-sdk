package transaction

import (
	"encoding/json"
	"strconv"
	"time"

	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type SubmitProposalResult struct {
	tx.TxCommitResult
	ProposalId int64 `json:"proposal_id"`
}

func (c *client) SubmitListPairProposal(title string, param msg.ListTradingPairParams, initialDeposit int64, votingPeriod time.Duration, sync bool, memo string, source int64) (*SubmitProposalResult, error) {
	bz, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}
	return c.SubmitProposal(title, string(bz), msg.ProposalTypeListTradingPair, initialDeposit, votingPeriod, sync, memo, source)
}

func (c *client) SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit int64, votingPeriod time.Duration, sync bool, memo string, source int64) (*SubmitProposalResult, error) {
	fromAddr := c.keyManager.GetAddr()
	coins := ctypes.Coins{ctypes.Coin{Denom: types.NativeSymbol, Amount: initialDeposit}}
	proposalMsg := msg.NewMsgSubmitProposal(title, description, proposalType, fromAddr, coins, votingPeriod)
	err := proposalMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(proposalMsg, sync, memo, source)
	if err != nil {
		return nil, err
	}
	var proposalId int64
	if commit.Ok && sync {
		// Todo since ap do not return proposal id now, do not return err
		proposalId, err = strconv.ParseInt(string(commit.Data), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &SubmitProposalResult{*commit, proposalId}, err

}
