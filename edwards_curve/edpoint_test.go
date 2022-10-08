package edwards_curve

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/test"
	// "github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type OnCurveTest[T, S emulated.FieldParams] struct {
	P AffinePoint[T]
}

func (c *OnCurveTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api)
	if err != nil {
		return err
	}
	cr.AssertIsOnCurve(c.P)
	return nil
}

func TestGeneratorIsOnCurve(t *testing.T) {
	assert := test.NewAssert(t)
	circuit := OnCurveTest[Ed25519, Ed25519Scalars]{}
	witness := OnCurveTest[Ed25519, Ed25519Scalars]{
		P: AffinePoint[Ed25519]{
			X: emulated.NewElement[Ed25519](newBigInt("216936D3CD6E53FEC0A4E231FDD6DC5C692CC7609525A7B2C9562D608F25D51A")),
			Y: emulated.NewElement[Ed25519](newBigInt("6666666666666666666666666666666666666666666666666666666666666658")),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

// s1*x1 + s2*x2 = y
type MulAddTest[T, S emulated.FieldParams] struct {
	X1, X2 AffinePoint[T]
	S1, S2 emulated.Element[S]
	Y AffinePoint[T]
}

func (c *MulAddTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api)
	if err != nil {
		return err
	}
	X1S1 := cr.ScalarMul(c.X1, c.S1)
	X2S2 := cr.ScalarMul(c.X2, c.S2)
	sum := cr.Add(X1S1, X2S2)
	cr.AssertIsEqual(sum, c.Y)
	cr.scalarApi.AssertIsEqual(
		cr.scalarApi.Add(c.S1, c.S2),
		emulated.NewElement[S](big.NewInt(1)),
	)
	return nil
}

func TestGeneratorGeneratesCurveOfCorrectOrder(t *testing.T) {
	assert := test.NewAssert(t)
	Gx, Gy := Ed25519{}.Generator()
	G := AffinePoint[Ed25519]{
		X: emulated.NewElement[Ed25519](Gx),
		Y: emulated.NewElement[Ed25519](Gy),
	}
	for i := 2; i <= 3; i++ {
		S1 := new(big.Int).Sub(rEd25519, big.NewInt(int64(i - 1)))
		S2 := big.NewInt(int64(i))
		circuit := MulAddTest[Ed25519, Ed25519Scalars]{}
		witness := MulAddTest[Ed25519, Ed25519Scalars]{
			X1: G,
			X2: G,
			S1: emulated.NewElement[Ed25519Scalars](S1),
			S2: emulated.NewElement[Ed25519Scalars](S2),
			Y: G,
		}
		err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
		assert.NoError(err)
	}
}

var testCurve = ecc.BN254

// type NegTest[T, S emulated.FieldParams] struct {
// 	P, Q AffinePoint[T]
// }

// func (c *NegTest[T, S]) Define(api frontend.API) error {
// 	cr, err := New[T, S](api)
// 	if err != nil {
// 		return err
// 	}
// 	res := cr.Neg(c.P)
// 	cr.AssertIsEqual(res, c.Q)
// 	return nil
// }

// func TestNeg(t *testing.T) {
// 	assert := test.NewAssert(t)
// 	secpCurve := secp256k1.S256()
// 	yn := new(big.Int).Sub(secpCurve.P, secpCurve.Gy)
// 	circuit := NegTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{}
// 	witness := NegTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{
// 		P: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](secpCurve.Gx),
// 			Y: emulated.NewElement[emulated.Secp256k1](secpCurve.Gy),
// 		},
// 		Q: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](secpCurve.Gx),
// 			Y: emulated.NewElement[emulated.Secp256k1](yn),
// 		},
// 	}
// 	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
// 	assert.NoError(err)
// }

// type AddTest[T, S emulated.FieldParams] struct {
// 	P, Q, R AffinePoint[T]
// }

// func (c *AddTest[T, S]) Define(api frontend.API) error {
// 	cr, err := New[T, S](api)
// 	if err != nil {
// 		return err
// 	}
// 	res := cr.Add(c.P, c.Q)
// 	cr.AssertIsEqual(res, c.R)
// 	return nil
// }

// func TestAdd(t *testing.T) {
// 	assert := test.NewAssert(t)
// 	secpCurve := secp256k1.S256()
// 	xd, yd := secpCurve.Double(secpCurve.Gx, secpCurve.Gy)
// 	xa, ya := secpCurve.Add(xd, yd, secpCurve.Gx, secpCurve.Gy)
// 	circuit := AddTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{}
// 	witness := AddTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{
// 		P: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](secpCurve.Gx),
// 			Y: emulated.NewElement[emulated.Secp256k1](secpCurve.Gy),
// 		},
// 		Q: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](xd),
// 			Y: emulated.NewElement[emulated.Secp256k1](yd),
// 		},
// 		R: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](xa),
// 			Y: emulated.NewElement[emulated.Secp256k1](ya),
// 		},
// 	}
// 	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
// 	assert.NoError(err)
// }

// type DoubleTest[T, S emulated.FieldParams] struct {
// 	P, Q AffinePoint[T]
// }

// func (c *DoubleTest[T, S]) Define(api frontend.API) error {
// 	cr, err := New[T, S](api)
// 	if err != nil {
// 		return err
// 	}
// 	res := cr.Double(c.P)
// 	cr.AssertIsEqual(res, c.Q)
// 	return nil
// }

// func TestDouble(t *testing.T) {
// 	assert := test.NewAssert(t)
// 	secpCurve := secp256k1.S256()
// 	xd, yd := secpCurve.Double(secpCurve.Gx, secpCurve.Gy)
// 	circuit := DoubleTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{}
// 	witness := DoubleTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{
// 		P: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](secpCurve.Gx),
// 			Y: emulated.NewElement[emulated.Secp256k1](secpCurve.Gy),
// 		},
// 		Q: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](xd),
// 			Y: emulated.NewElement[emulated.Secp256k1](yd),
// 		},
// 	}
// 	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
// 	assert.NoError(err)
// }

// type ScalarMulTest[T, S emulated.FieldParams] struct {
// 	P, Q AffinePoint[T]
// 	S    emulated.Element[S]
// }

// func (c *ScalarMulTest[T, S]) Define(api frontend.API) error {
// 	cr, err := New[T, S](api)
// 	if err != nil {
// 		return err
// 	}
// 	res := cr.ScalarMul(c.P, c.S)
// 	cr.AssertIsEqual(res, c.Q)
// 	return nil
// }

// func TestScalarMul(t *testing.T) {
// 	assert := test.NewAssert(t)
// 	secpCurve := secp256k1.S256()
// 	s, ok := new(big.Int).SetString("44693544921776318736021182399461740191514036429448770306966433218654680512345", 10)
// 	assert.True(ok)
// 	sx, sy := secpCurve.ScalarMult(secpCurve.Gx, secpCurve.Gy, s.Bytes())

// 	circuit := ScalarMulTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{}
// 	witness := ScalarMulTest[emulated.Secp256k1, emulated.Secp256k1Scalars]{
// 		S: emulated.NewElement[emulated.Secp256k1Scalars](s),
// 		P: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](secpCurve.Gx),
// 			Y: emulated.NewElement[emulated.Secp256k1](secpCurve.Gy),
// 		},
// 		Q: AffinePoint[emulated.Secp256k1]{
// 			X: emulated.NewElement[emulated.Secp256k1](sx),
// 			Y: emulated.NewElement[emulated.Secp256k1](sy),
// 		},
// 	}
// 	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
// 	assert.NoError(err)
// 	// _, err = frontend.Compile(testCurve.ScalarField(), r1cs.NewBuilder, &circuit)
// 	// assert.NoError(err)
// }