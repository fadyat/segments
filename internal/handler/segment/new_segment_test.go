package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/mocks"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
)

type newSegmentTestCase struct {
	name           string
	body           *dto.CreateSegment
	pre            func(s *segmentHandlerSuite, tc *newSegmentTestCase)
	expected       any
	expectedStatus int
	expErr         error
}

func (s *segmentHandlerSuite) TestNewSegmentHandler() {
	testCases := []newSegmentTestCase{
		{
			name: "success",
			body: &dto.CreateSegment{Slug: "test"},
			expected: &dto.SegmentCreated{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001").String(),
				Slug: "test",
			},
			pre: func(s *segmentHandlerSuite, tc *newSegmentTestCase) {
				s.h.segmentService.(*mocks.ISegment).
					On("NewSegment", context.Background(), tc.body).Return(tc.expected, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "validation error",
			body:     &dto.CreateSegment{Slug: "te"},
			expected: api.NewUnprocessableEntityError("slug is too short, min length is 3\n"),
			pre: func(s *segmentHandlerSuite, tc *newSegmentTestCase) {
				s.h.segmentService.(*mocks.ISegment).
					On("NewSegment", context.Background(), tc.body).Return(nil, tc.expected)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			w := httptest.NewRecorder()

			body, err := json.Marshal(tc.body)
			s.Require().NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/segment", io.NopCloser(bytes.NewReader(body)))
			s.Require().NoError(err)

			tc.pre(s, &tc)
			s.h.newSegment(w, req)

			var actualRawBody = w.Body.Bytes()
			var expectedRawBody []byte
			if tc.expected != nil {
				expectedRawBody, err = json.Marshal(tc.expected)
				s.Require().NoError(err)
			}

			s.Require().Equal(tc.expectedStatus, w.Code)
			s.Require().JSONEq(string(expectedRawBody), string(actualRawBody))
		})
	}
}
