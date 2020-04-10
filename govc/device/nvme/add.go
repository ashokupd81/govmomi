package nvme

import (
	"context"
	"flag"
	"fmt"
	//"strings"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	//"github.com/vmware/govmomi/object"
	//"github.com/vmware/govmomi/vim25/types"
)

type add struct {
	*flags.VirtualMachineFlag

	controller   string
	//sharedBus    string
	//hotAddRemove bool
}

func init() {
	cli.Register("device.nvme.add", &add{})
}

func (cmd *add) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)

}

func (cmd *add) Description() string {
	return `Add NVME controller to VM.

Examples:
  govc device.scsi.add -vm $vm`
}

func (cmd *add) Process(ctx context.Context) error {
	if err := cmd.VirtualMachineFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *add) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	if vm == nil {
		return flag.ErrHelp
	}

	devices, err := vm.Device(ctx)
	if err != nil {
		return err
	}

	nvme, err := devices.CreateNVMEController()
	if err != nil {
		return err
	}
	err = vm.AddDevice(ctx, nvme)
	if err != nil {
		return err
	}

	// output name of device we just created
	devices, err = vm.Device(ctx)
	if err != nil {
		return err
	}

	devices = devices.SelectByType(nvme)

	name := devices.Name(devices[len(devices)-1])

	fmt.Println(name)

	return nil
}
