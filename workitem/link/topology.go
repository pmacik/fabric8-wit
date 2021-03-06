package link

import (
	"database/sql"
	"database/sql/driver"

	"github.com/fabric8-services/fabric8-wit/errors"
	errs "github.com/pkg/errors"
)

// Topology determines the way that links can be created
type Topology string

// String implements the Stringer interface
func (t Topology) String() string { return string(t) }

// Scan implements the https://golang.org/pkg/database/sql/#Scanner interface
// See also https://stackoverflow.com/a/25374979/835098
// See also https://github.com/jinzhu/gorm/issues/302#issuecomment-80566841
func (t *Topology) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*t = Topology(v)
	case string:
		*t = Topology(v)
	case Topology:
		*t = v
	default:
		return errs.Errorf("failed to convert value of type %[1]T to topology: %[1]v", value)
	}
	return t.CheckValid()
}

// Ensure Topology implements the Scanner and Valuer interfaces
var _ sql.Scanner = (*Topology)(nil)
var _ driver.Valuer = (*Topology)(nil)

// Value implements the https://golang.org/pkg/database/sql/driver/#Valuer interface
func (t Topology) Value() (driver.Value, error) { return string(t), nil }

const (
	TopologyNetwork         Topology = "network"
	TopologyDirectedNetwork Topology = "directed_network"
	TopologyDependency      Topology = "dependency"
	TopologyTree            Topology = "tree"
)

// CheckValid returns nil if the given topology is valid; otherwise a
// BadParameterError is returned.
func (t Topology) CheckValid() error {
	switch t {
	case TopologyNetwork, TopologyDirectedNetwork, TopologyDependency, TopologyTree:
		return nil
	default:
		return errors.NewBadParameterError("topolgy", t).Expected(TopologyNetwork + "|" + TopologyDirectedNetwork + "|" + TopologyDependency + "|" + TopologyTree)
	}
}
