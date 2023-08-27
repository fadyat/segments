package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type segmentHandlerSuite struct {
	suite.Suite
	h *Handler
}

func (s *segmentHandlerSuite) SetupSuite() {
	s.h = NewHandler(mocks.NewISegment(s.T()), api.NewRenderer())
}

func TestSegmentHandlerSuite(t *testing.T) {
	suite.Run(t, new(segmentHandlerSuite))
}
