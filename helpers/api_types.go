// Copyright (c) Technofish Consulting Pty Ltd.
// SPDX-License-Identifier: MPL-2.0

package helpers

type ListSegmentPortsRequest struct {
	SegmentId string `json:"segment_id"`
}

type ListSegmentPortsResponse struct {
	Results       []SegmentPort `json:"results"`
	ResultCount   int           `json:"result_count"`
	SortBy        string        `json:"sort_by"`
	SortAscending bool          `json:"sort_ascending"`
}

type PatchSegmentPortRequest struct {
	SegmentId   string      `json:"segment_id"`
	PortId      string      `json:"port_id"`
	SegmentPort SegmentPort `json:"segment_port"`
}
