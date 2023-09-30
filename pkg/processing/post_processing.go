package processing

import (
	"log"
	"os"

	"srtor/pkg/fsutil"
)

func archive(zipGroups map[string][]string) {
	for zipPath, filePaths := range zipGroups {
		_, err := fsutil.ZipCreate(zipPath, filePaths)
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

func swap(files []*fileInfo) {
	for _, file := range files {
		err := fsutil.FileSwap(file.sourceFullPath, file.targetFullPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
