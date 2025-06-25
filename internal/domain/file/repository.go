package file


type FileRepository interface {
	GetFileURL(key, bucket string) (string, error)
}