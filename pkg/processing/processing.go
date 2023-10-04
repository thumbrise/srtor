package processing

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/schollz/progressbar/v3"
	"srtor/pkg/fsutil"
	"srtor/pkg/util"
)

type Processor struct {
	translator    Translator
	langSource    string
	langTarget    string
	numThreads    int
	resultDirName string
	needReplace   bool
	fileInfos     []*fileInfo
	zipGroups     map[string][]string
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

type Translator interface {
	Translate(source string, sourceLang string, targetLang string) (string, error)
}

func NewProcessor(translator Translator, langSource string, langTarget string, resultDirName string) *Processor {
	return &Processor{
		translator:    translator,
		langSource:    langSource,
		langTarget:    langTarget,
		numThreads:    runtime.NumCPU(),
		resultDirName: resultDirName,
		fileInfos:     make([]*fileInfo, 0),
		zipGroups:     make(map[string][]string),
	}
}

func (p *Processor) WithReplace(v bool) *Processor {
	p.needReplace = v
	return p
}

func (p *Processor) Process(paths []string) {
	if len(paths) == 0 {
		return
	}

	p.evaluate(paths)

	chunks, err := util.SliceSplit(p.fileInfos, p.numThreads)
	if err != nil {
		log.Fatal(err)
	}

	p.bar = newProgressBar(len(paths))

	p.processChunks(chunks)

	if p.needReplace {
		swap(p.fileInfos)
		archive(p.zipGroups)
	}
}

func (p *Processor) processChunks(chunks [][]*fileInfo) {
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

func (p *Processor) evaluate(paths []string) {
	for _, path := range paths {
		file := p.evaluateFile(path)

		p.fileInfos = append(p.fileInfos, file)
		p.zipGroups[file.zipFullPath] = append(p.zipGroups[file.zipFullPath], file.targetFullPath)
	}
}

func (p *Processor) evaluateFile(path string) *fileInfo {
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

	return r
}

func (p *Processor) processFile(file *fileInfo) error {
	originalText, err := fsutil.FileReadAsString(file.sourceFullPath)
	if err != nil {
		log.Println(err)
		return err
	}

	translatedText, err := p.translator.Translate(originalText, p.langSource, p.langTarget)
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
