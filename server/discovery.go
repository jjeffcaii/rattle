package server

import (
	"errors"
	"log"
	"sort"
	"sync"

	"github.com/hashicorp/memberlist"
)

type Node struct {
	id         string
	seeds      []string
	mlist      *memberlist.Memberlist
	broadcasts *memberlist.TransmitLimitedQueue
	lock       sync.Mutex
	port       int
}

func NewNode(id string, port int, seeds []string) *Node {
	return &Node{
		id:    id,
		seeds: seeds,
		port:  port,
	}
}

func (n *Node) Start() error {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.mlist != nil {
		return errors.New("started already")
	}
	c := memberlist.DefaultLocalConfig()
	c.BindPort = n.port
	c.Name = n.id
	//c.LogOutput = ioutil.Discard
	m, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	_, err = m.Join(n.seeds)
	if err != nil {
		return err
	}
	n.mlist = m

	n.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return n.mlist.NumMembers()
		},
		RetransmitMult: 3,
	}

	log.Println("discovery node started:", m.LocalNode().Address())

	return nil
}

func (n *Node) Members() []string {
	var s []string
	for _, node := range n.mlist.Members() {
		s = append(s, node.Address())
	}
	sort.Strings(s)
	return s
}
