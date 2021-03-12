package DCP

import (
	"crypto/rand"
	"github.com/didiercrunch/paillier"
	"math/big"
	"time"
)

type ICalculationObject interface {
	Add()
	Mul()
	Encrypt()
	Decrypt()
	KeyGen() error
}

type CalculationObjectPaillier struct {
	Counter    int
	PrivateKey *paillier.PrivateKey
	PublicKey  paillier.PublicKey
	Cipher     *paillier.Cypher
}

func (cop *CalculationObjectPaillier) KeyGen() error {
	p1, p2, e := paillier.GenerateSafePrime(128, 1, 1 * time.Second, rand.Reader)
	if e != nil {
		return e
	}

	cop.PrivateKey = paillier.CreatePrivateKey(p1, p2)
	cop.PublicKey = cop.PrivateKey.PublicKey

	return nil
}

func (cop *CalculationObjectPaillier) Encrypt() {
	cop.Cipher, _ = cop.PublicKey.Encrypt(big.NewInt(0), rand.Reader)
}

func (cop CalculationObjectPaillier) Decrypt() {

}

