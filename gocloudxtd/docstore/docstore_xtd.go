package docstorextd

import (
	"context"
	"net/url"

	"gocloud.dev/docstore"
)

type customCollectionOpenerFn = func(ctx context.Context, url string) (*docstore.Collection, error)

var (
	registeredCustomOpeners = map[string]customCollectionOpenerFn{}
)

// RegisterCustomOpener will register a new custom opener as being available for use to open a custom docstore URL with
// credentials.
func RegisterCustomOpener(key string, fn customCollectionOpenerFn) {
	registeredCustomOpeners[key] = fn
}

// OpenCollection can be used to open a new docstore collection using one of the custom openers that are loaded. This
// will fall back to the original docstore URL openers if there is no custom opener with the given scheme loaded.
func OpenCollection(ctx context.Context, urlStr string) (*docstore.Collection, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	fn, hasOpener := registeredCustomOpeners[parsed.Scheme]
	if !hasOpener {
		return docstore.OpenCollection(ctx, urlStr)
	}
	return fn(ctx, urlStr)
}
