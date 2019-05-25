package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type BlockData struct {
	Product		string
	Price		float32
	User		string
}

type Block struct {
	Id 			int
	Data		BlockData
	Hash		string
	PrevHash	string
	Difficulty  int
	Nonce		int
	Timestamp	time.Time
}

func (this *Block) create(data BlockData) {
	this.Data = data
	this.Timestamp = time.Now()
	this.Nonce = 0
}

func (this *Block) validateHash(hash string) bool {
	this.generateHash()

	if this.Hash != hash {
		return false
	}

	return true
}

func (this *Block) generateHash() {
	record := string(this.Id) + this.Timestamp.String() + string(this.Nonce) + this.PrevHash
	hash := sha256.Sum256([]byte(record))
	this.Hash = fmt.Sprintf("%x", hash)
}

func (this *Block) Mine() {
	prefix := getPrefix(this.Difficulty)
	finish := false

	for finish == false {
		this.Nonce++
		this.generateHash()

		if strings.HasPrefix(this.Hash, prefix) {
			finish = true
		}
	}
}

type Blockchain struct {
	Blocks     	[]Block
	CantBlocks 	int
	Difficulty 	int
}

func (this *Blockchain) createBlockchain() {
	this.CantBlocks = 0
	this.Difficulty = 3
	this.initialize()
}

func (this *Blockchain) initialize() {
	var genesisBlock Block
	data := BlockData{}

	genesisBlock.Id = 0
	genesisBlock.Difficulty = this.Difficulty
	hash := sha256.Sum256([]byte("Initialization Data"))
	genesisBlock.PrevHash = fmt.Sprintf("%x", hash)

	genesisBlock.create(data)
	genesisBlock.Mine()

	this.Append(genesisBlock)
}

func (this *Blockchain) AddBlock(data BlockData)  {
	var block Block
	prevBlock := this.GetLastBlock()

	block.Id = this.CantBlocks
	block.PrevHash = prevBlock.Hash
	block.Difficulty = this.Difficulty
	block.create(data)
	block.Mine()

	if validBlock(&block, &prevBlock) {
		this.Append(block)
	}
}

func (this *Blockchain) Append(block Block) {
	this.CantBlocks = this.CantBlocks + 1
	this.Blocks = append(this.Blocks, block)
}

func (this *Blockchain) GetLastBlock() Block {
	return this.Blocks[this.CantBlocks-1]
}

func getPrefix(length int) string {
	letterBytes := "0"
	b := make([]byte, length)

	for i := range b {
		b[i] = letterBytes[0]
	}

	return string(b)
}

func validBlock(block, prevBlock *Block) bool {
	// Confirm the hashes
	if block.PrevHash != prevBlock.Hash {
		return false
	}

	// confirm the block's hash is valid
	if !block.validateHash(block.Hash) {
		return false
	}

	// Check the position to confirm its been incremented
	if block.Id != prevBlock.Id+1 {
		return false
	}

	return true
}

func main() {
	var blockchain Blockchain

	blockchain.createBlockchain()

	data := BlockData {
		Product: "T-Shirt",
		Price: 225.00,
		User: "John Doe",
	}

	blockchain.AddBlock(data)

	for i := 0; i < blockchain.CantBlocks; i++ {
		fmt.Println(blockchain.Blocks[i])
	}
}