package sqlboiler

// import (
// 	"context"
// 	"crypto/rand"
// 	"database/sql"
// 	"flag"
// 	"fmt"
// 	"io"
// 	"kolresource/database"
// 	"kolresource/internal/kol"
// 	"kolresource/internal/kol/domain"

// 	"kolresource/internal/kol/domain/entities"
// 	"kolresource/pkg/database/postgreinit"
// 	"log"
// 	"net/url"
// 	"os"
// 	"reflect"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/golang-migrate/migrate/v4"
// 	_ "github.com/golang-migrate/migrate/v4/database/postgres"
// 	"github.com/golang-migrate/migrate/v4/source/iofs"
// 	"github.com/google/uuid"
// 	"github.com/ory/dockertest/v3"
// 	"github.com/ory/dockertest/v3/docker"
// 	"github.com/rs/zerolog"
// 	"github.com/volatiletech/sqlboiler/v4/boil"
// 	"go.uber.org/goleak"
// )

// type suitUp struct {
// 	pgHostPort string
// 	stdConn    *sql.DB
// }

// var suitUpInstance *suitUp

// func TestMain(m *testing.M) {
// 	boil.DebugMode = true

// 	leak := flag.Bool("leak", false, "use leak detector")
// 	flag.Parse()

// 	if *leak {
// 		goleak.VerifyTestMain(m)

// 		return
// 	}

// 	var err error

// 	dockerTestPool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	dockerTestPool.MaxWait = 120 * time.Second

// 	dockerTestResource, err := dockerTestPool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "17.0",
// 		Env: []string{
// 			"POSTGRES_PASSWORD=postgres",
// 			"POSTGRES_USER=postgres",
// 			"POSTGRES_DB=test",
// 			"listen_addresses = '*'",
// 		},
// 		Cmd: []string{
// 			"postgres",
// 			"-N",
// 			"200",
// 		},
// 	}, func(config *docker.HostConfig) {
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not start testRefDockertestResource: %s", err)
// 	}

// 	_ = dockerTestResource.Expire(180)

// 	pgHostPort := getHostPort(dockerTestResource, "5432/tcp")

// 	var pgStdConn *sql.DB

// 	if err = dockerTestPool.Retry(func() error {
// 		var stdConnErr error
// 		dbPort := strings.Split(pgHostPort, ":")[1]
// 		pgStdConn, stdConnErr = newStdConn(dbPort, "test")

// 		if stdConnErr != nil {
// 			return stdConnErr
// 		}

// 		if errPing := pgStdConn.Ping(); err != nil {
// 			return errPing
// 		}

// 		return nil
// 	}); err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	var teardown func()
// 	suitUpInstance, teardown = newSuitUp(pgStdConn, pgHostPort)

// 	code := m.Run()
// 	teardown()
// 	// You can't defer this because os.Exit doesn't care for defer
// 	if err := dockerTestPool.Purge(dockerTestResource); err != nil {
// 		log.Fatalf("Could not purge testRefDockertestResource: %s", err)
// 	}

// 	os.Exit(code)
// }

// func getHostPort(resource *dockertest.Resource, id string) string {
// 	dockerURL := os.Getenv("DOCKER_HOST")
// 	if dockerURL == "" {
// 		hostAndPort := resource.GetHostPort("5432/tcp")
// 		hp := strings.Split(hostAndPort, ":")
// 		testRefHost := hp[0]
// 		testRefPort := hp[1]

// 		return testRefHost + ":" + testRefPort
// 	}

// 	u, err := url.Parse(dockerURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return u.Hostname() + ":" + resource.GetPort(id)
// }

// func newStdConn(dbPort string, database string) (*sql.DB, error) {
// 	logger := zerolog.New(io.Discard)

// 	pgInit, err := postgreinit.New(
// 		&postgreinit.Config{
// 			Host:         "localhost",
// 			Port:         dbPort,
// 			User:         "postgres",
// 			Password:     "postgres",
// 			Database:     database,
// 			MaxConns:     10,
// 			MaxIdleConns: 10,
// 			MaxLifeTime:  1 * time.Minute,
// 		},
// 		postgreinit.WithLogLevel(zerolog.WarnLevel),
// 		postgreinit.WithLogger(&logger, "request-id"),
// 	)
// 	if err != nil {
// 		log.Fatalf("Could not init pginit: %s", err)
// 	}

// 	stdConn, err := pgInit.StdConn(context.Background())

// 	return stdConn, err
// }

// func newSuitUp(stdConn *sql.DB, pgHostPort string) (*suitUp, func()) {
// 	s := &suitUp{
// 		stdConn:    stdConn,
// 		pgHostPort: pgHostPort,
// 	}

// 	s.stdConn = s.generateNewStdConn()
// 	return s, func() {
// 		s.stdConn.Close()
// 		// s.defaultPool.Close()
// 		// s.changePool.Close()
// 	}
// }

// // generateNewConnPool If you want to have an isolated testing environment, you can create one from this function.
// func (s *suitUp) generateNewStdConn() *sql.DB {
// 	dbname := s.createDatabase(s.stdConn)

// 	dbPort := strings.Split(s.pgHostPort, ":")[1]

// 	stdConn, err := newStdConn(dbPort, dbname)
// 	if err != nil {
// 		log.Fatalf("Could not connect to pool: %s", err)
// 	}

// 	databaseURL := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", "postgres", s.pgHostPort, dbname)
// 	if err := s.runMigrations(databaseURL); err != nil {
// 		log.Fatalf("Could not run migrations: %s", err)
// 	}

// 	return stdConn
// }

// func (s *suitUp) createDatabase(stdConn *sql.DB) string {
// 	randStr := randomString(10)
// 	dbName := fmt.Sprintf("test_env_%s", randStr)

// 	sql := fmt.Sprintf("CREATE DATABASE %s", dbName)
// 	if _, err := stdConn.Exec(sql); err != nil {
// 		log.Fatalf("Could not create database: %v", err)
// 	}

// 	return dbName
// }

// func (s *suitUp) runMigrations(dbURL string) error {
// 	d, err := iofs.New(database.MigrationFiles, "migrations")
// 	if err != nil {
// 		return err
// 	}

// 	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
// 	if err != nil {
// 		return err
// 	}

// 	err = m.Up()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func randomString(length int) string {
// 	b := make([]byte, length)
// 	rand.Read(b)

// 	return fmt.Sprintf("%x", b)[:length]
// }

// func TestNewKolRepository(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *KolRepository
// 	}{
// 		{
// 			name: "happy",
// 			want: &KolRepository{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := NewKolRepository(suitUpInstance.stdConn)

// 			if reflect.TypeOf(repo) != reflect.TypeOf(tt.want) {
// 				t.Errorf("returned %v is not want %v", repo, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCreateKol(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		param   domain.CreateKolParams
// 		wantErr bool
// 	}{
// 		{
// 			name: "happy path",
// 			param: domain.CreateKolParams{
// 				Name:           "Test KOL",
// 				Email:          "test@example.com",
// 				Description:    "Test description",
// 				Sex:            "m",
// 				Enable:         true,
// 				UpdatedAdminID: uuid.New(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "dumplicate email error",
// 			param: domain.CreateKolParams{
// 				Name:           "Test KOL",
// 				Email:          "test@example.com",
// 				Description:    "Test description",
// 				Sex:            "m",
// 				Enable:         true,
// 				UpdatedAdminID: uuid.New(),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "very long name",
// 			param: domain.CreateKolParams{
// 				Name:           strings.Repeat("a", 256),
// 				Email:          "test@example.com",
// 				Description:    "Test description",
// 				Sex:            "f",
// 				Enable:         true,
// 				UpdatedAdminID: uuid.New(),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid sex",
// 			param: domain.CreateKolParams{
// 				Name:           "Test KOL",
// 				Email:          "test@example.com",
// 				Description:    "Test description",
// 				Sex:            "invalid",
// 				Enable:         true,
// 				UpdatedAdminID: uuid.New(),
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := NewKolRepository(suitUpInstance.stdConn)
// 			ctx := context.Background()

// 			kol, err := repo.CreateKol(ctx, tt.param)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CreateKol() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if !tt.wantErr {
// 				if kol == nil {
// 					t.Error("CreateKol() returned nil KOL when no error was expected")
// 					return
// 				}

// 				// Verify the created KOL
// 				if kol.Name != tt.param.Name {
// 					t.Errorf("Created KOL name = %v, want %v", kol.Name, tt.param.Name)
// 				}
// 				if kol.Email != tt.param.Email {
// 					t.Errorf("Created KOL email = %v, want %v", kol.Email, tt.param.Email)
// 				}
// 				if kol.Description != tt.param.Description {
// 					t.Errorf("Created KOL description = %v, want %v", kol.Description, tt.param.Description)
// 				}
// 				if kol.Sex != tt.param.Sex {
// 					t.Errorf("Created KOL sex = %v, want %v", kol.Sex, tt.param.Sex)
// 				}
// 				if kol.Enable != tt.param.Enable {
// 					t.Errorf("Created KOL enable = %v, want %v", kol.Enable, tt.param.Enable)
// 				}
// 				if kol.UpdatedAdminID != tt.param.UpdatedAdminID {
// 					t.Errorf("Created KOL updatedAdminID = %v, want %v", kol.UpdatedAdminID, tt.param.UpdatedAdminID)
// 				}
// 			}
// 		})
// 	}
// }

// func TestGetKolByID(t *testing.T) {
// 	repo := NewKolRepository(suitUpInstance.stdConn)
// 	ctx := context.Background()

// 	newKol, err := repo.CreateKol(ctx, domain.CreateKolParams{
// 		Name:           "stanley",
// 		Email:          "stanley@gmail.com",
// 		Description:    "description",
// 		Sex:            "m",
// 		Enable:         true,
// 		UpdatedAdminID: uuid.Must(uuid.NewV7()),
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to create test KOL: %v", err)
// 	}

// 	tests := []struct {
// 		name    string
// 		id      uuid.UUID
// 		wantKol *entities.Kol
// 		wantErr bool
// 	}{
// 		{
// 			name:    "happy",
// 			id:      newKol.ID,
// 			wantKol: newKol,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "data not found",
// 			id:      uuid.Must(uuid.NewV7()),
// 			wantKol: nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			gotKol, err := repo.GetKolByID(ctx, tt.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetKolByID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if !tt.wantErr {
// 				if gotKol == nil {
// 					t.Error("GetKolByID() returned nil KOL when no error was expected")
// 					return
// 				}

// 				if !reflect.DeepEqual(gotKol, tt.wantKol) {
// 					t.Errorf("got %v, but want %v", gotKol, tt.wantKol)
// 				}
// 			} else if gotKol != nil {
// 				t.Errorf("GetKolByID() = %v, want nil", gotKol)
// 			}
// 		})
// 	}
// }

// func TestUpdateKol(t *testing.T) {
// 	ctx := context.Background()
// 	repo := NewKolRepository(suitUpInstance.stdConn)

// 	// Create a test KOL to update
// 	initialKol, err := repo.CreateKol(ctx, domain.CreateKolParams{
// 		Name:           "Initial Name",
// 		Email:          "initial@example.com",
// 		Description:    "Initial description",
// 		Sex:            kol.SexMale,
// 		Enable:         true,
// 		UpdatedAdminID: uuid.New(),
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to create initial test KOL: %v", err)
// 	}

// 	tests := []struct {
// 		name    string
// 		param   domain.UpdateKolParams
// 		wantKol *entities.Kol
// 		wantErr bool
// 	}{
// 		{
// 			name: "happy path",
// 			param: domain.UpdateKolParams{
// 				ID:             initialKol.ID,
// 				Name:           "Updated Name",
// 				Email:          "updated@example.com",
// 				Description:    "Updated description",
// 				Sex:            kol.SexFemale,
// 				Enable:         false,
// 				UpdatedAdminID: uuid.MustParse("6a9f3135-8415-4002-8011-f8c1d283964e"),
// 			},
// 			wantKol: &entities.Kol{
// 				ID:             initialKol.ID,
// 				Name:           "Updated Name",
// 				Email:          "updated@example.com",
// 				Description:    "Updated description",
// 				Sex:            kol.SexFemale,
// 				Enable:         false,
// 				UpdatedAdminID: uuid.MustParse("6a9f3135-8415-4002-8011-f8c1d283964e"),
// 				CreatedAt:      initialKol.CreatedAt,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "non-existent KOL",
// 			param: domain.UpdateKolParams{
// 				ID:             uuid.New(),
// 				Name:           "Non-existent",
// 				Email:          "nonexistent@example.com",
// 				Description:    "This KOL doesn't exist",
// 				Sex:            kol.SexMale,
// 				Enable:         true,
// 				UpdatedAdminID: uuid.New(),
// 			},
// 			wantKol: nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotKol, err := repo.UpdateKol(ctx, tt.param)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateKol() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !tt.wantErr {
// 				if !reflect.DeepEqual(gotKol, tt.wantKol) {
// 					t.Errorf("UpdateKol() = %v, want %v", gotKol, tt.wantKol)
// 				}
// 			}
// 		})
// 	}
// }
