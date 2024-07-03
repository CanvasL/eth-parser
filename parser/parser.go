package parser

import (
	"context"
	"eth_parser/config"
	"eth_parser/eth"
	"fmt"
	"strings"
	"sync"
	"time"
)

type IParser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []eth.Transaction
}

type Parser struct {
	currentBlock int
	transactions map[string][]eth.Transaction
	subscribed   map[string]context.CancelFunc
	mu           sync.Mutex
}

func NewParser() *Parser {
	return &Parser{
		currentBlock: 0,
		transactions: make(map[string][]eth.Transaction),
		subscribed:   make(map[string]context.CancelFunc),
	}
}

func (p *Parser) GetCurrentBlock() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentBlock
}

func (p *Parser) Subscribe(address_ string) bool {
	address := strings.ToLower(address_)

	p.mu.Lock()
	if cancel, found := p.subscribed[address]; found {
		fmt.Printf("%v already subscribed.\n", address)
		cancel()
		delete(p.transactions, address)
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.subscribed[address] = cancel
	p.mu.Unlock()

	go func() {
		blockNumber, err := eth.GetCurrentBlockNumber()
		if err != nil {
			fmt.Printf("Failed to get the current block number: %v\n", err)
			return
		}
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context canceled.")
				return
			default:
				block, err := eth.GetBlockByNumber(blockNumber)
				if err != nil {
					fmt.Printf("Failed to get block details: %v\n", err)
					time.Sleep(config.POLLING_ERROR_TOLERANCE)
					continue
				}

				count := 0
				p.mu.Lock()
				for _, tx := range block.Transactions {
					if tx.From == address || tx.To == address {
						p.transactions[address] = append(p.transactions[address], tx)
						count++
					}
				}
				fmt.Printf("[%v] %v transactions related to %v.\n", blockNumber, count, address)

				p.currentBlock = blockNumber
				p.mu.Unlock()

				blockNumber++
				time.Sleep(config.POLLING_INTERVAL)
			}
		}
	}()

	return true
}

func (p *Parser) GetTransactions(address_ string) []eth.Transaction {
	address := strings.ToLower(address_)
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.transactions[address]
}
