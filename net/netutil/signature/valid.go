package signature

func VerifySign(content, sign string, pubKeys []string) (err error) {
	verifier := NewVerifier(nil)
	for _, v := range pubKeys {
		err = verifier.Verify(content, sign, v)
		// 验签成功，跳出循环
		if err == nil {
			break
		}
	}

	return
}
