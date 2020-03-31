package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DenisGubenko/ideasymbols/db"
	"github.com/DenisGubenko/ideasymbols/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	maxOrder = 50
	letters1 = `qwertyuiiop`
	letter2  = `sdfghjkl`
	letter3  = `xcvbnm`
)

type basicGenerator struct {
	storage  db.Storage
	letters  []string
	active   bool
	stop     chan interface{}
}

func NewBasicGenerator(storage db.Storage) Generator {
	return &basicGenerator{
		storage:  storage,
		letters: []string{
			letters1, letter2, letter3,
		},
	}
}

func (b *basicGenerator) Start() error {
	err := b.firstGen()
	if err != nil {
		return errors.WithStack(err)
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer func() {
		ticker.Stop()
		ticker = nil
		b.active = false
	}()
	for {
		select {
		case <-b.stop:
			return nil
		case <-ticker.C:
			err = b.storage.InactiveRandomOrder()
			if err != nil {
				logrus.Errorf(`%+v`, errors.WithStack(err))
			}

			err = b.oneGen(true)
			if err != nil {
				logrus.Errorf(`%+v`, errors.WithStack(err))
			}

			if !b.active {
				b.active = true
			}
		}
	}
}

func (b *basicGenerator) Stop() {
	if b.active {
		b.stop <- nil
	}
}

func (b *basicGenerator) firstGen() error {
	var err error
	for i := 0; i < maxOrder; i++ {
		if err = b.oneGen(true); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (b *basicGenerator) oneGen(active bool) error {
	rand.Seed(time.Now().UnixNano())
	firstSymbols := b.letters[rand.Intn(3)]
	secondSymbols := b.letters[rand.Intn(3)]
	firstSymbol := string(firstSymbols[rand.Intn(len(firstSymbols))])
	secondSymbol := string(secondSymbols[rand.Intn(len(secondSymbols))])

	err := b.storage.CreateOrder(&models.Order{
		Content: fmt.Sprintf(`%s%s`, firstSymbol, secondSymbol),
		Active:  active,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
