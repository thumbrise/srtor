package processing

import "github.com/schollz/progressbar/v3"

func newProgressBar(length int) *progressbar.ProgressBar {
	return progressbar.Default(int64(length))
}
