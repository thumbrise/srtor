package processing

import (
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"srtor/pkg/fsutil"
	"srtor/pkg/transl"
	"srtor/pkg/util"
	"sync"
)

type Processor struct {
	langSource    string
	langTarget    string
	numThreads    int
	targetDirName string
	needReplace   bool
	needArchive   bool
}

func NewProcessor(langSource string, langTarget string, destination string) *Processor {
	return &Processor{
		langSource:    langSource,
		langTarget:    langTarget,
		numThreads:    runtime.NumCPU(),
		targetDirName: destination,
	}
}
func (p *Processor) WithReplace(v bool) *Processor {
	p.needReplace = v
	return p
}
func (p *Processor) WithArchive(v bool) *Processor {
	p.needArchive = v
	return p
}
func (p *Processor) Process(files []string) {
	filesLen := len(files)

	if filesLen == 0 {
		return
	}

	bar := newProgressBar(filesLen)
	chunks, err := util.SplitSlice(files, p.numThreads)
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
func newProgressBar(length int) *progressbar.ProgressBar {
	return progressbar.Default(int64(length))
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

	translated, err := transl.Translate(string(source), p.langSource, p.langTarget)
	if err != nil {
		log.Println(err)
		return err
	}

	sourceName := filepath.Base(path)
	sourceDir := filepath.Dir(path)
	destination := filepath.Join(sourceDir, p.targetDirName)

	err = fsutil.MkdirOrIgnore(destination)
	if err != nil {
		log.Println(err)
		return err
	}

	resultPath := filepath.Join(destination, sourceName)
	resultBytes := util.FixUTF8(translated)

	err = os.WriteFile(resultPath, resultBytes, os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
