package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	tools := Tools{}
	s := tools.RandomString(10)
	if len(s) != 10 {
		t.Errorf("Expected length of 10, got %d", len(s))
	}
}

var uploadTests = []struct {
	name             string
	allowedFileTypes []string
	renameFile       bool
	errorExpected    bool
}{
	{"allow no rename", []string{"image/jpeg", "image/png"}, false, false},
	{"allow rename", []string{"image/jpeg", "image/png"}, true, false},
	{"not allowed filetype", []string{"image/jpeg"}, false, true},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer func(writer *multipart.Writer) {
				err := writer.Close()
				if err != nil {
					t.Error(err)
				}
			}(writer)
			defer wg.Done()

			part, err := writer.CreateFormFile("file", "testdata/img.png")
			if err != nil {
				t.Error(err)
				return
			}
			f, err := os.Open("testdata/img.png")
			if err != nil {
				t.Error(err)
				return
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					t.Error(err)
				}
			}(f)
			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("Error decoding image", err)
				return
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error("Error encoding image", err)
				return
			}
		}()

		// read from the pipe
		req, err := http.NewRequest("POST", "/", pr)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		var tools Tools
		tools.AllowedFileTypes = e.allowedFileTypes

		uploadedFiles, err := tools.UploadFiles(req, "./testdata/uploads/", e.renameFile)
		if e.errorExpected {
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		} else {
			if err != nil {
				t.Error(err)
			}
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: File not uploaded", e.name)
			}
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		wg.Wait()

	}
}

func TestTools_UploadOneFile(t *testing.T) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer func(writer *multipart.Writer) {
			err := writer.Close()
			if err != nil {
				t.Error(err)
			}
		}(writer)

		part, err := writer.CreateFormFile("file", "testdata/img.png")
		if err != nil {
			t.Error(err)
			return
		}
		f, err := os.Open("testdata/img.png")
		if err != nil {
			t.Error(err)
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				t.Error(err)
			}
		}(f)
		img, _, err := image.Decode(f)
		if err != nil {
			t.Error("Error decoding image", err)
			return
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error("Error encoding image", err)
			return
		}
	}()

	// read from the pipe
	req, err := http.NewRequest("POST", "/", pr)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	var tools Tools

	uploadedFile, err := tools.UploadOneFile(req, "./testdata/uploads/")
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {
		t.Error("File not uploaded")
	}
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))

}

func TestTools_CreateDirIfNotExists(t *testing.T) {
	var tools Tools
	err := tools.CreateDirIfNotExists("./testdata/test-dir")
	if err != nil {
		t.Error(err)
	}
	err = tools.CreateDirIfNotExists("./testdata/test-dir")
	if err != nil {
		t.Error(err)
	}
	_ = os.Remove("./testdata/test-dir")
}

var slugTests = []struct {
	name          string
	input         string
	expected      string
	errorExpected bool
}{
	{"valid string", "now 123", "now-123", false},
	{"valid string", "Now 123!+ &^123", "now-123-123", false},
	{"empty", "", "", true},
	{"japanese string", "今日は", "", true},
	{"japanese string and roman characters", "Hello world今日は", "hello-world", false},
}

func TestTools_Slugify(t *testing.T) {
	var tools Tools
	for _, e := range slugTests {
		slug, err := tools.Slugify(e.input)
		if e.errorExpected {
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		} else {
			if err != nil {
				t.Error(err)
			}
			if slug != e.expected {
				t.Errorf("Expected %s, got %s", e.expected, slug)
			}
		}
	}
}
