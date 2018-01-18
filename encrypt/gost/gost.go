package gost

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mtfelian/utils"
	"github.com/mtfelian/utils/encrypt"
)

type SSLParams struct {
	OpenSSLPath         string // путь к исполняемому файлу openSSL
	OurCertFilePath     string // путь к нашему сертификату
	ForeignCertFilePath string // путь к чужому сертификату
	OurPrivateKey       string // путь к нашему приватному ключу
}

type gostSSL struct{ params SSLParams }

func (s gostSSL) getParams() SSLParams { return s.params }

// NewGostSignerEncryptor создаёт новый объект для подписи
func NewGostSignerEncryptor(params SSLParams) encrypt.SignerEncryptor {
	if !utils.FileExists(params.OpenSSLPath) {
		fmt.Println("Не удаётся получить доступ к openSSL по заданному пути.")
		os.Exit(1)
	}
	return &gostSSL{params: params}
}

// SignDER подписывает файл с путём fileIn в формате DER
// содержимое вместе с цифровой подписью записывается в файл с путём fileOut
func (s gostSSL) SignDER(fileIn string, fileOut string) error {
	// Пример командной строки:
	// /gost-ssl/bin/openssl smime -sign -nodetach -signer certs/dirname/dirname.cer -inkey private.pem \
	// -engine gost -gost89 -binary -noattr -outform DER -in test.xml -out test.xml.sgn
	p := s.getParams()
	cmdParams := []string{
		"smime", "-sign", "-nodetach", "-signer", p.OurCertFilePath, "-inkey", p.OurPrivateKey,
		"-engine", "gost", "-gost89", "-binary", "-noattr", "-outform", "DER",
		"-in", fileIn,
		"-out", fileOut,
	}

	cmd := exec.Command(s.getParams().OpenSSLPath, cmdParams...)
	// в первом параметре вывод команды
	_, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	return err
}

// Encrypt шифрует файл с путём fileIn в формат DER
// выход - зашифрованный записывается в путь fileOut
func (s gostSSL) Encrypt(fileIn string, fileOut string) error {
	// Пример командной строки:
	// /gost-ssl/bin/openssl smime -encrypt -engine gost -gost89 -in test.xml.sgn -binary -outform der
	// -out test.xml.sgn.enc certs/equifax1617/Боевой\ сервер/Prod_Equifax_2016-2017.cer
	p := s.getParams()
	cmdParams := []string{
		"smime", "-encrypt", "-engine", "gost", "-gost89", "-in", fileIn, "-binary", "-outform", "der",
		"-out", fileOut, p.ForeignCertFilePath,
	}

	cmd := exec.Command(s.getParams().OpenSSLPath, cmdParams...)
	// в первом параметре вывод команды
	_, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	return err
}
