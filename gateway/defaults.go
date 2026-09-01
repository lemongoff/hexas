package gateway

import "github.com/lemongoff/hexas/rest"

// DefaultGatewayConf returns Hexas gateway defaults.
func DefaultGatewayConf() GatewayConf { return GatewayConf{RestConf: rest.DefaultRestConf()} }

// DefaultHTTPClientConf returns Hexas gateway HTTP upstream defaults.
func DefaultHTTPClientConf() HttpClientConf { return HttpClientConf{Timeout: 3000} }
