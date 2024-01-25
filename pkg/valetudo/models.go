package valetudo

type RobotStateRemainingAttribute struct {
	Value *int    `json:"value"`
	Unit  *string `json:"unit"`
}

type RobotStateAttribute struct {
	Class       string                        `json:"__class"`
	Type        *string                       `json:"type"`
	SubType     *string                       `json:"subType"`
	Value       *string                       `json:"value"`
	CustomValue *string                       `json:"customValue"`
	Attached    *bool                         `json:"attached,omitempty"`
	Level       *int                          `json:"level,omitempty"`
	Flag        *string                       `json:"flag,omitempty"`
	Remaining   *RobotStateRemainingAttribute `json:"remaining,omitempty"`
}

type RobotStateMap struct {
	PixelSize int                  `json:"pixelSize"`
	Layers    []RobotStateMapLayer `json:"layers"`
}

type RobotStateMapLayer struct {
	Type     string                     `json:"type"`
	Metadata RobotStateMapLayerMetadata `json:"metaData"`
}

type RobotStateMapLayerMetadata struct {
	Area      *int    `json:"area"`
	SegmentId *string `json:"segmentId"`
	Active    *bool   `json:"active"`
	Name      *string `json:"name"`
}

type RobotState struct {
	Attributes []RobotStateAttribute `json:"attributes"`
	Map        RobotStateMap         `json:"map"`
}

type MapSegmentationCapabilityPutRequest struct {
	Action      string   `json:"action"`
	SegmentIds  []string `json:"segment_ids"`
	Iterations  *int     `json:"iterations,omitempty"`
	CustomOrder *bool    `json:"custom_order,omitempty"`
}

type BasicControlCapabilityRequest struct {
	Action string `json:"action"`
}
