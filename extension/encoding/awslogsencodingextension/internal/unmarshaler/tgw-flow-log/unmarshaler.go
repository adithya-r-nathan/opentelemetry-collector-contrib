// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tgwflowlog // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/unmarshaler/tgw-flow-log"

import (
	"fmt"
	"io"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// Config defines the configuration for Transit Gateway flow log unmarshaling.
type Config struct {
	// FileFormat specifies the file format of the TGW flow logs.
	// Supported values: plain-text, parquet (parquet not yet implemented).
	// Default: plain-text.
	FileFormat string `mapstructure:"file_format"`
}

// tgwFlowLogUnmarshaler unmarshals Transit Gateway flow log records.
type tgwFlowLogUnmarshaler struct {
	buildInfo  component.BuildInfo
	logger     *zap.Logger
	fileFormat string
}

// NewTGWFlowLogUnmarshaler creates a new unmarshaler for Transit Gateway flow logs.
func NewTGWFlowLogUnmarshaler(cfg Config, buildInfo component.BuildInfo, logger *zap.Logger) (*tgwFlowLogUnmarshaler, error) {
	fileFormat := cfg.FileFormat
	if fileFormat == "" {
		fileFormat = "plain-text"
	}
	return &tgwFlowLogUnmarshaler{
		buildInfo:  buildInfo,
		logger:     logger,
		fileFormat: fileFormat,
	}, nil
}

// UnmarshalAWSLogs implements awsunmarshaler.AWSUnmarshaler.
func (u *tgwFlowLogUnmarshaler) UnmarshalAWSLogs(reader io.Reader) (plog.Logs, error) {
	switch u.fileFormat {
	case "plain-text":
		return unmarshalPlainText(reader, u.buildInfo)
	case "parquet":
		return plog.Logs{}, fmt.Errorf("parquet format not yet supported for Transit Gateway flow logs")
	default:
		return plog.Logs{}, fmt.Errorf("unsupported file format %q", u.fileFormat)
	}
}
