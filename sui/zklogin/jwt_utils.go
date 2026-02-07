package zklogin

func DecodeJWT(jwt string) (map[string]any, error) {
	v, err := JWTDecode(jwt, DecodeOptions{})
	if err != nil {
		return nil, err
	}
	return v.(map[string]any), nil
}

func ExtractClaimValue[R any](claim map[string]any, claimName string) (R, bool) {
	v, ok := claim[claimName]
	if !ok {
		var z R
		return z, false
	}
	r, ok := v.(R)
	if !ok {
		var z R
		return z, false
	}
	return r, true
}
