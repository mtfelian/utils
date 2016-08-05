package encrypt

// Signer умеет подписать файл цифровой подписью в формате DER, которая туда же и помещается
type Signer interface {
	SignDER(fileIn string, fileOut string) error
}

// Encryptor умеет зашифровать файл, на выходе используется формат DER
type Encryptor interface {
	Encrypt(fileIn string, fileOut string) error
}

// SignerEncryptor включает интерфейсы DERSigner и Encryptor
type SignerEncryptor interface {
	Signer
	Encryptor
}
