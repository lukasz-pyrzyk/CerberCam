// Code generated by protoc-gen-go.
// source: response.proto
// DO NOT EDIT!

package main

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Response struct {
	Email            *string  `protobuf:"bytes,1,req" json:"Email,omitempty"`
	Label            *string  `protobuf:"bytes,2,req" json:"Label,omitempty"`
	Probability      *float32 `protobuf:"fixed32,3,req" json:"Probability,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}

func (m *Response) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *Response) GetLabel() string {
	if m != nil && m.Label != nil {
		return *m.Label
	}
	return ""
}

func (m *Response) GetProbability() float32 {
	if m != nil && m.Probability != nil {
		return *m.Probability
	}
	return 0
}

func init() {
}
