package material

import (
	"time"
	"unicode"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const (
	siadsLength          = 7
	catmatLength         = 6
	ecampusLength        = 20
	descriptionMaxLength = 2048
)

type Material struct {
	uuid        uuid.UUID
	name        string
	siads       string
	catmat      string
	ecampus     string
	description string
	unit        string
	minQuantity float64
	createdAt   time.Time
	updatedAt   time.Time
}

func New(name, siads, catmat, ecampus, description, unit string, minQuantity float64) (Material, error) {
	var m Material
	err := errors.Join(
		m.SetName(name),
		m.SetSIADS(siads),
		m.SetCATMAT(catmat),
		m.SetECAMPUS(ecampus),
		m.SetDescription(description),
		m.SetUnit(unit),
		m.SetMinQuantity(minQuantity),
	)
	if err != nil {
		return Material{}, xerrors.ErrMaterialCreation.New(err)
	}
	m.uuid = uuid.NewUUIDv7()
	m.createdAt = time.Now()
	m.updatedAt = time.Now()
	return m, nil
}

func (m *Material) IsBelowMinimum(currentQuantity float64) bool {
	return currentQuantity < m.minQuantity
}

func (m *Material) UUID() uuid.UUID      { return m.uuid }
func (m *Material) Name() string         { return m.name }
func (m *Material) SIADS() string        { return m.siads }
func (m *Material) CATMAT() string       { return m.catmat }
func (m *Material) ECampus() string      { return m.ecampus }
func (m *Material) Description() string  { return m.description }
func (m *Material) Unit() string         { return m.unit }
func (m *Material) MinQuantity() float64 { return m.minQuantity }
func (m *Material) CreatedAt() time.Time { return m.createdAt }
func (m *Material) UpdatedAt() time.Time { return m.updatedAt }

func (m *Material) SetName(name string) error {
	return entity.SetWithUpdate(&m.name, name, ProcessName, &m.updatedAt)
}

func (m *Material) SetSIADS(siads string) error {
	return entity.SetWithUpdate(&m.siads, siads, ProcessSIADS, &m.updatedAt)
}

func (m *Material) SetCATMAT(catmat string) error {
	return entity.SetWithUpdate(&m.catmat, catmat, ProcessCATMAT, &m.updatedAt)
}

func (m *Material) SetECAMPUS(ecampus string) error {
	return entity.SetWithUpdate(&m.ecampus, ecampus, ProcessECAMPUS, &m.updatedAt)
}

func (m *Material) SetDescription(description string) error {
	return entity.SetWithUpdate(&m.description, description, ProcessDescription, &m.updatedAt)
}

func (m *Material) SetUnit(unit string) error {
	return entity.SetWithUpdate(&m.unit, unit, ProcessUnit, &m.updatedAt)
}

func (m *Material) SetMinQuantity(minQuantity float64) error {
	return entity.SetWithUpdate(&m.minQuantity, minQuantity, ProcessMinQuantity, &m.updatedAt)
}

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", xerrors.ErrNameEmpty
	}
	return name, nil
}

func ProcessSIADS(siads string) (string, error) {
	return processIdNumber(siads, siadsLength)
}

func ProcessCATMAT(catmat string) (string, error) {
	return processIdNumber(catmat, catmatLength)
}

func ProcessECAMPUS(ecampus string) (string, error) {
	return processIdNumber(ecampus, ecampusLength)
}

func ProcessDescription(description string) (string, error) {
	if len(description) > descriptionMaxLength {
		return "", xerrors.ErrTODO
	}
	return description, nil
}

func ProcessUnit(unit string) (string, error) {
	if unit == "" {
		return "", xerrors.ErrTODO
	}
	return unit, nil
}

func ProcessMinQuantity(minQuantity float64) (float64, error) {
	if minQuantity < 0 {
		return 0, xerrors.ErrNegativeMinQuantity
	}
	return minQuantity, nil
}

func processIdNumber(id string, expectedLength int) (string, error) {
	if len(id) > expectedLength {
		return "", xerrors.ErrInvalidIdLength
	}

	if !isAllDigits(id) {
		return "", xerrors.ErrIdContainsNonDigits
	}
	return id, nil
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
