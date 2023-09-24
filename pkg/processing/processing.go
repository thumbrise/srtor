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
	resultDirName string
	needReplace   bool
	needArchive   bool
	forArchive    map[string][]string
}

func NewProcessor(langSource string, langTarget string, resultDirName string) *Processor {
	return &Processor{
		langSource:    langSource,
		langTarget:    langTarget,
		numThreads:    runtime.NumCPU(),
		resultDirName: resultDirName,
		forArchive:    make(map[string][]string),
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
	chunks, err := util.SliceSplit(files, p.numThreads)
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

	for zipPath, filePaths := range p.forArchive {
		err := fsutil.ZipCreate(zipPath, filePaths)
		if err != nil {
			log.Fatal(err)
		}

		for _, filePath := range filePaths {
			err := os.Remove(filePath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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
	originalText, err := fsutil.FileReadAsString(path)
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

	err = fsutil.FileWrite(translatedText, resultPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if p.needReplace {
		err := fsutil.FileSwap(path, resultPath)
		if err != nil {
			log.Println(err)
			return err
		}
		if p.needArchive {
			zipPath := filepath.Join(dir, p.resultDirName, "original.zip")
			p.forArchive[zipPath] = append(p.forArchive[zipPath], resultPath)
		}
	}

	return nil
}
