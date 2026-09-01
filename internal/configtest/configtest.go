// Package configtest contains test-only helpers for loading typed YAML configuration.
package configtest

import (
	"context"

	config "github.com/lemongoff/hexas-config"
)

// LoadYAML overlays YAML onto the defaults already stored in target.
func LoadYAML[T any](data []byte, target *T) error {
	manager, err := config.NewManager(*target, config.YAMLBytes("test", data))
	if err != nil { return err }
	if err := manager.Load(context.Background()); err != nil { return err }
	snapshot, ok := manager.Current()
	if ok { *target = snapshot.Value() }
	return nil
}
