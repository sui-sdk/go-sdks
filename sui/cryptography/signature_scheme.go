package cryptography

type SignatureScheme string

const (
	SchemeED25519   SignatureScheme = "ED25519"
	SchemeSecp256k1 SignatureScheme = "Secp256k1"
	SchemeSecp256r1 SignatureScheme = "Secp256r1"
	SchemeMultiSig  SignatureScheme = "MultiSig"
	SchemeZkLogin   SignatureScheme = "ZkLogin"
	SchemePasskey   SignatureScheme = "Passkey"
)

var SignatureSchemeToFlag = map[SignatureScheme]byte{
	SchemeED25519:   0x00,
	SchemeSecp256k1: 0x01,
	SchemeSecp256r1: 0x02,
	SchemeMultiSig:  0x03,
	SchemeZkLogin:   0x05,
	SchemePasskey:   0x06,
}

var SignatureSchemeToSize = map[SignatureScheme]int{
	SchemeED25519:   32,
	SchemeSecp256k1: 33,
	SchemeSecp256r1: 33,
	SchemePasskey:   33,
}

var SignatureFlagToScheme = map[byte]SignatureScheme{
	0x00: SchemeED25519,
	0x01: SchemeSecp256k1,
	0x02: SchemeSecp256r1,
	0x03: SchemeMultiSig,
	0x05: SchemeZkLogin,
	0x06: SchemePasskey,
}
