package processing

import (
	"bytes"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"srtor/pkg/transl"
	"srtor/pkg/util"
	"sync"
)

var unicodeReplacement = []byte{0xef, 0xbf, 0xbd}

type Processor struct {
	langSource    string
	langTarget    string
	numThreads    int
	targetDirName string
}

func NewProcessor(langSource string, langTarget string, destination string) Processor {
	return Processor{
		langSource:    langSource,
		langTarget:    langTarget,
		numThreads:    runtime.NumCPU(),
		targetDirName: destination,
	}
}
func (p *Processor) Process(files []string) {
	filesLen := len(files)
	bar := progressbar.Default(int64(filesLen))
	numGoroutines := util.Max(p.numThreads, 0)
	numGoroutines = util.Min(numGoroutines, filesLen)
	chunkSize := filesLen / numGoroutines
	chunks, err := util.ChunkSlice(files, chunkSize)
	if err != nil {
		log.Println(err)
	}

	wg := sync.WaitGroup{}
	for i := range chunks {
		wg.Add(1)
		go func(paths []string, bar *progressbar.ProgressBar) {
			defer wg.Done()
			err := p.iteratePaths(paths, bar)
			if err != nil {
				log.Println(err)
			}
		}(chunks[i], bar)
	}
	wg.Wait()
}

func (p *Processor) iteratePaths(paths []string, bar *progressbar.ProgressBar) error {
	for _, path := range paths {
		err := p.processFile(path)
		if err != nil {
			return err
		}
		err = bar.Add(1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) processFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	target, err := transl.Translate(string(source), p.langSource, p.langTarget)
	if err != nil {
		return err
	}

	sourceName := filepath.Base(path)
	targetPath := filepath.Join(p.targetDirName, sourceName)
	targetBytes := []byte(target)

	targetBytes = bytes.ToValidUTF8(targetBytes, unicodeReplacement)

	err = os.WriteFile(targetPath, targetBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
