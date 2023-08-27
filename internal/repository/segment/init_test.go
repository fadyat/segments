package segment

import (
	"avito-internship-2023/internal/repository"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
	"time"
)

type SegmentRepoSuite struct {
	suite.Suite
	r  repository.ISegment
	db *sqlx.DB
}

func newTestDatabase() *sqlx.DB {
	dsn, driver := os.Getenv("TEST_DATABASE_DSN"), os.Getenv("TEST_DATABASE_DRIVER")
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	return db
}

func applyMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrationsPath := os.Getenv("TEST_DATABASE_MIGRATIONS_PATH")
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *SegmentRepoSuite) SetupSuite() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	s.db = newTestDatabase()
	s.r = NewRepository(repository.NewTransactor(s.db))
	err := applyMigrations(s.db.DB)
	if err != nil {
		log.Fatal("failed to apply migrations: ", err)
	}
}

func (s *SegmentRepoSuite) clean() {
	_, err := s.db.Exec("TRUNCATE TABLE segment CASCADE")
	if err != nil {
		s.Require().NoError(err)
	}
}

func (s *SegmentRepoSuite) TearDownTest() {
	s.clean()
}

func (s *SegmentRepoSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		log.Fatal("failed to close database connection: ", err)
	}
}

func TestSegmentRepository(t *testing.T) {
	suite.Run(t, new(SegmentRepoSuite))
}

func equalTimes(a, b *time.Time) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	ay, am, ad := a.Date()
	by, bm, bd := b.Date()

	ah, amin, _ := a.Clock()
	bh, bmin, _ := b.Clock()

	return ay == by && am == bm && ad == bd && ah == bh && amin == bmin
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}
