package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"squzy/internal/database/convertion"
	"testing"
)

//docker run -d --rm --name postgres -e POSTGRES_USER="user" -e POSTGRES_PASSWORD="password" -e POSTGRES_DB="database" -p 5432:5432 postgres
var (
	postgr      = &postgres{}
	postgrWrong = &postgres{}
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	postgr.db = s.DB

	s.DB.LogMode(true)
}

func TestPostgres_NewClient(t *testing.T) {
	t.Run("wrongPostgress", func(t *testing.T) {
		err := postgrWrong.newClient(func() (db *gorm.DB, e error) {
			return gorm.Open(
				"postgres",
				fmt.Sprintf("host=lkl port=00 user=us dbname=dbn password=ps connect_timeout=10 sslmode=disable"))
		})
		assert.Error(t, err)
	})
}

func (s *Suite) Test_InsertMetaData() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbSchedulerCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := postgr.InsertSnapshot(&apiPb.SchedulerResponse{})
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetMetaData() {
	var (
		id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbSchedulerCollection, dbSchedulerCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetSnapshots(id)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_InsertStatRequest() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s"`, dbStatRequestCollection)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	err := postgr.InsertStatRequest(convertion.ConvertFromPostgressStatRequest(&StatRequest{}))
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetStatRequest() {
	var (
		id = "1"
	)
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s"."deleted_at" IS NULL`, dbStatRequestCollection, dbStatRequestCollection)
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	_, err := postgr.GetStatRequest(id)
	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func TestPostgres_Migrate(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.Migrate()
		assert.Error(t, err)
	})
}

func TestPostgres_InsertMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.InsertSnapshot(&apiPb.SchedulerResponse{})
		assert.Error(t, err)
	})
}

func TestPostgres_GetMetaData(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetSnapshots("")
		assert.Error(t, err)
	})
}

func TestPostgres_InsertStatRequest(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		err := postgrWrong.InsertStatRequest(convertion.ConvertFromPostgressStatRequest(&StatRequest{}))
		assert.Error(t, err)
	})
}

func TestPostgres_GetStatRequest(t *testing.T) {
	t.Run("Should: return error", func(t *testing.T) {
		_, err := postgrWrong.GetStatRequest("")
		assert.Error(t, err)
	})
}
