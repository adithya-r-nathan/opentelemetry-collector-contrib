// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tgwflowlog // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/awslogsencodingextension/internal/unmarshaler/tgw-flow-log"

// tgwFieldMap maps Transit Gateway flow log field names to OpenTelemetry attribute keys.
// Fields shared with VPC flow logs reuse the same semconv keys for consistency.
// TGW-specific fields use the aws.tgw.* namespace.
//
// Reference: https://docs.aws.amazon.com/vpc/latest/tgw/tgw-flow-logs.html
var tgwFieldMap = map[string]string{
	// Fields shared with VPC flow logs
	"version":        "aws.tgw.flow.version",
	"account-id":     "cloud.account.id",
	"srcaddr":        "source.address",
	"dstaddr":        "destination.address",
	"srcport":        "source.port",
	"dstport":        "destination.port",
	"protocol":       "network.transport",
	"packets":        "aws.tgw.flow.packets",
	"bytes":          "aws.tgw.flow.bytes",
	"log-status":     "aws.tgw.flow.log_status",
	"type":           "network.type",
	"tcp-flags":      "aws.tgw.flow.tcp_flags",
	"flow-direction": "network.io.direction",
	"region":         "cloud.region",
	"resource-type":  "aws.tgw.resource_type",

	// TGW-specific fields
	"tgw-id":                 "aws.tgw.id",
	"tgw-attachment-id":      "aws.tgw.attachment.id",
	"tgw-src-vpc-account-id": "aws.tgw.src.vpc.account.id",
	"tgw-dst-vpc-account-id": "aws.tgw.dst.vpc.account.id",
	"tgw-src-vpc-id":         "aws.tgw.src.vpc.id",
	"tgw-dst-vpc-id":         "aws.tgw.dst.vpc.id",
	"tgw-src-subnet-id":      "aws.tgw.src.subnet.id",
	"tgw-dst-subnet-id":      "aws.tgw.dst.subnet.id",
	"tgw-src-eni":            "aws.tgw.src.eni",
	"tgw-dst-eni":            "aws.tgw.dst.eni",
	"tgw-src-az-id":          "aws.tgw.src.az.id",
	"tgw-dst-az-id":          "aws.tgw.dst.az.id",
	"tgw-pair-attachment-id": "aws.tgw.pair.attachment.id",

	// Packet loss fields (TGW-specific)
	"packets-lost-no-route":     "aws.tgw.flow.packets_lost.no_route",
	"packets-lost-blackhole":    "aws.tgw.flow.packets_lost.blackhole",
	"packets-lost-mtu-exceeded": "aws.tgw.flow.packets_lost.mtu_exceeded",
	"packets-lost-ttl-expired":  "aws.tgw.flow.packets_lost.ttl_expired",
}
