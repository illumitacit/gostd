package chistd

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"gitea.com/go-chi/session"
	"gocloud.dev/docstore"
)

// sessionDocument represents a stored session in the document store.
type sessionDocument struct {
	ID                   string
	SerializedAttributes string
	TTL                  time.Duration
	DocstoreRevision     any
}

// SessionDocumentPKeyField is the primary key of the session document entry. This should be used when opening the
// docstore collection.
const SessionDocumentPKeyField = "ID"

// since we do not use context define global once
var ctx = context.TODO()

// SessionDocumentStore represents a document store based session store implementation.
type SessionDocumentStore struct {
	coll        *docstore.Collection
	prefix, sid string
	duration    time.Duration
	data        *sync.Map
}

// NewSessionDocumentStore creates and returns a docstore based session store.
func NewSessionDocumentStore(
	coll *docstore.Collection,
	prefix, sid string,
	dur time.Duration,
	kv *sync.Map,
) *SessionDocumentStore {
	return &SessionDocumentStore{
		coll:   coll,
		prefix: prefix,
		sid:    sid,
		data:   kv,
	}
}

// Set sets value to given key in session.
func (s *SessionDocumentStore) Set(key, val any) error {
	s.data.Store(key, val)
	return nil
}

// Get gets value by given key in session.
func (s *SessionDocumentStore) Get(key any) any {
	val, _ := s.data.Load(key)
	return val
}

// Delete delete a key from session.
func (s *SessionDocumentStore) Delete(key any) error {
	s.data.LoadAndDelete(key)
	return nil
}

// ID returns current session ID.
func (s *SessionDocumentStore) ID() string {
	return s.sid
}

// Release releases resource and save data to provider.
func (s *SessionDocumentStore) Release() error {
	// Create a snapshot of the data so that it can be encoded into persistent storage.
	snapshot := map[any]any{}
	s.data.Range(func(key, val any) bool {
		snapshot[key] = val
		return true
	})

	// Skip encoding if the data is empty
	if len(snapshot) == 0 {
		return nil
	}

	data, err := session.EncodeGob(snapshot)
	if err != nil {
		return err
	}

	sess := &sessionDocument{
		ID:                   s.prefix + s.sid,
		SerializedAttributes: string(data),
		TTL:                  s.duration,
	}
	return s.coll.Put(ctx, sess)
}

// Flush deletes all session data.
func (s *SessionDocumentStore) Flush() error {
	s.data = &sync.Map{}
	return nil
}

// DocStoreProvider represents a redis session provider implementation.
type DocStoreProvider struct {
	coll     *docstore.Collection
	duration time.Duration
	prefix   string
}

// Init initializes docstore session provider. Config should be the doc store collection connection URL, as supported by
// gocloud.dev. Refer to the docs for more information and supported providers:
// https://gocloud.dev/howto/docstore/
//
// NOTE: You do not need to include the primary key field in the collection URL: it will be automatically appended in
// this function.
// NOTE: The collection schema should support storage of the sessionDocument.
func (p *DocStoreProvider) Init(maxlifetime int64, configs string) (err error) {
	p.duration, err = time.ParseDuration(fmt.Sprintf("%ds", maxlifetime))
	if err != nil {
		return err
	}

	parsed, err := url.Parse(configs)
	if err != nil {
		return err
	}
	qp, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return err
	}

	switch parsed.Scheme {
	case "firestore":
		qp.Set("name_field", SessionDocumentPKeyField)
	case "dynamodb":
		qp.Set("partition_key", SessionDocumentPKeyField)
	case "mongo":
		qp.Set("id_field", SessionDocumentPKeyField)
	case "mem":
		parsed.Path = "/" + SessionDocumentPKeyField
	}
	parsed.RawQuery = qp.Encode()

	coll, err := docstore.OpenCollection(ctx, parsed.String())
	if err != nil {
		return err
	}

	p.coll = coll
	return nil
}

// Read returns raw session store by session ID.
func (p *DocStoreProvider) Read(sid string) (session.RawStore, error) {
	psid := p.prefix + sid

	var attrs sync.Map
	if p.Exist(sid) {
		sess := &sessionDocument{ID: psid}
		if err := p.coll.Get(ctx, sess); err != nil {
			return nil, err
		}
		satts := sess.SerializedAttributes
		datts, err := session.DecodeGob([]byte(satts))
		if err != nil {
			return nil, err
		}

		for k, v := range datts {
			attrs.Store(k, v)
		}
	}

	return NewSessionDocumentStore(p.coll, p.prefix, sid, p.duration, &attrs), nil
}

// Exist returns true if session with given ID exists.
func (p *DocStoreProvider) Exist(sid string) bool {
	sess := &sessionDocument{ID: p.prefix + sid}
	err := p.coll.Get(ctx, sess, SessionDocumentPKeyField)
	return err == nil
}

// Destroy deletes a session by session ID.
func (p *DocStoreProvider) Destroy(sid string) error {
	sess := &sessionDocument{ID: p.prefix + sid}
	return p.coll.Delete(ctx, sess)
}

// Regenerate regenerates a stored session from old session ID to new one.
func (p *DocStoreProvider) Regenerate(oldsid, sid string) (_ session.RawStore, err error) {
	poldsid := p.prefix + oldsid
	psid := p.prefix + sid

	if p.Exist(sid) {
		return nil, fmt.Errorf("new sid '%s' already exists", sid)
	}

	if p.Exist(oldsid) {
		// TODO: wrap in a transaction so that this all happens atomically
		sess := &sessionDocument{ID: poldsid}
		if err := p.coll.Get(ctx, sess); err != nil {
			return nil, err
		}

		// Update the ID and call Create to replicate it under the new ID.
		sess.ID = psid
		sess.DocstoreRevision = nil
		if err := p.coll.Create(ctx, sess); err != nil {
			return nil, err
		}

		// Finally, delete the old stored session.
		oldSess := &sessionDocument{ID: poldsid}
		if err := p.coll.Delete(ctx, oldSess); err != nil {
			return nil, err
		}
	}

	// Read out the session using the new ID
	return p.Read(sid)
}

// Count counts and returns number of sessions.
// NOTE: because of how document stores work, this operation is expensive since it has to scan every document. This
// should be used sparingly!
func (p *DocStoreProvider) Count() int {
	fullScan := p.coll.Query().Get(ctx)
	defer fullScan.Stop()

	count := 0
	for {
		var s sessionDocument
		err := fullScan.Next(ctx, &s)
		if err == io.EOF {
			break
		} else if err != nil {
			return -1
		}
		count++
	}
	return count
}

// GC calls GC to clean expired sessions.
func (_ *DocStoreProvider) GC() {}

func init() {
	session.Register("docstore", &DocStoreProvider{})
}
