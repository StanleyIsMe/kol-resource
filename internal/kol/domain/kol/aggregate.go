package kol

import (
	"kolresource/internal/kol/entities"
)

type Kol struct {
	kol            *entities.Kol
	tags           []*entities.Tag
	UpdatedAdminID string
}

func NewKol(kol *entities.Kol) *Kol {
	return &Kol{
		kol:  kol,
		tags: make([]*entities.Tag, 0),
	}
}

func (k *Kol) AppendTag(tag *entities.Tag) {
	k.tags = append(k.tags, tag)
}

type SendEmailLog struct {
	kols     []*entities.Kol
	products []*entities.Product
}

func NewSendEmailLog(kols []*entities.Kol, products []*entities.Product) SendEmailLog {
	return SendEmailLog{
		kols:     kols,
		products: products,
	}
}
