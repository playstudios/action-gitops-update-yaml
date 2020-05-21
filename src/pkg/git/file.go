package git

import (
	"errors"
	"github.com/go-git/go-billy/v5"
	"os"
)

func fileOpenOrCreate(fs billy.Filesystem, name string) (billy.File, error) {
	if _, err := fs.Stat(name); errors.Is(err, os.ErrNotExist) {
		return fs.Create(name)
	}
	return fs.OpenFile(name, os.O_RDWR, 0666)
}
