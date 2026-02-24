// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tgwflowlog

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/constants"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
)

func compressToGZIPReader(t *testing.T, buf []byte) io.Reader {
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(buf)
	require.NoError(t, err)
	err = gzipWriter.Close()
	require.NoError(t, err)

	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedData.Bytes()))
	require.NoError(t, err)
	return gzipReader
}

func readAndCompressLogFile(t *testing.T, dir, file string) io.Reader {
	data, err := os.ReadFile(filepath.Join(dir, file))
	require.NoError(t, err)
	return compressToGZIPReader(t, data)
}

func readLogFile(t *testing.T, dir, file string) io.Reader {
	data, err := os.ReadFile(filepath.Join(dir, file))
	require.NoError(t, err)
	return bytes.NewReader(data)
}

func TestUnmarshalLogs_PlainText(t *testing.T) {
	t.Parallel()

	dir := "testdata"

	tests := []struct {
		name                 string
		logInputReader       io.Reader
		logsExpectedFilename string
		expectedErr          string
	}{
		{
			name:                 "Valid TGW flow log - single log",
			logInputReader:       readAndCompressLogFile(t, dir, "valid_tgw_flow_log_single.log"),
			logsExpectedFilename: "valid_tgw_flow_log_single_expected.yaml",
		},
		{
			name:                 "Valid TGW flow log - multi logs",
			logInputReader:       readAndCompressLogFile(t, dir, "valid_tgw_flow_log_multi.log"),
			logsExpectedFilename: "valid_tgw_flow_log_multi_expected.yaml",
		},
		{
			name:                 "Valid TGW flow log uncompressed",
			logInputReader:       readLogFile(t, dir, "valid_tgw_flow_log_single.log"),
			logsExpectedFilename: "valid_tgw_flow_log_single_expected.yaml",
		},
		{
			name:                 "Empty input",
			logInputReader:       bytes.NewReader([]byte{}),
			logsExpectedFilename: "valid_tgw_flow_log_empty_expected.yaml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := NewTGWFlowLogUnmarshaler(
				Config{FileFormat: constants.FileFormatPlainText},
				component.BuildInfo{},
				zap.NewNop(),
			)
			require.NoError(t, err)

			logs, err := u.UnmarshalAWSLogs(test.logInputReader)

			if test.expectedErr != "" {
				require.ErrorContains(t, err, test.expectedErr)
				return
			}

			// To generate golden files, uncomment:
			// golden.WriteLogsToFile(filepath.Join(dir, test.logsExpectedFilename), logs)

			require.NoError(t, err)

			expectedLogs, err := golden.ReadLogs(filepath.Join(dir, test.logsExpectedFilename))
			require.NoError(t, err)

			require.NoError(t, plogtest.CompareLogs(expectedLogs, logs))
		})
	}
}
