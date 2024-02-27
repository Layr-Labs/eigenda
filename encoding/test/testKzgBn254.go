package main

// "fmt"
// "log"
// "math"
// bn "github.com/Layr-Labs/eigenda/encoding/kzg/bn254"
// bnKzg "github.com/Layr-Labs/eigenda/encoding/kzgFFT/go-kzg-bn254"

//func testLinComb() {
//polyFr := GenPoly(1,2,3,4,5)
////fmt.Println("polyFr")
////fmt.Println(polyFr)
//s1, _ := GenerateTestingSetup("1927409816240961209460912649125", 16+1)

//var dr bn.G1Point
//for i := 0 ; i < 5 ; i++ {
//var product bn.G1Point
//var b big.Int
//polyFr[i].BigInt(b)
////product.ScalarMultiplication(&s1[i], &b)
//dr.AddG1()
//}

//}

//func SingleProof() {
//fftSetting := bnKzg.NewFFTSettings(4)
//s1, s2 := GenerateTestingSetup("1927409816240961209460912649125", 16+1)

//fmt.Println("srs single proof")
////for i, v := range s1 {
////fmt.Printf("%v: %v \n", i, v.String())
////}
////input := []uint64{1,2,3,4,5}
//polyFr := GenPoly(1,2,3,4,5)
////fmt.Println("polyFr")
////fmt.Println(polyFr)

//kzgSetting := bnKzg.NewKZGSettings(fftSetting, s1, s2)
////commit := kzgSetting.CommitToPolyUnoptimized(polyFr)
//commit := kzgSetting.CommitToPoly(polyFr)

//fmt.Println("bn254 commit")
//fmt.Printf("%v\n", commit)

////text := []byte("hello   hello   ")
////polyFr := GenTextPoly(text)

//proof := kzgSetting.ComputeProofSingle(polyFr, 3)

//fmt.Printf("proof\n")
//fmt.Printf("%v\n", proof.String())

//var x bn.Fr
//bn.AsFr(&x, 3)
//var value bn.Fr
//bn.EvalPolyAt(&value, polyFr, &x)
//fmt.Println("value")
//fmt.Println(value.String())

//y := value

//var xG2 bn.G2Point
//bn.MulG2(&xG2, &bn.GenG2, &x)
//var sMinuxX bn.G2Point
//bn.SubG2(&sMinuxX, &s2[1], &xG2)
//var yG1 bn.G1Point
//bn.MulG1(&yG1, &bn.GenG1, &y)
//var commitmentMinusY bn.G1Point
//bn.SubG1(&commitmentMinusY, commit, &yG1)

//// at verifier, it does not need the orig polyFr
//payProof, _ := proof.MarshalText()
//payCommit, _ := commit.MarshalText()
//payValue := value.String()

//payload := len(payProof) +len(payCommit) + len(payValue)
//if !kzgSetting.CheckProofSingle(commit, proof, &x, &value) {
//fmt.Println("not eval to result")
//} else {
//fmt.Println("eval to result", payload)
//fmt.Println("PayValye", payValue)
//}
//}

//func testKzgBnFr() {
//a := [4]uint64{0, 1, 2, 3}
//aFr := makeBnFr(a[:])
//fmt.Println("BnFr")
//printBnFr(aFr)

//diff := make([]bn.Fr, 4)

//for i := 0 ; i < 4 ; i++ {
//a := aFr[i]
//var b bn.Fr
//bn.CopyFr(&a, &b)

//bn.SubModFr(&diff[i], &a, &b)
//}
//fmt.Println("diff")
//printBnFr(diff)
//}

// func makeBn254Fr(input []uint64) []bn.Fr {
// 	inputFr := make([]bn.Fr, len(input))
// 	for i := 0; i < len(input); i++ {
// 		bn.AsFr(&inputFr[i], input[i])
// 	}
// 	return inputFr
// }

//func compareDataBn254(inputFr, dataFr []bn.Fr) {
//if len(inputFr) != len(dataFr) {
//log.Fatalf("Error. Diff length. input %v, data %v\n", len(inputFr), len(dataFr))
//}

//for i := 0 ; i < len(inputFr) ; i++ {
//if !bn.EqualFr(&inputFr[i], &dataFr[i]) {
//log.Fatalf("Error. Diff value at %v. input %v, data %v\n",
//i, inputFr[i].String(), dataFr[i].String())
//}
//}
//}

////func initBnData(size int) ([]bn.Fr, []bn.Fr) {
////v := make([]uint64, size)
////for i := 0 ; i < size ; i++ {
////v[i] = uint64(i + 1)
////}
////dataFr := makeBnFr(v)
////order := multiprover.CeilIntPowerOf2Num(size)
////fs := bnKzg.NewFFTSettings(uint8(order))
////polyFr, err := fs.FFT(dataFr, true)
////if err != nil {
////log.Fatal(err)
////}
////return polyFr, dataFr
////}
