//go:build !amd64 || purego
// +build !amd64 purego

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

package fp

import "math/bits"

// MulBy3 x *= 3 (mod q)
func MulBy3(x *Element) {
	_x := *x
	x.Double(x).Add(x, &_x)
}

// MulBy5 x *= 5 (mod q)
func MulBy5(x *Element) {
	_x := *x
	x.Double(x).Double(x).Add(x, &_x)
}

// MulBy13 x *= 13 (mod q)
func MulBy13(x *Element) {
	var y = Element{
		4345973640412121648,
		16340807117537158706,
		14673764841507373218,
		5587754667198343811,
		12846753860245084942,
		4041391838244625385,
		8324122986343791677,
		8773809490091176420,
		5465994123296109449,
		6649773564661156048,
		9147430723089113754,
		54281803719730243,
	}
	x.Mul(x, &y)
}

// Butterfly sets
//
//	a = a + b (mod q)
//	b = a - b (mod q)
func Butterfly(a, b *Element) {
	_butterflyGeneric(a, b)
}

func fromMont(z *Element) {
	_fromMontGeneric(z)
}

func reduce(z *Element) {
	_reduceGeneric(z)
}

// Mul z = x * y (mod q)
//
// x and y must be strictly inferior to q
func (z *Element) Mul(x, y *Element) *Element {
	// TODO @gbotrel restore doc

	var t [13]uint64
	var D uint64
	var m, C uint64
	// -----------------------------------
	// First loop

	C, t[0] = bits.Mul64(y[0], x[0])
	C, t[1] = madd1(y[0], x[1], C)
	C, t[2] = madd1(y[0], x[2], C)
	C, t[3] = madd1(y[0], x[3], C)
	C, t[4] = madd1(y[0], x[4], C)
	C, t[5] = madd1(y[0], x[5], C)
	C, t[6] = madd1(y[0], x[6], C)
	C, t[7] = madd1(y[0], x[7], C)
	C, t[8] = madd1(y[0], x[8], C)
	C, t[9] = madd1(y[0], x[9], C)
	C, t[10] = madd1(y[0], x[10], C)
	C, t[11] = madd1(y[0], x[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[1], x[0], t[0])
	C, t[1] = madd2(y[1], x[1], t[1], C)
	C, t[2] = madd2(y[1], x[2], t[2], C)
	C, t[3] = madd2(y[1], x[3], t[3], C)
	C, t[4] = madd2(y[1], x[4], t[4], C)
	C, t[5] = madd2(y[1], x[5], t[5], C)
	C, t[6] = madd2(y[1], x[6], t[6], C)
	C, t[7] = madd2(y[1], x[7], t[7], C)
	C, t[8] = madd2(y[1], x[8], t[8], C)
	C, t[9] = madd2(y[1], x[9], t[9], C)
	C, t[10] = madd2(y[1], x[10], t[10], C)
	C, t[11] = madd2(y[1], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[2], x[0], t[0])
	C, t[1] = madd2(y[2], x[1], t[1], C)
	C, t[2] = madd2(y[2], x[2], t[2], C)
	C, t[3] = madd2(y[2], x[3], t[3], C)
	C, t[4] = madd2(y[2], x[4], t[4], C)
	C, t[5] = madd2(y[2], x[5], t[5], C)
	C, t[6] = madd2(y[2], x[6], t[6], C)
	C, t[7] = madd2(y[2], x[7], t[7], C)
	C, t[8] = madd2(y[2], x[8], t[8], C)
	C, t[9] = madd2(y[2], x[9], t[9], C)
	C, t[10] = madd2(y[2], x[10], t[10], C)
	C, t[11] = madd2(y[2], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[3], x[0], t[0])
	C, t[1] = madd2(y[3], x[1], t[1], C)
	C, t[2] = madd2(y[3], x[2], t[2], C)
	C, t[3] = madd2(y[3], x[3], t[3], C)
	C, t[4] = madd2(y[3], x[4], t[4], C)
	C, t[5] = madd2(y[3], x[5], t[5], C)
	C, t[6] = madd2(y[3], x[6], t[6], C)
	C, t[7] = madd2(y[3], x[7], t[7], C)
	C, t[8] = madd2(y[3], x[8], t[8], C)
	C, t[9] = madd2(y[3], x[9], t[9], C)
	C, t[10] = madd2(y[3], x[10], t[10], C)
	C, t[11] = madd2(y[3], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[4], x[0], t[0])
	C, t[1] = madd2(y[4], x[1], t[1], C)
	C, t[2] = madd2(y[4], x[2], t[2], C)
	C, t[3] = madd2(y[4], x[3], t[3], C)
	C, t[4] = madd2(y[4], x[4], t[4], C)
	C, t[5] = madd2(y[4], x[5], t[5], C)
	C, t[6] = madd2(y[4], x[6], t[6], C)
	C, t[7] = madd2(y[4], x[7], t[7], C)
	C, t[8] = madd2(y[4], x[8], t[8], C)
	C, t[9] = madd2(y[4], x[9], t[9], C)
	C, t[10] = madd2(y[4], x[10], t[10], C)
	C, t[11] = madd2(y[4], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[5], x[0], t[0])
	C, t[1] = madd2(y[5], x[1], t[1], C)
	C, t[2] = madd2(y[5], x[2], t[2], C)
	C, t[3] = madd2(y[5], x[3], t[3], C)
	C, t[4] = madd2(y[5], x[4], t[4], C)
	C, t[5] = madd2(y[5], x[5], t[5], C)
	C, t[6] = madd2(y[5], x[6], t[6], C)
	C, t[7] = madd2(y[5], x[7], t[7], C)
	C, t[8] = madd2(y[5], x[8], t[8], C)
	C, t[9] = madd2(y[5], x[9], t[9], C)
	C, t[10] = madd2(y[5], x[10], t[10], C)
	C, t[11] = madd2(y[5], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[6], x[0], t[0])
	C, t[1] = madd2(y[6], x[1], t[1], C)
	C, t[2] = madd2(y[6], x[2], t[2], C)
	C, t[3] = madd2(y[6], x[3], t[3], C)
	C, t[4] = madd2(y[6], x[4], t[4], C)
	C, t[5] = madd2(y[6], x[5], t[5], C)
	C, t[6] = madd2(y[6], x[6], t[6], C)
	C, t[7] = madd2(y[6], x[7], t[7], C)
	C, t[8] = madd2(y[6], x[8], t[8], C)
	C, t[9] = madd2(y[6], x[9], t[9], C)
	C, t[10] = madd2(y[6], x[10], t[10], C)
	C, t[11] = madd2(y[6], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[7], x[0], t[0])
	C, t[1] = madd2(y[7], x[1], t[1], C)
	C, t[2] = madd2(y[7], x[2], t[2], C)
	C, t[3] = madd2(y[7], x[3], t[3], C)
	C, t[4] = madd2(y[7], x[4], t[4], C)
	C, t[5] = madd2(y[7], x[5], t[5], C)
	C, t[6] = madd2(y[7], x[6], t[6], C)
	C, t[7] = madd2(y[7], x[7], t[7], C)
	C, t[8] = madd2(y[7], x[8], t[8], C)
	C, t[9] = madd2(y[7], x[9], t[9], C)
	C, t[10] = madd2(y[7], x[10], t[10], C)
	C, t[11] = madd2(y[7], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[8], x[0], t[0])
	C, t[1] = madd2(y[8], x[1], t[1], C)
	C, t[2] = madd2(y[8], x[2], t[2], C)
	C, t[3] = madd2(y[8], x[3], t[3], C)
	C, t[4] = madd2(y[8], x[4], t[4], C)
	C, t[5] = madd2(y[8], x[5], t[5], C)
	C, t[6] = madd2(y[8], x[6], t[6], C)
	C, t[7] = madd2(y[8], x[7], t[7], C)
	C, t[8] = madd2(y[8], x[8], t[8], C)
	C, t[9] = madd2(y[8], x[9], t[9], C)
	C, t[10] = madd2(y[8], x[10], t[10], C)
	C, t[11] = madd2(y[8], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[9], x[0], t[0])
	C, t[1] = madd2(y[9], x[1], t[1], C)
	C, t[2] = madd2(y[9], x[2], t[2], C)
	C, t[3] = madd2(y[9], x[3], t[3], C)
	C, t[4] = madd2(y[9], x[4], t[4], C)
	C, t[5] = madd2(y[9], x[5], t[5], C)
	C, t[6] = madd2(y[9], x[6], t[6], C)
	C, t[7] = madd2(y[9], x[7], t[7], C)
	C, t[8] = madd2(y[9], x[8], t[8], C)
	C, t[9] = madd2(y[9], x[9], t[9], C)
	C, t[10] = madd2(y[9], x[10], t[10], C)
	C, t[11] = madd2(y[9], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[10], x[0], t[0])
	C, t[1] = madd2(y[10], x[1], t[1], C)
	C, t[2] = madd2(y[10], x[2], t[2], C)
	C, t[3] = madd2(y[10], x[3], t[3], C)
	C, t[4] = madd2(y[10], x[4], t[4], C)
	C, t[5] = madd2(y[10], x[5], t[5], C)
	C, t[6] = madd2(y[10], x[6], t[6], C)
	C, t[7] = madd2(y[10], x[7], t[7], C)
	C, t[8] = madd2(y[10], x[8], t[8], C)
	C, t[9] = madd2(y[10], x[9], t[9], C)
	C, t[10] = madd2(y[10], x[10], t[10], C)
	C, t[11] = madd2(y[10], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(y[11], x[0], t[0])
	C, t[1] = madd2(y[11], x[1], t[1], C)
	C, t[2] = madd2(y[11], x[2], t[2], C)
	C, t[3] = madd2(y[11], x[3], t[3], C)
	C, t[4] = madd2(y[11], x[4], t[4], C)
	C, t[5] = madd2(y[11], x[5], t[5], C)
	C, t[6] = madd2(y[11], x[6], t[6], C)
	C, t[7] = madd2(y[11], x[7], t[7], C)
	C, t[8] = madd2(y[11], x[8], t[8], C)
	C, t[9] = madd2(y[11], x[9], t[9], C)
	C, t[10] = madd2(y[11], x[10], t[10], C)
	C, t[11] = madd2(y[11], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)

	if t[12] != 0 {
		// we need to reduce, we have a result on 13 words
		var b uint64
		z[0], b = bits.Sub64(t[0], q0, 0)
		z[1], b = bits.Sub64(t[1], q1, b)
		z[2], b = bits.Sub64(t[2], q2, b)
		z[3], b = bits.Sub64(t[3], q3, b)
		z[4], b = bits.Sub64(t[4], q4, b)
		z[5], b = bits.Sub64(t[5], q5, b)
		z[6], b = bits.Sub64(t[6], q6, b)
		z[7], b = bits.Sub64(t[7], q7, b)
		z[8], b = bits.Sub64(t[8], q8, b)
		z[9], b = bits.Sub64(t[9], q9, b)
		z[10], b = bits.Sub64(t[10], q10, b)
		z[11], _ = bits.Sub64(t[11], q11, b)
		return z
	}

	// copy t into z
	z[0] = t[0]
	z[1] = t[1]
	z[2] = t[2]
	z[3] = t[3]
	z[4] = t[4]
	z[5] = t[5]
	z[6] = t[6]
	z[7] = t[7]
	z[8] = t[8]
	z[9] = t[9]
	z[10] = t[10]
	z[11] = t[11]

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], b = bits.Sub64(z[2], q2, b)
		z[3], b = bits.Sub64(z[3], q3, b)
		z[4], b = bits.Sub64(z[4], q4, b)
		z[5], b = bits.Sub64(z[5], q5, b)
		z[6], b = bits.Sub64(z[6], q6, b)
		z[7], b = bits.Sub64(z[7], q7, b)
		z[8], b = bits.Sub64(z[8], q8, b)
		z[9], b = bits.Sub64(z[9], q9, b)
		z[10], b = bits.Sub64(z[10], q10, b)
		z[11], _ = bits.Sub64(z[11], q11, b)
	}
	return z
}

// Square z = x * x (mod q)
//
// x must be strictly inferior to q
func (z *Element) Square(x *Element) *Element {
	// see Mul for algorithm documentation

	var t [13]uint64
	var D uint64
	var m, C uint64
	// -----------------------------------
	// First loop

	C, t[0] = bits.Mul64(x[0], x[0])
	C, t[1] = madd1(x[0], x[1], C)
	C, t[2] = madd1(x[0], x[2], C)
	C, t[3] = madd1(x[0], x[3], C)
	C, t[4] = madd1(x[0], x[4], C)
	C, t[5] = madd1(x[0], x[5], C)
	C, t[6] = madd1(x[0], x[6], C)
	C, t[7] = madd1(x[0], x[7], C)
	C, t[8] = madd1(x[0], x[8], C)
	C, t[9] = madd1(x[0], x[9], C)
	C, t[10] = madd1(x[0], x[10], C)
	C, t[11] = madd1(x[0], x[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[1], x[0], t[0])
	C, t[1] = madd2(x[1], x[1], t[1], C)
	C, t[2] = madd2(x[1], x[2], t[2], C)
	C, t[3] = madd2(x[1], x[3], t[3], C)
	C, t[4] = madd2(x[1], x[4], t[4], C)
	C, t[5] = madd2(x[1], x[5], t[5], C)
	C, t[6] = madd2(x[1], x[6], t[6], C)
	C, t[7] = madd2(x[1], x[7], t[7], C)
	C, t[8] = madd2(x[1], x[8], t[8], C)
	C, t[9] = madd2(x[1], x[9], t[9], C)
	C, t[10] = madd2(x[1], x[10], t[10], C)
	C, t[11] = madd2(x[1], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[2], x[0], t[0])
	C, t[1] = madd2(x[2], x[1], t[1], C)
	C, t[2] = madd2(x[2], x[2], t[2], C)
	C, t[3] = madd2(x[2], x[3], t[3], C)
	C, t[4] = madd2(x[2], x[4], t[4], C)
	C, t[5] = madd2(x[2], x[5], t[5], C)
	C, t[6] = madd2(x[2], x[6], t[6], C)
	C, t[7] = madd2(x[2], x[7], t[7], C)
	C, t[8] = madd2(x[2], x[8], t[8], C)
	C, t[9] = madd2(x[2], x[9], t[9], C)
	C, t[10] = madd2(x[2], x[10], t[10], C)
	C, t[11] = madd2(x[2], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[3], x[0], t[0])
	C, t[1] = madd2(x[3], x[1], t[1], C)
	C, t[2] = madd2(x[3], x[2], t[2], C)
	C, t[3] = madd2(x[3], x[3], t[3], C)
	C, t[4] = madd2(x[3], x[4], t[4], C)
	C, t[5] = madd2(x[3], x[5], t[5], C)
	C, t[6] = madd2(x[3], x[6], t[6], C)
	C, t[7] = madd2(x[3], x[7], t[7], C)
	C, t[8] = madd2(x[3], x[8], t[8], C)
	C, t[9] = madd2(x[3], x[9], t[9], C)
	C, t[10] = madd2(x[3], x[10], t[10], C)
	C, t[11] = madd2(x[3], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[4], x[0], t[0])
	C, t[1] = madd2(x[4], x[1], t[1], C)
	C, t[2] = madd2(x[4], x[2], t[2], C)
	C, t[3] = madd2(x[4], x[3], t[3], C)
	C, t[4] = madd2(x[4], x[4], t[4], C)
	C, t[5] = madd2(x[4], x[5], t[5], C)
	C, t[6] = madd2(x[4], x[6], t[6], C)
	C, t[7] = madd2(x[4], x[7], t[7], C)
	C, t[8] = madd2(x[4], x[8], t[8], C)
	C, t[9] = madd2(x[4], x[9], t[9], C)
	C, t[10] = madd2(x[4], x[10], t[10], C)
	C, t[11] = madd2(x[4], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[5], x[0], t[0])
	C, t[1] = madd2(x[5], x[1], t[1], C)
	C, t[2] = madd2(x[5], x[2], t[2], C)
	C, t[3] = madd2(x[5], x[3], t[3], C)
	C, t[4] = madd2(x[5], x[4], t[4], C)
	C, t[5] = madd2(x[5], x[5], t[5], C)
	C, t[6] = madd2(x[5], x[6], t[6], C)
	C, t[7] = madd2(x[5], x[7], t[7], C)
	C, t[8] = madd2(x[5], x[8], t[8], C)
	C, t[9] = madd2(x[5], x[9], t[9], C)
	C, t[10] = madd2(x[5], x[10], t[10], C)
	C, t[11] = madd2(x[5], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[6], x[0], t[0])
	C, t[1] = madd2(x[6], x[1], t[1], C)
	C, t[2] = madd2(x[6], x[2], t[2], C)
	C, t[3] = madd2(x[6], x[3], t[3], C)
	C, t[4] = madd2(x[6], x[4], t[4], C)
	C, t[5] = madd2(x[6], x[5], t[5], C)
	C, t[6] = madd2(x[6], x[6], t[6], C)
	C, t[7] = madd2(x[6], x[7], t[7], C)
	C, t[8] = madd2(x[6], x[8], t[8], C)
	C, t[9] = madd2(x[6], x[9], t[9], C)
	C, t[10] = madd2(x[6], x[10], t[10], C)
	C, t[11] = madd2(x[6], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[7], x[0], t[0])
	C, t[1] = madd2(x[7], x[1], t[1], C)
	C, t[2] = madd2(x[7], x[2], t[2], C)
	C, t[3] = madd2(x[7], x[3], t[3], C)
	C, t[4] = madd2(x[7], x[4], t[4], C)
	C, t[5] = madd2(x[7], x[5], t[5], C)
	C, t[6] = madd2(x[7], x[6], t[6], C)
	C, t[7] = madd2(x[7], x[7], t[7], C)
	C, t[8] = madd2(x[7], x[8], t[8], C)
	C, t[9] = madd2(x[7], x[9], t[9], C)
	C, t[10] = madd2(x[7], x[10], t[10], C)
	C, t[11] = madd2(x[7], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[8], x[0], t[0])
	C, t[1] = madd2(x[8], x[1], t[1], C)
	C, t[2] = madd2(x[8], x[2], t[2], C)
	C, t[3] = madd2(x[8], x[3], t[3], C)
	C, t[4] = madd2(x[8], x[4], t[4], C)
	C, t[5] = madd2(x[8], x[5], t[5], C)
	C, t[6] = madd2(x[8], x[6], t[6], C)
	C, t[7] = madd2(x[8], x[7], t[7], C)
	C, t[8] = madd2(x[8], x[8], t[8], C)
	C, t[9] = madd2(x[8], x[9], t[9], C)
	C, t[10] = madd2(x[8], x[10], t[10], C)
	C, t[11] = madd2(x[8], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[9], x[0], t[0])
	C, t[1] = madd2(x[9], x[1], t[1], C)
	C, t[2] = madd2(x[9], x[2], t[2], C)
	C, t[3] = madd2(x[9], x[3], t[3], C)
	C, t[4] = madd2(x[9], x[4], t[4], C)
	C, t[5] = madd2(x[9], x[5], t[5], C)
	C, t[6] = madd2(x[9], x[6], t[6], C)
	C, t[7] = madd2(x[9], x[7], t[7], C)
	C, t[8] = madd2(x[9], x[8], t[8], C)
	C, t[9] = madd2(x[9], x[9], t[9], C)
	C, t[10] = madd2(x[9], x[10], t[10], C)
	C, t[11] = madd2(x[9], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[10], x[0], t[0])
	C, t[1] = madd2(x[10], x[1], t[1], C)
	C, t[2] = madd2(x[10], x[2], t[2], C)
	C, t[3] = madd2(x[10], x[3], t[3], C)
	C, t[4] = madd2(x[10], x[4], t[4], C)
	C, t[5] = madd2(x[10], x[5], t[5], C)
	C, t[6] = madd2(x[10], x[6], t[6], C)
	C, t[7] = madd2(x[10], x[7], t[7], C)
	C, t[8] = madd2(x[10], x[8], t[8], C)
	C, t[9] = madd2(x[10], x[9], t[9], C)
	C, t[10] = madd2(x[10], x[10], t[10], C)
	C, t[11] = madd2(x[10], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)
	// -----------------------------------
	// First loop

	C, t[0] = madd1(x[11], x[0], t[0])
	C, t[1] = madd2(x[11], x[1], t[1], C)
	C, t[2] = madd2(x[11], x[2], t[2], C)
	C, t[3] = madd2(x[11], x[3], t[3], C)
	C, t[4] = madd2(x[11], x[4], t[4], C)
	C, t[5] = madd2(x[11], x[5], t[5], C)
	C, t[6] = madd2(x[11], x[6], t[6], C)
	C, t[7] = madd2(x[11], x[7], t[7], C)
	C, t[8] = madd2(x[11], x[8], t[8], C)
	C, t[9] = madd2(x[11], x[9], t[9], C)
	C, t[10] = madd2(x[11], x[10], t[10], C)
	C, t[11] = madd2(x[11], x[11], t[11], C)

	t[12], D = bits.Add64(t[12], C, 0)

	// m = t[0]n'[0] mod W
	m = t[0] * qInvNeg

	// -----------------------------------
	// Second loop
	C = madd0(m, q0, t[0])
	C, t[0] = madd2(m, q1, t[1], C)
	C, t[1] = madd2(m, q2, t[2], C)
	C, t[2] = madd2(m, q3, t[3], C)
	C, t[3] = madd2(m, q4, t[4], C)
	C, t[4] = madd2(m, q5, t[5], C)
	C, t[5] = madd2(m, q6, t[6], C)
	C, t[6] = madd2(m, q7, t[7], C)
	C, t[7] = madd2(m, q8, t[8], C)
	C, t[8] = madd2(m, q9, t[9], C)
	C, t[9] = madd2(m, q10, t[10], C)
	C, t[10] = madd2(m, q11, t[11], C)

	t[11], C = bits.Add64(t[12], C, 0)
	t[12], _ = bits.Add64(0, D, C)

	if t[12] != 0 {
		// we need to reduce, we have a result on 13 words
		var b uint64
		z[0], b = bits.Sub64(t[0], q0, 0)
		z[1], b = bits.Sub64(t[1], q1, b)
		z[2], b = bits.Sub64(t[2], q2, b)
		z[3], b = bits.Sub64(t[3], q3, b)
		z[4], b = bits.Sub64(t[4], q4, b)
		z[5], b = bits.Sub64(t[5], q5, b)
		z[6], b = bits.Sub64(t[6], q6, b)
		z[7], b = bits.Sub64(t[7], q7, b)
		z[8], b = bits.Sub64(t[8], q8, b)
		z[9], b = bits.Sub64(t[9], q9, b)
		z[10], b = bits.Sub64(t[10], q10, b)
		z[11], _ = bits.Sub64(t[11], q11, b)
		return z
	}

	// copy t into z
	z[0] = t[0]
	z[1] = t[1]
	z[2] = t[2]
	z[3] = t[3]
	z[4] = t[4]
	z[5] = t[5]
	z[6] = t[6]
	z[7] = t[7]
	z[8] = t[8]
	z[9] = t[9]
	z[10] = t[10]
	z[11] = t[11]

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], b = bits.Sub64(z[2], q2, b)
		z[3], b = bits.Sub64(z[3], q3, b)
		z[4], b = bits.Sub64(z[4], q4, b)
		z[5], b = bits.Sub64(z[5], q5, b)
		z[6], b = bits.Sub64(z[6], q6, b)
		z[7], b = bits.Sub64(z[7], q7, b)
		z[8], b = bits.Sub64(z[8], q8, b)
		z[9], b = bits.Sub64(z[9], q9, b)
		z[10], b = bits.Sub64(z[10], q10, b)
		z[11], _ = bits.Sub64(z[11], q11, b)
	}
	return z
}
