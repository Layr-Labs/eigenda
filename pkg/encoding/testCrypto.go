package main

import (
	"log"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	gkzg "github.com/consensys/gnark-crypto/ecc/bn254/fr/kzg"
)

// func testBnBytes() {
// 	//poly := []uint64{ 1, 2, 3, 4, 5}
// 	//polyFr := make([]fr.Element, len(poly))
// 	//for i := 0 ; i < len(poly) ; i++ {
// 	//polyFr[i] = fr.NewElement(poly[i])
// 	//}
//
// 	//for i := 0 ; i < len(polyFr) ; i++ {
// 	//var b fr.Element
// 	//a := polyFr[i].Bytes()
// 	//b.SetBytes(a[:])
// 	//fmt.Println("i", i)
// 	//fmt.Println(polyFr[i])
// 	//fmt.Println(b)
// 	//}
//
// 	text := genText(32 * 4)
// 	for i := 0; i < 4; i++ {
// 		var b fr.Element
// 		t := text[i*31 : (i+1)*31]
// 		b.SetBytes(t)
// 		a := b.Bytes()
//
// 		fmt.Println("i", i)
// 		fmt.Println(t)
// 		fmt.Println(a[1:])
// 	}
// }
//
// func gnarkKzgCommit() {
// 	srs, err := newKZGSRS(16)
// 	if err != nil {
// 		log.Fatal("cannot gen srs")
// 	}
// 	fmt.Printf("SRS\n")
// 	//for i, v := range srs.G1 {
// 	//fmt.Printf("%v: %v \n", i, v.String())
// 	//}
//
// 	poly := []uint64{1, 2, 3, 4, 5}
// 	polyFr := make([]fr.Element, len(poly))
// 	for i := 0; i < len(poly); i++ {
// 		polyFr[i] = fr.NewElement(poly[i])
// 	}
//
// 	//fmt.Println(polyFr)
//
// 	digest, err := gkzg.Commit(polyFr, srs)
//
// 	point := fr.NewElement(uint64(3))
//
// 	op, err := gkzg.Open(polyFr, point, srs)
// 	if err != nil {
// 		log.Fatal()
// 	}
//
// 	fmt.Println("gnark commit")
// 	fmt.Printf("%v\n", digest.String())
//
// 	fmt.Println("open proof, quotient")
// 	fmt.Printf("%v\n", op.H.String())
//
// 	fmt.Println("ClaimedValue")
// 	fmt.Printf("%v\n", op.ClaimedValue.String())
//
// 	fmt.Println("Verify")
// 	err = Verify(&digest, &op, point, srs)
// 	assert.Nil(t, err)
//
// }

// Verify verifies a KZG opening proof at a single point
func Verify(commitment *gkzg.Digest, proof *gkzg.OpeningProof, point fr.Element, srs *gkzg.SRS) error {
	// [f(a)]G₁
	var claimedValueG1Aff bn254.G1Affine
	var claimedValueBigInt big.Int
	proof.ClaimedValue.BigInt(&claimedValueBigInt)
	claimedValueG1Aff.ScalarMultiplication(&srs.Vk.G1, &claimedValueBigInt)

	// [f(α) - f(a)]G₁
	var fminusfaG1Jac, tmpG1Jac bn254.G1Jac
	fminusfaG1Jac.FromAffine(commitment)
	tmpG1Jac.FromAffine(&claimedValueG1Aff)
	fminusfaG1Jac.SubAssign(&tmpG1Jac)

	// [-H(α)]G₁
	var negH bn254.G1Affine
	negH.Neg(&proof.H)

	// [α-a]G₂
	var alphaMinusaG2Jac, genG2Jac, alphaG2Jac bn254.G2Jac
	var pointBigInt big.Int
	point.BigInt(&pointBigInt)
	genG2Jac.FromAffine(&srs.Vk.G2[0])
	alphaG2Jac.FromAffine(&srs.Vk.G2[1])
	alphaMinusaG2Jac.ScalarMultiplication(&genG2Jac, &pointBigInt).
		Neg(&alphaMinusaG2Jac).
		AddAssign(&alphaG2Jac)

	// [α-a]G₂
	var xminusaG2Aff bn254.G2Affine
	xminusaG2Aff.FromJacobian(&alphaMinusaG2Jac)

	// [f(α) - f(a)]G₁
	var fminusfaG1Aff bn254.G1Affine
	fminusfaG1Aff.FromJacobian(&fminusfaG1Jac)

	// e([f(α) - f(a)]G₁, G₂).e([-H(α)]G₁, [α-a]G₂) ==? 1
	check, err := bn254.PairingCheck(
		[]bn254.G1Affine{fminusfaG1Aff, negH},
		[]bn254.G2Affine{srs.Vk.G2[0], xminusaG2Aff},
	)
	if err != nil {
		return err
	}
	if !check {
		log.Fatal("checkfail")
	}

	return nil
}

// func newKZGSRS(kzgSize uint64) (*gkzg.SRS, error) {
// 	alpha, err := new(big.Int).SetString("1927409816240961209460912649125", 10)
// 	if !err {
// 		log.Fatal(err)
// 	}
// 	return gkzg.NewSRS(kzgSize, alpha)
// }
//
// func testBnFr() {
//
// 	a := [4]uint64{0, 1, 2, 3}
// 	aFr := makeBnFr(a[:])
// 	fmt.Println("BnFr")
// 	printBnFr(aFr)
//
// 	diff := make([]bnl.Fr, 4)
//
// 	for i := 0; i < 4; i++ {
// 		a := aFr[i]
// 		var b bnl.Fr
// 		bnl.CopyFr(&a, &b)
//
// 		bnl.SubModFr(&diff[i], &a, &b)
// 	}
// 	fmt.Println("diff")
// 	printBnFr(diff)
// }

// func makeBnFr(input []uint64) []bnl.Fr {
// 	inputFr := make([]bnl.Fr, len(input))
// 	for i := 0; i < len(input); i++ {
// 		bnl.AsFr(&inputFr[i], input[i])
// 	}
// 	return inputFr
// }
//
// func printBnFr(d []bnl.Fr) {
// 	for _, e := range d {
// 		fmt.Printf("%v ", e.String())
// 	}
// 	fmt.Printf("\n")
// }
//
// func genText(M uint64) []byte {
// 	signal := make([]byte, M)
// 	rand.Seed(time.Now().UnixNano())
// 	for i := uint64(0); i < M; i++ {
// 		r := rand.Intn(128)
// 		signal[i] = byte(r)
// 	}
// 	return signal
// }
