package mongodb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gocloud.dev/docstore"
	"gocloud.dev/docstore/mongodocstore"

	docstorextd "github.com/fensak-io/gostd/gocloudxtd/docstore"
)

const openerKey = "mongodb"

func collectionOpener(ctx context.Context, urlStr string) (*docstore.Collection, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Extract the parameters from the url
	// ASSUMPTION: the url path is always 2 elements, containing db and collection name. Note that the url path always
	// starts with /, so if we use strings.split, the first element should always be "".
	splitPath := strings.Split(parsed.Path, "/")
	if len(splitPath) != 3 {
		return nil, fmt.Errorf("unexpected path in collection url: expected 3 elements, contained %d elements", len(splitPath))
	}
	if splitPath[0] != "" {
		return nil, errors.New("expected URL path to start with /. This is most likely a bug.")
	}
	dbname := splitPath[1]
	collname := splitPath[2]

	qp, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return nil, err
	}

	// Open connection to mongodb server
	serverURL := url.URL{
		Scheme: "mongodb",
		User:   parsed.User,
		Host:   parsed.Host,
	}
	client, err := mongodocstore.Dial(ctx, serverURL.String())
	if err != nil {
		return nil, err
	}

	// Open collection
	mcoll := client.Database(dbname).Collection(collname)
	return mongodocstore.OpenCollection(mcoll, qp.Get("id_field"), nil)
}

func init() {
	docstorextd.RegisterCustomOpener(openerKey, collectionOpener)
}
