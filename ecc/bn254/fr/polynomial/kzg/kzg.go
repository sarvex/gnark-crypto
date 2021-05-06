package kzg

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
	bn254_pol "github.com/consensys/gnark-crypto/ecc/bn254/fr/polynomial"
	fiatshamir "github.com/consensys/gnark-crypto/fiat-shamir"
	"github.com/consensys/gnark-crypto/internal/parallel"
	"github.com/consensys/gnark-crypto/polynomial"
)

// Digest commitment of a polynomial
type Digest = bn254.G1Affine

// Scheme stores KZG data
type Scheme struct {

	// Size of the Srs.G1. Polynomial of degree up to Size-1 can be committed.
	// For Plonk use case, Size should be a power of 2.
	Size uint64

	// Domain to perform polynomial division. The size of the domain is the lowest power of 2 greater than Size.
	Domain *fft.Domain

	// Srs stores the result of the MPC
	Srs struct {
		G1 []bn254.G1Affine  // [gen [alpha]gen , [alpha**2]gen, ... ]
		G2 [2]bn254.G2Affine // [gen, [alpha]gen ]
	}
}

// Proof KZG proof for opening at a single point.
type Proof struct {

	// Point at which the polynomial is evaluated
	Point fr.Element

	// ClaimedValue purpoted value
	ClaimedValue fr.Element

	// H quotient polynomial (f - f(z))/(x-z)
	H bn254.G1Affine
}

func (mp *Proof) Marshal() []byte {
	panic("not implemented")
}

type BatchProofsSinglePoint struct {

	// Point at which the polynomials are evaluated
	Point fr.Element

	// ClaimedValues purpoted values
	ClaimedValues []fr.Element

	// H quotient polynomial Sum_i gamma**i*(f - f(z))/(x-z)
	H bn254.G1Affine
}

func (mp *BatchProofsSinglePoint) Marshal() []byte {
	panic("not implemented")
}

// Commit commits to a polynomial using a multi exponentiation with the SRS.
// It is assumed that the polynomial is in canonical form, in Montgomery form.
func (s *Scheme) Commit(p polynomial.Polynomial) polynomial.Digest {

	// TODO raise an error instead?
	if p.Degree() >= s.Size {
		panic("[Commit] Size of polynomial exceeds the size supported by the scheme")
	}

	var res Digest
	_p := p.(*bn254_pol.Polynomial)

	// ensure we don't modify p
	pCopy := make(bn254_pol.Polynomial, s.Domain.Cardinality)
	copy(pCopy, *_p)

	parallel.Execute(len(*_p), func(start, end int) {
		for i := start; i < end; i++ {
			pCopy[i].FromMont()
		}
	})
	res.MultiExp(s.Srs.G1, pCopy)

	return &res
}

// Open computes an opening proof of _p at _val.
// Returns a MockProof, which is an empty interface.
func (s *Scheme) Open(point interface{}, p polynomial.Polynomial) polynomial.OpeningProof {

	if p.Degree() >= s.Domain.Cardinality {
		panic("[Open] Size of polynomial exceeds the size supported by the scheme")
	}

	// build the proof
	var res Proof
	claimedValue := p.Eval(point)
	res.Point = bn254_pol.FromInterface(point)
	res.ClaimedValue = bn254_pol.FromInterface(claimedValue)

	// compute H
	_p := p.(*bn254_pol.Polynomial)
	h := dividePolyByXminusA(*s.Domain, *_p, res.ClaimedValue, res.Point)

	// commit to H
	c := s.Commit(&h)
	res.H.Set(c.(*bn254.G1Affine))

	return &res
}

// Verify verifies a KZG opening proof at a single point
func (s *Scheme) Verify(commitment polynomial.Digest, proof polynomial.OpeningProof) error {

	_commitment := commitment.(*bn254.G1Affine)
	_proof := proof.(*Proof)

	// comm(f(a))
	var claimedValueG1Aff bn254.G1Affine
	var claimedValueBigInt big.Int
	_proof.ClaimedValue.ToBigIntRegular(&claimedValueBigInt)
	claimedValueG1Aff.ScalarMultiplication(&s.Srs.G1[0], &claimedValueBigInt)

	// [f(alpha) - f(a)]G1Jac
	var fminusfaG1Jac, tmpG1Jac bn254.G1Jac
	fminusfaG1Jac.FromAffine(_commitment)
	tmpG1Jac.FromAffine(&claimedValueG1Aff)
	fminusfaG1Jac.SubAssign(&tmpG1Jac)

	// [-H(alpha)]G1Aff
	var negH bn254.G1Affine
	negH.Neg(&_proof.H)

	// [alpha-a]G2Jac
	var alphaMinusaG2Jac, genG2Jac, alphaG2Jac bn254.G2Jac
	var pointBigInt big.Int
	_proof.Point.ToBigIntRegular(&pointBigInt)
	genG2Jac.FromAffine(&s.Srs.G2[0])
	alphaG2Jac.FromAffine(&s.Srs.G2[1])
	alphaMinusaG2Jac.ScalarMultiplication(&genG2Jac, &pointBigInt).
		Neg(&alphaMinusaG2Jac).
		AddAssign(&alphaG2Jac)

	// [alpha-a]G2Aff
	var xminusaG2Aff bn254.G2Affine
	xminusaG2Aff.FromJacobian(&alphaMinusaG2Jac)

	// [f(alpha) - f(a)]G1Aff
	var fminusfaG1Aff bn254.G1Affine
	fminusfaG1Aff.FromJacobian(&fminusfaG1Jac)

	// e([-H(alpha)]G1Aff, G2gen).e([-H(alpha)]G1Aff, [alpha-a]G2Aff) ==? 1
	check, err := bn254.PairingCheck(
		[]bn254.G1Affine{fminusfaG1Aff, negH},
		[]bn254.G2Affine{s.Srs.G2[0], xminusaG2Aff},
	)
	if err != nil {
		return err
	}
	if !check {
		return polynomial.ErrVerifyOpeningProof
	}
	return nil
}

// BatchOpenSinglePoint creates a batch opening proof of several polynomials at a single point
func (s *Scheme) BatchOpenSinglePoint(point interface{}, digests []polynomial.Digest, polynomials []polynomial.Polynomial) polynomial.BatchOpeningProofSinglePoint {

	nbDigests := len(digests)
	if nbDigests != len(polynomials) {
		panic("The number of polynomials and digests don't match")
	}

	var res BatchProofsSinglePoint

	// compute the purported values
	res.ClaimedValues = make([]fr.Element, len(polynomials))
	for i := 0; i < len(polynomials); i++ {
		_val := polynomials[i].Eval(point)
		res.ClaimedValues[i] = bn254_pol.FromInterface(_val)
	}

	// set the point at which the evaluation is done
	res.Point = bn254_pol.FromInterface(point)

	// derive the challenge gamma, binded to the point and the commitments
	fs := fiatshamir.NewTranscript(fiatshamir.SHA256, "gamma")
	fs.Bind("gamma", res.Point.Marshal())
	for i := 0; i < len(digests); i++ {
		fs.Bind("gamma", digests[i].Marshal())
	}
	gammaByte, _ := fs.ComputeChallenge("gamma") // TODO handle error
	var gamma fr.Element
	gamma.SetBytes(gammaByte)

	// compute sum_i gamma**i*f and sum_i gamma**i*f(a)
	var sumGammaiTimesEval fr.Element
	sumGammaiTimesEval.Set(&res.ClaimedValues[nbDigests-1])
	sumGammaiTimesPol := polynomials[nbDigests-1].Clone()
	for i := nbDigests - 2; i >= 0; i-- {
		sumGammaiTimesEval.Mul(&sumGammaiTimesEval, &gamma).
			Add(&sumGammaiTimesEval, &res.ClaimedValues[i])
		sumGammaiTimesPol.ScaleInPlace(&gamma)
		sumGammaiTimesPol.Add(polynomials[i], sumGammaiTimesPol)
	}

	// compute H
	_sumGammaiTimesPol := sumGammaiTimesPol.(*bn254_pol.Polynomial)
	h := dividePolyByXminusA(*s.Domain, *_sumGammaiTimesPol, sumGammaiTimesEval, res.Point)
	c := s.Commit(&h)
	res.H.Set(c.(*bn254.G1Affine))

	return &res
}

func (s *Scheme) BatchVerifySinglePoint(digests []polynomial.Digest, batchOpeningProof polynomial.BatchOpeningProofSinglePoint) error {

	nbDigests := len(digests)

	_batchOpeningProof := batchOpeningProof.(*BatchProofsSinglePoint)

	if len(digests) != len(_batchOpeningProof.ClaimedValues) {
		panic("size of digests and size of proof don't match")
	}

	// derive the challenge gamma, binded to the point and the commitments
	fs := fiatshamir.NewTranscript(fiatshamir.SHA256, "gamma")
	fs.Bind("gamma", _batchOpeningProof.Point.Marshal())
	for i := 0; i < len(digests); i++ {
		fs.Bind("gamma", digests[i].Marshal())
	}
	gammaByte, _ := fs.ComputeChallenge("gamma") // TODO handle error
	var gamma fr.Element
	gamma.SetBytes(gammaByte)

	var sumGammaiTimesEval fr.Element
	sumGammaiTimesEval.Set(&_batchOpeningProof.ClaimedValues[nbDigests-1])
	for i := nbDigests - 2; i >= 0; i-- {
		sumGammaiTimesEval.Mul(&sumGammaiTimesEval, &gamma).
			Add(&sumGammaiTimesEval, &_batchOpeningProof.ClaimedValues[i])
	}

	var sumGammaiTimesEvalBigInt big.Int
	sumGammaiTimesEval.ToBigIntRegular(&sumGammaiTimesEvalBigInt)
	var sumGammaiTimesEvalG1Aff bn254.G1Affine
	sumGammaiTimesEvalG1Aff.ScalarMultiplication(&s.Srs.G1[0], &sumGammaiTimesEvalBigInt)

	var acc fr.Element
	acc.SetOne()
	gammai := make([]fr.Element, len(digests))
	gammai[0].SetOne().FromMont()
	for i := 1; i < len(digests); i++ {
		acc.Mul(&acc, &gamma)
		gammai[i].Set(&acc).FromMont()
	}
	var sumGammaiTimesDigestsG1Aff bn254.G1Affine
	_digests := make([]bn254.G1Affine, len(digests))
	for i := 0; i < len(digests); i++ {
		_digests[i].Set(digests[i].(*bn254.G1Affine))
	}

	sumGammaiTimesDigestsG1Aff.MultiExp(_digests, gammai)

	// sum_i [gamma**i * (f-f(a))]G1
	var sumGammiDiffG1Aff bn254.G1Affine
	var t1, t2 bn254.G1Jac
	t1.FromAffine(&sumGammaiTimesDigestsG1Aff)
	t2.FromAffine(&sumGammaiTimesEvalG1Aff)
	t1.SubAssign(&t2)
	sumGammiDiffG1Aff.FromJacobian(&t1)

	// [alpha-a]G2Jac
	var alphaMinusaG2Jac, genG2Jac, alphaG2Jac bn254.G2Jac
	var pointBigInt big.Int
	_batchOpeningProof.Point.ToBigIntRegular(&pointBigInt)
	genG2Jac.FromAffine(&s.Srs.G2[0])
	alphaG2Jac.FromAffine(&s.Srs.G2[1])
	alphaMinusaG2Jac.ScalarMultiplication(&genG2Jac, &pointBigInt).
		Neg(&alphaMinusaG2Jac).
		AddAssign(&alphaG2Jac)

	// [alpha-a]G2Aff
	var xminusaG2Aff bn254.G2Affine
	xminusaG2Aff.FromJacobian(&alphaMinusaG2Jac)

	// [-H(alpha)]G1Aff
	var negH bn254.G1Affine
	negH.Neg(&_batchOpeningProof.H)

	// check the pairing equation
	check, err := bn254.PairingCheck(
		[]bn254.G1Affine{sumGammiDiffG1Aff, negH},
		[]bn254.G2Affine{s.Srs.G2[0], xminusaG2Aff},
	)
	if err != nil {
		return err
	}
	if !check {
		return polynomial.ErrVerifyBatchOpeningSinglePoint
	}
	return nil
}
