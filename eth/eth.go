package eth

import (
	"bytes"
	"encoding/json"
	"eth_parser/config"
	"eth_parser/utils"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Block struct {
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func GetCurrentBlockNumber() (int, error) {
	requestBody, _ := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	})

	resp, err := http.Post(config.RPC_URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var rpcResponse RPCResponse
	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return -1, err
	}

	if rpcResponse.Error != nil {
		return -1, fmt.Errorf("RPC error: %s", rpcResponse.Error.Message)
	}

	var blockNumberHex string
	if err := json.Unmarshal(rpcResponse.Result, &blockNumberHex); err != nil {
		return -1, err
	}

	blockNumber, err := utils.HexStringToInt(blockNumberHex)
	if err != nil {
		return -1, err
	}

	return blockNumber, nil
}

func GetBlockByNumber(blockNumber int) (*Block, error) {
	blockNumberHex := utils.IntToHexString(blockNumber)

	requestBody, _ := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNumberHex, true},
		ID:      1,
	})

	resp, err := http.Post(config.RPC_URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var rpcResponse RPCResponse
	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return nil, err
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResponse.Error.Message)
	}

	var block Block
	if err := json.Unmarshal(rpcResponse.Result, &block); err != nil {
		return nil, err
	}

	return &block, nil
}
