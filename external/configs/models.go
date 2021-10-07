package configs

import (
	"encoding/json"
	"github.com/dysnix/ai-scale-libs/external/enums"
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"

	tc "github.com/dysnix/ai-scale-libs/external/types_convertation"
)

type Base struct {
	IsDebugMode  bool `yaml:"debugMode" json:"debug_mode"`
	UseProfiling bool `yaml:"useProfiling" json:"use_profiling"`
}

type Informer struct {
	Resource string        `yaml:"resource" json:"resource" validate:"required"`
	Interval time.Duration `yaml:"interval" json:"interval" validate:"required,gt=0"`
}

type K8sCloudWatcher struct {
	CtxPath   string     `yaml:"kubeConfigPath" json:"kube_config_path" validate:"required,file"`
	Informers []Informer `yaml:"informers" json:"informers" validate:"required,gt=0"`
}

type GRPC struct {
	Enabled       bool
	UseReflection bool        `yaml:"useReflection" json:"use_reflection"`
	Compression   Compression `yaml:"compression" json:"compression"`
	Conn          *Connection `yaml:"connection" json:"connection" validate:"required"`
	Keepalive     *Keepalive  `yaml:"keepalive" json:"keepalive"`
}

type Compression struct {
	Enabled bool                  `yaml:"enabled" json:"enabled"`
	Type    enums.CompressionType `yaml:"type" json:"type"`
}

type Connection struct {
	Host            string        `yaml:"host" json:"host"`
	Port            uint16        `yaml:"port" json:"port" validate:"required,gt=0"`
	ReadBufferSize  uint          `yaml:"readBufferSize" json:"read_buffer_size" validate:"required,gte=4096"`
	WriteBufferSize uint          `yaml:"writeBufferSize" json:"write_buffer_size" validate:"required,gte=4096"`
	MaxMessageSize  uint          `yaml:"maxMessageSize" json:"max_message_size" validate:"required,gte=2048"`
	Insecure        bool          `yaml:"insecure" json:"insecure" validate:"required"`
	Timeout         time.Duration `yaml:"timeout" json:"timeout" validate:"gte=0"`
}

func (c *Connection) MarshalJSON() ([]byte, error) {
	type alias struct {
		Host            string  `yaml:"host" json:"host"`
		Port            uint16  `yaml:"port" json:"port"`
		ReadBufferSize  string  `yaml:"readBufferSize" json:"read_buffer_size"`
		WriteBufferSize string  `yaml:"writeBufferSize" json:"write_buffer_size"`
		MaxMessageSize  string  `yaml:"maxMessageSize" json:"max_message_size"`
		Insecure        bool    `yaml:"insecure" json:"insecure"`
		Timeout         *string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	}

	if c == nil {
		*c = Connection{}
	}

	return json.Marshal(alias{
		Host:            c.Host,
		Port:            c.Port,
		ReadBufferSize:  tc.BytesSize(float64(c.ReadBufferSize)),
		WriteBufferSize: tc.BytesSize(float64(c.WriteBufferSize)),
		MaxMessageSize:  tc.BytesSize(float64(c.MaxMessageSize)),
		Insecure:        c.Insecure,
		Timeout:         tc.String(ConvertDurationToStr(c.Timeout)),
	})
}

func (c *Connection) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		Host            string  `yaml:"host" json:"host"`
		Port            uint16  `yaml:"port" json:"port"`
		ReadBufferSize  string  `yaml:"readBufferSize" json:"read_buffer_size"`
		WriteBufferSize string  `yaml:"writeBufferSize" json:"write_buffer_size"`
		MaxMessageSize  string  `yaml:"maxMessageSize" json:"max_message_size"`
		Insecure        bool    `yaml:"insecure" json:"insecure"`
		Timeout         *string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if c == nil {
		*c = Connection{}
	}

	if tmp.Timeout != nil {
		c.Timeout, err = str2duration.ParseDuration(*tmp.Timeout)
		if err != nil {
			return err
		}
	}

	c.Host = tmp.Host
	c.Port = tmp.Port

	var tmpB int64
	if tmpB, err = tc.RAMInBytes(tmp.ReadBufferSize); err != nil {
		return err
	}

	c.ReadBufferSize = uint(tmpB)

	if tmpB, err = tc.RAMInBytes(tmp.WriteBufferSize); err != nil {
		return err
	}

	c.WriteBufferSize = uint(tmpB)

	if tmpB, err = tc.RAMInBytes(tmp.MaxMessageSize); err != nil {
		return err
	}

	c.MaxMessageSize = uint(tmpB)
	c.Insecure = tmp.Insecure

	return nil
}

func (c *Connection) MarshalYAML() (interface{}, error) {
	type alias struct {
		Host            string  `yaml:"host" json:"host"`
		Port            uint16  `yaml:"port" json:"port"`
		ReadBufferSize  string  `yaml:"readBufferSize" json:"read_buffer_size"`
		WriteBufferSize string  `yaml:"writeBufferSize" json:"write_buffer_size"`
		MaxMessageSize  string  `yaml:"maxMessageSize" json:"max_message_size"`
		Insecure        bool    `yaml:"insecure" json:"insecure"`
		Timeout         *string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	}

	if c == nil {
		*c = Connection{}
	}

	return alias{
		Host:            c.Host,
		Port:            c.Port,
		ReadBufferSize:  tc.BytesSize(float64(c.ReadBufferSize)),
		WriteBufferSize: tc.BytesSize(float64(c.WriteBufferSize)),
		MaxMessageSize:  tc.BytesSize(float64(c.MaxMessageSize)),
		Insecure:        c.Insecure,
		Timeout:         tc.String(ConvertDurationToStr(c.Timeout)),
	}, nil
}

func (c *Connection) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Host            string  `yaml:"host" json:"host"`
		Port            uint16  `yaml:"port" json:"port"`
		ReadBufferSize  string  `yaml:"readBufferSize" json:"read_buffer_size"`
		WriteBufferSize string  `yaml:"writeBufferSize" json:"write_buffer_size"`
		MaxMessageSize  string  `yaml:"maxMessageSize" json:"max_message_size"`
		Insecure        bool    `yaml:"insecure" json:"insecure"`
		Timeout         *string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	}
	var tmp alias
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	if c == nil {
		*c = Connection{}
	}

	if tmp.Timeout != nil {
		c.Timeout, err = str2duration.ParseDuration(*tmp.Timeout)
		if err != nil {
			return err
		}
	}

	c.Host = tmp.Host
	c.Port = tmp.Port

	var tmpB int64
	if tmpB, err = tc.RAMInBytes(tmp.ReadBufferSize); err != nil {
		return err
	}

	c.ReadBufferSize = uint(tmpB)

	if tmpB, err = tc.RAMInBytes(tmp.WriteBufferSize); err != nil {
		return err
	}

	c.WriteBufferSize = uint(tmpB)

	if tmpB, err = tc.RAMInBytes(tmp.MaxMessageSize); err != nil {
		return err
	}

	c.MaxMessageSize = uint(tmpB)
	c.Insecure = tmp.Insecure

	return nil
}

type Keepalive struct {
	Time              time.Duration      `yaml:"time" json:"time" validate:"required,gt=0"`
	Timeout           time.Duration      `yaml:"timeout" json:"timeout" validate:"required,gt=0"`
	EnforcementPolicy *EnforcementPolicy `yaml:"enforcementPolicy" json:"enforcement_policy"`
}

func (ka *Keepalive) MarshalJSON() ([]byte, error) {
	type alias struct {
		Time              string             `yaml:"time" json:"time"`
		Timeout           string             `yaml:"timeout" json:"timeout"`
		EnforcementPolicy *EnforcementPolicy `yaml:"enforcementPolicy" json:"enforcement_policy"`
	}

	if ka == nil {
		*ka = Keepalive{}
	}

	return json.Marshal(alias{
		Time:              ConvertDurationToStr(ka.Time),
		Timeout:           ConvertDurationToStr(ka.Timeout),
		EnforcementPolicy: ka.EnforcementPolicy,
	})
}

func (ka *Keepalive) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		Time              string             `yaml:"time" json:"time"`
		Timeout           string             `yaml:"timeout" json:"timeout"`
		EnforcementPolicy *EnforcementPolicy `yaml:"enforcementPolicy" json:"enforcement_policy"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if ka == nil {
		*ka = Keepalive{}
	}

	ka.Time, err = str2duration.ParseDuration(tmp.Time)
	if err != nil {
		return err
	}

	ka.Timeout, err = str2duration.ParseDuration(tmp.Timeout)
	if err != nil {
		return err
	}

	ka.EnforcementPolicy = tmp.EnforcementPolicy

	return nil
}

func (ka *Keepalive) MarshalYAML() (interface{}, error) {
	type alias struct {
		Time              string             `yaml:"time" json:"time"`
		Timeout           string             `yaml:"timeout" json:"timeout"`
		EnforcementPolicy *EnforcementPolicy `yaml:"enforcementPolicy" json:"enforcement_policy"`
	}

	if ka == nil {
		*ka = Keepalive{}
	}

	return alias{
		Time:              ConvertDurationToStr(ka.Time),
		Timeout:           ConvertDurationToStr(ka.Timeout),
		EnforcementPolicy: ka.EnforcementPolicy,
	}, nil
}

func (ka *Keepalive) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Time              string             `yaml:"time" json:"time"`
		Timeout           string             `yaml:"timeout" json:"timeout"`
		EnforcementPolicy *EnforcementPolicy `yaml:"enforcementPolicy" json:"enforcement_policy"`
	}
	var tmp alias
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	if ka == nil {
		*ka = Keepalive{}
	}

	ka.Time, err = str2duration.ParseDuration(tmp.Time)
	if err != nil {
		return err
	}

	ka.Timeout, err = str2duration.ParseDuration(tmp.Timeout)
	if err != nil {
		return err
	}

	ka.EnforcementPolicy = tmp.EnforcementPolicy

	return nil
}

type EnforcementPolicy struct {
	MinTime             time.Duration `yaml:"minTime" json:"min_time" validate:"required,gt=0"`
	PermitWithoutStream bool          `yaml:"permitWithoutStream" json:"permit_without_stream"`
}

func (ep *EnforcementPolicy) MarshalJSON() ([]byte, error) {
	type alias struct {
		MinTime             string `yaml:"minTime" json:"min_time"`
		PermitWithoutStream bool   `yaml:"permitWithoutStream" json:"permit_without_stream"`
	}

	if ep == nil {
		*ep = EnforcementPolicy{}
	}

	return json.Marshal(alias{
		MinTime:             ConvertDurationToStr(ep.MinTime),
		PermitWithoutStream: ep.PermitWithoutStream,
	})
}

func (ep *EnforcementPolicy) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		MinTime             string `yaml:"minTime" json:"min_time"`
		PermitWithoutStream bool   `yaml:"permitWithoutStream" json:"permit_without_stream"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if ep == nil {
		*ep = EnforcementPolicy{}
	}

	ep.PermitWithoutStream = tmp.PermitWithoutStream
	ep.MinTime, err = str2duration.ParseDuration(tmp.MinTime)

	return err
}

func (ep *EnforcementPolicy) MarshalYAML() (interface{}, error) {
	type alias struct {
		MinTime             string `yaml:"minTime" json:"min_time"`
		PermitWithoutStream bool   `yaml:"permitWithoutStream" json:"permit_without_stream"`
	}

	if ep == nil {
		*ep = EnforcementPolicy{}
	}

	return alias{
		MinTime:             ConvertDurationToStr(ep.MinTime),
		PermitWithoutStream: ep.PermitWithoutStream,
	}, nil
}

func (ep *EnforcementPolicy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		MinTime             string `yaml:"minTime" json:"min_time"`
		PermitWithoutStream bool   `yaml:"permitWithoutStream" json:"permit_without_stream"`
	}
	var tmp alias
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	if ep == nil {
		*ep = EnforcementPolicy{}
	}

	ep.PermitWithoutStream = tmp.PermitWithoutStream
	ep.MinTime, err = str2duration.ParseDuration(tmp.MinTime)

	return err
}
