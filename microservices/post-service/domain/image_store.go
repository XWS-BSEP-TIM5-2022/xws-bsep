package domain

type ImageStore interface {
	UploadImage(image []byte) (string, error)
	GetImage(filename string) []byte
}
