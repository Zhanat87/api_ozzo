// Code generated by protoc-gen-go.
// source: google.golang.org/genproto/googleapis/api/serviceconfig/http.proto
// DO NOT EDIT!

package google_api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Defines the HTTP configuration for a service. It contains a list of
// [HttpRule][google.api.HttpRule], each specifying the mapping of an RPC method
// to one or more HTTP REST API methods.
type Http struct {
	// A list of HTTP configuration rules that apply to individual API methods.
	//
	// **NOTE:** All service configuration rules follow "last one wins" order.
	Rules []*HttpRule `protobuf:"bytes,1,rep,name=rules" json:"rules,omitempty"`
}

func (m *Http) Reset()                    { *m = Http{} }
func (m *Http) String() string            { return proto.CompactTextString(m) }
func (*Http) ProtoMessage()               {}
func (*Http) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{0} }

func (m *Http) GetRules() []*HttpRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

// `HttpRule` defines the mapping of an RPC method to one or more HTTP
// REST APIs.  The mapping determines what portions of the request
// message are populated from the path, query parameters, or body of
// the HTTP request.  The mapping is typically specified as an
// `google.api.http` annotation, see "google/api/annotations.proto"
// for details.
//
// The mapping consists of a field specifying the path template and
// method kind.  The path template can refer to fields in the request
// message, as in the example below which describes a REST GET
// operation on a resource collection of messages:
//
// ```proto
// service Messaging {
//   rpc GetMessage(GetMessageRequest) returns (Message) {
//     option (google.api.http).get = "/v1/messages/{message_id}/{sub.subfield}";
//   }
// }
// message GetMessageRequest {
//   message SubMessage {
//     string subfield = 1;
//   }
//   string message_id = 1; // mapped to the URL
//   SubMessage sub = 2;    // `sub.subfield` is url-mapped
// }
// message Message {
//   string text = 1; // content of the resource
// }
// ```
//
// This definition enables an automatic, bidrectional mapping of HTTP
// JSON to RPC. Example:
//
// HTTP | RPC
// -----|-----
// `GET /v1/messages/123456/foo`  | `GetMessage(message_id: "123456" sub: SubMessage(subfield: "foo"))`
//
// In general, not only fields but also field paths can be referenced
// from a path pattern. Fields mapped to the path pattern cannot be
// repeated and must have a primitive (non-message) type.
//
// Any fields in the request message which are not bound by the path
// pattern automatically become (optional) HTTP query
// parameters. Assume the following definition of the request message:
//
// ```proto
// message GetMessageRequest {
//   message SubMessage {
//     string subfield = 1;
//   }
//   string message_id = 1; // mapped to the URL
//   int64 revision = 2;    // becomes a parameter
//   SubMessage sub = 3;    // `sub.subfield` becomes a parameter
// }
// ```
//
// This enables a HTTP JSON to RPC mapping as below:
//
// HTTP | RPC
// -----|-----
// `GET /v1/messages/123456?revision=2&sub.subfield=foo` | `GetMessage(message_id: "123456" revision: 2 sub: SubMessage(subfield: "foo"))`
//
// Note that fields which are mapped to HTTP parameters must have a
// primitive type or a repeated primitive type. Message types are not
// allowed. In the case of a repeated type, the parameter can be
// repeated in the URL, as in `...?param=A&param=B`.
//
// For HTTP method kinds which allow a request body, the `body` field
// specifies the mapping. Consider a REST update method on the
// message resource collection:
//
// ```proto
// service Messaging {
//   rpc UpdateMessage(UpdateMessageRequest) returns (Message) {
//     option (google.api.http) = {
//       put: "/v1/messages/{message_id}"
//       body: "message"
//     };
//   }
// }
// message UpdateMessageRequest {
//   string message_id = 1; // mapped to the URL
//   Message message = 2;   // mapped to the body
// }
// ```
//
// The following HTTP JSON to RPC mapping is enabled, where the
// representation of the JSON in the request body is determined by
// protos JSON encoding:
//
// HTTP | RPC
// -----|-----
// `PUT /v1/messages/123456 { "text": "Hi!" }` | `UpdateMessage(message_id: "123456" message { text: "Hi!" })`
//
// The special name `*` can be used in the body mapping to define that
// every field not bound by the path template should be mapped to the
// request body.  This enables the following alternative definition of
// the update method:
//
// ```proto
// service Messaging {
//   rpc UpdateMessage(Message) returns (Message) {
//     option (google.api.http) = {
//       put: "/v1/messages/{message_id}"
//       body: "*"
//     };
//   }
// }
// message Message {
//   string message_id = 1;
//   string text = 2;
// }
// ```
//
// The following HTTP JSON to RPC mapping is enabled:
//
// HTTP | RPC
// -----|-----
// `PUT /v1/messages/123456 { "text": "Hi!" }` | `UpdateMessage(message_id: "123456" text: "Hi!")`
//
// Note that when using `*` in the body mapping, it is not possible to
// have HTTP parameters, as all fields not bound by the path end in
// the body. This makes this option more rarely used in practice of
// defining REST APIs. The common usage of `*` is in custom methods
// which don't use the URL at all for transferring data.
//
// It is possible to define multiple HTTP methods for one RPC by using
// the `additional_bindings` option. Example:
//
// ```proto
// service Messaging {
//   rpc GetMessage(GetMessageRequest) returns (Message) {
//     option (google.api.http) = {
//       get: "/v1/messages/{message_id}"
//       additional_bindings {
//         get: "/v1/users/{user_id}/messages/{message_id}"
//       }
//     };
//   }
// }
// message GetMessageRequest {
//   string message_id = 1;
//   string user_id = 2;
// }
// ```
//
// This enables the following two alternative HTTP JSON to RPC
// mappings:
//
// HTTP | RPC
// -----|-----
// `GET /v1/messages/123456` | `GetMessage(message_id: "123456")`
// `GET /v1/users/me/messages/123456` | `GetMessage(user_id: "me" message_id: "123456")`
//
// # Rules for HTTP mapping
//
// The rules for mapping HTTP path, query parameters, and body fields
// to the request message are as follows:
//
// 1. The `body` field specifies either `*` or a field path, or is
//    omitted. If omitted, it assumes there is no HTTP body.
// 2. Leaf fields (recursive expansion of nested messages in the
//    request) can be classified into three types:
//     (a) Matched in the URL template.
//     (b) Covered by body (if body is `*`, everything except (a) fields;
//         else everything under the body field)
//     (c) All other fields.
// 3. URL query parameters found in the HTTP request are mapped to (c) fields.
// 4. Any body sent with an HTTP request can contain only (b) fields.
//
// The syntax of the path template is as follows:
//
//     Template = "/" Segments [ Verb ] ;
//     Segments = Segment { "/" Segment } ;
//     Segment  = "*" | "**" | LITERAL | Variable ;
//     Variable = "{" FieldPath [ "=" Segments ] "}" ;
//     FieldPath = IDENT { "." IDENT } ;
//     Verb     = ":" LITERAL ;
//
// The syntax `*` matches a single path segment. It follows the semantics of
// [RFC 6570](https://tools.ietf.org/html/rfc6570) Section 3.2.2 Simple String
// Expansion.
//
// The syntax `**` matches zero or more path segments. It follows the semantics
// of [RFC 6570](https://tools.ietf.org/html/rfc6570) Section 3.2.3 Reserved
// Expansion.
//
// The syntax `LITERAL` matches literal text in the URL path.
//
// The syntax `Variable` matches the entire path as specified by its template;
// this nested template must not contain further variables. If a variable
// matches a single path segment, its template may be omitted, e.g. `{var}`
// is equivalent to `{var=*}`.
//
// NOTE: the field paths in variables and in the `body` must not refer to
// repeated fields or map fields.
//
// Use CustomHttpPattern to specify any HTTP method that is not included in the
// `pattern` field, such as HEAD, or "*" to leave the HTTP method unspecified for
// a given URL path rule. The wild-card rule is useful for services that provide
// content to Web (HTML) clients.
type HttpRule struct {
	// Selects methods to which this rule applies.
	//
	// Refer to [selector][google.api.DocumentationRule.selector] for syntax details.
	Selector string `protobuf:"bytes,1,opt,name=selector" json:"selector,omitempty"`
	// Determines the URL pattern is matched by this rules. This pattern can be
	// used with any of the {get|put|post|delete|patch} methods. A custom method
	// can be defined using the 'custom' field.
	//
	// Types that are valid to be assigned to Pattern:
	//	*HttpRule_Get
	//	*HttpRule_Put
	//	*HttpRule_Post
	//	*HttpRule_Delete
	//	*HttpRule_Patch
	//	*HttpRule_Custom
	Pattern isHttpRule_Pattern `protobuf_oneof:"pattern"`
	// The name of the request field whose value is mapped to the HTTP body, or
	// `*` for mapping all fields not captured by the path pattern to the HTTP
	// body. NOTE: the referred field must not be a repeated field and must be
	// present at the top-level of response message type.
	Body string `protobuf:"bytes,7,opt,name=body" json:"body,omitempty"`
	// Additional HTTP bindings for the selector. Nested bindings must
	// not contain an `additional_bindings` field themselves (that is,
	// the nesting may only be one level deep).
	AdditionalBindings []*HttpRule `protobuf:"bytes,11,rep,name=additional_bindings,json=additionalBindings" json:"additional_bindings,omitempty"`
}

func (m *HttpRule) Reset()                    { *m = HttpRule{} }
func (m *HttpRule) String() string            { return proto.CompactTextString(m) }
func (*HttpRule) ProtoMessage()               {}
func (*HttpRule) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{1} }

type isHttpRule_Pattern interface {
	isHttpRule_Pattern()
}

type HttpRule_Get struct {
	Get string `protobuf:"bytes,2,opt,name=get,oneof"`
}
type HttpRule_Put struct {
	Put string `protobuf:"bytes,3,opt,name=put,oneof"`
}
type HttpRule_Post struct {
	Post string `protobuf:"bytes,4,opt,name=post,oneof"`
}
type HttpRule_Delete struct {
	Delete string `protobuf:"bytes,5,opt,name=delete,oneof"`
}
type HttpRule_Patch struct {
	Patch string `protobuf:"bytes,6,opt,name=patch,oneof"`
}
type HttpRule_Custom struct {
	Custom *CustomHttpPattern `protobuf:"bytes,8,opt,name=custom,oneof"`
}

func (*HttpRule_Get) isHttpRule_Pattern()    {}
func (*HttpRule_Put) isHttpRule_Pattern()    {}
func (*HttpRule_Post) isHttpRule_Pattern()   {}
func (*HttpRule_Delete) isHttpRule_Pattern() {}
func (*HttpRule_Patch) isHttpRule_Pattern()  {}
func (*HttpRule_Custom) isHttpRule_Pattern() {}

func (m *HttpRule) GetPattern() isHttpRule_Pattern {
	if m != nil {
		return m.Pattern
	}
	return nil
}

func (m *HttpRule) GetGet() string {
	if x, ok := m.GetPattern().(*HttpRule_Get); ok {
		return x.Get
	}
	return ""
}

func (m *HttpRule) GetPut() string {
	if x, ok := m.GetPattern().(*HttpRule_Put); ok {
		return x.Put
	}
	return ""
}

func (m *HttpRule) GetPost() string {
	if x, ok := m.GetPattern().(*HttpRule_Post); ok {
		return x.Post
	}
	return ""
}

func (m *HttpRule) GetDelete() string {
	if x, ok := m.GetPattern().(*HttpRule_Delete); ok {
		return x.Delete
	}
	return ""
}

func (m *HttpRule) GetPatch() string {
	if x, ok := m.GetPattern().(*HttpRule_Patch); ok {
		return x.Patch
	}
	return ""
}

func (m *HttpRule) GetCustom() *CustomHttpPattern {
	if x, ok := m.GetPattern().(*HttpRule_Custom); ok {
		return x.Custom
	}
	return nil
}

func (m *HttpRule) GetAdditionalBindings() []*HttpRule {
	if m != nil {
		return m.AdditionalBindings
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*HttpRule) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _HttpRule_OneofMarshaler, _HttpRule_OneofUnmarshaler, _HttpRule_OneofSizer, []interface{}{
		(*HttpRule_Get)(nil),
		(*HttpRule_Put)(nil),
		(*HttpRule_Post)(nil),
		(*HttpRule_Delete)(nil),
		(*HttpRule_Patch)(nil),
		(*HttpRule_Custom)(nil),
	}
}

func _HttpRule_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*HttpRule)
	// pattern
	switch x := m.Pattern.(type) {
	case *HttpRule_Get:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Get)
	case *HttpRule_Put:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Put)
	case *HttpRule_Post:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Post)
	case *HttpRule_Delete:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Delete)
	case *HttpRule_Patch:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Patch)
	case *HttpRule_Custom:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Custom); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("HttpRule.Pattern has unexpected type %T", x)
	}
	return nil
}

func _HttpRule_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*HttpRule)
	switch tag {
	case 2: // pattern.get
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Pattern = &HttpRule_Get{x}
		return true, err
	case 3: // pattern.put
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Pattern = &HttpRule_Put{x}
		return true, err
	case 4: // pattern.post
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Pattern = &HttpRule_Post{x}
		return true, err
	case 5: // pattern.delete
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Pattern = &HttpRule_Delete{x}
		return true, err
	case 6: // pattern.patch
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Pattern = &HttpRule_Patch{x}
		return true, err
	case 8: // pattern.custom
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CustomHttpPattern)
		err := b.DecodeMessage(msg)
		m.Pattern = &HttpRule_Custom{msg}
		return true, err
	default:
		return false, nil
	}
}

func _HttpRule_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*HttpRule)
	// pattern
	switch x := m.Pattern.(type) {
	case *HttpRule_Get:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Get)))
		n += len(x.Get)
	case *HttpRule_Put:
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Put)))
		n += len(x.Put)
	case *HttpRule_Post:
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Post)))
		n += len(x.Post)
	case *HttpRule_Delete:
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Delete)))
		n += len(x.Delete)
	case *HttpRule_Patch:
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Patch)))
		n += len(x.Patch)
	case *HttpRule_Custom:
		s := proto.Size(x.Custom)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// A custom pattern is used for defining custom HTTP verb.
type CustomHttpPattern struct {
	// The name of this custom HTTP verb.
	Kind string `protobuf:"bytes,1,opt,name=kind" json:"kind,omitempty"`
	// The path matched by this custom verb.
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
}

func (m *CustomHttpPattern) Reset()                    { *m = CustomHttpPattern{} }
func (m *CustomHttpPattern) String() string            { return proto.CompactTextString(m) }
func (*CustomHttpPattern) ProtoMessage()               {}
func (*CustomHttpPattern) Descriptor() ([]byte, []int) { return fileDescriptor9, []int{2} }

func init() {
	proto.RegisterType((*Http)(nil), "google.api.Http")
	proto.RegisterType((*HttpRule)(nil), "google.api.HttpRule")
	proto.RegisterType((*CustomHttpPattern)(nil), "google.api.CustomHttpPattern")
}

func init() {
	proto.RegisterFile("google.golang.org/genproto/googleapis/api/serviceconfig/http.proto", fileDescriptor9)
}

var fileDescriptor9 = []byte{
	// 348 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x51, 0xcd, 0x4a, 0x2b, 0x31,
	0x14, 0xbe, 0xd3, 0x4e, 0xa7, 0xed, 0x29, 0x5c, 0xb8, 0xb9, 0x45, 0x82, 0x20, 0x48, 0x37, 0x16,
	0x17, 0x33, 0x50, 0x17, 0x2e, 0x5c, 0x39, 0x22, 0xd6, 0x5d, 0x99, 0x17, 0x90, 0xe9, 0x4c, 0x4c,
	0x83, 0xd3, 0x24, 0x4c, 0x4e, 0x05, 0x5f, 0xc7, 0x77, 0xf0, 0xdd, 0x5c, 0x9a, 0x64, 0x52, 0x5b,
	0x10, 0xdc, 0x84, 0xf3, 0xfd, 0x9c, 0x9f, 0x9c, 0x03, 0x39, 0x57, 0x8a, 0x37, 0x2c, 0xe5, 0xaa,
	0x29, 0x25, 0x4f, 0x55, 0xcb, 0x33, 0xce, 0xa4, 0x6e, 0x15, 0xaa, 0xac, 0x93, 0x4a, 0x2d, 0x4c,
	0x66, 0x9f, 0xcc, 0xb0, 0xf6, 0x55, 0x54, 0xac, 0x52, 0xf2, 0x59, 0xf0, 0x6c, 0x83, 0xa8, 0x53,
	0xef, 0x23, 0x10, 0x6a, 0x58, 0xd3, 0x6c, 0x01, 0xf1, 0xd2, 0x2a, 0xe4, 0x12, 0x06, 0xed, 0xae,
	0x61, 0x86, 0x46, 0xe7, 0xfd, 0xf9, 0x64, 0x31, 0x4d, 0x0f, 0x9e, 0xd4, 0x19, 0x0a, 0x2b, 0x16,
	0x9d, 0x65, 0xf6, 0xd1, 0x83, 0xd1, 0x9e, 0x23, 0xa7, 0x30, 0x32, 0xac, 0x61, 0x15, 0xaa, 0xd6,
	0xe6, 0x46, 0xf3, 0x71, 0xf1, 0x8d, 0x09, 0x81, 0x3e, 0x67, 0x48, 0x7b, 0x8e, 0x5e, 0xfe, 0x29,
	0x1c, 0x70, 0x9c, 0xde, 0x21, 0xed, 0xef, 0x39, 0x0b, 0xc8, 0x14, 0x62, 0xad, 0x0c, 0xd2, 0x38,
	0x90, 0x1e, 0x11, 0x0a, 0x49, 0x6d, 0x2b, 0x21, 0xa3, 0x83, 0xc0, 0x07, 0x4c, 0x4e, 0x60, 0xa0,
	0x4b, 0xac, 0x36, 0x34, 0x09, 0x42, 0x07, 0xc9, 0x35, 0x24, 0xd5, 0xce, 0xa0, 0xda, 0xd2, 0x91,
	0x15, 0x26, 0x8b, 0xb3, 0xe3, 0x5f, 0xdc, 0x79, 0xc5, 0xcd, 0xbd, 0x2a, 0x11, 0x59, 0x2b, 0x5d,
	0xc1, 0xce, 0x6e, 0x87, 0x8a, 0xd7, 0xaa, 0x7e, 0xa3, 0x43, 0xff, 0x01, 0x1f, 0x93, 0x7b, 0xf8,
	0x5f, 0xd6, 0xb5, 0x40, 0xa1, 0x64, 0xd9, 0x3c, 0xad, 0x85, 0xac, 0x85, 0xe4, 0x86, 0x4e, 0x7e,
	0xd9, 0x0f, 0x39, 0x24, 0xe4, 0xc1, 0x9f, 0x8f, 0x61, 0xa8, 0xbb, 0x7e, 0xb3, 0x1b, 0xf8, 0xf7,
	0x63, 0x08, 0xd7, 0xfa, 0xc5, 0x7a, 0xc3, 0xee, 0x7c, 0xec, 0x38, 0x9b, 0xb3, 0xe9, 0x16, 0x57,
	0xf8, 0x38, 0xbf, 0x80, 0xbf, 0x95, 0xda, 0x1e, 0xb5, 0xcd, 0xc7, 0xbe, 0x8c, 0xbb, 0xe8, 0x2a,
	0xfa, 0x8c, 0xa2, 0xf7, 0x5e, 0xfc, 0x70, 0xbb, 0x7a, 0x5c, 0x27, 0xfe, 0xc8, 0x57, 0x5f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x32, 0x48, 0x5c, 0x87, 0x2a, 0x02, 0x00, 0x00,
}
