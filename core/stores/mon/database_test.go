package mon

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestDatabaseCollectionConcurrent(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	require.NoError(t, err)
	db := newDatabase("test", client, client.Database("test"), nil)

	const workers = 32
	collections := make([]Collection, workers)
	var group sync.WaitGroup
	for i := range collections {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			collections[index] = db.Collection("players")
		}(i)
	}
	group.Wait()

	for _, collection := range collections[1:] {
		assert.Same(t, collections[0], collection)
	}
}
