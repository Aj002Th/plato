package discovery

import "github.com/cloudwego/hertz/pkg/common/json"

// EndpointInfo etcd 中存储的 ipconfig endpoint 信息
type EndpointInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`

	// 可以灵活扩展
	// 现在记录的参数有: connect_num连接数, message_bytes传输字节数
	MetaData map[string]any `json:"meta"`
}

// 暂时先用 json 编解码, 可以据此扩展出其他编码类型

func UnMarshal(data []byte) (*EndpointInfo, error) {
	edi := &EndpointInfo{}
	err := json.Unmarshal(data, edi)
	if err != nil {
		return nil, err
	}
	return edi, nil
}

func (edi *EndpointInfo) Marshal() string {
	data, err := json.Marshal(edi)
	if err != nil {
		// 不应该会有问题
		panic(err)
	}
	return string(data)
}
