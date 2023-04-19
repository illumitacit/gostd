package render

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"go.uber.org/zap"
)

// FileRenderer is a renderer that can be used to render files (as opposed to HTML pages) in a embedded FS. Note that
// this only supports individual file templates and does not support componentization.
type FileRenderer struct {
	templateMap map[string]*template.Template
}

type FileRendererOpts struct {
	ViewsFS         embed.FS
	GlobPath        string
	CustomFunctions template.FuncMap
}

// NewFileRenderer will load the file templates from the views FS and load each template as a named template by its
// path.
func NewFileRenderer(logger *zap.Logger, opts FileRendererOpts) (*FileRenderer, error) {
	sugar := logger.Sugar()
	sugar.Debug("Loading file templates")

	fm := sprig.TxtFuncMap()

	// Load custom functions
	for k, v := range opts.CustomFunctions {
		fm[k] = v
	}

	filePaths, err := fs.Glob(opts.ViewsFS, opts.GlobPath)
	if err != nil {
		return nil, err
	}

	tm := map[string]*template.Template{}
	for _, path := range filePaths {
		t, err := loadEmbeddedTextTemplate(opts.ViewsFS, fm, path)
		if err != nil {
			return nil, err
		}
		tm[path] = t
	}

	sugar.Debug("Successfully loaded templates")
	return &FileRenderer{templateMap: tm}, nil
}

// RenderFile will render the file template at the given path using the provided data.
func (r FileRenderer) RenderFile(tplPath string, data interface{}) (string, error) {
	tpl, hasTpl := r.templateMap[tplPath]
	if !hasTpl {
		return "", fmt.Errorf("No template %s in embedded filesystem", tplPath)
	}

	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	return buf.String(), err
}

func loadEmbeddedTextTemplate(
	fs embed.FS,
	fm template.FuncMap,
	path string,
) (*template.Template, error) {
	newT := template.New(path).Funcs(fm)
	data, err := fs.ReadFile(path)
	if err != nil {
		return newT, err
	}
	return newT.Parse(string(data))
}
