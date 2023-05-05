package ovms_capi

type LogLevel string
type OVMS_DataType int
type OVMS_BufferType int

const (
	OVMS_LOG_TRACE   LogLevel = "TRACE"
	OVMS_LOG_DEBUG            = "DEBUG"
	OVMS_LOG_INFO             = "INFO"
	OVMS_LOG_WARNING          = "WARNING"
	OVMS_LOG_ERROR            = "ERROR"
)

const (
	OVMS_DATATYPE_BF16 OVMS_DataType = iota
	OVMS_DATATYPE_FP64
	OVMS_DATATYPE_FP32
	OVMS_DATATYPE_FP16
	OVMS_DATATYPE_I64
	OVMS_DATATYPE_I32
	OVMS_DATATYPE_I16
	OVMS_DATATYPE_I8
	OVMS_DATATYPE_I4
	OVMS_DATATYPE_U64
	OVMS_DATATYPE_U32
	OVMS_DATATYPE_U16
	OVMS_DATATYPE_U8
	OVMS_DATATYPE_U4
	OVMS_DATATYPE_U1
	OVMS_DATATYPE_BOOL
	OVMS_DATATYPE_CUSTOM
	OVMS_DATATYPE_UNDEFINED
	OVMS_DATATYPE_DYNAMIC
	OVMS_DATATYPE_MIXED
	OVMS_DATATYPE_Q78
	OVMS_DATATYPE_BIN
	OVMS_DATATYPE_END
)

const (
	OVMS_BUFFERTYPE_CPU OVMS_BufferType = iota
	OVMS_BUFFERTYPE_CPU_PINNED
	OVMS_BUFFERTYPE_GPU
	OVMS_BUFFERTYPE_HDDL
)

func (level LogLevel) LogLevelInt() int {
	switch level {
	case OVMS_LOG_TRACE:
		return 0
	case OVMS_LOG_DEBUG:
		return 1
	case OVMS_LOG_INFO:
		return 2
	case OVMS_LOG_WARNING:
		return 3
	case OVMS_LOG_ERROR:
		return 4
	default:
		return 224
	}
}