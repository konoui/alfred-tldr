package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/konoui/go-alfred"
)

func printUpdateResults(err error) (_ error) {
	if err != nil {
		fmt.Fprintf(outStream, "update failed due to %s", err)
	} else {
		fmt.Fprintf(outStream, "update succeeded")
	}
	return
}

func (cfg *config) updateTLDRWorkflow() error {
	if cfg.confirm {
		awf.Logger().Infoln("updating tldr workflow...")
		ctx, cancel := context.WithTimeout(context.Background(), updateWorkflowTimeout)
		defer cancel()
		err := awf.Updater().Update(ctx)
		return printUpdateResults(err)
	}

	return errors.New("direct update via flag is not supported")
}

func (cfg *config) updateDB() error {
	if cfg.confirm {
		// update explicitly
		awf.Logger().Infoln("updating tldr database...")
		ctx, cancel := context.WithTimeout(context.Background(), updateDBTimeout)
		defer cancel()
		err := cfg.tldrClient.Update(ctx)
		return printUpdateResults(err)
	}

	awf.Append(
		alfred.NewItem().
			Title("Please Enter if update tldr database").
			Arg(fmt.Sprintf("--%s --%s", longUpdateFlag, confirmFlag)),
	).
		Variable(nextActionKey, nextActionShell).
		Output()

	return nil
}
