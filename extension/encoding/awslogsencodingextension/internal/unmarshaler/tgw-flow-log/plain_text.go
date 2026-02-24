// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tgwflowlog // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/unmarshaler/tgw-flow-log"

import (
	"bufio"
	"io"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/constants"
)

// unmarshalPlainText parses plain-text Transit Gateway flow logs from an io.Reader.
// The first line is a space-separated header; each subsequent line is a record.
// See https://docs.aws.amazon.com/vpc/latest/tgw/tgw-flow-logs.html.
func unmarshalPlainText(reader io.Reader, buildInfo component.BuildInfo) (plog.Logs, error) {
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	rl.Resource().Attributes().PutStr(constants.FormatIdentificationTag, constants.FormatTransitGatewayFlowLog)
	sl := rl.ScopeLogs().AppendEmpty()
	sl.Scope().SetName(buildInfo.Command)
	sl.Scope().SetVersion(buildInfo.Version)

	scanner := bufio.NewScanner(reader)

	// First line is the header.
	if !scanner.Scan() {
		return ld, scanner.Err()
	}
	fields := strings.Fields(scanner.Text())

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		lr := sl.LogRecords().AppendEmpty()
		parseRecord(fields, line, lr)
	}

	return ld, scanner.Err()
}

// parseRecord maps a single TGW flow log line onto a LogRecord.
// Uses strings.Cut for performance, consistent with the VPC flow log parser.
func parseRecord(fields []string, line string, lr plog.LogRecord) {
	attrs := lr.Attributes()
	for _, field := range fields {
		if line == "" {
			break
		}
		var value string
		value, line, _ = strings.Cut(line, " ")

		// AWS uses "-" for not-applicable or missing fields; skip them.
		if value == "-" {
			continue
		}

		switch field {
		case "end":
			// "end" maps to the log record timestamp, consistent with vpc_flow_log convention.
			setTimestampFromUnix(value, lr)
		case "start":
			attrs.PutStr("aws.tgw.flow.start", value)
		default:
			if otelKey, ok := tgwFieldMap[field]; ok {
				attrs.PutStr(otelKey, value)
			}
			// Unknown/future fields are silently ignored for forward compatibility.
		}
	}
}

// setTimestampFromUnix converts a Unix epoch string (seconds) to a pcommon.Timestamp.
func setTimestampFromUnix(value string, lr plog.LogRecord) {
	var ts int64
	for _, c := range value {
		if c < '0' || c > '9' {
			return
		}
		ts = ts*10 + int64(c-'0')
	}
	lr.SetTimestamp(pcommon.Timestamp(ts * 1_000_000_000))
}
