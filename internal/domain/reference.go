package domain

import (
	"bytes"
	"embed"
	"maps"
	"slices"

	"github.com/alan-b-lima/almodon/internal/support/resource/doc"
)

//go:embed */resource/http.go
var resource embed.FS

func Reference() (*doc.Ref, error) {
	refs := map[string]string{
		"Auth Reference":      "auth/resource/http.go",
		"Items Reference":     "item/resource/http.go",
		"Material Reference":  "material/resource/http.go",
		"Promotion Reference": "promotion/resource/http.go",
		"User Reference":      "user/resource/http.go",
	}

	titles := slices.Sorted(maps.Keys(refs))
	docs := make([]*doc.Doc, 0, len(titles))

	for _, title := range titles {
		file, err := resource.ReadFile(refs[title])
		if err != nil {
			return nil, err
		}

		doc, err := doc.New(title, bytes.NewBuffer(file))
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	ref, err := doc.NewRef("Almodon Reference", docs)
	if err != nil {
		return nil, err
	}

	return ref, nil
}
