// Copyright (c) Technofish Consulting Pty Ltd
// SPDX-License-Identifier: MPL-2.0

package helpers

import "encoding/json"

func ConvertToTFSegmentPort(jsonData string) (SegmentPort, error) {
	var segmentPort SegmentPort
	err := json.Unmarshal([]byte(jsonData), &segmentPort)
	if err != nil {
		return segmentPort, err
	}

	return segmentPort, nil
}
