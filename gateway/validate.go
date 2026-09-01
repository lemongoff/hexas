package gateway

import "fmt"

// Validate checks gateway and upstream configuration.
func (c GatewayConf) Validate() error {
	if err := c.RestConf.Validate(); err != nil { return err }
	for i, upstream := range c.Upstreams {
		if upstream.Grpc == nil && upstream.Http == nil { return fmt.Errorf("upstream %d requires grpc or http target", i) }
		if upstream.Grpc != nil && upstream.Http != nil { return fmt.Errorf("upstream %d cannot define both grpc and http", i) }
		if upstream.Http != nil && upstream.Http.Timeout <= 0 { return fmt.Errorf("upstream %d http timeout must be positive", i) }
	}
	return nil
}
