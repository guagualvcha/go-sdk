package main

import (
	"flag"
	"fmt"
	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/monitor/tg"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	testnetNodeAddr = "tcp://data-seed-pre-1-s3.binance.org:80"
	prodNodeAddr    = "tcp://dataseed5.ninicoin.io:80"
	tgBot           = tg.NewRoRob("bot715771123:AAG1VRqO2D_cD6u1ZCO2Ezz1DNvbEy78KOA", "-342791328")
)

func main() {
	network := flag.String("net", "prod", "network")
	flag.Parse()
	if *network == "prod" {
		fmt.Println("starting watch prod")
		go monitTx(types.ProdNetwork)
	} else if *network == "testnet" {
		fmt.Println("starting watch testnet")
		go monitTx(types.TestNetwork)
	}
	select {}
}

func monitTx(network types.ChainNetwork) {
	var c *rpc.HTTP
	var net string
	if network == types.ProdNetwork {
		c = rpc.NewRPCClient(prodNodeAddr, network)
		net = "prod"
	} else {
		c = rpc.NewRPCClient(testnetNodeAddr, network)
		net = "testnet"
	}
	query := "tm.event = 'Tx'"
	_, err := tmquery.New(query)
	if err != nil {
		fmt.Println(err)
	}
	out, err := c.Subscribe(query, 100)
	var height int64
	for o := range out {
		txResult, ok := o.Data.(tmtypes.EventDataTx)
		if !ok {
			fmt.Println("error, not a event datatx")
			continue
		}
		if height!=txResult.Height{
			height = txResult.Height
			if height %100 ==0{
				fmt.Printf("%s receive height %d \n",net, txResult.Height)
			}
		}
		tx, err := rpc.ParseTx(tx.Cdc, txResult.Tx)
		if err != nil {
			fmt.Printf("parse tx error: %v\n", err)
			continue
		}
		msgs := tx.GetMsgs()
		for _, m := range msgs {
			if issue, ok := m.(msg.TokenIssueMsg); ok {
				tgBot.SentMessage(fmt.Sprintf("issue from: %v ,token: %v, supply: %d\n, net %s", issue.From.String(), issue.Symbol, issue.TotalSupply, net))
			} else if list, ok := m.(msg.DexListMsg); ok {
				tgBot.SentMessage(fmt.Sprintf("list from: %v ,pair: %v\n, net %s", list.From.String(), list.BaseAssetSymbol, net))
			} else if proposal, ok := m.(msg.SubmitProposalMsg); ok {
				if proposal.ProposalType == msg.ProposalTypeListTradingPair {
					tgBot.SentMessage(fmt.Sprintf("list proposal: from %v, content %v, net %v", proposal.Proposer.String(), proposal.Description, net))
				}
			}
		}
	}
}
