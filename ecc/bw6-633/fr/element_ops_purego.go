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

package fr

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
		8178485296672800069,
		8476448362227282520,
		14180928431697993131,
		4308307642551989706,
		120359802761433421,
	}
	x.Mul(x, &y)
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

	// Implements CIOS multiplication -- section 2.3.2 of Tolga Acar's thesis
	// https://www.microsoft.com/en-us/research/wp-content/uploads/1998/06/97Acar.pdf
	//
	// The algorithm:
	//
	// for i=0 to N-1
	// 		C := 0
	// 		for j=0 to N-1
	// 			(C,t[j]) := t[j] + x[j]*y[i] + C
	// 		(t[N+1],t[N]) := t[N] + C
	//
	// 		C := 0
	// 		m := t[0]*q'[0] mod D
	// 		(C,_) := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 			(C,t[j-1]) := t[j] + m*q[j] + C
	//
	// 		(C,t[N-1]) := t[N] + C
	// 		t[N] := t[N+1] + C
	//
	// → N is the number of machine words needed to store the modulus q
	// → D is the word size. For example, on a 64-bit architecture D is 2	64
	// → x[i], y[i], q[i] is the ith word of the numbers x,y,q
	// → q'[0] is the lowest word of the number -q⁻¹ mod r. This quantity is pre-computed, as it does not depend on the inputs.
	// → t is a temporary array of size N+2
	// → C, S are machine words. A pair (C,S) refers to (hi-bits, lo-bits) of a two-word number
	//
	// As described here https://hackmd.io/@gnark/modular_multiplication we can get rid of one carry chain and simplify:
	// (also described in https://eprint.iacr.org/2022/1400.pdf annex)
	//
	// for i=0 to N-1
	// 		(A,t[0]) := t[0] + x[0]*y[i]
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 			(A,t[j])  := t[j] + x[j]*y[i] + A
	// 			(C,t[j-1]) := t[j] + m*q[j] + C
	//
	// 		t[N-1] = C + A
	//
	// This optimization saves 5N + 2 additions in the algorithm, and can be used whenever the highest bit
	// of the modulus is zero (and not all of the remaining bits are set).

	var t0, t1, t2, t3, t4 uint64
	var u0, u1, u2, u3, u4 uint64
	{
		var c0, c1, c2 uint64
		v := x[0]
		u0, t0 = bits.Mul64(v, y[0])
		u1, t1 = bits.Mul64(v, y[1])
		u2, t2 = bits.Mul64(v, y[2])
		u3, t3 = bits.Mul64(v, y[3])
		u4, t4 = bits.Mul64(v, y[4])
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, 0, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[1]
		u0, c1 = bits.Mul64(v, y[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, y[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, y[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, y[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, y[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[2]
		u0, c1 = bits.Mul64(v, y[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, y[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, y[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, y[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, y[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[3]
		u0, c1 = bits.Mul64(v, y[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, y[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, y[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, y[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, y[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[4]
		u0, c1 = bits.Mul64(v, y[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, y[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, y[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, y[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, y[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	z[0] = t0
	z[1] = t1
	z[2] = t2
	z[3] = t3
	z[4] = t4

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], b = bits.Sub64(z[2], q2, b)
		z[3], b = bits.Sub64(z[3], q3, b)
		z[4], _ = bits.Sub64(z[4], q4, b)
	}
	return z
}

// Square z = x * x (mod q)
//
// x must be strictly inferior to q
func (z *Element) Square(x *Element) *Element {
	// see Mul for algorithm documentation

	var t0, t1, t2, t3, t4 uint64
	var u0, u1, u2, u3, u4 uint64
	{
		var c0, c1, c2 uint64
		v := x[0]
		u0, t0 = bits.Mul64(v, x[0])
		u1, t1 = bits.Mul64(v, x[1])
		u2, t2 = bits.Mul64(v, x[2])
		u3, t3 = bits.Mul64(v, x[3])
		u4, t4 = bits.Mul64(v, x[4])
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, 0, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[1]
		u0, c1 = bits.Mul64(v, x[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, x[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, x[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, x[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, x[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[2]
		u0, c1 = bits.Mul64(v, x[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, x[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, x[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, x[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, x[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[3]
		u0, c1 = bits.Mul64(v, x[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, x[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, x[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, x[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, x[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	{
		var c0, c1, c2 uint64
		v := x[4]
		u0, c1 = bits.Mul64(v, x[0])
		t0, c0 = bits.Add64(c1, t0, 0)
		u1, c1 = bits.Mul64(v, x[1])
		t1, c0 = bits.Add64(c1, t1, c0)
		u2, c1 = bits.Mul64(v, x[2])
		t2, c0 = bits.Add64(c1, t2, c0)
		u3, c1 = bits.Mul64(v, x[3])
		t3, c0 = bits.Add64(c1, t3, c0)
		u4, c1 = bits.Mul64(v, x[4])
		t4, c0 = bits.Add64(c1, t4, c0)

		c2, _ = bits.Add64(0, 0, c0)
		t1, c0 = bits.Add64(u0, t1, 0)
		t2, c0 = bits.Add64(u1, t2, c0)
		t3, c0 = bits.Add64(u2, t3, c0)
		t4, c0 = bits.Add64(u3, t4, c0)
		c2, _ = bits.Add64(u4, c2, c0)

		m := qInvNeg * t0

		u0, c1 = bits.Mul64(m, q0)
		_, c0 = bits.Add64(t0, c1, 0)
		u1, c1 = bits.Mul64(m, q1)
		t0, c0 = bits.Add64(t1, c1, c0)
		u2, c1 = bits.Mul64(m, q2)
		t1, c0 = bits.Add64(t2, c1, c0)
		u3, c1 = bits.Mul64(m, q3)
		t2, c0 = bits.Add64(t3, c1, c0)
		u4, c1 = bits.Mul64(m, q4)

		t3, c0 = bits.Add64(0, c1, c0)
		u4, _ = bits.Add64(u4, 0, c0)
		t0, c0 = bits.Add64(u0, t0, 0)
		t1, c0 = bits.Add64(u1, t1, c0)
		t2, c0 = bits.Add64(u2, t2, c0)
		t3, c0 = bits.Add64(u3, t3, c0)
		c2, _ = bits.Add64(c2, 0, c0)
		t3, c0 = bits.Add64(t4, t3, 0)
		t4, _ = bits.Add64(u4, c2, c0)

	}
	z[0] = t0
	z[1] = t1
	z[2] = t2
	z[3] = t3
	z[4] = t4

	// if z ⩾ q → z -= q
	if !z.smallerThanModulus() {
		var b uint64
		z[0], b = bits.Sub64(z[0], q0, 0)
		z[1], b = bits.Sub64(z[1], q1, b)
		z[2], b = bits.Sub64(z[2], q2, b)
		z[3], b = bits.Sub64(z[3], q3, b)
		z[4], _ = bits.Sub64(z[4], q4, b)
	}
	return z
}
