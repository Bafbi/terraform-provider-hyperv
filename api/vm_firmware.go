package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConsoleModeType int

const (
	ConsoleModeType_Default ConsoleModeType = 0
	ConsoleModeType_Com1    ConsoleModeType = 1
	ConsoleModeType_Com2    ConsoleModeType = 2
	ConsoleModeType_None    ConsoleModeType = 3
)

var ConsoleModeType_name = map[ConsoleModeType]string{
	ConsoleModeType_Default: "Default",
	ConsoleModeType_Com1:    "COM1",
	ConsoleModeType_Com2:    "COM2",
	ConsoleModeType_None:    "None",
}

var ConsoleModeType_value = map[string]ConsoleModeType{
	"default": ConsoleModeType_Default,
	"com1":    ConsoleModeType_Com1,
	"com2":    ConsoleModeType_Com2,
	"none":    ConsoleModeType_None,
}

func (x ConsoleModeType) String() string {
	return ConsoleModeType_name[x]
}

func ToConsoleModeType(x string) ConsoleModeType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return ConsoleModeType(integerValue)
	}
	return ConsoleModeType_value[strings.ToLower(x)]
}

func (d *ConsoleModeType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *ConsoleModeType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = ConsoleModeType(i)
			return nil
		}

		return err
	}
	*d = ToConsoleModeType(s)
	return nil
}

type IPProtocolPreference int

const (
	IPProtocolPreference_IPv4 IPProtocolPreference = 0
	IPProtocolPreference_IPv6 IPProtocolPreference = 1
)

var IPProtocolPreference_name = map[IPProtocolPreference]string{
	IPProtocolPreference_IPv4: "IPv4",
	IPProtocolPreference_IPv6: "IPv6",
}

var IPProtocolPreference_value = map[string]IPProtocolPreference{
	"ipv4": IPProtocolPreference_IPv4,
	"ipv6": IPProtocolPreference_IPv6,
}

func (x IPProtocolPreference) String() string {
	return IPProtocolPreference_name[x]
}

func ToIPProtocolPreference(x string) IPProtocolPreference {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return IPProtocolPreference(integerValue)
	}
	return IPProtocolPreference_value[strings.ToLower(x)]
}

func (d *IPProtocolPreference) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *IPProtocolPreference) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = IPProtocolPreference(i)
			return nil
		}

		return err
	}
	*d = ToIPProtocolPreference(s)
	return nil
}

type Gen2BootType int

const (
	Gen2BootType_HardDiskDrive  Gen2BootType = 1
	Gen2BootType_DvdDrive       Gen2BootType = 2
	Gen2BootType_NetworkAdapter Gen2BootType = 3
)

var Gen2BootType_name = map[Gen2BootType]string{
	Gen2BootType_HardDiskDrive:  "HardDiskDrive",
	Gen2BootType_DvdDrive:       "DvdDrive",
	Gen2BootType_NetworkAdapter: "NetworkAdapter",
}

var Gen2BootType_value = map[string]Gen2BootType{
	"harddiskdrive":  Gen2BootType_HardDiskDrive,
	"dvddrive":       Gen2BootType_DvdDrive,
	"networkadapter": Gen2BootType_NetworkAdapter,
}

func (x Gen2BootType) String() string {
	return Gen2BootType_name[x]
}

func ToGen2BootType(x string) Gen2BootType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return Gen2BootType(integerValue)
	}
	return Gen2BootType_value[strings.ToLower(x)]
}

func (d *Gen2BootType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *Gen2BootType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = Gen2BootType(i)
			return nil
		}

		return err
	}
	*d = ToGen2BootType(s)
	return nil
}

type Gen2BootOrder struct {
	Type Gen2BootType

	NetworkAdapterName string
	SwitchName         string
	MacAddress         string

	Path               string
	ControllerNumber   int
	ControllerLocation int
}

type VmFirmware struct {
	VmName                       string
	BootOrders                   []Gen2BootOrder
	EnableSecureBoot             OnOffState
	SecureBootTemplate           string
	PreferredNetworkBootProtocol IPProtocolPreference
	ConsoleMode                  ConsoleModeType
	PauseAfterBootFailure        OnOffState
}

func DefaultVmFirmwares() (interface{}, error) {
	result := make([]VmFirmware, 0)
	vmFirmware := VmFirmware{
		BootOrders:                   []Gen2BootOrder{},
		EnableSecureBoot:             OnOffState_On,
		SecureBootTemplate:           "MicrosoftWindows",
		PreferredNetworkBootProtocol: IPProtocolPreference_IPv4,
		ConsoleMode:                  ConsoleModeType_Default,
		PauseAfterBootFailure:        OnOffState_Off,
	}

	result = append(result, vmFirmware)
	return result, nil
}

func ExpandGen2BootOrder(bootOrders []interface{}) ([]Gen2BootOrder, error) {
	gen2bootOrders := make([]Gen2BootOrder, 0)
	for _, bootOrder := range bootOrders {
		bootOrder, ok := bootOrder.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] boot_order should be a Hash - was '%+v'", bootOrder)
		}

		log.Printf("[DEBUG] bootOrder = [%+v]", bootOrder)

		bootType, ok := bootOrder["boot_type"].(string)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] boot_type should be a string")
		}
		networkAdapterName, ok := bootOrder["network_adapter_name"].(string)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] network_adapter_name should be a string")
		}
		switchName, ok := bootOrder["switch_name"].(string)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] switch_name should be a string")
		}
		macAddress, ok := bootOrder["mac_address"].(string)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] mac_address should be a string")
		}
		path, ok := bootOrder["path"].(string)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] path should be a string")
		}
		controllerNumber, ok := bootOrder["controller_number"].(int)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] controller_number should be an int")
		}
		controllerLocation, ok := bootOrder["controller_location"].(int)
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] controller_location should be an int")
		}

		expandedGen2BootOrder := Gen2BootOrder{
			Type: ToGen2BootType(bootType),

			NetworkAdapterName: networkAdapterName,
			SwitchName:         switchName,
			MacAddress:         macAddress,

			Path:               path,
			ControllerNumber:   controllerNumber,
			ControllerLocation: controllerLocation,
		}

		gen2bootOrders = append(gen2bootOrders, expandedGen2BootOrder)
	}

	return gen2bootOrders, nil
}

func ExpandVmFirmwares(d *schema.ResourceData) ([]VmFirmware, error) {
	expandedVmFirmwares := make([]VmFirmware, 0)

	if v, ok := d.GetOk("vm_firmware"); ok {
		vmFirmwares, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] vm_firmware should be a list - was '%+v'", v)
		}
		for _, firmware := range vmFirmwares {
			firmware, ok := firmware.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vm_firmware should be a Hash - was '%+v'", firmware)
			}

			log.Printf("[DEBUG] firmware = [%+v]", firmware)

			bootOrderInterface, ok := firmware["boot_order"].([]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] boot_order should be a list")
			}
			bootOrders, err := ExpandGen2BootOrder(bootOrderInterface)
			if err != nil {
				return nil, err
			}

			enableSecureBoot, ok := firmware["enable_secure_boot"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] enable_secure_boot should be a string")
			}
			secureBootTemplate, ok := firmware["secure_boot_template"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] secure_boot_template should be a string")
			}
			preferredNetworkBootProtocol, ok := firmware["preferred_network_boot_protocol"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] preferred_network_boot_protocol should be a string")
			}
			consoleMode, ok := firmware["console_mode"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] console_mode should be a string")
			}
			pauseAfterBootFailure, ok := firmware["pause_after_boot_failure"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] pause_after_boot_failure should be a string")
			}

			expandedVmFirmware := VmFirmware{
				BootOrders:                   bootOrders,
				EnableSecureBoot:             ToOnOffState(enableSecureBoot),
				SecureBootTemplate:           secureBootTemplate,
				PreferredNetworkBootProtocol: ToIPProtocolPreference(preferredNetworkBootProtocol),
				ConsoleMode:                  ToConsoleModeType(consoleMode),
				PauseAfterBootFailure:        ToOnOffState(pauseAfterBootFailure),
			}

			expandedVmFirmwares = append(expandedVmFirmwares, expandedVmFirmware)
		}
	}

	if len(expandedVmFirmwares) < 1 {
		vmFirmware := VmFirmware{
			BootOrders:                   []Gen2BootOrder{},
			EnableSecureBoot:             OnOffState_On,
			SecureBootTemplate:           "MicrosoftWindows",
			PreferredNetworkBootProtocol: IPProtocolPreference_IPv4,
			ConsoleMode:                  ConsoleModeType_Default,
			PauseAfterBootFailure:        OnOffState_Off,
		}
		expandedVmFirmwares = append(expandedVmFirmwares, vmFirmware)
	}

	return expandedVmFirmwares, nil
}

func FlattenGen2BootOrders(bootOrders []Gen2BootOrder) []interface{} {
	if bootOrders == nil || len(bootOrders) < 1 {
		return nil
	}

	flattenedGen2BootOrders := make([]interface{}, 0)

	for _, bootOrder := range bootOrders {
		flattenedGen2BootOrder := make(map[string]interface{})
		flattenedGen2BootOrder["boot_type"] = bootOrder.Type.String()

		flattenedGen2BootOrder["network_adapter_name"] = bootOrder.NetworkAdapterName
		flattenedGen2BootOrder["switch_name"] = bootOrder.SwitchName
		flattenedGen2BootOrder["mac_address"] = bootOrder.MacAddress

		flattenedGen2BootOrder["path"] = bootOrder.Path
		flattenedGen2BootOrder["controller_number"] = bootOrder.ControllerNumber
		flattenedGen2BootOrder["controller_location"] = bootOrder.ControllerLocation
		flattenedGen2BootOrders = append(flattenedGen2BootOrders, flattenedGen2BootOrder)
	}

	return flattenedGen2BootOrders
}

func FlattenVmFirmwares(vmFirmwares *[]VmFirmware) []interface{} {
	if vmFirmwares == nil || len(*vmFirmwares) < 1 {
		return nil
	}

	flattenedVmFirmwares := make([]interface{}, 0)

	for _, vmFirmware := range *vmFirmwares {
		flattenedGen2BootOrder := FlattenGen2BootOrders(vmFirmware.BootOrders)
		flattenedVmFirmware := make(map[string]interface{})
		flattenedVmFirmware["boot_order"] = flattenedGen2BootOrder
		flattenedVmFirmware["enable_secure_boot"] = vmFirmware.EnableSecureBoot.String()
		flattenedVmFirmware["secure_boot_template"] = vmFirmware.SecureBootTemplate
		flattenedVmFirmware["preferred_network_boot_protocol"] = vmFirmware.PreferredNetworkBootProtocol.String()
		flattenedVmFirmware["console_mode"] = vmFirmware.ConsoleMode.String()
		flattenedVmFirmware["pause_after_boot_failure"] = vmFirmware.PauseAfterBootFailure.String()
		flattenedVmFirmwares = append(flattenedVmFirmwares, flattenedVmFirmware)
	}

	return flattenedVmFirmwares
}

type HypervVmFirmwareClient interface {
	CreateOrUpdateVmFirmware(
		ctx context.Context,
		vmName string,
		bootOrders []Gen2BootOrder,
		enableSecureBoot OnOffState,
		secureBootTemplate string,
		preferredNetworkBootProtocol IPProtocolPreference,
		consoleMode ConsoleModeType,
		pauseAfterBootFailure OnOffState,
	) (err error)
	GetVmFirmware(ctx context.Context, vmName string) (result VmFirmware, err error)
	GetNoVmFirmwares(ctx context.Context) (result []VmFirmware)
	GetVmFirmwares(ctx context.Context, vmName string) (result []VmFirmware, err error)
	CreateOrUpdateVmFirmwares(ctx context.Context, vmName string, vmFirmwares []VmFirmware) (err error)
}
