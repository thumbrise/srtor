package processing

import (
	"errors"
	"fmt"
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
	bar           *progressbar.ProgressBar
	sync.Mutex
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
	}
}

func (p *Processor) WithReplace(v bool) *Processor {
	p.needReplace = v
	return p
}

func (p *Processor) Process(paths []string) {
	filesLen := len(paths)

	if filesLen == 0 {
		return
	}

	p.bar = newProgressBar(filesLen)

	files := make([]*fileInfo, 0)
	for _, path := range paths {
		file, err := p.evaluateFile(path)
		if err != nil {
			log.Fatal(err)
		}

		files = append(files, file)
	}

	chunks, err := util.SliceSplit(files, p.numThreads)
	if err != nil {
		log.Println(err)
	}
	p.threadChunks(chunks)

	forArchive := make(map[string][]string)
	if p.needReplace {
		for _, file := range files {

			err := fsutil.FileSwap(file.sourceFullPath, file.targetFullPath)
			if err != nil {
				log.Fatal(err)
			}

			forArchive[file.zipFullPath] = append(forArchive[file.zipFullPath], file.targetFullPath)
		}
	}
	for zipPath, filePaths := range forArchive {
		_, err = fsutil.ZipCreate(zipPath, filePaths)
		if err != nil {
			log.Fatal(err)
		}

		for _, filePath := range filePaths {
			err = os.Remove(filePath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (p *Processor) threadChunks(chunks [][]*fileInfo) {
	wg := sync.WaitGroup{}

	for i := range chunks {
		wg.Add(1)

		chunk := chunks[i]
		go func(files []*fileInfo) {
			defer wg.Done()

			for _, file := range files {
				err := p.processFile(file)
				if err != nil {
					err = fmt.Errorf("error on file %s\n%v", file.sourceFullPath, err)
					log.Println(err)
				}
			}

			p.Lock()
			err := p.bar.Add(1)
			if err != nil {
				log.Println(err)
			}
			p.Unlock()
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

	return nil
}
