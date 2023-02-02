package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/Masterminds/sprig/v3"
	"go.uber.org/zap"
)

type Renderer struct {
	templateMap map[string]*template.Template
}

type RendererOpts struct {
	ViewsFS           embed.FS
	LayoutGlobPath    string
	ComponentGlobPath string
	PageGlobPath      string
	CustomFunctions   template.FuncMap
}

// NewRenderer will load the templates from the views FS and load each page template as a named template by its path.
// Note that every template will include all the templates from the layouts and components folder so that cross
// references work, but NOT individual pages. This means that a page template can reference components and layouts with
// a template directive, but will not be able to load templates from another page.
func NewRenderer(logger *zap.Logger, opts RendererOpts) (*Renderer, error) {
	sugar := logger.Sugar()
	sugar.Debug("Loading templates")

	fm := sprig.FuncMap()

	// Load custom functions
	for k, v := range opts.CustomFunctions {
		fm[k] = v
	}

	layoutPaths, err := fs.Glob(opts.ViewsFS, opts.LayoutGlobPath)
	if err != nil {
		return nil, err
	}

	componentPaths, err := fs.Glob(opts.ViewsFS, opts.ComponentGlobPath)
	if err != nil {
		return nil, err
	}

	includedPaths := append(layoutPaths, componentPaths...)

	pagePaths, err := fs.Glob(opts.ViewsFS, opts.PageGlobPath)
	if err != nil {
		return nil, err
	}

	tm := map[string]*template.Template{}
	for _, path := range append(componentPaths, pagePaths...) {
		t, err := loadAllEmbeddedTemplates(opts.ViewsFS, fm, path, includedPaths)
		if err != nil {
			return nil, err
		}
		tm[path] = t
	}

	sugar.Debug("Successfully loaded templates")
	return &Renderer{templateMap: tm}, nil
}

// RenderHTML will render the html template at the given path using the provided data.
func (r Renderer) RenderHTML(tplPath string, data interface{}) (string, error) {
	tpl, hasTpl := r.templateMap[tplPath]
	if !hasTpl {
		return "", fmt.Errorf("No template %s in embedded filesystem", tplPath)
	}

	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	return buf.String(), err
}

func loadAllEmbeddedTemplates(
	fs embed.FS,
	fm template.FuncMap,
	rootPath string,
	paths []string,
) (*template.Template, error) {
	root, err := loadEmbeddedTemplate(fs, fm, nil, rootPath)
	if err != nil {
		return root, err
	}

	for _, path := range paths {
		// Skip the root path. The root path will be included in the paths slice when loading in each component as an
		// individual template.
		if path == rootPath {
			continue
		}

		if _, err := loadEmbeddedTemplate(fs, fm, root, path); err != nil {
			return root, err
		}
	}
	return root, nil
}

func loadEmbeddedTemplate(
	fs embed.FS,
	fm template.FuncMap,
	maybeRoot *template.Template,
	path string,
) (*template.Template, error) {
	var newT *template.Template
	if maybeRoot == nil {
		newT = template.New(path).Funcs(fm)
	} else {
		newT = maybeRoot.New(path).Funcs(fm)
	}

	data, err := fs.ReadFile(path)
	if err != nil {
		return newT, err
	}

	return newT.Parse(string(data))
}
