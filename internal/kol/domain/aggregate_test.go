package domain

import (
	"flag"
	"kolresource/internal/kol/domain/entities"
	"os"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestNewKol(t *testing.T) {
	t.Parallel()

	entity := &entities.Kol{
		ID:    uuid.New(),
		Name:  "test-kol",
		Email: "test@example.com",
	}

	k := NewKol(entity)

	if k.GetKol() != entity {
		t.Error("GetKol() did not return the expected entity")
	}

	if tags := k.GetTags(); len(tags) != 0 {
		t.Errorf("GetTags() length = %d, want 0", len(tags))
	}
}

func TestKol_AppendTag(t *testing.T) {
	t.Parallel()

	entity := &entities.Kol{
		ID:   uuid.New(),
		Name: "test-kol",
	}

	k := NewKol(entity)

	tag := &entities.Tag{
		ID:   uuid.New(),
		Name: "test-tag",
	}

	k.AppendTag(tag)

	tags := k.GetTags()
	if len(tags) != 1 {
		t.Fatalf("GetTags() length = %d, want 1", len(tags))
	}

	if tags[0] != tag {
		t.Error("GetTags()[0] did not return the expected tag")
	}
}

func TestNewSendEmailLog(t *testing.T) {
	t.Parallel()

	kols := []*entities.Kol{
		{ID: uuid.New(), Name: "kol-1"},
		{ID: uuid.New(), Name: "kol-2"},
	}

	products := []*entities.Product{
		{ID: uuid.New(), Name: "product-1"},
	}

	log := NewSendEmailLog(kols, products)

	_ = log
}
