//go:build integration

package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/homework7/test/fixtures"
	"net/http"
	"testing"
)

func Test_getUser(t *testing.T) {
	t.Run("success,db", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student().Valid().P()
		jsonStudent, err := json.Marshal(student)
		_, err = http.Post(Ts.URL+"/student", "", bytes.NewReader(jsonStudent))

		res, err := http.Get(Ts.URL + "/student?studentId=1")

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		res, err := http.Get(Ts.URL + "/student?studentId")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func Test_createUser(t *testing.T) {
	t.Run("success,db", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student().Valid().P()
		jsonStudent, err := json.Marshal(student)

		res, err := http.Post(Ts.URL+"/student", "", bytes.NewReader(jsonStudent))

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

	})

	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student()
		jsonStudent, err := json.Marshal(student)

		res, err := http.Post(Ts.URL+"/student", "", bytes.NewReader(jsonStudent))

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
