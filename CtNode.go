package DCP

import (
	"fmt"
	"github.com/google/uuid"
)

type ICtNode interface {
	InitRoutine(Prepare) error
	Broadcast()
	Listen()
	HandleCalculationObject(interface{})
	Print()
}

type CtNode struct {
	Id             uuid.UUID
	Co             *CalculationObjectPaillier
	Ids            []string
	ReachableNodes map[chan *CalculationObjectPaillier]struct{}
	Channel        chan *CalculationObjectPaillier
	HandledCoIds   map[uuid.UUID]struct{}
}

func InitRoutine(fn Prepare, node *CtNode) error {
	e := fn(node)
	return e
}

func (node *CtNode) Broadcast(externalCo *CalculationObjectPaillier) {
	fmt.Printf("Broadcasting triggered in node %s\n", node.Id)

	var objToBroadcast *CalculationObjectPaillier
	if externalCo != nil {
		objToBroadcast = externalCo
	} else {
		objToBroadcast = node.Co
	}

	for rn, _ := range node.ReachableNodes {
		go func(rn chan *CalculationObjectPaillier) {
			rn <- objToBroadcast
		}(rn)
	}
}

func (node *CtNode) Listen() {
	go func() {
		for {
			co := <-node.Channel
			fmt.Printf("Listen triggered in node %s\n", node.Id)
			node.HandleCalculationObject(co)
		}
	}()
}

func (node *CtNode) HandleCalculationObject(co *CalculationObjectPaillier) {
	// Run Eval
	// Broadcast

	if node.Co.PublicKey.N.Cmp(co.PublicKey.N) == 0 {
		fmt.Println("Public key match")
		if co.Counter > nodeVisitDecryptThreshold {
			fmt.Println("Calculation process finished, updating internal CalculationObject")
			node.Co = co
			close(node.Channel)
		}

		fmt.Println("Too few participants to satisfy privacy. Still listening")
		// Too few participants to satisfy privacy, abort Calculation process
		return
	}

	if _, exist := node.HandledCoIds[co.Id]; exist {
		fmt.Printf("Calculation object with ID: %s already handled\n", co.Id.String())
		return
	}

	idLen := len(node.Ids)

	cipher, e := co.Encrypt(idLen)
	if e != nil {
		// No-op
		fmt.Println(e.Error())
		return
	}

	// Add to co cipher
	co.Add(cipher)
	co.Counter = co.Counter + 1

	node.HandledCoIds[co.Id] = struct{}{}
	node.Broadcast(co)
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.privateKey, node.Co.PublicKey)
}
