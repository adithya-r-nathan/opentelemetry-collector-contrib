// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tgwflowlog // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/unmarshaler/tgw-flow-log"

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/constants"
)

func createTGWFlowLogContent(b *testing.B, filename string, nLogs int) []byte {
	data, err := os.ReadFile(filename)
	require.NoError(b, err)

	lines := bytes.Split(data, []byte{'\n'})
	if len(lines) < 2 {
		require.Fail(b, "file should have at least 1 line for fields and 1 line for the TGW flow log")
	}

	fieldLine := lines[0]
	flowLog := lines[1]

	result := make([][]byte, nLogs+1)
	result[0] = fieldLine
	for i := range nLogs {
		result[i+1] = flowLog
	}

	return bytes.Join(result, []byte{'\n'})
}

func BenchmarkUnmarshalPlainTextLogs(b *testing.B) {
	// each log line of this file is around 200B
	filename := "./testdata/valid_tgw_flow_log_single.log"

	tests := map[string]struct {
		nLogs int
	}{
		"1_log": {
			nLogs: 1,
		},
		"1000_logs": {
			nLogs: 1_000,
		},
	}

	u := tgwFlowLogUnmarshaler{
		fileFormat: constants.FileFormatPlainText,
		buildInfo:  component.BuildInfo{},
		logger:     zap.NewNop(),
	}

	for name, test := range tests {
		data := createTGWFlowLogContent(b, filename, test.nLogs)

		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, err := unmarshalPlainText(bytes.NewReader(data), u.buildInfo)
				require.NoError(b, err)
			}
		})
	}
}
