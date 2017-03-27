package zip

import (
	"fmt"
	"github.com/mihteh/utils"
	"os"
	"os/exec"
)

const (
	zipPath   = "zip"
	unzipPath = "unzip"
)

// Compress создаёт архив с путём и именем fileOut из файлов и папок, переданных в pathIn
func Compress(fileOut string, pathIn ...string) error {
	if utils.FileExists(fileOut) {
		if err := os.Remove(fileOut); err != nil {
			return err
		}
	}

	cmdParams := []string{"-j", fileOut}
	cmdParams = append(cmdParams, pathIn...)

	cmd := exec.Command(zipPath, cmdParams...)
	output, err := cmd.CombinedOutput()
	fmt.Printf("\n### OUTPUT COMPRESS: %v\n", string(output))
	if err != nil {
		fmt.Println("Ошибка выполнения команды zip: ", err.Error())
		return err
	}

	return nil
}

// Decompress извлекает из архива с путём и именем fileIn содержимое и помещает его в pathOut
func Decompress(fileIn string, pathOut string) error {
	cmdParams := []string{
		fileIn,
		"-d",
		pathOut,
	}
	cmd := exec.Command(unzipPath, cmdParams...)
	output, err := cmd.CombinedOutput()
	fmt.Printf("\n### OUTPUT DECOMPRESS: %v\n", string(output))
	if err != nil {
		fmt.Println("Ошибка выполнения команды unzip: ", err.Error())
		return err
	}

	return nil
}
