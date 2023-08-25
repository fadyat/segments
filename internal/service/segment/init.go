package segment

import (
	"avito-internship-2023/internal/repository"
	"avito-internship-2023/internal/service"
)

type svc struct {
	segmentRepository repository.ISegment
}

func NewService(repo repository.ISegment) service.ISegment {
	return &svc{segmentRepository: repo}
}
