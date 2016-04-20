package fs

type FileExplorer interface {
	Init() error
	ListDir(path string) ([]ListDirEntry, error)
	Move(path string, newPath string) error
	Copy(path string, newPath string) error
	Delete(path string) error
	Rename(oldPath string, newPath string) error
	// Chmod(path string, code string) error
	Mkdir(path string, name string) error
	Close() error
}
