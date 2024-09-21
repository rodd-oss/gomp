// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.1
// source: events.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Direction int32

const (
	Direction_left  Direction = 0
	Direction_right Direction = 1
	Direction_up    Direction = 2
	Direction_down  Direction = 3
)

// Enum value maps for Direction.
var (
	Direction_name = map[int32]string{
		0: "left",
		1: "right",
		2: "up",
		3: "down",
	}
	Direction_value = map[string]int32{
		"left":  0,
		"right": 1,
		"up":    2,
		"down":  3,
	}
)

func (x Direction) Enum() *Direction {
	p := new(Direction)
	*p = x
	return p
}

func (x Direction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Direction) Descriptor() protoreflect.EnumDescriptor {
	return file_events_proto_enumTypes[0].Descriptor()
}

func (Direction) Type() protoreflect.EnumType {
	return &file_events_proto_enumTypes[0]
}

func (x Direction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Direction.Descriptor instead.
func (Direction) EnumDescriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0}
}

type Event_Type int32

const (
	Event_type_init    Event_Type = 0
	Event_type_connect Event_Type = 1
	Event_type_exit    Event_Type = 2
	Event_type_idle    Event_Type = 3
	Event_type_move    Event_Type = 4
	Event_type_empty   Event_Type = 5
)

// Enum value maps for Event_Type.
var (
	Event_Type_name = map[int32]string{
		0: "type_init",
		1: "type_connect",
		2: "type_exit",
		3: "type_idle",
		4: "type_move",
		5: "type_empty",
	}
	Event_Type_value = map[string]int32{
		"type_init":    0,
		"type_connect": 1,
		"type_exit":    2,
		"type_idle":    3,
		"type_move":    4,
		"type_empty":   5,
	}
)

func (x Event_Type) Enum() *Event_Type {
	p := new(Event_Type)
	*p = x
	return p
}

func (x Event_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Event_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_events_proto_enumTypes[1].Descriptor()
}

func (Event_Type) Type() protoreflect.EnumType {
	return &file_events_proto_enumTypes[1]
}

func (x Event_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Event_Type.Descriptor instead.
func (Event_Type) EnumDescriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{1, 0}
}

type Unit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	X         float64   `protobuf:"fixed64,2,opt,name=x,proto3" json:"x,omitempty"`
	Y         float64   `protobuf:"fixed64,3,opt,name=y,proto3" json:"y,omitempty"`
	Frame     int32     `protobuf:"varint,4,opt,name=frame,proto3" json:"frame,omitempty"`
	Skin      string    `protobuf:"bytes,5,opt,name=skin,proto3" json:"skin,omitempty"`
	Action    string    `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	Speed     float64   `protobuf:"fixed64,7,opt,name=speed,proto3" json:"speed,omitempty"`
	Direction Direction `protobuf:"varint,8,opt,name=direction,proto3,enum=Direction" json:"direction,omitempty"`
	Side      Direction `protobuf:"varint,9,opt,name=side,proto3,enum=Direction" json:"side,omitempty"`
}

func (x *Unit) Reset() {
	*x = Unit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unit) ProtoMessage() {}

func (x *Unit) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unit.ProtoReflect.Descriptor instead.
func (*Unit) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0}
}

func (x *Unit) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Unit) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Unit) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Unit) GetFrame() int32 {
	if x != nil {
		return x.Frame
	}
	return 0
}

func (x *Unit) GetSkin() string {
	if x != nil {
		return x.Skin
	}
	return ""
}

func (x *Unit) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *Unit) GetSpeed() float64 {
	if x != nil {
		return x.Speed
	}
	return 0
}

func (x *Unit) GetDirection() Direction {
	if x != nil {
		return x.Direction
	}
	return Direction_left
}

func (x *Unit) GetSide() Direction {
	if x != nil {
		return x.Side
	}
	return Direction_left
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type Event_Type `protobuf:"varint,1,opt,name=type,proto3,enum=Event_Type" json:"type,omitempty"`
	// Types that are assignable to Data:
	//
	//	*Event_Init
	//	*Event_Connect
	//	*Event_Exit
	//	*Event_Idle
	//	*Event_Move
	Data isEvent_Data `protobuf_oneof:"data"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{1}
}

func (x *Event) GetType() Event_Type {
	if x != nil {
		return x.Type
	}
	return Event_type_init
}

func (m *Event) GetData() isEvent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Event) GetInit() *EventInit {
	if x, ok := x.GetData().(*Event_Init); ok {
		return x.Init
	}
	return nil
}

func (x *Event) GetConnect() *EventConnect {
	if x, ok := x.GetData().(*Event_Connect); ok {
		return x.Connect
	}
	return nil
}

func (x *Event) GetExit() *EventExit {
	if x, ok := x.GetData().(*Event_Exit); ok {
		return x.Exit
	}
	return nil
}

func (x *Event) GetIdle() *EventIdle {
	if x, ok := x.GetData().(*Event_Idle); ok {
		return x.Idle
	}
	return nil
}

func (x *Event) GetMove() *EventMove {
	if x, ok := x.GetData().(*Event_Move); ok {
		return x.Move
	}
	return nil
}

type isEvent_Data interface {
	isEvent_Data()
}

type Event_Init struct {
	Init *EventInit `protobuf:"bytes,2,opt,name=init,proto3,oneof"`
}

type Event_Connect struct {
	Connect *EventConnect `protobuf:"bytes,3,opt,name=connect,proto3,oneof"`
}

type Event_Exit struct {
	Exit *EventExit `protobuf:"bytes,4,opt,name=exit,proto3,oneof"`
}

type Event_Idle struct {
	Idle *EventIdle `protobuf:"bytes,5,opt,name=idle,proto3,oneof"`
}

type Event_Move struct {
	Move *EventMove `protobuf:"bytes,6,opt,name=move,proto3,oneof"`
}

func (*Event_Init) isEvent_Data() {}

func (*Event_Connect) isEvent_Data() {}

func (*Event_Exit) isEvent_Data() {}

func (*Event_Idle) isEvent_Data() {}

func (*Event_Move) isEvent_Data() {}

type EventInit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string           `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Units    map[string]*Unit `protobuf:"bytes,2,rep,name=units,proto3" json:"units,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EventInit) Reset() {
	*x = EventInit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventInit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventInit) ProtoMessage() {}

func (x *EventInit) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventInit.ProtoReflect.Descriptor instead.
func (*EventInit) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{2}
}

func (x *EventInit) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *EventInit) GetUnits() map[string]*Unit {
	if x != nil {
		return x.Units
	}
	return nil
}

type EventConnect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unit *Unit `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
}

func (x *EventConnect) Reset() {
	*x = EventConnect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventConnect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventConnect) ProtoMessage() {}

func (x *EventConnect) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventConnect.ProtoReflect.Descriptor instead.
func (*EventConnect) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{3}
}

func (x *EventConnect) GetUnit() *Unit {
	if x != nil {
		return x.Unit
	}
	return nil
}

type EventExit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
}

func (x *EventExit) Reset() {
	*x = EventExit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventExit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventExit) ProtoMessage() {}

func (x *EventExit) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventExit.ProtoReflect.Descriptor instead.
func (*EventExit) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{4}
}

func (x *EventExit) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

type EventIdle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
}

func (x *EventIdle) Reset() {
	*x = EventIdle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventIdle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventIdle) ProtoMessage() {}

func (x *EventIdle) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventIdle.ProtoReflect.Descriptor instead.
func (*EventIdle) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{5}
}

func (x *EventIdle) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

type EventMove struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId  string    `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Direction Direction `protobuf:"varint,2,opt,name=direction,proto3,enum=Direction" json:"direction,omitempty"`
}

func (x *EventMove) Reset() {
	*x = EventMove{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventMove) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventMove) ProtoMessage() {}

func (x *EventMove) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventMove.ProtoReflect.Descriptor instead.
func (*EventMove) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{6}
}

func (x *EventMove) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *EventMove) GetDirection() Direction {
	if x != nil {
		return x.Direction
	}
	return Direction_left
}

var File_events_proto protoreflect.FileDescriptor

var file_events_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd4,
	0x01, 0x0a, 0x04, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x01, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6b, 0x69,
	0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6b, 0x69, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x12, 0x28, 0x0a, 0x09, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a,
	0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x04, 0x73, 0x69, 0x64, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x04, 0x73, 0x69, 0x64, 0x65, 0x22, 0xc9, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x20, 0x0a, 0x04, 0x69, 0x6e, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x69, 0x74, 0x48, 0x00, 0x52, 0x04, 0x69, 0x6e,
	0x69, 0x74, 0x12, 0x29, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x20, 0x0a,
	0x04, 0x65, 0x78, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x45, 0x78, 0x69, 0x74, 0x48, 0x00, 0x52, 0x04, 0x65, 0x78, 0x69, 0x74, 0x12,
	0x20, 0x0a, 0x04, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x04, 0x69, 0x64, 0x6c,
	0x65, 0x12, 0x20, 0x0a, 0x04, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0a, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x76, 0x65, 0x48, 0x00, 0x52, 0x04, 0x6d,
	0x6f, 0x76, 0x65, 0x22, 0x64, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x74,
	0x79, 0x70, 0x65, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x74, 0x79,
	0x70, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09,
	0x74, 0x79, 0x70, 0x65, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x74,
	0x79, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x6c, 0x65, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x74, 0x79,
	0x70, 0x65, 0x5f, 0x6d, 0x6f, 0x76, 0x65, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x74, 0x79, 0x70,
	0x65, 0x5f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x10, 0x05, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x96, 0x01, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x69, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x05,
	0x75, 0x6e, 0x69, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x49, 0x6e, 0x69, 0x74, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x1a, 0x3f, 0x0a, 0x0a, 0x55, 0x6e, 0x69,
	0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x29, 0x0a, 0x0c, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x19, 0x0a, 0x04, 0x75, 0x6e,
	0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x52,
	0x04, 0x75, 0x6e, 0x69, 0x74, 0x22, 0x28, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x45, 0x78,
	0x69, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x28, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x09, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x4d, 0x6f, 0x76, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x32, 0x0a,
	0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x08, 0x0a, 0x04, 0x6c, 0x65,
	0x66, 0x74, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x72, 0x69, 0x67, 0x68, 0x74, 0x10, 0x01, 0x12,
	0x06, 0x0a, 0x02, 0x75, 0x70, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x64, 0x6f, 0x77, 0x6e, 0x10,
	0x03, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_events_proto_rawDescOnce sync.Once
	file_events_proto_rawDescData = file_events_proto_rawDesc
)

func file_events_proto_rawDescGZIP() []byte {
	file_events_proto_rawDescOnce.Do(func() {
		file_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_proto_rawDescData)
	})
	return file_events_proto_rawDescData
}

var file_events_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_events_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_events_proto_goTypes = []any{
	(Direction)(0),       // 0: Direction
	(Event_Type)(0),      // 1: Event.Type
	(*Unit)(nil),         // 2: Unit
	(*Event)(nil),        // 3: Event
	(*EventInit)(nil),    // 4: EventInit
	(*EventConnect)(nil), // 5: EventConnect
	(*EventExit)(nil),    // 6: EventExit
	(*EventIdle)(nil),    // 7: EventIdle
	(*EventMove)(nil),    // 8: EventMove
	nil,                  // 9: EventInit.UnitsEntry
}
var file_events_proto_depIdxs = []int32{
	0,  // 0: Unit.direction:type_name -> Direction
	0,  // 1: Unit.side:type_name -> Direction
	1,  // 2: Event.type:type_name -> Event.Type
	4,  // 3: Event.init:type_name -> EventInit
	5,  // 4: Event.connect:type_name -> EventConnect
	6,  // 5: Event.exit:type_name -> EventExit
	7,  // 6: Event.idle:type_name -> EventIdle
	8,  // 7: Event.move:type_name -> EventMove
	9,  // 8: EventInit.units:type_name -> EventInit.UnitsEntry
	2,  // 9: EventConnect.unit:type_name -> Unit
	0,  // 10: EventMove.direction:type_name -> Direction
	2,  // 11: EventInit.UnitsEntry.value:type_name -> Unit
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_events_proto_init() }
func file_events_proto_init() {
	if File_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_events_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Unit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_events_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_events_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EventInit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_events_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*EventConnect); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_events_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EventExit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_events_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*EventIdle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_events_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*EventMove); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_events_proto_msgTypes[1].OneofWrappers = []any{
		(*Event_Init)(nil),
		(*Event_Connect)(nil),
		(*Event_Exit)(nil),
		(*Event_Idle)(nil),
		(*Event_Move)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_events_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_proto_goTypes,
		DependencyIndexes: file_events_proto_depIdxs,
		EnumInfos:         file_events_proto_enumTypes,
		MessageInfos:      file_events_proto_msgTypes,
	}.Build()
	File_events_proto = out.File
	file_events_proto_rawDesc = nil
	file_events_proto_goTypes = nil
	file_events_proto_depIdxs = nil
}
