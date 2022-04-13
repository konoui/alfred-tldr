package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/konoui/go-alfred"
)

func printUpdateResults(w io.Writer, err error) (_ error) {
	if err != nil {
		fmt.Fprintf(w, "update failed due to %s", err)
	} else {
		fmt.Fprintf(w, "update succeeded")
	}
	return
}

func updateTLDRWorkflow(c *client) error {
	if c.cfg.confirm {
		c.Logger().Infoln("updating tldr workflow...")
		ctx, cancel := context.WithTimeout(context.Background(), updateWorkflowTimeout)
		defer cancel()
		err := c.Updater().Update(ctx)
		return printUpdateResults(c.OutWriter(), err)
	}

	return errors.New("update workflow flag is not supported")
}

func updateDB(c *client) error {
	if c.cfg.confirm {
		// update explicitly
		c.Logger().Infoln("updating tldr database...")
		ctx, cancel := context.WithTimeout(context.Background(), updateDBTimeout)
		defer cancel()
		err := c.tldrClient.Update(ctx)
		return printUpdateResults(c.OutWriter(), err)
	}

	c.Append(
		alfred.NewItem().
			Title("Please Enter if update tldr database").
			Arg(fmt.Sprintf("--%s --%s", longUpdateFlag, confirmFlag)),
	).
		Variable(nextActionKey, nextActionShell).
		Output()

	return nil
}
