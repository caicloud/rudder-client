package universal

func convertBoolToPointer(b bool) *bool {
	if b {
		return &b
	}
	return nil
}
