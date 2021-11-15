package api

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	motypes "github.com/smartbch/moeingevm/types"
	"github.com/smartbch/smartbch/internal/bigutils"
	"github.com/smartbch/smartbch/internal/testutils"
)

var (
	// ../testdata/basic/contracts/InternalTxs.sol
	contract1CreationBytecode = testutils.HexToBytes(`0x608060405234801561001057600080fd5b50604051610825380380610825833981810160405281019061003291906100d0565b81600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050610155565b6000815190506100ca8161013e565b92915050565b600080604083850312156100e357600080fd5b60006100f1858286016100bb565b9250506020610102858286016100bb565b9150509250929050565b60006101178261011e565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6101478161010c565b811461015257600080fd5b50565b6106c1806101646000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80634af6b5c21461005c57806361bc221a1461007a578063a11a1f3614610098578063a27eaec2146100c8578063c952aa9c146100f8575b600080fd5b610064610116565b6040516100719190610534565b60405180910390f35b61008261013c565b60405161008f919061054f565b60405180910390f35b6100b260048036038101906100ad91906104c4565b610142565b6040516100bf919061054f565b60405180910390f35b6100e260048036038101906100dd91906104c4565b6102db565b6040516100ef919061054f565b60405180910390f35b610100610474565b60405161010d9190610534565b60405180910390f35b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b6000806000815480929190610156906105fc565b9190505550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a27eaec26001846101a6919061056a565b6040518263ffffffff1660e01b81526004016101c2919061054f565b602060405180830381600087803b1580156101dc57600080fd5b505af11580156101f0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021491906104ed565b50600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a27eaec2600584610260919061056a565b6040518263ffffffff1660e01b815260040161027c919061054f565b602060405180830381600087803b15801561029657600080fd5b505af11580156102aa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102ce91906104ed565b50604082901b9050919050565b60008060008154809291906102ef906105fc565b9190505550600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e73620c360018461033f919061056a565b6040518263ffffffff1660e01b815260040161035b919061054f565b602060405180830381600087803b15801561037557600080fd5b505af1158015610389573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ad91906104ed565b50600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e73620c36005846103f9919061056a565b6040518263ffffffff1660e01b8152600401610415919061054f565b602060405180830381600087803b15801561042f57600080fd5b505af1158015610443573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061046791906104ed565b50604082901b9050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000813590506104a981610674565b92915050565b6000815190506104be81610674565b92915050565b6000602082840312156104d657600080fd5b60006104e48482850161049a565b91505092915050565b6000602082840312156104ff57600080fd5b600061050d848285016104af565b91505092915050565b61051f816105c0565b82525050565b61052e816105f2565b82525050565b60006020820190506105496000830184610516565b92915050565b60006020820190506105646000830184610525565b92915050565b6000610575826105f2565b9150610580836105f2565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156105b5576105b4610645565b5b828201905092915050565b60006105cb826105d2565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610607826105f2565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561063a57610639610645565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b61067d816105f2565b811461068857600080fd5b5056fea2646970667358221220abdd9e3c194feeb5f4281cd4cb41cd99cccf592dd984b690aca2c81c088564c964736f6c63430008000033`)
	contract2CreationBytecode = testutils.HexToBytes(`0x608060405234801561001057600080fd5b506040516105ab3803806105ab8339818101604052810190610032919061008e565b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610100565b600081519050610088816100e9565b92915050565b6000602082840312156100a057600080fd5b60006100ae84828501610079565b91505092915050565b60006100c2826100c9565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6100f2816100b7565b81146100fd57600080fd5b50565b61049c8061010f6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806361bc221a14610046578063a27eaec214610064578063c952aa9c14610094575b600080fd5b61004e6100b2565b60405161005b919061032a565b60405180910390f35b61007e6004803603810190610079919061029f565b6100b8565b60405161008b919061032a565b60405180910390f35b61009c61024f565b6040516100a9919061030f565b60405180910390f35b60005481565b60008060008154809291906100cc906103d7565b9190505550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e73620c360018461011c9190610345565b6040518263ffffffff1660e01b8152600401610138919061032a565b602060405180830381600087803b15801561015257600080fd5b505af1158015610166573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061018a91906102c8565b50600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166309010e936002846101d69190610345565b6040518263ffffffff1660e01b81526004016101f2919061032a565b60206040518083038186803b15801561020a57600080fd5b505afa15801561021e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061024291906102c8565b50604082901b9050919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000813590506102848161044f565b92915050565b6000815190506102998161044f565b92915050565b6000602082840312156102b157600080fd5b60006102bf84828501610275565b91505092915050565b6000602082840312156102da57600080fd5b60006102e88482850161028a565b91505092915050565b6102fa8161039b565b82525050565b610309816103cd565b82525050565b600060208201905061032460008301846102f1565b92915050565b600060208201905061033f6000830184610300565b92915050565b6000610350826103cd565b915061035b836103cd565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156103905761038f610420565b5b828201905092915050565b60006103a6826103ad565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b60006103e2826103cd565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561041557610414610420565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b610458816103cd565b811461046357600080fd5b5056fea264697066735822122051045f8546c4c69d329f3419221cdbcfb9e6875b3f255dcad533e5e07ac0575064736f6c63430008000033`)
	contract3CreationBytecode = testutils.HexToBytes(`0x608060405234801561001057600080fd5b50610234806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806309010e931461004657806361bc221a14610076578063e73620c314610094575b600080fd5b610060600480360381019061005b9190610112565b6100c4565b60405161006d919061014a565b60405180910390f35b61007e6100d2565b60405161008b919061014a565b60405180910390f35b6100ae60048036038101906100a99190610112565b6100d8565b6040516100bb919061014a565b60405180910390f35b6000604082901b9050919050565b60005481565b60008060008154809291906100ec9061016f565b9190505550604082901b9050919050565b60008135905061010c816101e7565b92915050565b60006020828403121561012457600080fd5b6000610132848285016100fd565b91505092915050565b61014481610165565b82525050565b600060208201905061015f600083018461013b565b92915050565b6000819050919050565b600061017a82610165565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156101ad576101ac6101b8565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6101f081610165565b81146101fb57600080fd5b5056fea26469706673582212203fe00ed2f97dd3fa1c5c00f7a0615a64ea011ae5d3debe28007bd7d1d1edd14b64736f6c63430008000033`)
	methodIdCall2             = "0xa11a1f36"
	//methodIdCall3 := "0xa27eaec2"
	//methodIdCallMe := "0xe73620c3"
)

type CallStackNode struct {
	depth      int32
	From       gethcmn.Address
	To         gethcmn.Address
	Input      hexutil.Bytes
	Output     hexutil.Bytes
	StatusCode int
	GasLeft    int64
	Calls      []*CallStackNode
}

func buildCallStack(tx *motypes.Transaction) *CallStackNode {
	var nodes []*CallStackNode

	for _, call := range tx.InternalTxCalls {
		newNode := &CallStackNode{
			depth: call.Depth,
			From:  call.Sender,
			To:    call.Destination,
			Input: call.Input,
		}
		if len(nodes) == 0 { // first call
			nodes = append(nodes, newNode)
			continue
		}

		lastNode := nodes[len(nodes)-1]
		if call.Depth == lastNode.depth+1 { // new call
			lastNode.Calls = append(lastNode.Calls, newNode)
			nodes = append(nodes, newNode)
			continue
		}

		if call.Depth <= lastNode.depth { // last calls return
			n := lastNode.depth - call.Depth
			for i := int32(0); i <= n; i++ {
				lastRet := tx.InternalTxReturns[0]
				tx.InternalTxReturns = tx.InternalTxReturns[1:]
				lastNode.Output = lastRet.Output
				lastNode.GasLeft = lastRet.GasLeft
				lastNode.StatusCode = lastRet.StatusCode
				nodes = nodes[:len(nodes)-1]
				lastNode = nodes[len(nodes)-1]
			}

			lastNode.Calls = append(lastNode.Calls, newNode)
			nodes = append(nodes, newNode)
			continue
		}

		panic(fmt.Errorf("lastNode.depth:%d, call.Depth:%d",
			lastNode.depth, call.Depth))
	}

	node0 := nodes[0]
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		ret := tx.InternalTxReturns[0]
		tx.InternalTxReturns = tx.InternalTxReturns[1:]
		node.Output = ret.Output
		node.GasLeft = ret.GasLeft
		node.StatusCode = ret.StatusCode
	}

	return node0
}

func TestInternalTxCalls(t *testing.T) {
	key := "8d0eb0baad6ea91b33c148698372bc2e220ea6cb841112577f93c8194c0c8f11"
	_app := testutils.CreateTestApp(key)
	defer _app.Destroy()

	tx1, _, contract3Addr := _app.DeployContractInBlock(key, contract3CreationBytecode)
	_app.EnsureTxSuccess(tx1.Hash())
	//println(contract3Addr.String())

	tx2, _, contract2Addr := _app.DeployContractInBlock(key,
		testutils.JoinBytes(contract2CreationBytecode, make([]byte, 12), contract3Addr[:]))
	_app.EnsureTxSuccess(tx2.Hash())
	//println(contract2Addr.String())

	tx3, _, contract1Addr := _app.DeployContractInBlock(key,
		testutils.JoinBytes(contract1CreationBytecode, make([]byte, 12), contract2Addr[:], make([]byte, 12), contract3Addr[:]))
	_app.EnsureTxSuccess(tx3.Hash())
	//println(contract1Addr.String())

	/*
		contract1.call2()
			=> contract2.call3()
				=> contract3.callMe()
				=> contract3.callMe()
			=> contract2.call3()
				=> contract3.callMe()
				=> contract3.callMe()
	*/
	callData := testutils.JoinBytes(testutils.HexToBytes(methodIdCall2), bigutils.NewU256(0x100).PaddedBytes(32))
	tx4, _ := _app.MakeAndExecTxInBlock(key, contract1Addr, 0, callData)
	_app.EnsureTxSuccess(tx4.Hash())
	moTx4 := _app.GetTx(tx4.Hash())
	require.Len(t, moTx4.InternalTxCalls, 7)
	require.Len(t, moTx4.InternalTxReturns, 7)
	for _, call := range moTx4.InternalTxCalls {
		fmt.Println(hex.EncodeToString(call.Input))
	}
	for _, ret := range moTx4.InternalTxReturns {
		fmt.Println(hex.EncodeToString(ret.Output))
	}
	cs := `{
  "From": "0x6db26a33492ccc4006599ed88b569c0b13c5d17a",
  "To": "0xe32d21f68654d87a4aad8c80616db99d95dde0f1",
  "Input": "0xa11a1f360000000000000000000000000000000000000000000000000000000000000100",
  "Output": "0x0000000000000000000000000000000000000000000001000000000000000000",
  "StatusCode": 0,
  "GasLeft": 890031,
  "Calls": [
    {
      "From": "0xe32d21f68654d87a4aad8c80616db99d95dde0f1",
      "To": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
      "Input": "0xa27eaec20000000000000000000000000000000000000000000000000000000000000101",
      "Output": "0x0000000000000000000000000000000000000000000001010000000000000000",
      "StatusCode": 0,
      "GasLeft": 890032,
      "Calls": [
        {
          "From": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
          "To": "0x0eefec15be847ced628df09459cb9b8492337210",
          "Input": "0xe73620c30000000000000000000000000000000000000000000000000000000000000102",
          "Output": "0x0000000000000000000000000000000000000000000001020000000000000000",
          "StatusCode": 0,
          "GasLeft": 879854,
          "Calls": null
        },
        {
          "From": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
          "To": "0x0eefec15be847ced628df09459cb9b8492337210",
          "Input": "0x09010e930000000000000000000000000000000000000000000000000000000000000103",
          "Output": "0x0000000000000000000000000000000000000000000001030000000000000000",
          "StatusCode": 0,
          "GasLeft": 876640,
          "Calls": null
        }
      ]
    },
    {
      "From": "0xe32d21f68654d87a4aad8c80616db99d95dde0f1",
      "To": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
      "Input": "0xa27eaec20000000000000000000000000000000000000000000000000000000000000105",
      "Output": "0x0000000000000000000000000000000000000000000001050000000000000000",
      "StatusCode": 0,
      "GasLeft": 876471,
      "Calls": [
        {
          "From": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
          "To": "0x0eefec15be847ced628df09459cb9b8492337210",
          "Input": "0xe73620c30000000000000000000000000000000000000000000000000000000000000106",
          "Output": "0x0000000000000000000000000000000000000000000001060000000000000000",
          "StatusCode": 0,
          "GasLeft": 866805,
          "Calls": null
        },
        {
          "From": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
          "To": "0x0eefec15be847ced628df09459cb9b8492337210",
          "Input": "0x09010e930000000000000000000000000000000000000000000000000000000000000107",
          "Output": "0x0000000000000000000000000000000000000000000001070000000000000000",
          "StatusCode": 0,
          "GasLeft": 863291,
          "Calls": null
        }
      ]
    }
  ]
}`
	require.Equal(t, cs, testutils.ToPrettyJSON(buildCallStack(moTx4)))
}

func TestGetTransactionReceipt(t *testing.T) {
	key := "8d0eb0baad6ea91b33c148698372bc2e220ea6cb841112577f93c8194c0c8f11"
	_app := testutils.CreateTestApp(key)
	defer _app.Destroy()
	_api := createSbchAPI(_app)

	tx1, _, contract3Addr := _app.DeployContractInBlock(key, contract3CreationBytecode)
	_app.EnsureTxSuccess(tx1.Hash())
	println("contract3Addr:", contract3Addr.String())

	tx2, _, contract2Addr := _app.DeployContractInBlock(key,
		testutils.JoinBytes(contract2CreationBytecode, make([]byte, 12), contract3Addr[:]))
	_app.EnsureTxSuccess(tx2.Hash())
	println("contract2Addr:", contract2Addr.String())

	tx3, _, contract1Addr := _app.DeployContractInBlock(key,
		testutils.JoinBytes(contract1CreationBytecode, make([]byte, 12), contract2Addr[:], make([]byte, 12), contract3Addr[:]))
	_app.EnsureTxSuccess(tx3.Hash())
	println("contract1Addr:", contract1Addr.String())

	/*
		contract1.call2()
			=> contract2.call3()
				=> contract3.callMe()
				=> contract3.callMe()
			=> contract2.call3()
				=> contract3.callMe()
				=> contract3.callMe()
	*/
	callData := testutils.JoinBytes(testutils.HexToBytes(methodIdCall2), bigutils.NewU256(0x100).PaddedBytes(32))
	tx4, h4 := _app.MakeAndExecTxInBlock(key, contract1Addr, 0, callData)
	_app.EnsureTxSuccess(tx4.Hash())
	moTx4 := _app.GetTx(tx4.Hash())
	require.Len(t, moTx4.InternalTxCalls, 7)
	require.Len(t, moTx4.InternalTxReturns, 7)
	for _, call := range moTx4.InternalTxCalls {
		fmt.Println(hex.EncodeToString(call.Input))
	}
	for _, ret := range moTx4.InternalTxReturns {
		fmt.Println(hex.EncodeToString(ret.Output))
	}

	callList := `[
  {
    "callPath": "call_0",
    "from": "0x6db26a33492ccc4006599ed88b569c0b13c5d17a",
    "to": "0xe32d21f68654d87a4aad8c80616db99d95dde0f1",
    "gas": "0xeef6c",
    "value": "0x0",
    "input": "0xa11a1f360000000000000000000000000000000000000000000000000000000000000100",
    "status": "0x0",
    "gasUsed": "0x15abd",
    "output": "0x0000000000000000000000000000000000000000000001000000000000000000"
  },
  {
    "callPath": "call_0_0",
    "from": "0xe32d21f68654d87a4aad8c80616db99d95dde0f1",
    "to": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
    "gas": "0xe571e",
    "value": "0x0",
    "input": "0xa27eaec20000000000000000000000000000000000000000000000000000000000000101",
    "status": "0x0",
    "gasUsed": "0xc26e",
    "output": "0x0000000000000000000000000000000000000000000001010000000000000000"
  },
  {
    "callPath": "call_0_0_0",
    "from": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
    "to": "0x0eefec15be847ced628df09459cb9b8492337210",
    "gas": "0xdc147",
    "value": "0x0",
    "input": "0xe73620c30000000000000000000000000000000000000000000000000000000000000102",
    "status": "0x0",
    "gasUsed": "0x5459",
    "output": "0x0000000000000000000000000000000000000000000001020000000000000000"
  },
  {
    "callPath": "staticcall_0_0_1",
    "from": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
    "to": "0x0eefec15be847ced628df09459cb9b8492337210",
    "gas": "0xd62ad",
    "value": "0x0",
    "input": "0x09010e930000000000000000000000000000000000000000000000000000000000000103",
    "status": "0x0",
    "gasUsed": "0x24d",
    "output": "0x0000000000000000000000000000000000000000000001030000000000000000"
  },
  {
    "callPath": "call_0_1",
    "from": "0xe32d21f68654d87a4aad8c80616db99d95dde0f1",
    "to": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
    "gas": "0xd8c25",
    "value": "0x0",
    "input": "0xa27eaec20000000000000000000000000000000000000000000000000000000000000105",
    "status": "0x0",
    "gasUsed": "0x2c6e",
    "output": "0x0000000000000000000000000000000000000000000001050000000000000000"
  },
  {
    "callPath": "call_0_1_0",
    "from": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
    "to": "0x0eefec15be847ced628df09459cb9b8492337210",
    "gas": "0xd434e",
    "value": "0x0",
    "input": "0xe73620c30000000000000000000000000000000000000000000000000000000000000106",
    "status": "0x0",
    "gasUsed": "0x959",
    "output": "0x0000000000000000000000000000000000000000000001060000000000000000"
  },
  {
    "callPath": "staticcall_0_1_1",
    "from": "0xa8115c4df61f9fb1e686d1692cd53fa4d4ced237",
    "to": "0x0eefec15be847ced628df09459cb9b8492337210",
    "gas": "0xd2e88",
    "value": "0x0",
    "input": "0x09010e930000000000000000000000000000000000000000000000000000000000000107",
    "status": "0x0",
    "gasUsed": "0x24d",
    "output": "0x0000000000000000000000000000000000000000000001070000000000000000"
  }
]`
	ret, err := _api.GetTransactionReceipt(tx4.Hash())
	require.NoError(t, err)
	//println(testutils.ToPrettyJSON(ret["internalTransactions"]))
	require.Equal(t, callList, testutils.ToPrettyJSON(ret["internalTransactions"]))

	retTxs, err := _api.GetTxListByHeight(gethrpc.BlockNumber(h4))
	require.NoError(t, err)
	require.Len(t, retTxs, 1)
	require.Equal(t, ret, retTxs[0])
}
