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
func (p *Processor) evaluateFile(path string) (*FileInfo, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrPathNotAbsolute
	}

	originalDir := filepath.Dir(path)
	originalBase := filepath.Base(path)
	targetFullPath := filepath.Join(originalDir, p.resultDirName, originalBase)
	zipPath := filepath.Join(originalDir, p.resultDirName, "original.zip")

	return &FileInfo{
		sourceFullPath: path,
		sourceBase:     originalBase,
		sourceDir:      originalDir,
		targetFullPath: targetFullPath,
		zipFullPath:    zipPath,
	}, nil
}
func (p *Processor) processReal(paths []string, onFileProcessed func() error) error {
	fileInfos := make([]*FileInfo, 0)
	for _, path := range paths {
		fileInfo, err := p.evaluateFile(path)
		if err != nil {
			return err
		}

		fileInfos = append(fileInfos, fileInfo)
	}
	for _, fileInfo := range fileInfos {
		err := p.processFile(fileInfo)
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

type FileInfo struct {
	sourceFullPath string
	sourceBase     string
	sourceDir      string
	targetFullPath string
	zipFullPath    string
}

func (p *Processor) processFile(fileInfo *FileInfo) error {
	originalText, err := fsutil.FileReadAsString(fileInfo.sourceFullPath)
	if err != nil {
		log.Println(err)
		return err
	}

	translatedText, err := transl.Translate(originalText, p.langSource, p.langTarget)
	if err != nil {
		log.Println(err)
		return err
	}

	err = fsutil.FileWrite(translatedText, fileInfo.targetFullPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if p.needReplace {
		err := fsutil.FileSwap(fileInfo.sourceFullPath, fileInfo.targetFullPath)
		if err != nil {
			log.Println(err)
			return err
		}

		p.forArchive[fileInfo.zipFullPath] = append(p.forArchive[fileInfo.zipFullPath], fileInfo.targetFullPath)
	}

	return nil
}
