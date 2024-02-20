package kzgEncoder

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/utils"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type SubTable struct {
	//SizeLow   uint64
	//SizeUp    uint64
	FilePath string
}

type TableParam struct {
	DimE      uint64
	CosetSize uint64
}

type SRSTable struct {
	Tables    map[TableParam]SubTable
	TableDir  string
	NumWorker uint64
	s1        []bls.G1Point
}

func NewSRSTable(tableDir string, s1 []bls.G1Point, numWorker uint64) (*SRSTable, error) {

	err := os.MkdirAll(tableDir, os.ModePerm)
	if err != nil {
		log.Println("NEWSRSTABLE.ERR.1", err)
		return nil, err
	}

	files, err := os.ReadDir(tableDir)
	if err != nil {
		log.Println("NEWSRSTABLE.ERR.2", err)
		return nil, err
	}

	tables := make(map[TableParam]SubTable)
	for _, file := range files {
		filename := file.Name()

		tokens := strings.Split(filename, ".")

		dimEValue, err := strconv.Atoi(tokens[0][4:])
		if err != nil {
			log.Println("NEWSRSTABLE.ERR.3", err)
			return nil, err
		}
		cosetSizeValue, err := strconv.Atoi(tokens[1][5:])
		if err != nil {
			log.Println("NEWSRSTABLE.ERR.4", err)
			return nil, err
		}

		param := TableParam{
			DimE:      uint64(dimEValue),
			CosetSize: uint64(cosetSizeValue),
		}

		filePath := path.Join(tableDir, filename)
		tables[param] = SubTable{FilePath: filePath}
	}

	return &SRSTable{
		Tables:    tables,
		TableDir:  tableDir,
		NumWorker: numWorker,
		s1:        s1, // g1 points
	}, nil
}

func (p *SRSTable) GetSubTables(
	numChunks uint64,
	chunkLen uint64,
) ([][]bls.G1Point, error) {
	cosetSize := chunkLen
	dimE := numChunks
	m := numChunks*chunkLen - 1
	dim := m / cosetSize

	param := TableParam{
		DimE:      dimE,
		CosetSize: cosetSize,
	}

	start := time.Now()
	table, ok := p.Tables[param]
	if !ok {
		log.Printf("Table with params: DimE=%v CosetSize=%v does not exist\n", dimE, cosetSize)
		log.Printf("Generating the table. May take a while\n")
		log.Printf("... ...\n")
		filename := fmt.Sprintf("dimE%v.coset%v", dimE, cosetSize)
		dstFilePath := path.Join(p.TableDir, filename)
		fftPoints := p.Precompute(dim, dimE, cosetSize, m, dstFilePath, p.NumWorker)

		elapsed := time.Since(start)
		log.Printf("    Precompute finishes using %v\n", elapsed)

		return fftPoints, nil
	} else {
		log.Printf("Detected Precomputed FFT sliced G1 table\n")
		fftPoints, err := p.TableReaderThreads(table.FilePath, dimE, cosetSize, p.NumWorker)
		if err != nil {
			log.Println("GetSubTables.ERR.0", err)
			return nil, err
		}

		elapsed := time.Since(start)
		log.Printf("    Loading Table uses %v\n", elapsed)

		return fftPoints, nil
	}
}

type DispatchReturn struct {
	points []bls.G1Point
	j      uint64
}

// m = len(poly) - 1, which is deg
func (p *SRSTable) Precompute(dim, dimE, l, m uint64, filePath string, numWorker uint64) [][]bls.G1Point {
	order := dimE * l
	if l == 1 {
		order = dimE * 2
	}
	// TODO, create function only read g1 points
	//s1 := ReadG1Points(p.SrsFilePath, order)
	n := uint8(math.Log2(float64(order)))
	fs := kzg.NewFFTSettings(n)

	fftPoints := make([][]bls.G1Point, l)

	numJob := l
	jobChan := make(chan uint64, numJob)
	results := make(chan DispatchReturn, l)

	for w := uint64(0); w < numWorker; w++ {
		go p.precomputeWorker(fs, m, dim, dimE, jobChan, l, results)
	}

	for j := uint64(0); j < l; j++ {
		jobChan <- j
	}
	close(jobChan)

	for w := uint64(0); w < l; w++ {
		computeResult := <-results
		fftPoints[computeResult.j] = computeResult.points
	}

	err := p.TableWriter(fftPoints, dimE, filePath)
	if err != nil {
		log.Println("Precompute error:", err)
	}
	return fftPoints
}

func (p *SRSTable) precomputeWorker(fs *kzg.FFTSettings, m, dim, dimE uint64, jobChan <-chan uint64, l uint64, results chan DispatchReturn) {
	for j := range jobChan {
		dr, err := p.PrecomputeSubTable(fs, m, dim, dimE, j, l)
		if err != nil {
			log.Println("precomputeWorker.ERR.1", err)
			return
		}
		results <- dr
	}
}

func (p *SRSTable) PrecomputeSubTable(fs *kzg.FFTSettings, m, dim, dimE, j, l uint64) (DispatchReturn, error) {
	// there is a constant term
	points := make([]bls.G1Point, 2*dimE)
	k := m - l - j

	for i := uint64(0); i < dim; i++ {
		bls.CopyG1(&points[i], &p.s1[k])
		k -= l
	}
	for i := dim; i < 2*dimE; i++ {
		bls.CopyG1(&points[i], &bls.ZeroG1)
	}

	y, err := fs.FFTG1(points, false)
	if err != nil {
		log.Println("PrecomputeSubTable.ERR.1", err)
		return DispatchReturn{}, err
	}

	return DispatchReturn{
		points: y,
		j:      j,
	}, nil

}

type Boundary struct {
	start   uint64
	end     uint64 // informational
	sliceAt uint64
}

func (p *SRSTable) TableReaderThreads(filePath string, dimE, l uint64, numWorker uint64) ([][]bls.G1Point, error) {
	g1f, err := os.Open(filePath)
	if err != nil {
		log.Println("TableReaderThreads.ERR.0", err)
		return nil, err
	}
	//todo: resolve panic
	defer func() {
		if err := g1f.Close(); err != nil {
			panic(err)
		}
	}()

	// 2 due to circular FFT  mul
	subTableSize := dimE * 2 * utils.G1PointBytes
	totalSubTableSize := subTableSize * l

	if numWorker > l {
		numWorker = l
	}

	reader := bufio.NewReaderSize(g1f, int(totalSubTableSize+l))
	buf := make([]byte, totalSubTableSize+l)
	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("TableReaderThreads.ERR.1", err, "file path:", filePath)
		return nil, err
	}

	boundaries := make([]Boundary, l)
	for i := uint64(0); i < uint64(l); i++ {
		start := (subTableSize + 1) * i
		end := (subTableSize+1)*(i+1) - 1 // exclude \n
		boundary := Boundary{
			start:   start,
			end:     end,
			sliceAt: i,
		}
		boundaries[i] = boundary
	}

	fftPoints := make([][]bls.G1Point, l)

	jobChan := make(chan Boundary, l)

	var wg sync.WaitGroup
	wg.Add(int(numWorker))
	for i := uint64(0); i < numWorker; i++ {
		go p.readWorker(buf, fftPoints, jobChan, dimE, &wg)
	}

	for i := uint64(0); i < l; i++ {
		jobChan <- boundaries[i]
	}
	close(jobChan)
	wg.Wait()
	return fftPoints, nil
}

func (p *SRSTable) readWorker(
	buf []byte,
	fftPoints [][]bls.G1Point,
	jobChan <-chan Boundary,
	dimE uint64,
	wg *sync.WaitGroup,
) {
	for b := range jobChan {
		slicePoints := make([]bls.G1Point, dimE*2)
		for i := uint64(0); i < dimE*2; i++ {
			g1 := buf[b.start+i*utils.G1PointBytes : b.start+(i+1)*utils.G1PointBytes]
			err := slicePoints[i].UnmarshalText(g1[:])
			if err != nil {
				log.Printf("Error. From %v to %v. %v", b.start, b.end, err)
				log.Println()
				log.Println("readWorker.ERR.0", err)
				return
			}
		}
		fftPoints[b.sliceAt] = slicePoints
	}
	wg.Done()
}

func (p *SRSTable) TableWriter(fftPoints [][]bls.G1Point, dimE uint64, filePath string) error {
	wf, err := os.Create(filePath)
	if err != nil {
		log.Println("TableWriter.ERR.0", err)
		return err
	}

	writer := bufio.NewWriter(wf)
	l := uint64(len(fftPoints))

	delimiter := [1]byte{'\n'}

	for j := uint64(0); j < l; j++ {
		for i := uint64(0); i < dimE*2; i++ {

			g1Bytes := fftPoints[j][i].MarshalText()
			if _, err := writer.Write(g1Bytes); err != nil {
				log.Println("TableWriter.ERR.2", err)
				return err
			}
		}
		// every line for each slice
		if _, err := writer.Write(delimiter[:]); err != nil {
			log.Println("TableWriter.ERR.3", err)
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		log.Println("TableWriter.ERR.4", err)
		return err
	}
	return nil
}
