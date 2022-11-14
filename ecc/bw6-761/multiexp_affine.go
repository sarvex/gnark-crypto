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

package bw6761

import (
	"github.com/consensys/gnark-crypto/ecc/bw6-761/fp"
)

type batchOpG1Affine struct {
	bucketID uint16
	point    G1Affine
}

func (o batchOpG1Affine) isNeg() bool {
	return o.bucketID&1 == 1
}

// processChunkG1BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG1BatchAffine[B ibG1Affine, BS bitSet, TP pG1Affine, TPP ppG1Affine, TQ qOpsG1Affine, TC cG1Affine](
	chunk uint64,
	chRes chan<- g1JacExtended,
	c uint64,
	points []G1Affine,
	digits []uint16) {

	// init the buckets
	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// setup for the batch affine;
	// we do that instead of a separate object to give enough hints to the compiler to..
	var bucketIds BS // bitSet to signify presence of a bucket in current batch
	cptAdd := 0      // count the number of bucket + point added to current batch

	var R TPP // bucket references
	var P TP  // points to be added to R (buckets); it is beneficial to store them on the stack (ie copy)

	batchSize := len(P)

	canAdd := func(bID uint16) bool {
		return !bucketIds[bID]
	}

	isFull := func() bool {
		return (cptAdd) == batchSize
	}

	executeAndReset := func() {
		if (cptAdd) == 0 {
			return
		}
		batchAddG1Affine[TP, TPP, TC](&R, &P, cptAdd)

		var tmp BS
		bucketIds = tmp
		cptAdd = 0
	}

	addFromQueue := func(op batchOpG1Affine) {
		// CanAdd must be called before --> ensures bucket is not "used" in current batch

		BK := &buckets[op.bucketID]
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			BK.Set(&op.point)
			return
		}
		if BK.X.Equal(&op.point.X) {
			if BK.Y.Equal(&op.point.Y) {
				// P + P: doubling, which should be quite rare --
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				BK.Add(BK, BK)
				return
			}
			BK.setInfinity()
			return
		}

		bucketIds[op.bucketID] = true
		R[cptAdd] = BK
		P[cptAdd] = op.point
		cptAdd++
	}

	add := func(bucketID uint16, PP *G1Affine, isAdd bool) {
		// CanAdd must be called before --> ensures bucket is not "used" in current batch

		BK := &buckets[bucketID]
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			if isAdd {
				BK.Set(PP)
			} else {
				BK.Neg(PP)
			}
			return
		}
		if BK.X.Equal(&PP.X) {
			if BK.Y.Equal(&PP.Y) {
				// P + P: doubling, which should be quite rare --
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				if isAdd {
					BK.Add(BK, BK)
				} else {
					BK.setInfinity()
				}

				return
			}
			if isAdd {
				BK.setInfinity()
			} else {
				BK.Add(BK, BK)
			}
			return
		}

		bucketIds[bucketID] = true
		R[cptAdd] = BK
		if isAdd {
			P[cptAdd].Set(PP)
		} else {
			P[cptAdd].Neg(PP)
		}
		cptAdd++
	}

	var queue TQ
	qID := 0

	processQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if !canAdd(queue[i].bucketID) {
				continue
			}
			addFromQueue(queue[i])
			if isFull() {
				executeAndReset()
			}
			queue[i] = queue[qID-1]
			qID--
		}
	}

	processTopQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if !canAdd(queue[i].bucketID) {
				return
			}
			addFromQueue(queue[i])
			if isFull() {
				executeAndReset()
			}
			qID--
		}
	}

	for i, digit := range digits {

		if digit == 0 || points[i].IsInfinity() {
			continue
		}

		bucketID := uint16((digit >> 1))
		isAdd := digit&1 == 0
		if isAdd {
			// add
			bucketID -= 1
		}

		if !canAdd(bucketID) {
			// put it in queue
			queue[qID].bucketID = bucketID
			if isAdd {
				queue[qID].point = points[i]
			} else {
				queue[qID].point.Neg(&points[i])
			}
			qID++

			// queue is full, flush it.
			if qID == len(queue)-1 {
				executeAndReset()
				processQueue()
			}
			continue
		}

		// we add the point to the batch.
		add(bucketID, &points[i], isAdd)
		if isFull() {
			executeAndReset()
			processTopQueue()
		}
	}

	// empty the queue
	for qID != 0 {
		processQueue()
		executeAndReset()
	}

	// flush items in batch.
	executeAndReset()

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g1JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].IsInfinity() {
			runningSum.addMixed(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total

}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketG1AffineC16 [1 << (16 - 1)]G1Affine

// buckets: array of G1Affine points of size 1 << (c-1)
type ibG1Affine interface {
	bucketG1AffineC16
}

// array of coordinates fp.Element
type cG1Affine interface {
	cG1AffineC16
}

// buckets: array of G1Affine points (for the batch addition)
type pG1Affine interface {
	pG1AffineC16
}

// buckets: array of *G1Affine points (for the batch addition)
type ppG1Affine interface {
	ppG1AffineC16
}

// buckets: array of G1Affine queue operations (for the batch addition)
type qOpsG1Affine interface {
	qOpsG1AffineC16
}
type cG1AffineC16 [640]fp.Element
type pG1AffineC16 [640]G1Affine
type ppG1AffineC16 [640]*G1Affine
type qOpsG1AffineC16 [640]batchOpG1Affine

type batchOpG2Affine struct {
	bucketID uint16
	point    G2Affine
}

func (o batchOpG2Affine) isNeg() bool {
	return o.bucketID&1 == 1
}

// processChunkG2BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG2BatchAffine[B ibG2Affine, BS bitSet, TP pG2Affine, TPP ppG2Affine, TQ qOpsG2Affine, TC cG2Affine](
	chunk uint64,
	chRes chan<- g2JacExtended,
	c uint64,
	points []G2Affine,
	digits []uint16) {

	// init the buckets
	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// setup for the batch affine;
	// we do that instead of a separate object to give enough hints to the compiler to..
	var bucketIds BS // bitSet to signify presence of a bucket in current batch
	cptAdd := 0      // count the number of bucket + point added to current batch

	var R TPP // bucket references
	var P TP  // points to be added to R (buckets); it is beneficial to store them on the stack (ie copy)

	batchSize := len(P)

	canAdd := func(bID uint16) bool {
		return !bucketIds[bID]
	}

	isFull := func() bool {
		return (cptAdd) == batchSize
	}

	executeAndReset := func() {
		if (cptAdd) == 0 {
			return
		}
		batchAddG2Affine[TP, TPP, TC](&R, &P, cptAdd)

		var tmp BS
		bucketIds = tmp
		cptAdd = 0
	}

	addFromQueue := func(op batchOpG2Affine) {
		// CanAdd must be called before --> ensures bucket is not "used" in current batch

		BK := &buckets[op.bucketID]
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			BK.Set(&op.point)
			return
		}
		if BK.X.Equal(&op.point.X) {
			if BK.Y.Equal(&op.point.Y) {
				// P + P: doubling, which should be quite rare --
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				BK.Add(BK, BK)
				return
			}
			BK.setInfinity()
			return
		}

		bucketIds[op.bucketID] = true
		R[cptAdd] = BK
		P[cptAdd] = op.point
		cptAdd++
	}

	add := func(bucketID uint16, PP *G2Affine, isAdd bool) {
		// CanAdd must be called before --> ensures bucket is not "used" in current batch

		BK := &buckets[bucketID]
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			if isAdd {
				BK.Set(PP)
			} else {
				BK.Neg(PP)
			}
			return
		}
		if BK.X.Equal(&PP.X) {
			if BK.Y.Equal(&PP.Y) {
				// P + P: doubling, which should be quite rare --
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				if isAdd {
					BK.Add(BK, BK)
				} else {
					BK.setInfinity()
				}

				return
			}
			if isAdd {
				BK.setInfinity()
			} else {
				BK.Add(BK, BK)
			}
			return
		}

		bucketIds[bucketID] = true
		R[cptAdd] = BK
		if isAdd {
			P[cptAdd].Set(PP)
		} else {
			P[cptAdd].Neg(PP)
		}
		cptAdd++
	}

	var queue TQ
	qID := 0

	processQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if !canAdd(queue[i].bucketID) {
				continue
			}
			addFromQueue(queue[i])
			if isFull() {
				executeAndReset()
			}
			queue[i] = queue[qID-1]
			qID--
		}
	}

	processTopQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if !canAdd(queue[i].bucketID) {
				return
			}
			addFromQueue(queue[i])
			if isFull() {
				executeAndReset()
			}
			qID--
		}
	}

	for i, digit := range digits {

		if digit == 0 || points[i].IsInfinity() {
			continue
		}

		bucketID := uint16((digit >> 1))
		isAdd := digit&1 == 0
		if isAdd {
			// add
			bucketID -= 1
		}

		if !canAdd(bucketID) {
			// put it in queue
			queue[qID].bucketID = bucketID
			if isAdd {
				queue[qID].point = points[i]
			} else {
				queue[qID].point.Neg(&points[i])
			}
			qID++

			// queue is full, flush it.
			if qID == len(queue)-1 {
				executeAndReset()
				processQueue()
			}
			continue
		}

		// we add the point to the batch.
		add(bucketID, &points[i], isAdd)
		if isFull() {
			executeAndReset()
			processTopQueue()
		}
	}

	// empty the queue
	for qID != 0 {
		processQueue()
		executeAndReset()
	}

	// flush items in batch.
	executeAndReset()

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g2JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].IsInfinity() {
			runningSum.addMixed(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total

}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketG2AffineC16 [1 << (16 - 1)]G2Affine

// buckets: array of G2Affine points of size 1 << (c-1)
type ibG2Affine interface {
	bucketG2AffineC16
}

// array of coordinates fp.Element
type cG2Affine interface {
	cG2AffineC16
}

// buckets: array of G2Affine points (for the batch addition)
type pG2Affine interface {
	pG2AffineC16
}

// buckets: array of *G2Affine points (for the batch addition)
type ppG2Affine interface {
	ppG2AffineC16
}

// buckets: array of G2Affine queue operations (for the batch addition)
type qOpsG2Affine interface {
	qOpsG2AffineC16
}
type cG2AffineC16 [640]fp.Element
type pG2AffineC16 [640]G2Affine
type ppG2AffineC16 [640]*G2Affine
type qOpsG2AffineC16 [640]batchOpG2Affine

type bitSetC4 [1 << (4 - 1)]bool
type bitSetC5 [1 << (5 - 1)]bool
type bitSetC8 [1 << (8 - 1)]bool
type bitSetC16 [1 << (16 - 1)]bool

type bitSet interface {
	bitSetC4 |
		bitSetC5 |
		bitSetC8 |
		bitSetC16
}
