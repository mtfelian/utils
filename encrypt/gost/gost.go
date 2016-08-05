package gost

import (
	"fmt"
	"github.com/mihteh/utils"
	"github.com/mihteh/utils/encrypt"
	"os"
	"os/exec"
)

type SSLParams struct {
	OpenSSLPath         string // путь к исполняемому файлу openSSL
	OurCertFilePath     string // путь к нашему сертификату
	ForeignCertFilePath string // путь к чужому сертификату
	OurPrivateKey       string // путь к нашему приватному ключу
}

type gostSSL struct {
	params SSLParams
}

func (s gostSSL) getParams() SSLParams {
	return s.params
}

// NewGostSigner создаёт новый объект для подписи
func NewGostSigner(params SSLParams) encrypt.Signer {
	if !utils.FileExists(params.OpenSSLPath) {
		fmt.Println("Не удаётся получить доступ к openSSL по заданному пути.")
		os.Exit(1)
	}

	return &gostSSL{
		params: params,
	}
}

// SignDER подписывает файл с путём fileIn в формате DER
// содержимое вместе с цифровой подписью записывается в файл с путём fileOut
func (s gostSSL) SignDER(fileIn string, fileOut string) error {
	// Пример командной строки:
	// /gost-ssl/bin/openssl smime -sign -nodetach -signer certs/kakunin/kakunin.cer -inkey private.pem \
	// -engine gost -gost89 -binary -noattr -outform DER -in test.xml -out test.xml.sgn
	p := s.getParams()
	cmdParams := []string{
		"smime",
		"-sign",
		"-nodetach",
		"-signer",
		p.OurCertFilePath,
		"-inkey",
		p.OurPrivateKey,
		"-engine",
		"gost",
		"-gost89",
		"-binary",
		"-noattr",
		"-outform",
		"DER",
		"-in",
		fileIn,
		"-out",
		fileOut,
	}

	cmd := exec.Command(s.getParams().OpenSSLPath, cmdParams...)
	// в первом параметре вывод команды
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return err
	}

	return nil
}
