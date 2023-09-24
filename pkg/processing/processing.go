package processing

import (
	"github.com/schollz/progressbar/v3"
	"log"
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
	resultDirName string
	needReplace   bool
	needArchive   bool
}

func NewProcessor(langSource string, langTarget string, resultDirName string) *Processor {
	return &Processor{
		langSource:    langSource,
		langTarget:    langTarget,
		numThreads:    runtime.NumCPU(),
		resultDirName: resultDirName,
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

	p.threadChunks(chunks, func() error {
		err = bar.Add(1)
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *Processor) threadChunks(chunks [][]string, onFileProcessed func() error) {
	wg := sync.WaitGroup{}

	for i := range chunks {

		wg.Add(1)

		chunk := chunks[i]
		go func(paths []string) {
			defer wg.Done()
			err := p.processReal(paths, onFileProcessed)
			if err != nil {
				log.Println(err)
			}
		}(chunk)
	}

	wg.Wait()
}

func newProgressBar(length int) *progressbar.ProgressBar {
	return progressbar.Default(int64(length))
}

func (p *Processor) processReal(paths []string, onFileProcessed func() error) error {
	for _, path := range paths {
		err := p.processFile(path)
		if err != nil {
			return err
		}

		err = onFileProcessed()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) processFile(path string) error {
	originalText, err := fsutil.ReadFileAsString(path)
	if err != nil {
		log.Println(err)
		return err
	}

	translatedText, err := transl.Translate(originalText, p.langSource, p.langTarget)
	if err != nil {
		log.Println(err)
		return err
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	resultPath := filepath.Join(dir, p.resultDirName, base)

	err = fsutil.WriteFileForced(translatedText, resultPath)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
