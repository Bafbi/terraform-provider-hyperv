package api

import (
	"encoding/json"
	"testing"
)

func TestSerializeVmNetworkAdapter(t *testing.T) {
	t.Parallel()

	vmNetworkAdapterJson, err := json.Marshal(VmNetworkAdapter{
		Name: "test",
	})

	if err != nil {
		t.Errorf("Unable to deserialize vmNetworkAdapter: %s", err.Error())
	}

	vmNetworkAdapterJsonString := string(vmNetworkAdapterJson)

	if vmNetworkAdapterJsonString == "" {
		t.Errorf("Unable to deserialize vmNetworkAdapter: %s", err.Error())
	}
}

func TestDeserializeVmNetworkAdapter(t *testing.T) {
	t.Parallel()

	var vmNetworkAdapterJson = `
{
    "Name":  "TestMachine"
}
`

	var vmNetworkAdapter VmNetworkAdapter
	err := json.Unmarshal([]byte(vmNetworkAdapterJson), &vmNetworkAdapter)
	if err != nil {
		t.Errorf("Unable to deserialize vmNetworkAdapter: %s", err.Error())
	}
}
