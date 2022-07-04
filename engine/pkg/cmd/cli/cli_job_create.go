// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/engine/client"
	"github.com/pingcap/tiflow/engine/enginepb"
	engineModel "github.com/pingcap/tiflow/engine/model"
	"github.com/pingcap/tiflow/engine/pkg/tenant"
	cmdcontext "github.com/pingcap/tiflow/pkg/cmd/context"
	"github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/security"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// createJobOptions defines common job flags.
type createJobOptions struct {
	jobTypeStr   string
	jobConfigStr string
	projectID    string
	tenantID     string
	rpcTimeout   time.Duration

	jobType   engineModel.JobType
	jobConfig []byte
	tenant    tenant.ProjectInfo

	credential *security.Credential
}

// newCreateJobOptions creates new job options.
func newCreateJobOptions() *createJobOptions {
	return &createJobOptions{}
}

// addFlags receives a *cobra.Command reference and binds
// flags related to template printing to it.
func (o *createJobOptions) addFlags(cmd *cobra.Command) {
	if o == nil {
		return
	}

	cmd.Flags().StringVar(&o.jobTypeStr, "job-type", "", "job type")
	cmd.Flags().StringVar(&o.jobConfigStr, "job-config", "", "path of config file for the job")
	cmd.Flags().StringVar(&o.tenantID, "tenant-id", "", "the tenant id")
	cmd.Flags().StringVar(&o.projectID, "project-id", "", "the project id")
	cmd.Flags().DurationVar(&o.rpcTimeout, "rpc-timeout", time.Second*30, "default rpc timeout")
}

// validate checks that the provided job options are valid.
func (o *createJobOptions) validate(ctx context.Context, cmd *cobra.Command) error {
	jobType, ok := engineModel.GetJobTypeByName(o.jobTypeStr)
	if !ok {
		return errors.ErrInvalidJobType.GenWithStackByArgs(o.jobType)
	}
	jobConfig, err := openFileAndReadString(o.jobConfigStr)
	if err != nil {
		return err
	}
	o.jobType = jobType
	o.jobConfig = jobConfig
	o.tenant = o.getProjectInfo()

	return nil
}

func openFileAndReadString(path string) (content []byte, err error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return io.ReadAll(fp)
}

func (o *createJobOptions) getProjectInfo() tenant.ProjectInfo {
	var tenantID, projectID string
	if o.tenantID != "" {
		tenantID = o.tenantID
	} else {
		log.Warn("tenant-id is empty, use default tenant id")
		tenantID = tenant.DefaultUserProjectInfo.TenantID()
	}
	if o.projectID != "" {
		projectID = o.projectID
	} else {
		log.Warn("project-id is empty, use default project id")
		projectID = tenant.DefaultUserProjectInfo.ProjectID()
	}
	return tenant.NewProjectInfo(tenantID, projectID)
}

// run the `cli job create` command.
func (o *createJobOptions) run(ctx context.Context, cmd *cobra.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), o.rpcTimeout)
	defer cancel()

	cliManager := client.NewClientManager()
	resp, err := cliManager.MasterClient().SubmitJob(ctx, &enginepb.SubmitJobRequest{
		Tp:     int32(o.jobType),
		Config: o.jobConfig,
		ProjectInfo: &enginepb.ProjectInfo{
			TenantId:  o.tenant.TenantID(),
			ProjectId: o.tenant.ProjectID(),
		},
	})
	if err != nil {
		return err
	}
	log.L().Info("create job successfully", zap.Any("resp", resp))
	return nil
}

// newCmdCreateJob creates the `cli job create` command.
func newCmdCreateJob() *cobra.Command {
	o := newCreateJobOptions()

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a new job",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmdcontext.GetDefaultContext()

			if err := o.validate(ctx, cmd); err != nil {
				return err
			}

			return o.run(ctx, cmd)
		},
	}

	o.addFlags(command)

	return command
}
