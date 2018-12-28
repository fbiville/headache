/*
 * Copyright 2018 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fs

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type FileSystem struct {
	FileWriter FileWriter
	FileReader FileReader
}

func DefaultFileSystem() FileSystem {
	return FileSystem{
		FileWriter: &OsFileWriter{},
		FileReader: &OsFileReader{},
	}
}

func (fi *FileSystem) IsFile(path string) bool {
	info, err := fi.FileReader.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

type FileWriter interface {
	Open(path string, mask int, permissions os.FileMode) (File, error)
	Write(path string, contents string, permissions os.FileMode) error
}

type OsFileWriter struct{}

func (*OsFileWriter) Open(name string, mask int, permissions os.FileMode) (File, error) {
	file, err := os.OpenFile(name, mask, permissions)
	return &OsFile{
		File: file,
	}, err
}

func (*OsFileWriter) Write(path string, contents string, permissions os.FileMode) error {
	return ioutil.WriteFile(path, []byte(contents), permissions)
}

type File interface {
	Write([]byte) error
	Close() error
}

type OsFile struct {
	File *os.File
}

func (of *OsFile) Write(contents []byte) error {
	_, err := of.File.Write(contents)
	return err
}

func (of *OsFile) Close() error {
	return of.File.Close()
}

type FileReader interface {
	http.FileSystem
	Read(path string) ([]byte, error)
	Stat(path string) (os.FileInfo, error)
}

type OsFileReader struct{}

func (*OsFileReader) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
func (*OsFileReader) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
func (*OsFileReader) Open(name string) (http.File, error) {
	return os.Open(name)
}

func UnsafeClose(file File) {
	err := file.Close()
	if err != nil {
		log.Fatalf("headache execution error, cannot close file %v\n\t%v", file, err)
	}
}


// test utility
type FakeFileInfo struct {
	FileMode os.FileMode
}

func (*FakeFileInfo) Name() string       { panic("not implemented") }
func (*FakeFileInfo) Size() int64        { panic("not implemented") }
func (*FakeFileInfo) ModTime() time.Time { panic("not implemented") }
func (*FakeFileInfo) IsDir() bool        { panic("not implemented") }
func (*FakeFileInfo) Sys() interface{}   { panic("not implemented") }
func (ffi *FakeFileInfo) Mode() os.FileMode {
	return ffi.FileMode
}
