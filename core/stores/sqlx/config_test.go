package sqlx

import (
	"testing"

	"github.com/lemongoff/hexas/internal/configtest"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	text := []byte(`DataSource: primary:password@tcp(127.0.0.1:3306)/primary_db
`)

	sc := DefaultSqlConf()
	err := configtest.LoadYAML(text, &sc)
	assert.Nil(t, err)
	assert.Equal(t, "mysql", sc.DriverName)
	assert.Equal(t, policyRoundRobin, sc.Policy)
	assert.Nil(t, sc.Validate())

	sc = SqlConf{}
	assert.Equal(t, errEmptyDatasource, sc.Validate())

	sc.DataSource = "primary:password@tcp(127.0.0.1:3306)/primary_db"
	assert.Equal(t, errEmptyDriverName, sc.Validate())

	sc.DriverName = "mysql"
	assert.Nil(t, sc.Validate())
}
