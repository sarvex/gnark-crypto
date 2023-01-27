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

package ecdsa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"hash"
	"io"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bls12-377"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fp"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	utils "github.com/consensys/gnark-crypto/ecc/bls12-377/signature"
	"github.com/consensys/gnark-crypto/signature"
)

const (
	sizeFr         = fr.Bytes
	sizeFp         = fp.Bytes
	sizePublicKey  = sizeFp
	sizePrivateKey = sizeFr + sizePublicKey
	sizeSignature  = 2 * sizeFr
)

var order = fr.Modulus()

// PublicKey represents an ECDSA public key
type PublicKey struct {
	A bls12377.G1Affine
}

// PrivateKey represents an ECDSA private key
type PrivateKey struct {
	PublicKey PublicKey
	scalar    [sizeFr]byte // secret scalar, in big Endian
}

// Signature represents an ECDSA signature
type Signature struct {
	R, S [sizeFr]byte
}

var one = new(big.Int).SetInt64(1)

// randFieldElement returns a random element of the order of the given
// curve using the procedure given in FIPS 186-4, Appendix B.5.1.
func randFieldElement(rand io.Reader) (k *big.Int, err error) {
	b := make([]byte, fr.Bits/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return
	}

	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(order, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}

// GenerateKey generates a public and private key pair.
func GenerateKey(rand io.Reader) (*PrivateKey, error) {

	k, err := randFieldElement(rand)
	if err != nil {
		return nil, err

	}
	_, _, g, _ := bls12377.Generators()

	privateKey := new(PrivateKey)
	k.FillBytes(privateKey.scalar[:sizeFr])
	privateKey.PublicKey.A.ScalarMultiplication(&g, k)
	return privateKey, nil
}

type zr struct{}

// Read replaces the contents of dst with zeros. It is safe for concurrent use.
func (zr) Read(dst []byte) (n int, err error) {
	for i := range dst {
		dst[i] = 0
	}
	return len(dst), nil
}

var zeroReader = zr{}

const (
	aesIV = "gnark-crypto IV." // must be 16 chars (equal block size)
)

func nonce(privateKey *PrivateKey, hash []byte) (csprng *cipher.StreamReader, err error) {
	// This implementation derives the nonce from an AES-CTR CSPRNG keyed by:
	//
	//    SHA2-512(privateKey.scalar ∥ entropy ∥ hash)[:32]
	//
	// The CSPRNG key is indifferentiable from a random oracle as shown in
	// [Coron], the AES-CTR stream is indifferentiable from a random oracle
	// under standard cryptographic assumptions (see [Larsson] for examples).
	//
	// [Coron]: https://cs.nyu.edu/~dodis/ps/merkle.pdf
	// [Larsson]: https://web.archive.org/web/20040719170906/https://www.nada.kth.se/kurser/kth/2D1441/semteo03/lecturenotes/assump.pdf

	// Get 256 bits of entropy from rand.
	entropy := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, entropy)
	if err != nil {
		return

	}

	// Initialize an SHA-512 hash context; digest...
	md := sha512.New()
	md.Write(privateKey.scalar[:sizeFr]) // the private key,
	md.Write(entropy)                    // the entropy,
	md.Write(hash)                       // and the input hash;
	key := md.Sum(nil)[:32]              // and compute ChopMD-256(SHA-512),
	// which is an indifferentiable MAC.

	// Create an AES-CTR instance to use as a CSPRNG.
	block, _ := aes.NewCipher(key)

	// Create a CSPRNG that xors a stream of zeros with
	// the output of the AES-CTR instance.
	csprng = &cipher.StreamReader{
		R: zeroReader,
		S: cipher.NewCTR(block, []byte(aesIV)),
	}

	return csprng, err
}

// Equal compares 2 public keys
func (pub *PublicKey) Equal(x signature.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	bpk := pub.Bytes()
	bxx := xx.Bytes()
	return subtle.ConstantTimeCompare(bpk, bxx) == 1
}

// Public returns the public key associated to the private key.
func (privKey *PrivateKey) Public() signature.PublicKey {
	var pub PublicKey
	pub.A.Set(&privKey.PublicKey.A)
	return &pub
}

// Sign performs the ECDSA signature
//
// k ← 𝔽r (random)
// P = k ⋅ g1Gen
// r = x_P (mod order)
// s = k⁻¹ . (m + sk ⋅ r)
// signature = {r, s}
//
// SEC 1, Version 2.0, Section 4.1.3
func (privKey *PrivateKey) Sign(message []byte, hFunc hash.Hash) ([]byte, error) {
	scalar, r, s, kInv := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	scalar.SetBytes(privKey.scalar[:sizeFr])
	for {
		for {
			csprng, err := nonce(privKey, message)
			if err != nil {
				return nil, err
			}
			k, err := randFieldElement(csprng)
			if err != nil {
				return nil, err
			}

			var P bls12377.G1Affine
			P.ScalarMultiplicationBase(k)
			kInv.ModInverse(k, order)

			P.X.BigInt(r)
			r.Mod(r, order)
			if r.Sign() != 0 {
				break
			}
		}
		s.Mul(r, scalar)

		var m *big.Int
		if hFunc != nil {
			// compute the hash of the message as an integer
			dataToHash := make([]byte, len(message))
			copy(dataToHash[:], message[:])
			hFunc.Reset()
			_, err := hFunc.Write(dataToHash[:])
			if err != nil {
				return nil, err
			}
			hramBin := hFunc.Sum(nil)
			m = utils.HashToInt(hramBin)
		} else {
			m = utils.HashToInt(message)
		}

		s.Add(m, s).
			Mul(kInv, s).
			Mod(s, order) // order != 0
		if s.Sign() != 0 {
			break
		}
	}

	var sig Signature
	r.FillBytes(sig.R[:sizeFr])
	s.FillBytes(sig.S[:sizeFr])

	return sig.Bytes(), nil
}

// Verify validates the ECDSA signature
//
// R ?= (s⁻¹ ⋅ m ⋅ Base + s⁻¹ ⋅ R ⋅ publiKey)_x
//
// SEC 1, Version 2.0, Section 4.1.4
func (publicKey *PublicKey) Verify(sigBin, message []byte, hFunc hash.Hash) (bool, error) {

	// Deserialize the signature
	var sig Signature
	if _, err := sig.SetBytes(sigBin); err != nil {
		return false, err
	}

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sig.R[:sizeFr])
	s.SetBytes(sig.S[:sizeFr])

	sInv := new(big.Int).ModInverse(s, order)

	var m *big.Int
	if hFunc != nil {
		// compute the hash of the message as an integer
		dataToHash := make([]byte, len(message))
		copy(dataToHash[:], message[:])
		hFunc.Reset()
		_, err := hFunc.Write(dataToHash[:])
		if err != nil {
			return false, err
		}
		hramBin := hFunc.Sum(nil)
		m = utils.HashToInt(hramBin)
	} else {
		m = utils.HashToInt(message)
	}

	u1 := new(big.Int).Mul(m, sInv)
	u1.Mod(u1, order)
	u2 := new(big.Int).Mul(r, sInv)
	u2.Mod(u2, order)
	var U bls12377.G1Jac
	U.JointScalarMultiplicationBase(&publicKey.A, u1, u2)

	var z big.Int
	U.Z.Square(&U.Z).
		Inverse(&U.Z).
		Mul(&U.Z, &U.X).
		BigInt(&z)

	z.Mod(&z, order)

	return z.Cmp(r) == 0, nil

}
