package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PortMirroring int

const (
	PortMirroring_None        PortMirroring = 0
	PortMirroring_Destination PortMirroring = 1
	PortMirroring_Source      PortMirroring = 2
)

var PortMirroring_name = map[PortMirroring]string{
	PortMirroring_None:        "None",
	PortMirroring_Destination: "Destination",
	PortMirroring_Source:      "Source",
}

var PortMirroring_value = map[string]PortMirroring{
	"none":        PortMirroring_None,
	"destination": PortMirroring_Destination,
	"source":      PortMirroring_Source,
}

func (x PortMirroring) String() string {
	return PortMirroring_name[x]
}

func ToPortMirroring(x string) PortMirroring {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return PortMirroring(integerValue)
	}
	return PortMirroring_value[strings.ToLower(x)]
}

func (d *PortMirroring) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *PortMirroring) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = PortMirroring(i)
			return nil
		}

		return err
	}
	*d = ToPortMirroring(s)
	return nil
}

type IovInterruptModerationValue int

const (
	IovInterruptModerationValue_Default  IovInterruptModerationValue = 0
	IovInterruptModerationValue_Adaptive IovInterruptModerationValue = 1
	IovInterruptModerationValue_Off      IovInterruptModerationValue = 2
	IovInterruptModerationValue_Low      IovInterruptModerationValue = 100
	IovInterruptModerationValue_Medium   IovInterruptModerationValue = 200
	IovInterruptModerationValue_High     IovInterruptModerationValue = 300
)

var IovInterruptModerationValue_name = map[IovInterruptModerationValue]string{
	IovInterruptModerationValue_Default:  "Default",
	IovInterruptModerationValue_Adaptive: "Adaptive",
	IovInterruptModerationValue_Off:      "Off",
	IovInterruptModerationValue_Low:      "Low",
	IovInterruptModerationValue_Medium:   "Medium",
	IovInterruptModerationValue_High:     "High",
}

var IovInterruptModerationValue_value = map[string]IovInterruptModerationValue{
	"default":  IovInterruptModerationValue_Default,
	"adaptive": IovInterruptModerationValue_Adaptive,
	"off":      IovInterruptModerationValue_Off,
	"low":      IovInterruptModerationValue_Low,
	"medium":   IovInterruptModerationValue_Medium,
	"high":     IovInterruptModerationValue_High,
}

func (x IovInterruptModerationValue) String() string {
	return IovInterruptModerationValue_name[x]
}

func ToIovInterruptModerationValue(x string) IovInterruptModerationValue {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return IovInterruptModerationValue(integerValue)
	}
	return IovInterruptModerationValue_value[strings.ToLower(x)]
}

func (d *IovInterruptModerationValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *IovInterruptModerationValue) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = IovInterruptModerationValue(i)
			return nil
		}

		return err
	}
	*d = ToIovInterruptModerationValue(s)
	return nil
}

func DiffSuppressVmStaticMacAddress(key, old, new string, d *schema.ResourceData) bool {
	// Static Mac Address has not been set, so we don't mind what ever value is automatically generated
	if new == "" {
		return true
	}

	return new == old
}

func ExpandNetworkAdapters(d *schema.ResourceData) ([]VmNetworkAdapter, error) {
	expandedNetworkAdapters := make([]VmNetworkAdapter, 0)

	if v, ok := d.GetOk("network_adaptors"); ok {
		networkAdapters, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a list - was '%+v'", v)
		}

		for _, networkAdapter := range networkAdapters {
			networkAdapter, ok := networkAdapter.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a Hash - was '%+v'", networkAdapter)
			}

			mandatoryFeatureIdVal := networkAdapter["mandatory_feature_id"]
			mandatoryFeatureIdSet, ok := mandatoryFeatureIdVal.(*schema.Set)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] mandatory_feature_id should be a Set - was '%+v'", mandatoryFeatureIdVal)
			}
			mandatoryFeatureIds := make([]string, 0)
			for _, mandatoryFeatureId := range mandatoryFeatureIdSet.List() {
				featureId, ok := mandatoryFeatureId.(string)
				if !ok {
					return nil, fmt.Errorf("[ERROR][hyperv] mandatory_feature_id item should be a string - was '%+v'", mandatoryFeatureId)
				}
				mandatoryFeatureIds = append(mandatoryFeatureIds, featureId)
			}

			ipAddressesVal := networkAdapter["ip_addresses"]
			ipAddressesSet, ok := ipAddressesVal.([]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] ip_addresses should be a list - was '%+v'", ipAddressesVal)
			}
			ipAddresses := make([]string, 0)
			for _, ipAddress := range ipAddressesSet {
				ipAddr, ok := ipAddress.(string)
				if !ok {
					return nil, fmt.Errorf("[ERROR][hyperv] ip_address item should be a string - was '%+v'", ipAddress)
				}
				ipAddresses = append(ipAddresses, ipAddr)
			}

			name, ok := networkAdapter["name"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] name should be a string")
			}
			switchName, ok := networkAdapter["switch_name"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] switch_name should be a string")
			}
			managementOs, ok := networkAdapter["management_os"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] management_os should be a bool")
			}
			isLegacy, ok := networkAdapter["is_legacy"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] is_legacy should be a bool")
			}
			dynamicMacAddress, ok := networkAdapter["dynamic_mac_address"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] dynamic_mac_address should be a bool")
			}
			staticMacAddress, ok := networkAdapter["static_mac_address"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] static_mac_address should be a string")
			}
			macAddressSpoofing, ok := networkAdapter["mac_address_spoofing"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] mac_address_spoofing should be a string")
			}
			dhcpGuard, ok := networkAdapter["dhcp_guard"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] dhcp_guard should be a string")
			}
			routerGuard, ok := networkAdapter["router_guard"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] router_guard should be a string")
			}
			portMirroring, ok := networkAdapter["port_mirroring"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] port_mirroring should be a string")
			}
			ieeePriorityTag, ok := networkAdapter["ieee_priority_tag"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] ieee_priority_tag should be a string")
			}
			vmqWeight, ok := networkAdapter["vmq_weight"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vmq_weight should be an int")
			}
			iovQueuePairsRequested, ok := networkAdapter["iov_queue_pairs_requested"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] iov_queue_pairs_requested should be an int")
			}
			iovInterruptModeration, ok := networkAdapter["iov_interrupt_moderation"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] iov_interrupt_moderation should be a string")
			}
			iovWeight, ok := networkAdapter["iov_weight"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] iov_weight should be an int")
			}
			ipsecOffloadMaxSA, ok := networkAdapter["ipsec_offload_maximum_security_association"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] ipsec_offload_maximum_security_association should be an int")
			}
			maximumBandwidth, ok := networkAdapter["maximum_bandwidth"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] maximum_bandwidth should be an int")
			}
			minimumBandwidthAbsolute, ok := networkAdapter["minimum_bandwidth_absolute"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] minimum_bandwidth_absolute should be an int")
			}
			minimumBandwidthWeight, ok := networkAdapter["minimum_bandwidth_weight"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] minimum_bandwidth_weight should be an int")
			}
			resourcePoolName, ok := networkAdapter["resource_pool_name"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] resource_pool_name should be a string")
			}
			testReplicaPoolName, ok := networkAdapter["test_replica_pool_name"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] test_replica_pool_name should be a string")
			}
			testReplicaSwitchName, ok := networkAdapter["test_replica_switch_name"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] test_replica_switch_name should be a string")
			}
			virtualSubnetId, ok := networkAdapter["virtual_subnet_id"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] virtual_subnet_id should be an int")
			}
			allowTeaming, ok := networkAdapter["allow_teaming"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] allow_teaming should be a string")
			}
			notMonitoredInCluster, ok := networkAdapter["not_monitored_in_cluster"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] not_monitored_in_cluster should be a bool")
			}
			stormLimit, ok := networkAdapter["storm_limit"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] storm_limit should be an int")
			}
			dynamicIpAddressLimit, ok := networkAdapter["dynamic_ip_address_limit"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] dynamic_ip_address_limit should be an int")
			}
			deviceNaming, ok := networkAdapter["device_naming"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] device_naming should be a string")
			}
			fixSpeed10G, ok := networkAdapter["fix_speed_10g"].(string)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] fix_speed_10g should be a string")
			}
			packetDirectNumProcs, ok := networkAdapter["packet_direct_num_procs"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] packet_direct_num_procs should be an int")
			}
			packetDirectModerationCount, ok := networkAdapter["packet_direct_moderation_count"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] packet_direct_moderation_count should be an int")
			}
			packetDirectModerationInterval, ok := networkAdapter["packet_direct_moderation_interval"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] packet_direct_moderation_interval should be an int")
			}
			vrssEnabled, ok := networkAdapter["vrss_enabled"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vrss_enabled should be a bool")
			}

			expandedNetworkAdapter := VmNetworkAdapter{
				Name:                                   name,
				SwitchName:                             switchName,
				ManagementOs:                           managementOs,
				IsLegacy:                               isLegacy,
				DynamicMacAddress:                      dynamicMacAddress,
				StaticMacAddress:                       staticMacAddress,
				MacAddressSpoofing:                     ToOnOffState(macAddressSpoofing),
				DhcpGuard:                              ToOnOffState(dhcpGuard),
				RouterGuard:                            ToOnOffState(routerGuard),
				PortMirroring:                          ToPortMirroring(portMirroring),
				IeeePriorityTag:                        ToOnOffState(ieeePriorityTag),
				VmqWeight:                              vmqWeight,
				IovQueuePairsRequested:                 iovQueuePairsRequested,
				IovInterruptModeration:                 ToIovInterruptModerationValue(iovInterruptModeration),
				IovWeight:                              iovWeight,
				IpsecOffloadMaximumSecurityAssociation: ipsecOffloadMaxSA,
				MaximumBandwidth:                       maximumBandwidth,
				MinimumBandwidthAbsolute:               minimumBandwidthAbsolute,
				MinimumBandwidthWeight:                 minimumBandwidthWeight,
				MandatoryFeatureId:                     mandatoryFeatureIds,
				ResourcePoolName:                       resourcePoolName,
				TestReplicaPoolName:                    testReplicaPoolName,
				TestReplicaSwitchName:                  testReplicaSwitchName,
				VirtualSubnetId:                        virtualSubnetId,
				AllowTeaming:                           ToOnOffState(allowTeaming),
				NotMonitoredInCluster:                  notMonitoredInCluster,
				StormLimit:                             stormLimit,
				DynamicIpAddressLimit:                  dynamicIpAddressLimit,
				DeviceNaming:                           ToOnOffState(deviceNaming),
				FixSpeed10G:                            ToOnOffState(fixSpeed10G),
				PacketDirectNumProcs:                   packetDirectNumProcs,
				PacketDirectModerationCount:            packetDirectModerationCount,
				PacketDirectModerationInterval:         packetDirectModerationInterval,
				VrssEnabled:                            vrssEnabled,
			}

			// Continue with remaining fields
			vmmqEnabled, ok := networkAdapter["vmmq_enabled"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vmmq_enabled should be a bool")
			}
			vmmqQueuePairs, ok := networkAdapter["vmmq_queue_pairs"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vmmq_queue_pairs should be an int")
			}
			vlanAccess, ok := networkAdapter["vlan_access"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vlan_access should be a bool")
			}
			vlanId, ok := networkAdapter["vlan_id"].(int)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vlan_id should be an int")
			}
			waitForIps, ok := networkAdapter["wait_for_ips"].(bool)
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] wait_for_ips should be a bool")
			}

			expandedNetworkAdapter.VmmqEnabled = vmmqEnabled
			expandedNetworkAdapter.VmmqQueuePairs = vmmqQueuePairs
			expandedNetworkAdapter.VlanAccess = vlanAccess
			expandedNetworkAdapter.VlanId = vlanId
			expandedNetworkAdapter.WaitForIps = waitForIps
			expandedNetworkAdapter.IpAddresses = ipAddresses

			expandedNetworkAdapters = append(expandedNetworkAdapters, expandedNetworkAdapter)
		}
	}

	return expandedNetworkAdapters, nil
}

func FlattenMandatoryFeatureIds(mandatoryFeatureIdStrings []string) *schema.Set {
	if len(mandatoryFeatureIdStrings) < 1 {
		return nil
	}

	var mandatoryFeatureIds []interface{}

	for _, mandatoryFeatureId := range mandatoryFeatureIdStrings {
		mandatoryFeatureIds = append(mandatoryFeatureIds, mandatoryFeatureId)
	}

	return schema.NewSet(schema.HashString, mandatoryFeatureIds)
}

func FlattenNetworkAdapters(networkAdapters *[]VmNetworkAdapter) []interface{} {
	if networkAdapters == nil || len(*networkAdapters) < 1 {
		return nil
	}

	flattenedNetworkAdapters := make([]interface{}, 0)
	for _, networkAdapter := range *networkAdapters {
		flattenedNetworkAdapter := make(map[string]interface{})
		flattenedNetworkAdapter["name"] = networkAdapter.Name
		flattenedNetworkAdapter["switch_name"] = networkAdapter.SwitchName
		flattenedNetworkAdapter["management_os"] = networkAdapter.ManagementOs
		flattenedNetworkAdapter["is_legacy"] = networkAdapter.IsLegacy
		flattenedNetworkAdapter["dynamic_mac_address"] = networkAdapter.DynamicMacAddress
		flattenedNetworkAdapter["static_mac_address"] = networkAdapter.StaticMacAddress
		flattenedNetworkAdapter["mac_address_spoofing"] = networkAdapter.MacAddressSpoofing.String()
		flattenedNetworkAdapter["dhcp_guard"] = networkAdapter.DhcpGuard.String()
		flattenedNetworkAdapter["router_guard"] = networkAdapter.RouterGuard.String()
		flattenedNetworkAdapter["port_mirroring"] = networkAdapter.PortMirroring.String()
		flattenedNetworkAdapter["ieee_priority_tag"] = networkAdapter.IeeePriorityTag.String()
		flattenedNetworkAdapter["vmq_weight"] = networkAdapter.VmqWeight
		flattenedNetworkAdapter["iov_queue_pairs_requested"] = networkAdapter.IovQueuePairsRequested
		flattenedNetworkAdapter["iov_interrupt_moderation"] = networkAdapter.IovInterruptModeration.String()
		flattenedNetworkAdapter["iov_weight"] = networkAdapter.IovWeight
		flattenedNetworkAdapter["ipsec_offload_maximum_security_association"] = networkAdapter.IpsecOffloadMaximumSecurityAssociation
		flattenedNetworkAdapter["maximum_bandwidth"] = networkAdapter.MaximumBandwidth
		flattenedNetworkAdapter["minimum_bandwidth_absolute"] = networkAdapter.MinimumBandwidthAbsolute
		flattenedNetworkAdapter["minimum_bandwidth_weight"] = networkAdapter.MinimumBandwidthWeight
		flattenedNetworkAdapter["mandatory_feature_id"] = FlattenMandatoryFeatureIds(networkAdapter.MandatoryFeatureId)
		flattenedNetworkAdapter["resource_pool_name"] = networkAdapter.ResourcePoolName
		flattenedNetworkAdapter["test_replica_pool_name"] = networkAdapter.TestReplicaPoolName
		flattenedNetworkAdapter["test_replica_switch_name"] = networkAdapter.TestReplicaSwitchName
		flattenedNetworkAdapter["virtual_subnet_id"] = networkAdapter.VirtualSubnetId
		flattenedNetworkAdapter["allow_teaming"] = networkAdapter.AllowTeaming.String()
		flattenedNetworkAdapter["not_monitored_in_cluster"] = networkAdapter.NotMonitoredInCluster
		flattenedNetworkAdapter["storm_limit"] = networkAdapter.StormLimit
		flattenedNetworkAdapter["dynamic_ip_address_limit"] = networkAdapter.DynamicIpAddressLimit
		flattenedNetworkAdapter["device_naming"] = networkAdapter.DeviceNaming.String()
		flattenedNetworkAdapter["fix_speed_10g"] = networkAdapter.FixSpeed10G.String()
		flattenedNetworkAdapter["packet_direct_num_procs"] = networkAdapter.PacketDirectNumProcs
		flattenedNetworkAdapter["packet_direct_moderation_count"] = networkAdapter.PacketDirectModerationCount
		flattenedNetworkAdapter["packet_direct_moderation_interval"] = networkAdapter.PacketDirectModerationInterval
		flattenedNetworkAdapter["vrss_enabled"] = networkAdapter.VrssEnabled
		flattenedNetworkAdapter["vmmq_enabled"] = networkAdapter.VmmqEnabled
		flattenedNetworkAdapter["vmmq_queue_pairs"] = networkAdapter.VmmqQueuePairs
		flattenedNetworkAdapter["vlan_access"] = networkAdapter.VlanAccess
		flattenedNetworkAdapter["vlan_id"] = networkAdapter.VlanId
		flattenedNetworkAdapter["wait_for_ips"] = networkAdapter.WaitForIps
		flattenedNetworkAdapter["ip_addresses"] = networkAdapter.IpAddresses

		flattenedNetworkAdapters = append(flattenedNetworkAdapters, flattenedNetworkAdapter)
	}

	return flattenedNetworkAdapters
}

type VmNetworkAdapterWaitForIp struct {
	Name       string
	WaitForIps bool
}

type VmNetworkAdapter struct {
	VmName                                 string
	Index                                  int
	Name                                   string
	SwitchName                             string
	ManagementOs                           bool
	IsLegacy                               bool
	DynamicMacAddress                      bool
	StaticMacAddress                       string
	MacAddressSpoofing                     OnOffState
	DhcpGuard                              OnOffState
	RouterGuard                            OnOffState
	PortMirroring                          PortMirroring
	IeeePriorityTag                        OnOffState
	VmqWeight                              int
	IovQueuePairsRequested                 int
	IovInterruptModeration                 IovInterruptModerationValue
	IovWeight                              int
	IpsecOffloadMaximumSecurityAssociation int
	MaximumBandwidth                       int
	MinimumBandwidthAbsolute               int
	MinimumBandwidthWeight                 int
	MandatoryFeatureId                     []string
	ResourcePoolName                       string
	TestReplicaPoolName                    string
	TestReplicaSwitchName                  string
	VirtualSubnetId                        int
	AllowTeaming                           OnOffState
	NotMonitoredInCluster                  bool
	StormLimit                             int
	DynamicIpAddressLimit                  int
	DeviceNaming                           OnOffState
	FixSpeed10G                            OnOffState
	PacketDirectNumProcs                   int
	PacketDirectModerationCount            int
	PacketDirectModerationInterval         int
	VrssEnabled                            bool
	VmmqEnabled                            bool
	VmmqQueuePairs                         int
	VlanAccess                             bool
	VlanId                                 int
	WaitForIps                             bool
	IpAddresses                            []string
}

func ExpandVmNetworkAdapterWaitForIps(d *schema.ResourceData) ([]VmNetworkAdapterWaitForIp, uint32, uint32, error) {
	expandVmNetworkAdapterWaitForIps := make([]VmNetworkAdapterWaitForIp, 0)

	timeoutVal := d.Get("wait_for_ips_timeout")
	timeout, ok := timeoutVal.(int)
	if !ok {
		return nil, 0, 0, fmt.Errorf("[ERROR][hyperv] wait_for_ips_timeout should be an int - was '%+v'", timeoutVal)
	}
	waitForIpsTimeout := uint32(timeout)
	if waitForIpsTimeout == 0 {
		waitForIpsTimeout = 300
	}

	pollPeriodVal := d.Get("wait_for_ips_poll_period")
	pollPeriod, ok := pollPeriodVal.(int)
	if !ok {
		return nil, 0, 0, fmt.Errorf("[ERROR][hyperv] wait_for_ips_poll_period should be an int - was '%+v'", pollPeriodVal)
	}
	waitForIpsPollPeriod := uint32(pollPeriod)

	if v, ok := d.GetOk("network_adaptors"); ok {
		networkAdapters, ok := v.([]interface{})
		if !ok {
			return nil, 0, 0, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a list - was '%+v'", v)
		}

		for _, networkAdapter := range networkAdapters {
			networkAdapter, ok := networkAdapter.(map[string]interface{})
			if !ok {
				return nil, waitForIpsTimeout, waitForIpsPollPeriod, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a Hash - was '%+v'", networkAdapter)
			}

			name, ok := networkAdapter["name"].(string)
			if !ok {
				return nil, waitForIpsTimeout, waitForIpsPollPeriod, fmt.Errorf("[ERROR][hyperv] name should be a string")
			}
			waitForIps, ok := networkAdapter["wait_for_ips"].(bool)
			if !ok {
				return nil, waitForIpsTimeout, waitForIpsPollPeriod, fmt.Errorf("[ERROR][hyperv] wait_for_ips should be a bool")
			}

			expandedNetworkAdapterWaitForIp := VmNetworkAdapterWaitForIp{
				Name:       name,
				WaitForIps: waitForIps,
			}

			expandVmNetworkAdapterWaitForIps = append(expandVmNetworkAdapterWaitForIps, expandedNetworkAdapterWaitForIp)
		}
	}

	return expandVmNetworkAdapterWaitForIps, waitForIpsTimeout, waitForIpsPollPeriod, nil
}

type HypervVmNetworkAdapterClient interface {
	CreateVmNetworkAdapter(
		ctx context.Context,
		vmName string,
		name string,
		switchName string,
		managementOs bool,
		isLegacy bool,
		dynamicMacAddress bool,
		staticMacAddress string,
		macAddressSpoofing OnOffState,
		dhcpGuard OnOffState,
		routerGuard OnOffState,
		portMirroring PortMirroring,
		ieeePriorityTag OnOffState,
		vmqWeight int,
		iovQueuePairsRequested int,
		iovInterruptModeration IovInterruptModerationValue,
		iovWeight int,
		ipsecOffloadMaximumSecurityAssociation int,
		maximumBandwidth int,
		minimumBandwidthAbsolute int,
		minimumBandwidthWeight int,
		mandatoryFeatureId []string,
		resourcePoolName string,
		testReplicaPoolName string,
		testReplicaSwitchName string,
		virtualSubnetId int,
		allowTeaming OnOffState,
		notMonitoredInCluster bool,
		stormLimit int,
		dynamicIpAddressLimit int,
		deviceNaming OnOffState,
		fixSpeed10G OnOffState,
		packetDirectNumProcs int,
		packetDirectModerationCount int,
		packetDirectModerationInterval int,
		vrssEnabled bool,
		vmmqEnabled bool,
		vmmqQueuePairs int,
		vlanAccess bool,
		vlanId int,
	) (err error)
	WaitForVmNetworkAdaptersIps(
		ctx context.Context,
		vmName string,
		timeout uint32,
		pollPeriod uint32,
		vmNetworkAdaptersWaitForIps []VmNetworkAdapterWaitForIp,
	) (err error)
	GetVmNetworkAdapters(ctx context.Context, vmName string, networkAdaptersWaitForIps []VmNetworkAdapterWaitForIp) (result []VmNetworkAdapter, err error)
	UpdateVmNetworkAdapter(
		ctx context.Context,
		vmName string,
		index int,
		name string,
		switchName string,
		managementOs bool,
		isLegacy bool,
		dynamicMacAddress bool,
		staticMacAddress string,
		macAddressSpoofing OnOffState,
		dhcpGuard OnOffState,
		routerGuard OnOffState,
		portMirroring PortMirroring,
		ieeePriorityTag OnOffState,
		vmqWeight int,
		iovQueuePairsRequested int,
		iovInterruptModeration IovInterruptModerationValue,
		iovWeight int,
		ipsecOffloadMaximumSecurityAssociation int,
		maximumBandwidth int,
		minimumBandwidthAbsolute int,
		minimumBandwidthWeight int,
		mandatoryFeatureId []string,
		resourcePoolName string,
		testReplicaPoolName string,
		testReplicaSwitchName string,
		virtualSubnetId int,
		allowTeaming OnOffState,
		notMonitoredInCluster bool,
		stormLimit int,
		dynamicIpAddressLimit int,
		deviceNaming OnOffState,
		fixSpeed10G OnOffState,
		packetDirectNumProcs int,
		packetDirectModerationCount int,
		packetDirectModerationInterval int,
		vrssEnabled bool,
		vmmqEnabled bool,
		vmmqQueuePairs int,
		vlanAccess bool,
		vlanId int,
	) (err error)
	DeleteVmNetworkAdapter(ctx context.Context, vmName string, index int) (err error)
	CreateOrUpdateVmNetworkAdapters(ctx context.Context, vmName string, networkAdapters []VmNetworkAdapter) (err error)
}
