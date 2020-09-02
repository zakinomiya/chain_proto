package miner

import (
	"go_chain/block"
	"go_chain/utils"
	"log"
	"runtime"
	"sync"
)

type Miner struct {
	quit      chan struct{}
	blockLock sync.Mutex
	wg        sync.WaitGroup
}

func New() *Miner {
	return &Miner{}
}

func (m *Miner) Start() {
	m.quit = make(chan struct{})
	m.wg.Add(1)
	go m.mining()

	log.Println("Mining process started")
}

func (m *Miner) Stop() {
	close(m.quit)
	m.wg.Wait()

	log.Println("Mining process stopped")
}

func (m *Miner) ServiceName() string {
	return "Miner"
}

func (m *Miner) mining() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go m.findNonce(utils.RandomUint32())
	}
}

func (m *Miner) findNonce(offset uint32, block *block.Block) {

}
