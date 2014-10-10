package softlayer

import (
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

type SoftLayer_Virtual_Guest_Service interface {
	Service

	GetObject(instanceId int) (datatypes.SoftLayer_Virtual_Guest, error)
	CreateObject(template datatypes.SoftLayer_Virtual_Guest_Template) (datatypes.SoftLayer_Virtual_Guest, error)
	EditObject(instanceId int, template datatypes.SoftLayer_Virtual_Guest) (bool, error)
	DeleteObject(instanceId int) (bool, error)

	GetPowerState(instanceId int) (datatypes.SoftLayer_Virtual_Guest_Power_State, error)

	GetSshKeys(instanceId int) ([]datatypes.SoftLayer_Security_Ssh_Key, error)

	GetActiveTransaction(instanceId int) (datatypes.SoftLayer_Provisioning_Version1_Transaction, error)
	GetActiveTransactions(instanceId int) ([]datatypes.SoftLayer_Provisioning_Version1_Transaction, error)

	RebootSoft(instanceId int) (bool, error)
	RebootHard(instanceId int) (bool, error)

	SetMetadata(instanceId int, metadata string) (bool, error)

	ConfigureMetadataDisk(instanceId int) (datatypes.SoftLayer_Provisioning_Version1_Transaction, error)
}
