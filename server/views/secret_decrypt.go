package views

type SecretDecryptView struct {
	Secret string `json:"secret"`
}

func NewSecretDecodeView(decryptedSecret string) SecretDecryptView {
	result := SecretDecryptView{
		Secret: decryptedSecret,
	}

	return result
}
