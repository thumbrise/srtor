package processing

import (
	"errors"
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

var ErrPathNotAbsolute = errors.New("file path is not absolute")

type Processor struct {
	langSource    string
	langTarget    string
	numThreads    int
	resultDirName string
	needReplace   bool
	forArchive    map[string][]string
}

type fileInfo struct {
	sourceFullPath string
	sourceBase     string
	sourceDir      string
	targetFullPath string
	zipFullPath    string
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

func (p *Processor) evaluateFile(path string) (*fileInfo, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrPathNotAbsolute
	}

	originalDir := filepath.Dir(path)
	originalBase := filepath.Base(path)
	targetFullPath := filepath.Join(originalDir, p.resultDirName, originalBase)
	zipPath := filepath.Join(originalDir, p.resultDirName, "original.zip")

	r := &fileInfo{
		sourceFullPath: path,
		sourceBase:     originalBase,
		sourceDir:      originalDir,
		targetFullPath: targetFullPath,
		zipFullPath:    zipPath,
	}

	return r, nil
}

func (p *Processor) processReal(paths []string, onFileProcessed func() error) error {
	files := make([]*fileInfo, 0)
	for _, path := range paths {
		file, err := p.evaluateFile(path)
		if err != nil {
			return err
		}

		files = append(files, file)
	}
	for _, file := range files {
		err := p.processFile(file)
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

func (p *Processor) processFile(file *fileInfo) error {
	originalText, err := fsutil.FileReadAsString(file.sourceFullPath)
	if err != nil {
		log.Println(err)
		return err
	}

	translatedText, err := transl.Translate(originalText, p.langSource, p.langTarget)
	if err != nil {
		log.Println(err)
		return err
	}

	err = fsutil.FileWrite(translatedText, file.targetFullPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if p.needReplace {
		err := fsutil.FileSwap(file.sourceFullPath, file.targetFullPath)
		if err != nil {
			log.Println(err)
			return err
		}

		p.forArchive[file.zipFullPath] = append(p.forArchive[file.zipFullPath], file.targetFullPath)
	}

	return nil
}
