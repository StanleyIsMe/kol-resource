package sqlboiler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"kolresource/database"

	"kolresource/pkg/database/postgreinit"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
)

type suitUp struct {
	pgHostPort string
	stdConn    *sql.DB
}

var suitUpInstance *suitUp

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	var err error

	dockerTestPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dockerTestPool.MaxWait = 120 * time.Second

	dockerTestResource, err := dockerTestPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17.0",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=test",
			"listen_addresses = '*'",
		},
		Cmd: []string{
			"postgres",
			"-N",
			"200",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start testRefDockertestResource: %s", err)
	}

	_ = dockerTestResource.Expire(180)

	pgHostPort := getHostPort(dockerTestResource, "5432/tcp")

	var pgStdConn *sql.DB

	if err = dockerTestPool.Retry(func() error {
		var stdConnErr error
		dbPort := strings.Split(pgHostPort, ":")[1]
		pgStdConn, stdConnErr = newStdConn(dbPort, "test")

		if stdConnErr != nil {
			return stdConnErr
		}

		if errPing := pgStdConn.Ping(); errPing != nil {
			pgStdConn.Close()

			return errPing
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	var teardown func()
	suitUpInstance, teardown = newSuitUp(pgStdConn, pgHostPort)

	code := m.Run()
	teardown()

	// You can't defer this because os.Exit doesn't care for defer
	if err := dockerTestPool.Purge(dockerTestResource); err != nil {
		log.Fatalf("Could not purge testRefDockertestResource: %s", err)
	}

	// Mark first due to the detector is always find the database/sql.(*DB).connectionOpener
	if *leak && code == 0 {
		// 	if err := goleak.Find(); err != nil {
		// 		log.Printf("goleak: Errors on successful test run: %v\n", err)

		// 		code = 1
		// }
	}

	os.Exit(code)
}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		hostAndPort := resource.GetHostPort("5432/tcp")
		hp := strings.Split(hostAndPort, ":")
		testRefHost := hp[0]
		testRefPort := hp[1]

		return testRefHost + ":" + testRefPort
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}

	return u.Hostname() + ":" + resource.GetPort(id)
}

func newStdConn(dbPort string, database string) (*sql.DB, error) {
	logger := zerolog.New(io.Discard)

	pgInit, err := postgreinit.New(
		&postgreinit.Config{
			Host:         "localhost",
			Port:         dbPort,
			User:         "postgres",
			Password:     "postgres",
			Database:     database,
			MaxConns:     1,
			MaxIdleConns: 1,
			MaxLifeTime:  1 * time.Minute,
		},
		postgreinit.WithLogLevel(zerolog.WarnLevel),
		postgreinit.WithLogger(&logger, "request-id"),
	)
	if err != nil {
		log.Fatalf("Could not init pginit: %s", err)
	}

	stdConn, err := pgInit.StdConn(context.Background())

	return stdConn, err
}

func newSuitUp(stdConn *sql.DB, pgHostPort string) (*suitUp, func()) {
	s := &suitUp{
		stdConn:    stdConn,
		pgHostPort: pgHostPort,
	}

	s.stdConn = s.generateNewStdConn()
	return s, func() {
		stdConn.Close()
		s.stdConn.Close()
	}
}

// generateNewStdConn If you want to have an isolated testing environment, you can create one from this function.
func (s *suitUp) generateNewStdConn() *sql.DB {
	dbname := s.createDatabase(s.stdConn)

	dbPort := strings.Split(s.pgHostPort, ":")[1]

	stdConn, err := newStdConn(dbPort, dbname)
	if err != nil {
		log.Fatalf("Could not connect to pool: %s", err)
	}

	databaseURL := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", "postgres", s.pgHostPort, dbname)
	if err := s.runMigrations(databaseURL); err != nil {
		stdConn.Close()

		log.Fatalf("Could not run migrations: %s", err)
	}

	return stdConn
}

func (s *suitUp) createDatabase(stdConn *sql.DB) string {
	randStr := randomString(10)
	dbName := fmt.Sprintf("test_env_%s", randStr)

	sql := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if _, err := stdConn.Exec(sql); err != nil {
		log.Fatalf("Could not create database: %v", err)
	}

	return dbName
}

func (s *suitUp) runMigrations(dbURL string) error {
	d, err := iofs.New(database.MigrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)

	return fmt.Sprintf("%x", b)[:length]
}

func TestNewKolRepository(t *testing.T) {
	tests := []struct {
		name string
		want *KolRepository
	}{
		{
			name: "happy",
			want: &KolRepository{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			repo := NewKolRepository(suitUpInstance.stdConn)

			if reflect.TypeOf(repo) != reflect.TypeOf(tt.want) {
				t.Errorf("returned %v is not want %v", repo, tt.want)
			}
		})
	}
}
