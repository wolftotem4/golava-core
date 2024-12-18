package session

import (
	"context"

	"github.com/wolftotem4/golava-core/util"
)

type Store struct {
	ID         string
	Handler    SessionHandler
	Attributes map[string]interface{}
}

func NewStore(id string, handler SessionHandler) *Store {
	return &Store{
		ID:         id,
		Handler:    handler,
		Attributes: make(map[string]interface{}),
	}
}

func (s *Store) Start(ctx context.Context) error {
	err := s.loadSession(ctx)
	if err != nil {
		return err
	}

	if !s.Has("_token") {
		s.RegenerateToken()
	}

	return nil
}

func (s *Store) loadSession(ctx context.Context) error {
	payload, err := s.Handler.Read(ctx, s.ID)
	if err != nil {
		return err
	}

	if payload == nil {
		return nil
	}

	return unmarshal(payload, &s.Attributes)
}

func (s *Store) Get(key string) (interface{}, bool) {
	val, ok := s.Attributes[key]
	return val, ok
}

func (s *Store) Put(key string, val interface{}) {
	s.Attributes[key] = val
}

func (s *Store) Has(key string) bool {
	_, ok := s.Attributes[key]
	return ok
}

func (s *Store) Forget(key string) {
	delete(s.Attributes, key)
}

func (s *Store) Flash(key string, val interface{}) {
	s.Put(key, val)

	s.pushToStringSlice("_flash.new", key)

	s.removeFromOldFlashData(key)
}

func (s *Store) pushToStringSlice(key string, val string) {
	if s.Attributes[key] == nil {
		s.Attributes[key] = []string{}
	}
	s.Attributes[key] = append(s.Attributes[key].([]string), val)
}

func (s *Store) removeFromOldFlashData(key string) {
	old, _ := s.Attributes["_flash.old"].([]string)
	for _, oldKey := range old {
		if oldKey == key {
			return
		}
	}

	s.Attributes["_flash.old"] = old
}

func (s *Store) getStringSlice(key string) []string {
	value, _ := s.Attributes[key].([]string)
	return value
}

func (s *Store) AgeFlashData() {
	for _, key := range s.getStringSlice("_flash.old") {
		s.Forget(key)
	}

	s.Attributes["_flash.old"] = s.getStringSlice("_flash.new")
	s.Attributes["_flash.new"] = []string{}
}

func (s *Store) compactForStorage() {
	for _, key := range []string{"_flash.new", "_flash.old"} {
		value, _ := s.Attributes[key].([]string)
		if len(value) == 0 {
			delete(s.Attributes, key)
		}
	}
}

func (s *Store) Save(ctx context.Context) error {
	s.AgeFlashData()

	s.compactForStorage()

	payload, err := marshal(s.Attributes)
	if err != nil {
		return err
	}

	return s.Handler.Write(ctx, s.ID, payload)
}

func (s *Store) FlashInput(value any) error {
	data, err := inputToMap(value)
	if err != nil {
		return err
	}
	s.Flash("_old_input", data)
	return nil
}

func (s *Store) GetOldInput() (map[string]interface{}, bool) {
	value, ok := s.Get("_old_input")
	if !ok {
		return nil, false
	}
	return value.(map[string]interface{}), true
}

func (s *Store) GetOldInputValue(key string) (interface{}, bool) {
	data, ok := s.GetOldInput()
	if !ok {
		return nil, false
	}

	value, ok := data[key]
	return value, ok
}

func (s *Store) Migrate(ctx context.Context, destroy bool) error {
	if destroy {
		err := s.Handler.Destroy(ctx, s.ID)
		if err != nil {
			return err
		}
	}

	s.ID = NewSessionId()

	return nil
}

func (s *Store) Remove(key string) {
	delete(s.Attributes, key)
}

func (s *Store) Flush() {
	s.Attributes = make(map[string]interface{})
}

func (s *Store) Invalidate(ctx context.Context) {
	s.Flush()
	s.Migrate(ctx, true)
}

func (s *Store) Token() string {
	value, _ := s.Get("_token")
	token, _ := value.(string)
	return token
}

func (s *Store) RegenerateToken() {
	s.Put("_token", util.RandomToken(30))
}
