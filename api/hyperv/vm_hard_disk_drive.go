package hyperv

import (
	"context"
	"encoding/json"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type createVmHardDiskDriveArgs struct {
	VmHardDiskDriveJson string
}

var createVmHardDiskDriveTemplate = template.Must(template.New("CreateVmHardDiskDrive").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmHardDiskDrive = '{{.VmHardDiskDriveJson}}' | ConvertFrom-Json

$NewVmHardDiskDriveArgs = @{
	VmName=$vmHardDiskDrive.VmName
	ControllerType=$vmHardDiskDrive.ControllerType
	ControllerNumber=$vmHardDiskDrive.ControllerNumber
	ControllerLocation=$vmHardDiskDrive.ControllerLocation
	Path=$vmHardDiskDrive.Path
	ResourcePoolName=$vmHardDiskDrive.ResourcePoolName
	SupportPersistentReservations=$vmHardDiskDrive.SupportPersistentReservations
	MaximumIops=$_.MaximumIops;
	MinimumIops=$_.MinimumIops;
	QosPolicyId=$_.QosPolicyId;
	OverrideCacheAttributes=$vmHardDiskDrive.OverrideCacheAttributes
	AllowUnverifiedPaths=$true
}

if ($vmHardDiskDrive.DiskNumber -lt 4294967295){
	$NewVmHardDiskDriveArgs.DiskNumber=$vmHardDiskDrive.DiskNumber
}

Add-VmHardDiskDrive @NewVmHardDiskDriveArgs
`))

func (c *ClientConfig) CreateVmHardDiskDrive(
	ctx context.Context,
	vmName string,
	controllerType api.ControllerType,
	controllerNumber int32,
	controllerLocation int32,
	path string,
	diskNumber uint32,
	resourcePoolName string,
	supportPersistentReservations bool,
	maximumIops uint64,
	minimumIops uint64,
	qosPolicyId string,
	overrideCacheAttributes api.CacheAttributes,

) (err error) {
	vmHardDiskDriveJson, err := json.Marshal(api.VmHardDiskDrive{
		VmName:                        vmName,
		ControllerType:                controllerType,
		ControllerNumber:              controllerNumber,
		ControllerLocation:            controllerLocation,
		Path:                          path,
		DiskNumber:                    diskNumber,
		ResourcePoolName:              resourcePoolName,
		SupportPersistentReservations: supportPersistentReservations,
		MaximumIops:                   maximumIops,
		MinimumIops:                   minimumIops,
		QosPolicyId:                   qosPolicyId,
		OverrideCacheAttributes:       overrideCacheAttributes,
	})

	if err != nil {
		return err
	}

	err = c.ScriptRunner.RunFireAndForgetScript(ctx, createVmHardDiskDriveTemplate, createVmHardDiskDriveArgs{
		VmHardDiskDriveJson: string(vmHardDiskDriveJson),
	})

	return err
}

type getVmHardDiskDrivesArgs struct {
	VmName string
}

var getVmHardDiskDrivesTemplate = template.Must(template.New("GetVmHardDiskDrives").Parse(`
$ErrorActionPreference = 'Stop'
$vmHardDiskDrivesObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMHardDiskDrive | %{ @{
	ControllerType=$_.ControllerType;
	ControllerNumber=$_.ControllerNumber;
	ControllerLocation=$_.ControllerLocation;
	Path=$_.Path;
	DiskNumber=if ($_.DiskNumber -eq $null) { 4294967295 } else { $_.DiskNumber };
	ResourcePoolName=$_.PoolName;
	SupportPersistentReservations=$_.SupportPersistentReservations;
	MaximumIops=$_.MaximumIops;
	MinimumIops=$_.MinimumIops;
	QosPolicyId=$_.QosPolicyId;	
	OverrideCacheAttributes=$_.WriteHardeningMethod;
}})

if ($vmHardDiskDrivesObject) {
	$vmHardDiskDrives = ConvertTo-Json -InputObject $vmHardDiskDrivesObject
	$vmHardDiskDrives
} else {
	"[]"
}
`))

func (c *ClientConfig) GetVmHardDiskDrives(ctx context.Context, vmName string) (result []api.VmHardDiskDrive, err error) {
	result = make([]api.VmHardDiskDrive, 0)

	err = c.ScriptRunner.RunScriptWithResult(ctx, getVmHardDiskDrivesTemplate, getVmHardDiskDrivesArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

type updateVmHardDiskDriveArgs struct {
	VmName              string
	ControllerNumber    int32
	ControllerLocation  int32
	ControllerType      api.ControllerType
	VmHardDiskDriveJson string
}

var updateVmHardDiskDriveTemplate = template.Must(template.New("UpdateVmHardDiskDrive").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmHardDiskDrive = '{{.VmHardDiskDriveJson}}' | ConvertFrom-Json

$vmHardDiskDrivesObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMHardDiskDrive -ControllerLocation {{.ControllerLocation}} -ControllerNumber {{.ControllerNumber}} -ControllerType {{.ControllerType}})

if (!$vmHardDiskDrivesObject){
	throw "VM hard disk drive does not exist - {{.ControllerLocation}} {{.ControllerNumber}} {{.ControllerType}}"
}

$SetVmHardDiskDriveArgs = @{}
$SetVmHardDiskDriveArgs.VmName=$vmHardDiskDrivesObject.VmName
$SetVmHardDiskDriveArgs.ControllerType=$vmHardDiskDrivesObject.ControllerType
$SetVmHardDiskDriveArgs.ControllerLocation=$vmHardDiskDrivesObject.ControllerLocation
$SetVmHardDiskDriveArgs.ControllerNumber=$vmHardDiskDrivesObject.ControllerNumber
$SetVmHardDiskDriveArgs.ToControllerLocation=$vmHardDiskDrive.ControllerLocation
$SetVmHardDiskDriveArgs.ToControllerNumber=$vmHardDiskDrive.ControllerNumber
$SetVmHardDiskDriveArgs.Path=$vmHardDiskDrive.Path
if ($vmHardDiskDrive.DiskNumber -lt 4294967295){
	$SetVmHardDiskDriveArgs.DiskNumber=$vmHardDiskDrive.DiskNumber
}
if ($vmHardDiskDrivesObject.ResourcePoolName -ne $vmHardDiskDrive.ResourcePoolName) {
	if ($vmHardDiskDrive.ResourcePoolName) {
		$SetVmHardDiskDriveArgs.ResourcePoolName=$vmHardDiskDrive.ResourcePoolName
	} else {
		throw "Unable to remove resource pool $($vmHardDiskDrive.ResourcePoolName) from hard disk drive $(ConvertTo-Json -InputObject $vmHardDiskDrivesObject)"
	}
}
$SetVmHardDiskDriveArgs.SupportPersistentReservations=$vmHardDiskDrive.SupportPersistentReservations
$SetVmHardDiskDriveArgs.MaximumIops=$vmHardDiskDrive.MaximumIops
$SetVmHardDiskDriveArgs.MinimumIops=$vmHardDiskDrive.MinimumIops
$SetVmHardDiskDriveArgs.QosPolicyId=$vmHardDiskDrive.QosPolicyId
$SetVmHardDiskDriveArgs.OverrideCacheAttributes=$vmHardDiskDrive.OverrideCacheAttributes	
$SetVmHardDiskDriveArgs.AllowUnverifiedPaths=$true

Set-VMHardDiskDrive @SetVmHardDiskDriveArgs

`))

func (c *ClientConfig) UpdateVmHardDiskDrive(
	ctx context.Context,
	vmName string,
	controllerNumber int32,
	controllerLocation int32,
	controllerType api.ControllerType,
	toControllerNumber int32,
	toControllerLocation int32,
	path string,
	diskNumber uint32,
	resourcePoolName string,
	supportPersistentReservations bool,
	maximumIops uint64,
	minimumIops uint64,
	qosPolicyId string,
	overrideCacheAttributes api.CacheAttributes,
) (err error) {
	vmHardDiskDriveJson, err := json.Marshal(api.VmHardDiskDrive{
		VmName:                        vmName,
		ControllerType:                controllerType,
		ControllerNumber:              toControllerNumber,
		ControllerLocation:            toControllerLocation,
		Path:                          path,
		DiskNumber:                    diskNumber,
		ResourcePoolName:              resourcePoolName,
		SupportPersistentReservations: supportPersistentReservations,
		MaximumIops:                   maximumIops,
		MinimumIops:                   minimumIops,
		QosPolicyId:                   qosPolicyId,
		OverrideCacheAttributes:       overrideCacheAttributes,
	})

	if err != nil {
		return err
	}

	err = c.ScriptRunner.RunFireAndForgetScript(ctx, updateVmHardDiskDriveTemplate, updateVmHardDiskDriveArgs{
		VmName:              vmName,
		ControllerNumber:    controllerNumber,
		ControllerLocation:  controllerLocation,
		ControllerType:      controllerType,
		VmHardDiskDriveJson: string(vmHardDiskDriveJson),
	})

	return err
}

type deleteVmHardDiskDriveArgs struct {
	VmName             string
	ControllerNumber   int32
	ControllerLocation int32
	ControllerType     api.ControllerType
}

var deleteVmHardDiskDriveTemplate = template.Must(template.New("DeleteVmHardDiskDrive").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VMHardDiskDrive -VmName '{{.VmName}}' -ControllerNumber {{.ControllerNumber}} -ControllerLocation {{.ControllerLocation}} -ControllerType {{.ControllerType}}) | Remove-VMHardDiskDrive
`))

func (c *ClientConfig) DeleteVmHardDiskDrive(ctx context.Context, vmname string, controllerNumber int32, controllerLocation int32, controllerType api.ControllerType) (err error) {
	err = c.ScriptRunner.RunFireAndForgetScript(ctx, deleteVmHardDiskDriveTemplate, deleteVmHardDiskDriveArgs{
		VmName:             vmname,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
		ControllerType:     controllerType,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmHardDiskDrives(ctx context.Context, vmName string, hardDiskDrives []api.VmHardDiskDrive) (err error) {
	currentHardDiskDrives, err := c.GetVmHardDiskDrives(ctx, vmName)
	if err != nil {
		return err
	}

	currentHardDiskDrivesLength := len(currentHardDiskDrives)
	desiredHardDiskDrivesLength := len(hardDiskDrives)

	// 1. Delete extra drives (index out of bounds of desired)
	for i := currentHardDiskDrivesLength - 1; i > desiredHardDiskDrivesLength-1; i-- {
		currentHardDiskDrive := currentHardDiskDrives[i]
		err = c.DeleteVmHardDiskDrive(ctx, vmName, currentHardDiskDrive.ControllerNumber, currentHardDiskDrive.ControllerLocation, currentHardDiskDrive.ControllerType)
		if err != nil {
			return err
		}
	}

	// 2. Identify drives that need to be moved/recreated and delete them
	// We only iterate up to the minimum length (intersection of current and desired)
	limit := currentHardDiskDrivesLength
	if desiredHardDiskDrivesLength < limit {
		limit = desiredHardDiskDrivesLength
	}

	indicesToCreate := make(map[int]bool)

	for i := 0; i < limit; i++ {
		currentHardDiskDrive := currentHardDiskDrives[i]
		hardDiskDrive := hardDiskDrives[i]

		// Check if identity/location changed
		if currentHardDiskDrive.ControllerType != hardDiskDrive.ControllerType ||
			currentHardDiskDrive.ControllerNumber != hardDiskDrive.ControllerNumber ||
			currentHardDiskDrive.ControllerLocation != hardDiskDrive.ControllerLocation {

			// Mismatch in location/type -> Delete and Recreate
			err = c.DeleteVmHardDiskDrive(ctx, vmName, currentHardDiskDrive.ControllerNumber, currentHardDiskDrive.ControllerLocation, currentHardDiskDrive.ControllerType)
			if err != nil {
				return err
			}
			indicesToCreate[i] = true
		} else {
			// Location matches -> Update in place
			err = c.UpdateVmHardDiskDrive(
				ctx,
				vmName,
				currentHardDiskDrive.ControllerNumber,
				currentHardDiskDrive.ControllerLocation,
				currentHardDiskDrive.ControllerType,
				hardDiskDrive.ControllerNumber,
				hardDiskDrive.ControllerLocation,
				hardDiskDrive.Path,
				hardDiskDrive.DiskNumber,
				hardDiskDrive.ResourcePoolName,
				hardDiskDrive.SupportPersistentReservations,
				hardDiskDrive.MaximumIops,
				hardDiskDrive.MinimumIops,
				hardDiskDrive.QosPolicyId,
				hardDiskDrive.OverrideCacheAttributes,
			)
			if err != nil {
				return err
			}
		}
	}

	// 3. Create new drives (either completely new, or those we deleted to move)
	for i := 0; i < desiredHardDiskDrivesLength; i++ {
		// Create if it's a new index (>= limit) OR if we marked it for recreation
		if i >= limit || indicesToCreate[i] {
			hardDiskDrive := hardDiskDrives[i]
			err = c.CreateVmHardDiskDrive(
				ctx,
				vmName,
				hardDiskDrive.ControllerType,
				hardDiskDrive.ControllerNumber,
				hardDiskDrive.ControllerLocation,
				hardDiskDrive.Path,
				hardDiskDrive.DiskNumber,
				hardDiskDrive.ResourcePoolName,
				hardDiskDrive.SupportPersistentReservations,
				hardDiskDrive.MaximumIops,
				hardDiskDrive.MinimumIops,
				hardDiskDrive.QosPolicyId,
				hardDiskDrive.OverrideCacheAttributes,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
