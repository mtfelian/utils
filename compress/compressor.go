package compress

// Compressor умеет упаковать файлы в архив
type Compressor interface {
	Compress(fileOut string, pathIn ...string) error
}

// Decompressor умеет распаковать файлы из архива
type Decompressor interface {
	Decompress(fileIn string, pathOut string) error
}

type CompressorDecompressor interface {
	Compressor
	Decompressor
}
