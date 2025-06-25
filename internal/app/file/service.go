package file

import (
	"backend/internal/domain/file"
)

type Service struct {
	repo file.FileRepository
}

func NewService(repo file.FileRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetFileURL(key, bucket string) (string, error){
	return s.repo.GetFileURL(key, bucket)
}