//go:build unit

package httpServer

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/homework8/internal/pkg/model"
	"gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository"
	studentService "gitlab.ozon.dev/homework8/internal/pkg/service/student"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_repository "gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository/mocks"
)

func Test_getUser(t *testing.T) {
	t.Parallel()
	var (
		id = 1
	)
	t.Run("success,db", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mock_repository.NewMockStudentsRepo(ctrl)
		mc := mock_repository.NewMockStudentsRepoCached(ctrl)

		service := studentService.New(m, mc)
		s := New(service)

		req, err := http.NewRequest(http.MethodGet, "student?studentId=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)
		res := httptest.NewRecorder()

		mc.EXPECT().Get(gomock.Any(), uint64(id)).Return(nil, studentsRepository.ErrObjectNotFound)
		m.EXPECT().GetById(gomock.Any(), uint64(id)).Return(&model.Student{ID: uint64(id), Name: "Andrew"}, nil)

		// act

		s.GetStudent(res, req)
		// assert
		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Parallel()
		tt := []struct {
			name    string
			request *url.URL
			isOk    bool
		}{
			{
				"without id",
				&url.URL{RawQuery: "student?studentId"},
				false,
			},
			{
				"wrong id",
				&url.URL{RawQuery: "student?studentId=AndrewAndrew"},
				false,
			},
			{
				"empty",
				&url.URL{RawQuery: ""},
				false,
			},
		}
		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				_, err := getStudentID(tc.request)
				require.Error(t, err)
				assert.EqualError(t, err, "invalid student.id")
			})
		}
	})
}

func Test_createUser(t *testing.T) {
	t.Parallel()

	t.Run("success,db", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mock_repository.NewMockStudentsRepo(ctrl)
		mc := mock_repository.NewMockStudentsRepoCached(ctrl)

		service := studentService.New(m, mc)
		s := New(service)
		student := model.Student{Name: "Andrew", CourseProgram: "Go"}
		cashedStudent := model.Student{Name: "Andrew", CourseProgram: "Go"}
		jsonBody, err := json.Marshal(student)

		req, err := http.NewRequest(http.MethodPost, "student", bytes.NewReader(jsonBody))
		require.NoError(t, err)
		res := httptest.NewRecorder()

		m.EXPECT().Add(gomock.Any(), &student).Return(&student, nil)
		mc.EXPECT().Add(gomock.Any(), &cashedStudent).Return(nil)

		// act

		s.CreateStudent(res, req)
		// assert
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Parallel()
		tt := []struct {
			name    string
			student *model.Student
		}{
			{
				"without name",
				&model.Student{CourseProgram: "Go"},
			},
			{
				"empty",
				nil,
			},
		}
		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				json, _ := json.Marshal(tc.student)
				_, err := getStudentFromBody(json)
				require.Error(t, err)
				assert.EqualError(t, err, "invalid body")
			})
		}
	})
}

func Test_updateUser(t *testing.T) {
	t.Parallel()
	var (
		id = 1
	)
	t.Run("success,db", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mock_repository.NewMockStudentsRepo(ctrl)
		mc := mock_repository.NewMockStudentsRepoCached(ctrl)

		service := studentService.New(m, mc)
		s := New(service)
		student := model.Student{ID: uint64(id), Name: "Andrew", CourseProgram: "Go"}
		jsonBody, err := json.Marshal(student)

		req, err := http.NewRequest(http.MethodPost, "student", bytes.NewReader(jsonBody))
		require.NoError(t, err)
		res := httptest.NewRecorder()

		m.EXPECT().Update(gomock.Any(), &student).Return(true, nil)
		mc.EXPECT().Add(gomock.Any(), &student).Return(nil)

		// act

		s.UpdateStudent(res, req)
		// assert
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mock_repository.NewMockStudentsRepo(ctrl)
		mc := mock_repository.NewMockStudentsRepoCached(ctrl)
		service := studentService.New(m, mc)
		s := New(service)
		student := model.Student{ID: uint64(id), Name: "Andrew", CourseProgram: "Go"}
		jsonBody, err := json.Marshal(student)

		req, err := http.NewRequest(http.MethodPost, "student", bytes.NewReader(jsonBody))
		require.NoError(t, err)
		res := httptest.NewRecorder()

		m.EXPECT().Update(gomock.Any(), &student).Return(false, studentsRepository.ErrObjectNotFound)

		// act
		s.UpdateStudent(res, req)
		// assert
		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}
