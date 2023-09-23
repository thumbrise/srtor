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

type Processor struct {
	LangSource    string
	LangTarget    string
	NumGoroutines int
	TargetDirName string
}

func NewProcessor() Processor {
	return Processor{
		LangSource:    "en",
		LangTarget:    "ru",
		NumGoroutines: runtime.NumCPU(),
		TargetDirName: "srtor",
	}
}
func (p *Processor) Process(files []string) {
	filesLen := len(files)
	bar := progressbar.Default(int64(filesLen))
	numGoroutines := util.Max(p.NumGoroutines, 0)
	numGoroutines = util.Min(numGoroutines, filesLen)
	chunkSize := filesLen / numGoroutines
	chunks := util.ChunkSlice(files, chunkSize)

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
	target, err := transl.Translate(string(source), p.LangSource, p.LangTarget)
	if err != nil {
		return err
	}

	sourceDir := filepath.Dir(path)
	sourceName := filepath.Base(path)
	targetPath := filepath.Join(sourceDir, p.TargetDirName, sourceName)
	targetBytes := []byte(target)
	unicodeReplacement := []byte{0xef, 0xbf, 0xbd}
	targetBytes = bytes.ToValidUTF8(targetBytes, unicodeReplacement)

	err = os.WriteFile(targetPath, targetBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
