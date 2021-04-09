package empty

func NewUnknown(unknown []byte)Empty{
	return Empty{unknownFields: unknown}
}