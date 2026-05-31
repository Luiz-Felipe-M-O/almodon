package domain

import (
	"embed"

	"github.com/alan-b-lima/almodon/internal/support/resource/doc"
	"github.com/alan-b-lima/almodon/ui/web"
)

//go:embed */resource/http.go
var resource embed.FS

func Reference(glob *web.Glob) (*doc.Ref, error) {
	refs := []struct{ Title, Path string }{
		{Title: "Auth Reference", Path: "auth/resource/http.go"},
		{Title: "Items Reference", Path: "item/resource/http.go"},
		{Title: "Material Reference", Path: "material/resource/http.go"},
		{Title: "Promotion Reference", Path: "promotion/resource/http.go"},
		{Title: "User Reference", Path: "user/resource/http.go"},
	}

	docs := make([]*doc.Doc, 0, len(refs))

	for _, ref := range refs {
		file, err := resource.Open(ref.Path)
		if err != nil {
			return nil, err
		}

		d, err := doc.NewDoc(glob, ref.Title, file)
		if err != nil {
			if err == doc.ErrDocNotFound {
				continue
			}

			return nil, err
		}

		docs = append(docs, d)
	}

	ref, err := doc.NewRef(glob, "Almodon Reference", docs)
	if err != nil {
		return nil, err
	}

	return ref, nil
}
