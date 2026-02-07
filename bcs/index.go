package bcs

func FromBase64(v string) ([]byte, error) { return DecodeStr(v, EncodingBase64) }
func FromHex(v string) ([]byte, error)    { return DecodeStr(v, EncodingHex) }
func ToBase64(v []byte) (string, error)   { return EncodeStr(v, EncodingBase64) }
func ToHex(v []byte) (string, error)      { return EncodeStr(v, EncodingHex) }
