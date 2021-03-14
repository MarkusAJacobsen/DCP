package DCP

import "fmt"

type ICtNode interface {
	InitRoutine(Prepare) error
	Broadcast()
	Listen()
	HandleCalculationObject(interface{})
	Print()
}

type CtNode struct {
	Co CalculationObjectPaillier
	Ids []string
	ReachableNodes []chan CalculationObjectPaillier
	channel chan CalculationObjectPaillier
}

func (node *CtNode) InitRoutine(fn Prepare) {
	fn(node)
}

func (node *CtNode) Broadcast() {
	for _, rn := range node.ReachableNodes {
		fmt.Printf("%T\n", rn)
		go func(rn chan CalculationObjectPaillier) {
			fmt.Println("Broadcasting")
			rn <- node.Co
		}(rn)
	}
}

func (node *CtNode) Listen() {
	go func() {
		for {
			co := <- node.channel
			fmt.Println("Listen triggered")
			node.HandleCalculationObject(co)
		}
	}()
}

func (node *CtNode) HandleCalculationObject(co CalculationObjectPaillier)  {
	// Run Eval
	// Broadcast

	// Check PK

	// Check counter

	//idLen := len(node.Ids)

	//cipher, e := co.Encrypt(idLen)
}

func (node CtNode) Print() {
	fmt.Printf("Counter %d, PK %s, SK %s\n", node.Co.Counter, node.Co.PrivateKey, node.Co.PublicKey)
}