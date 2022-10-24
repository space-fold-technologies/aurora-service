package security

// EncryptionParameters :
type EncryptionParameters struct {
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	Parallelism uint8  `yaml:"parallelism"`
	SaltLength  uint32 `yaml:"salt-length"`
	KeyLength   uint32 `yaml:"key-length"`
}

// DecodeDetails :
type DecodedDetails struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
	Hash        []byte
	Salt        []byte
}
