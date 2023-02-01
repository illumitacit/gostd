package webstd

// AppContextKey is a value for use with context.WithValue. This should be used to define context keys that are specific
// to the app component.
//
// NOTE: It's used as a pointer so it fits in an interface{} without allocation. This technique for defining context
// keys was copied from Go 1.7's new use of context in net/http.
type AppContextKey struct {
	app  string
	name string
}

func (k *AppContextKey) String() string {
	return k.app + " context value " + k.name
}

func NewAppContextKey(app, name string) *AppContextKey {
	return &AppContextKey{
		app:  app,
		name: name,
	}
}
