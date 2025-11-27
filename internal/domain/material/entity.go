package material

import (
	"strings"

	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/errors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const (
	siadsLength          = 7
	catmatLength         = 6
	ecampusLength        = 20
	descriptionMaxLength = 2048
)

type unitDefinition struct {
	canonical     string
	abbreviations []string
}

var predefinedUnits = map[string]unitDefinition{
	"unidade": {
		canonical:     "unidade",
		abbreviations: []string{"und", "unid", "un", "u"},
	},
	"caixa": {
		canonical:     "caixa",
		abbreviations: []string{"cx"},
	},
	"pacote": {
		canonical:     "pacote",
		abbreviations: []string{"pct"},
	},
	"frasco": {
		canonical:     "frasco",
		abbreviations: []string{"fr"},
	},
	"mililitro": {
		canonical:     "mililitro",
		abbreviations: []string{"ml"},
	},
	"grama": {
		canonical:     "grama",
		abbreviations: []string{"g", "gr"},
	},
	"par": {
		canonical:     "par",
		abbreviations: []string{},
	},
	"litro": {
		canonical:     "litro",
		abbreviations: []string{"l", "lt"},
	},
	"resma": {
		canonical:     "resma",
		abbreviations: []string{},
	},
	"rolo": {
		canonical:     "rolo",
		abbreviations: []string{},
	},
	"galão": {
		canonical:     "galão",
		abbreviations: []string{"gal"},
	},
	"cartela": {
		canonical:     "cartela",
		abbreviations: []string{"crt"},
	},
	"tira": {
		canonical:     "tira",
		abbreviations: []string{},
	},
}

var customUnits = make(map[string]unitDefinition)

type Material struct {
	uuid        uuid.UUID
	name        string
	siads       string
	catmat      string
	ecampus     string
	description string
	unit        string
	minQuantity float64
}

func New(name, siads, catmat, ecampus, description, unit string, minQuantity float64) (Material, error) {
	var m Material

	err := errors.Join(
		m.SetName(name),
		m.SetSIADS(siads),
		m.SetCATMAT(catmat),
		m.SetECampus(ecampus),
		m.SetDescription(description),
		m.SetUnit(unit),
		m.SetMinQuantity(minQuantity),
	)
	if err != nil {
		return Material{}, xerrors.ErrMaterialCreation.New(err)
	}

	m.uuid = uuid.NewUUIDv7()
	return m, nil
}

func (m *Material) IsBelowMinimum(quantity float64) bool {
	return quantity < m.minQuantity
}

func (m *Material) UUID() uuid.UUID      { return m.uuid }
func (m *Material) Name() string         { return m.name }
func (m *Material) SIADS() string        { return m.siads }
func (m *Material) CATMAT() string       { return m.catmat }
func (m *Material) ECampus() string      { return m.ecampus }
func (m *Material) Description() string  { return m.description }
func (m *Material) Unit() string         { return m.unit }
func (m *Material) MinQuantity() float64 { return m.minQuantity }

func (m *Material) SetName(name string) error {
	return entity.Set(&m.name, name, ProcessName)
}

func (m *Material) SetSIADS(siads string) error {
	return entity.Set(&m.siads, siads, ProcessSIADS)
}

func (m *Material) SetCATMAT(catmat string) error {
	return entity.Set(&m.catmat, catmat, ProcessCATMAT)
}

func (m *Material) SetECampus(ecampus string) error {
	return entity.Set(&m.ecampus, ecampus, ProcessECampus)
}

func (m *Material) SetDescription(description string) error {
	return entity.Set(&m.description, description, ProcessDescription)
}

func (m *Material) SetUnit(unit string) error {
	return entity.Set(&m.unit, unit, ProcessUnit)
}

func (m *Material) SetMinQuantity(minQuantity float64) error {
	return entity.Set(&m.minQuantity, minQuantity, ProcessMinQuantity)
}

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", xerrors.ErrNameEmpty
	}
	return name, nil
}

func ProcessSIADS(siads string) (string, error) {
	id, err := processIdNumber(siads, siadsLength, siadsLength)
	if err != nil {
		return "", xerrors.ErrSIADSInvalid.New(err)
	}
	return id, nil
}

func ProcessCATMAT(catmat string) (string, error) {
	id, err := processIdNumber(catmat, catmatLength, catmatLength)
	if err != nil {
		return "", xerrors.ErrCATMATInvalid.New(err)
	}
	return id, nil
}

func ProcessECampus(ecampus string) (string, error) {
	id, err := processIdNumber(ecampus, 1, ecampusLength)
	if err != nil {
		return "", xerrors.ErrECampusInvalid.New(err)
	}
	return id, nil
}

func ProcessDescription(description string) (string, error) {
	if len(description) > descriptionMaxLength {
		return "", xerrors.ErrDescriptionTooLong
	}
	return description, nil
}

func ProcessUnit(unit string) (string, error) {
	if unit == "" {
		return "", xerrors.ErrUnitEmpty
	}

	normalized, err := normalizeUnit(unit)
	if err != nil {
		// If unit not found, add it as a custom unit automatically
		if err == xerrors.ErrUnitNotFound {
			cleanUnit := strings.ToLower(strings.TrimSpace(unit))
			customUnits[cleanUnit] = unitDefinition{
				canonical:     cleanUnit,
				abbreviations: []string{},
			}
			return cleanUnit, nil
		}
		return "", err
	}

	return normalized, nil
}

func ProcessMinQuantity(minQuantity float64) (float64, error) {
	if minQuantity < 0 {
		return 0, xerrors.ErrMinQuantityNegative
	}
	return minQuantity, nil
}

func normalizeUnit(input string) (string, error) {
	if input == "" {
		return "", xerrors.ErrUnitEmpty
	}

	normalized := strings.ToLower(strings.TrimSpace(input))

	if unit, exists := predefinedUnits[normalized]; exists {
		return unit.canonical, nil
	}

	if unit, exists := customUnits[normalized]; exists {
		return unit.canonical, nil
	}

	for _, unit := range predefinedUnits {
		for _, abbr := range unit.abbreviations {
			if normalized == abbr {
				return unit.canonical, nil
			}
		}
	}

	for _, unit := range customUnits {
		for _, abbr := range unit.abbreviations {
			if normalized == abbr {
				return unit.canonical, nil
			}
		}
	}

	return "", xerrors.ErrUnitNotFound
}

func processIdNumber(id string, min, max int) (string, error) {
	if min == max && len(id) != min {
		return "", xerrors.ErrInvalidIdLength.New(min)
	}

	if min != max && (len(id) < min || max < len(id)) {
		return "", xerrors.ErrInvalidIdLengthRange.New(min, max)
	}

	if !numerical(id) {
		return "", xerrors.ErrIdContainsNonDigits
	}

	return id, nil
}

func numerical(s string) bool {
	for _, r := range s {
		if r < '0' || '9' < r {
			return false
		}
	}
	return true
}
