// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package polynomial

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func TestEval(t *testing.T) {

	// build polynomial
	f := make(Polynomial, 20)
	for i := 0; i < 20; i++ {
		f[i].SetOne()
	}

	// random value
	var point fr.Element
	point.SetRandom()

	// compute manually f(val)
	var expectedEval, one, den fr.Element
	var expo big.Int
	one.SetOne()
	expo.SetUint64(20)
	expectedEval.Exp(point, &expo).
		Sub(&expectedEval, &one)
	den.Sub(&point, &one)
	expectedEval.Div(&expectedEval, &den)

	// compute purpoted evaluation
	_purportedEval := f.Eval(&point)
	purportedEval := FromInterface(_purportedEval)

	// check
	if !purportedEval.Equal(&expectedEval) {
		t.Fatal("polynomial evaluation failed")
	}
}

func TestAddConstantInPlace(t *testing.T) {

	// build polynomial
	f := make(Polynomial, 20)
	for i := 0; i < 20; i++ {
		f[i].SetOne()
	}

	// constant to add
	var c fr.Element
	c.SetRandom()

	// add constant
	f.AddConstantInPlace(&c)

	// check
	var expectedCoeffs, one fr.Element
	one.SetOne()
	expectedCoeffs.Add(&one, &c)
	for i := 0; i < 20; i++ {
		if !f[i].Equal(&expectedCoeffs) {
			t.Fatal("AddConstantInPlace failed")
		}
	}
}

func TestSubConstantInPlace(t *testing.T) {

	// build polynomial
	f := make(Polynomial, 20)
	for i := 0; i < 20; i++ {
		f[i].SetOne()
	}

	// constant to sub
	var c fr.Element
	c.SetRandom()

	// sub constant
	f.SubConstantInPlace(&c)

	// check
	var expectedCoeffs, one fr.Element
	one.SetOne()
	expectedCoeffs.Sub(&one, &c)
	for i := 0; i < 20; i++ {
		if !f[i].Equal(&expectedCoeffs) {
			t.Fatal("SubConstantInPlace failed")
		}
	}
}

func TestScaleInPlace(t *testing.T) {

	// build polynomial
	f := make(Polynomial, 20)
	for i := 0; i < 20; i++ {
		f[i].SetOne()
	}

	// constant to scale by
	var c fr.Element
	c.SetRandom()

	// scale by constant
	f.ScaleInPlace(&c)

	// check
	for i := 0; i < 20; i++ {
		if !f[i].Equal(&c) {
			t.Fatal("ScaleInPlace failed")
		}
	}

}

func TestAdd(t *testing.T) {

	// build unbalanced polynomials
	f1 := make(Polynomial, 20)
	for i := 0; i < 20; i++ {
		f1[i].SetOne()
	}
	f2 := make(Polynomial, 10)
	for i := 0; i < 10; i++ {
		f2[i].SetOne()
	}

	// expected result
	var one, two fr.Element
	one.SetOne()
	two.Double(&one)
	expectedSum := make(Polynomial, 20)
	for i := 0; i < 10; i++ {
		expectedSum[i].Set(&two)
	}
	for i := 10; i < 20; i++ {
		expectedSum[i].Set(&one)
	}

	// set different combinations
	var g Polynomial
	g.Add(&f1, &f2)
	if !g.Equal(&expectedSum) {
		t.Fatal("add polynomials fails")
	}

	_f1 := f1.Clone()
	_f1.Add(&f1, &f2)
	if !_f1.Equal(&expectedSum) {
		t.Fatal("add polynomials fails")
	}

	_f2 := f2.Clone()
	_f2.Add(&f1, &f2)
	if !_f2.Equal(&expectedSum) {
		t.Fatal("add polynomials fails")
	}
}
