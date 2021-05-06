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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package mockcommitment

import (
	"io"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/polynomial"
	"github.com/consensys/gnark-crypto/polynomial"
)

// Scheme mock commitment, useful for testing polynomial based IOP
// like PLONK, where the scheme should not depend on which polynomial commitment scheme
// is used.
type Scheme struct{}

// WriteTo panics
func (s *Scheme) WriteTo(w io.Writer) (n int64, err error) {
	panic("not implemented")
}

// ReadFrom panics
func (s *Scheme) ReadFrom(r io.Reader) (n int64, err error) {
	panic("not implemented")
}

// Commit returns the first coefficient of p
func (s *Scheme) Commit(p polynomial.Polynomial) polynomial.Digest {
	_p := p.(*bn254.Polynomial)
	var res fr.Element
	res.Set(&(*_p)[0])
	return &res
}

// Open computes an opening proof of _p at _val.
// Returns a MockProof, which is an empty interface.
func (s *Scheme) Open(point interface{}, p polynomial.Polynomial) polynomial.OpeningProof {

	claimedValue := p.Eval(point)
	var _claimedValue fr.Element
	_claimedValue.Set(claimedValue.(*fr.Element))

	// ensure we get a copy of point
	var _point fr.Element
	_point.Set(point.(*fr.Element))

	return &MockProof{
		Point:        _point,
		ClaimedValue: _claimedValue,
	}
}

// Verify mock implementation of verify
func (s *Scheme) Verify(commitment polynomial.Digest, proof polynomial.OpeningProof) error {
	return nil
}

// BatchOpenSinglePoint computes a batch opening proof for _p at _val.
func (s *Scheme) BatchOpenSinglePoint(point interface{}, digests []polynomial.Digest, polynomials []polynomial.Polynomial) polynomial.BatchOpeningProofSinglePoint {

	var res MockBatchProofsSinglePoint
	res.ClaimedValues = make([]fr.Element, len(polynomials))
	res.Point.Set(point.(*fr.Element))

	for i := 0; i < len(polynomials); i++ {
		claimedValue := polynomials[i].Eval(point)
		res.ClaimedValues[i].Set(claimedValue.(*fr.Element))
	}

	return &res
}

// BatchVerifySinglePoint computes a batch opening proof for
func (s *Scheme) BatchVerifySinglePoint(digests []polynomial.Digest, batchOpeningProof polynomial.BatchOpeningProofSinglePoint) error {

	return nil

}
