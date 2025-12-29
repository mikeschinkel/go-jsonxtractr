package jsonxtractr

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"errors"
	"os"

	"github.com/mikeschinkel/go-dt"
)

var (
	ErrParsingInput     = errors.New("error parsing input")
	ErrGeneratingOutput = errors.New("error generating output")
	ErrReadingInputFile = errors.New("error reading input file")
)

func Load(path dt.Filepath, data any) (err error) {
	var raw []byte
	if path == "" {
		err = NewErr(dt.ErrFileNotExists)
		goto end
	}

	raw, err = path.ReadFile()
	if err != nil {
		err = NewErr(ErrReadingInputFile, path.ErrKV(), err)
		goto end
	}

	err = jsonv2.Unmarshal(raw, &data)
	if err != nil {
		err = NewErr(ErrParsingInput, path.ErrKV(), err)
		goto end
	}

end:
	return err
}

func Save(path dt.Filepath, data any) (err error) {
	var raw []byte
	var dir dt.DirPath

	dir, err = path.Dir().Expand()
	if err != nil {
		goto end
	}
	err = dir.EnsureExists(os.ModePerm)
	if err != nil {
		goto end
	}

	raw, err = jsonv2.Marshal(data, jsontext.WithIndent("\t"))
	if err != nil {
		err = NewErr(ErrGeneratingOutput, err)
		goto end
	}
	err = path.WriteFile(raw, 0o644)
	if err != nil {
		err = NewErr(dt.ErrFailedToSaveFile, err)
		goto end
	}

end:
	if err != nil {
		err = WithErr(err, path.ErrKV())
	}
	return err
}
