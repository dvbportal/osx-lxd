package main

import (
	"fmt"

	"github.com/chai2010/gettext-go/gettext"

	"github.com/dvbportal/osx-lxd"
	"github.com/dvbportal/osx-lxd/shared"
)

type deleteCmd struct{}

func (c *deleteCmd) showByDefault() bool {
	return true
}

func (c *deleteCmd) usage() string {
	return gettext.Gettext(
		"Delete containers or container snapshots.\n" +
			"\n" +
			"lxc delete <container>[/<snapshot>] [<container>[/<snapshot>]...]\n" +
			"\n" +
			"Destroy containers or snapshots with any attached data (configuration,\n" +
			"snapshots, ...).\n")
}

func (c *deleteCmd) flags() {}

func doDelete(d *lxd.Client, name string) error {
	resp, err := d.Delete(name)
	if err != nil {
		return err
	}

	return d.WaitForSuccess(resp.Operation)
}

func (c *deleteCmd) run(config *lxd.Config, args []string) error {
	if len(args) == 0 {
		return errArgs
	}

	for _, nameArg := range args {
		remote, name := config.ParseRemoteAndContainer(nameArg)

		d, err := lxd.NewClient(config, remote)
		if err != nil {
			return err
		}

		ct, err := d.ContainerStatus(name, false)

		if err != nil {
			// Could be a snapshot
			return doDelete(d, name)
		}

		if ct.State() != shared.STOPPED {
			resp, err := d.Action(name, shared.Stop, -1, true)
			if err != nil {
				return err
			}

			op, err := d.WaitFor(resp.Operation)
			if err != nil {
				return err
			}

			if op.StatusCode == shared.Failure {
				return fmt.Errorf(gettext.Gettext("Stopping container failed!"))
			}

			if ct.Ephemeral == true {
				return nil
			}
		}
		if err := doDelete(d, name); err != nil {
			return err
		}
	}
	return nil

}
