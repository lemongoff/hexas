package sqlx

import "errors"

var (
	errEmptyDatasource = errors.New("empty datasource")
	errEmptyDriverName = errors.New("empty driver name")
)

// SqlConf defines the configuration for sqlx.
type SqlConf struct {
	DataSource string
	DriverName string
	Replicas   []string `json:",optional"`
	Policy     string
}

// Validate validates the SqlxConf.
func (sc SqlConf) Validate() error {
	if len(sc.DataSource) == 0 {
		return errEmptyDatasource
	}

	if len(sc.DriverName) == 0 {
		return errEmptyDriverName
	}
	if sc.Policy != "" && sc.Policy != "round-robin" && sc.Policy != "random" {
		return errors.New("sql policy must be round-robin or random")
	}

	return nil
}
