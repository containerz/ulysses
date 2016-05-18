package framework

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
)

type StateStorage interface {
	SaveTaskState(map[string]string, map[string]string, ApplicationState) error
	SaveApplicationState(string, string, ApplicationState) error
	SaveStackState(string, ApplicationState) error
	GetStackState(string) (map[string]ApplicationState, error)
}

// NoopStateStorage does nothing and is used during bootstrapping when the real state storage is not yet available
type NoopStateStorage struct{}

func (*NoopStateStorage) SaveTaskState(map[string]string, map[string]string, ApplicationState) error {
	return nil
}

func (*NoopStateStorage) SaveApplicationState(string, string, ApplicationState) error {
	return nil
}

func (*NoopStateStorage) SaveStackState(string, ApplicationState) error {
	return nil
}

func (*NoopStateStorage) GetStackState(string) (map[string]ApplicationState, error) {
	return make(map[string]ApplicationState), nil
}

type CassandraStateStorage struct {
	connection *gocql.Session
	keyspace   string
}

func NewCassandraStateStorage(connection *gocql.Session, keyspace string) (StateStorage, error) {
	storage := &CassandraStateStorage{connection: connection, keyspace: keyspace}
	return storage, storage.Init()
}

func (css CassandraStateStorage) GetStackState(id string) (map[string]ApplicationState, error) {
	query := css.prepareQuery("SELECT id, state, props FROM %s.states WHERE parent = ?")
	iter := css.connection.Query(query, id).Iter()
	var (
		appid    string
		appstate ApplicationState
		appprops string
		res      map[string]ApplicationState
	)
	for iter.Scan(&appid, &appstate, &appprops) {
		res[appid] = appstate
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return res, nil
}

func (css CassandraStateStorage) SaveTaskState(task map[string]string, context map[string]string, state ApplicationState) error {
	props, err := json.Marshal(context)
	if err != nil {
		return err
	}
	return css.saveState(task["id"], "task", "", int(state), string(props))
}

func (css CassandraStateStorage) SaveApplicationState(id string, parent string, state ApplicationState) error {
	return css.saveState(id, "application", parent, int(state), "")
}

func (css CassandraStateStorage) SaveStackState(id string, state ApplicationState) error {
	return css.saveState(id, "stack", "", int(state), "")
}

func (css CassandraStateStorage) Init() error {
	query := css.prepareQuery("CREATE TABLE IF NOT EXISTS %s.states (id text, type text, parent text, state int, props text, PRIMARY KEY(id, type, state))")
	return css.connection.Query(query).Exec()
}

func (css CassandraStateStorage) saveState(id string, stateType string, parent string, state int, props string) error {
	query := css.prepareQuery("INSERT INTO %s.states (id, type, parent, state, props) VALUES (?, ?, ?, ?)")
	err := css.connection.Query(query, id, stateType, state, props).Exec()
	return err
}

func (css CassandraStateStorage) prepareQuery(query string) string {
	return fmt.Sprintf(query, css.keyspace)
}
