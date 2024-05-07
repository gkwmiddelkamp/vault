package views

type SecretDecodeView struct {
	Secret string `json:"secret"`
}

func NewSecretDecodeView(decryptedSecret string) SecretDecodeView {
	result := SecretDecodeView{
		Secret: decryptedSecret,
	}

	return result
}
