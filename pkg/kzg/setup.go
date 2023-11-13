//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"

	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// GenerateTestingSetup creates a setup of n values from the given secret. **for testing purposes only**
func GenerateTestingSetup(secret string, n uint64) ([]bls.G1Point, []bls.G2Point) {
	var s bls.Fr
	bls.SetFr(&s, secret)

	var sPow bls.Fr
	bls.CopyFr(&sPow, &bls.ONE)

	s1Out := make([]bls.G1Point, n)
	s2Out := make([]bls.G2Point, n)
	for i := uint64(0); i < n; i++ {
		bls.MulG1(&s1Out[i], &bls.GenG1, &sPow)
		bls.MulG2(&s2Out[i], &bls.GenG2, &sPow)
		var tmp bls.Fr
		bls.CopyFr(&tmp, &sPow)
		bls.MulModFr(&sPow, &tmp, &s)
	}
	return s1Out, s2Out
}

// func ReadGeneatorPoints(n uint64, g1FilePath, g2FilePath string) ([]bls.G1Point, []bls.G2Point, error) {
// 	g1f, err := os.Open(g1FilePath)
// 	if err != nil {
// 		log.Println("ReadGeneatorPoints.ERR.0", err)
// 		return nil, nil, err
// 	}
// 	g2f, err := os.Open(g2FilePath)
// 	if err != nil {
// 		log.Println("ReadGeneatorPoints.ERR.1", err)
// 		return nil, nil, err
// 	}
// 	//todo: handle panic
// 	defer func() {
// 		if err := g1f.Close(); err != nil {
// 			panic(err)
// 		}
// 		if err := g2f.Close(); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	start := time.Now()
// 	g1r := bufio.NewReaderSize(g1f, int(n*64))
// 	g2r := bufio.NewReaderSize(g2f, int(n*128))

// 	g1Bytes, _, err := g1r.ReadLine()
// 	if err != nil {
// 		log.Println("ReadGeneatorPoints.ERR.2", err)
// 		return nil, nil, err
// 	}

// 	g2Bytes, _, err := g2r.ReadLine()
// 	if err != nil {
// 		log.Println("ReadGeneatorPoints.ERR.3", err)
// 		return nil, nil, err
// 	}

// 	if uint64(len(g1Bytes)) < 64*n {
// 		log.Printf("Error. Insufficient G1 points. Only contains %v. Requesting %v\n", len(g1Bytes)/64, n)
// 		log.Println()
// 		log.Println("ReadGeneatorPoints.ERR.4", err)
// 		return nil, nil, err
// 	}
// 	if uint64(len(g2Bytes)) < 128*n {
// 		log.Printf("Error. Insufficient G2 points. Only contains %v. Requesting %v\n", len(g1Bytes)/128, n)
// 		log.Println()
// 		log.Println("ReadGeneatorPoints.ERR.5", err)
// 		return nil, nil, err
// 	}

// 	// measure reading time
// 	t := time.Now()
// 	elapsed := t.Sub(start)
// 	fmt.Println("    Reading G1 G2 raw points takes", elapsed)
// 	start = time.Now()

// 	s1Outs := make([]bls.G1Point, n, n)
// 	s2Outs := make([]bls.G2Point, n, n)
// 	for i := uint64(0); i < n; i++ {
// 		g1 := g1Bytes[i*64 : (i+1)*128]
// 		err := s1Outs[i].UnmarshalText(g1[:])
// 		if err != nil {
// 			log.Println("ReadGeneatorPoints.ERR.6", err)
// 			return nil, nil, err
// 		}

// 		g2 := g2Bytes[i*128 : (i+1)*128]
// 		err = s2Outs[i].UnmarshalText(g2[:])
// 		if err != nil {
// 			log.Println("ReadGeneatorPoints.ERR.7", err)
// 			return nil, nil, err
// 		}
// 	}

// 	// measure parsing time
// 	t = time.Now()
// 	elapsed = t.Sub(start)
// 	fmt.Println("    Parsing G1 G2 to crypto data struct takes", elapsed)
// 	return s1Outs, s2Outs, nil
// }

func WriteGeneratorPoints(n uint64) error {
	secret := "1927409816240961209460912649125"
	ns := strconv.Itoa(int(n))

	var s bls.Fr
	bls.SetFr(&s, secret)

	var sPow bls.Fr
	bls.CopyFr(&sPow, &bls.ONE)

	g1f, err := os.Create("g1.point." + ns)
	if err != nil {
		log.Println("WriteGeneratorPoints.ERR.0", err)
		return err
	}

	g1w := bufio.NewWriter(g1f)
	g2f, err := os.Create("g2.point." + ns)
	if err != nil {
		log.Println("WriteGeneratorPoints.ERR.1", err)
		return err
	}
	g2w := bufio.NewWriter(g2f)

	//delimiter := [1]byte{'\n'}

	start := time.Now()
	for i := uint64(0); i < n; i++ {
		var s1Out bls.G1Point
		var s2Out bls.G2Point
		bls.MulG1(&s1Out, &bls.GenG1, &sPow)
		bls.MulG2(&s2Out, &bls.GenG2, &sPow)

		g1Byte := s1Out.MarshalText()
		if _, err := g1w.Write(g1Byte); err != nil {
			log.Println("WriteGeneratorPoints.ERR.3", err)
			return err
		}

		g2Byte := s2Out.MarshalText()
		if _, err := g2w.Write(g2Byte); err != nil {
			log.Println("WriteGeneratorPoints.ERR.5", err)
			return err
		}

		var tmp bls.Fr
		bls.CopyFr(&tmp, &sPow)
		bls.MulModFr(&sPow, &tmp, &s)
	}

	if err = g1w.Flush(); err != nil {
		log.Println("WriteGeneratorPoints.ERR.6", err)
		return err
	}
	if err = g2w.Flush(); err != nil {
		log.Println("WriteGeneratorPoints.ERR.7", err)
		return err
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Generating takes", elapsed)
	return nil
}
