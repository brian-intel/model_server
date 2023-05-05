package ovms_capi

type OVMSOutput struct {
	Name       string
	DataType   OVMS_DataType
	Shape      []int64
	DimCount   int
	BufferType OVMS_BufferType
	DeviceId   int
	Data       []byte
	ByteSize   int
}