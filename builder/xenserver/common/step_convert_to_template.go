package common

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	xsclient "github.com/xenserver/go-xenserver-client"
)

type StepConvertToTemplate struct{}

func (self *StepConvertToTemplate) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("commonconfig").(CommonConfig)
	ui := state.Get("ui").(packer.Ui)
	client := state.Get("client").(xsclient.XenAPIClient)
	instance_uuid := state.Get("instance_uuid").(string)

	instance, err := client.GetVMByUuid(instance_uuid)
	if err != nil {
		ui.Error(fmt.Sprintf("Could not get VM with UUID '%s': %s", instance_uuid, err.Error()))
		return multistep.ActionHalt
	}

	ui.Say("Step: convert VM to template")

	if config.Convert == false {
		ui.Say("Skipping conversion")
		return multistep.ActionContinue
	}

	success := func() bool {
		ui.Message("Converting VM to template")
		err := instance.SetIsATemplate(true)
		if err != nil {
			return false
		}
		return true
	}()

	if !success {
		ui.Error(fmt.Sprintf("Could convert VM to template"))
		return multistep.ActionHalt
	}

	ui.Message("Successfully converted VM to template")
	return multistep.ActionContinue
}

func (StepConvertToTemplate) Cleanup(state multistep.StateBag) {}
