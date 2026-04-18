package material

const (
	NameMaxLen = 128

	DescriptionMaxLen = 4096
)

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", ErrNameEmpty
	}

	if len(name) >= NameMaxLen {
		return "", ErrNameEmpty
	}

	return name, nil
}

func ProcessECampus(ecampus int) (int, error) {
	return ecampus, nil
}

func ProcessCATMAT(catmat int) (int, error) {
	return catmat, nil
}

func ProcessSIADS(siads int) (int, error) {
	return siads, nil
}

func ProcessDescription(description string) (string, error) {
	if len(description) >= DescriptionMaxLen {
		return "", ErrDescriptionTooLong
	}

	return description, nil
}

func ProcessUnit(unit string) (string, error) {
	return unit, nil
}

func ProcessMin(quantity float64) (float64, error) {
	if quantity < 0 {
		return 0, ErrMinNegative
	}

	return quantity, nil
}
