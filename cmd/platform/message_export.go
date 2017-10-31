// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/spf13/cobra"
)

var messageExportCmd = &cobra.Command{
	Use:     "export",
	Short:   "Export data from Mattermost",
	Long:    "Export data from Mattermost in a format suitable for import into a third-party application",
	Example: "export --format=actiance --exportFrom=12345",
	RunE:    messageExportCmdF,
}

func init() {
	messageExportCmd.Flags().String("format", "actiance", "The format to export data in")
	messageExportCmd.Flags().Int64("exportFrom", -1, "The timestamp of the earliest post to export, expressed in seconds since the unix epoch.")
}

func messageExportCmdF(cmd *cobra.Command, args []string) error {
	a, err := initDBCommandContextCobra(cmd)
	if err != nil {
		return err
	}

	if !*a.Config().MessageExportSettings.EnableExport {
		CommandPrintErrorln("ERROR: The message export feature is not enabled.")
		return nil
	}

	// for now, format is hard-coded to actiance. In time, we'll have to support other formats and inject them into job data
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return errors.New("format flag error")
	} else if format != "actiance" {
		return errors.New("unsupported export format")
	}

	startTime, err := cmd.Flags().GetInt64("exportFrom")
	if err != nil {
		return errors.New("exportFrom flag error")
	}

	if messageExportI := a.MessageExport; messageExportI != nil {
		job, err := messageExportI.StartSynchronizeJob(true, startTime)
		if err != nil || job.Status == model.JOB_STATUS_ERROR || job.Status == model.JOB_STATUS_CANCELED {
			CommandPrintErrorln("ERROR: Message export job failed. Please check the server logs")
		} else {
			CommandPrettyPrintln("SUCCESS: Message export job complete")
		}
	}

	return nil
}
