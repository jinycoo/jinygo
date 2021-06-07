package signature

type Signer struct {
	S Sign
}

func (s *Signer) Sign(content, privateKey string) (sign string, err error) {
	return s.S(content, privateKey)
}
