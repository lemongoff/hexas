package rpc

import "testing"

func TestProtocCommandRegistersSkipScaffold(t *testing.T) {
	flag := protocCmd.Flags().Lookup("skip-scaffold")
	if flag == nil {
		t.Fatal("rpc protoc does not register --skip-scaffold")
	}
	if flag.DefValue != "false" {
		t.Fatalf("skip-scaffold default = %q", flag.DefValue)
	}
}
