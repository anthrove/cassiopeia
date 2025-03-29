package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anthrove/identity/pkg/object"
)

var (
	storage       Provider
	randomPath    string
	fileName      string
	fileName2     string
	sampleFile    string
	exceptObjects = 2
)

func setup(t *testing.T) {
	provider := object.Provider{
		TenantID:  "testTenant",
		Parameter: []byte(`{"base_path": "/test/path"}`),
	}

	var err error
	storage, err = newLocalProvider(provider)
	if err != nil {
		t.Fatalf("Failed to initialize local provider: %v", err)
	}

	randomPath = strings.Replace(time.Now().Format("20060102150506.000"), ".", "", -1)
	fmt.Printf("testing file in %v\n", filepath.Join(storage.GetEndpoint(), randomPath))

	fileName = "/" + filepath.Join(randomPath, "alexander-andrews-mEdKuPYJe1I-unsplash.jpg")
	fileName2 = "/" + filepath.Join(randomPath, "sample2", "alexander-andrews-mEdKuPYJe1I-unsplash.jpg")
	sampleFile, _ = filepath.Abs("../../../test/testdata/alexander-andrews-mEdKuPYJe1I-unsplash.jpg")
}

func TestLocalProvider(t *testing.T) {
	setup(t)

	t.Run("PutFile", testPutFile)
	t.Run("PutFileAgain", testPutFileAgain)
	t.Run("PutFile2", testPutFile2)
	t.Run("GetFile", testGetFile)
	t.Run("GetURL", testGetURL)
	t.Run("GetStream", testGetStream)
	t.Run("ListObjects", testListObjects)
	t.Run("DeleteFile", testDeleteFile)
	t.Run("GetFileAfterDelete", testGetFileAfterDelete)
	t.Run("GetFile2AfterDelete", testGetFile2AfterDelete)
}

func testPutFile(t *testing.T) {
	file, err := os.Open(sampleFile)
	if err != nil {
		t.Errorf("No error should happen when open sample file, but got %v", err)
		return
	}
	defer file.Close()

	object, err := storage.Put(fileName, file)
	if err != nil {
		t.Errorf("No error should happen when save sample file, but got %v", err)
	} else if object.Path == "" || object.StorageInterface == nil {
		t.Errorf("returned object should have necessary information")
	}
}

func testPutFileAgain(t *testing.T) {
	file, err := storage.GetStream(fileName)
	if err != nil {
		t.Errorf("No error should happen when open sample file, but got %v", err)
		return
	}

	object, err := storage.Put(fileName, file)
	if err != nil {
		t.Errorf("No error should happen when save sample file, but got %v", err)
	} else if object.Path == "" || object.StorageInterface == nil {
		t.Errorf("returned object should have necessary information")
	}
}

func testPutFile2(t *testing.T) {
	file, err := os.Open(sampleFile)
	if err != nil {
		t.Errorf("No error should happen when open sample file, but got %v", err)
		return
	}
	defer file.Close()

	object, err := storage.Put(fileName2, file)
	if err != nil {
		t.Errorf("No error should happen when save sample file, but got %v", err)
	} else if object.Path == "" || object.StorageInterface == nil {
		t.Errorf("returned object should have necessary information")
	}
}

func testGetFile(t *testing.T) {
	file, err := storage.Get(fileName)
	if err != nil {
		t.Errorf("No error should happen when get sample file, but got %v", err)
		return
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		t.Errorf("No error should happen when read downloaded file, but got %v", err)
	} else if string(buffer) == "sample" {
		t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
	}
}

func testGetURL(t *testing.T) {
	url, err := storage.GetURL(fileName)
	if err != nil {
		t.Errorf("No error should happen when GetURL for sample file, but got %v", err)
		return
	}

	if strings.HasPrefix(url, "http") {
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("No error should happen when get file with public URL")
			return
		}
		defer resp.Body.Close()

		buffer, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("No error should happen when read downloaded file, but got %v", err)
		} else if string(buffer) == "sample" {
			t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
		}
	}
}

func testGetStream(t *testing.T) {
	stream, err := storage.GetStream(fileName)
	if err != nil {
		t.Errorf("No error should happen when get sample file, but got %v", err)
		return
	}
	defer stream.Close()

	buffer, err := io.ReadAll(stream)
	if err != nil {
		t.Errorf("No error should happen when read downloaded file, but got %v", err)
	} else if string(buffer) == "sample" {
		t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
	}
}

func testListObjects(t *testing.T) {
	objects, err := storage.List(randomPath)
	if err != nil {
		t.Errorf("No error should happen when list objects, but got %v", err)
		return
	}

	if len(objects) != exceptObjects {
		t.Errorf("Should found %v objects, but got %v", exceptObjects, len(objects))
	}
}

func testDeleteFile(t *testing.T) {
	err := storage.Delete(fileName)
	if err != nil {
		t.Errorf("No error should happen when delete sample file, but got %v", err)
	}
}

func testGetFileAfterDelete(t *testing.T) {
	_, err := storage.Get(fileName)
	if err == nil {
		t.Errorf("There should be an error when get deleted sample file")
	}
}

func testGetFile2AfterDelete(t *testing.T) {
	_, err := storage.Get(fileName2)
	if err != nil {
		t.Errorf("Sample file 2 should not have been deleted")
	}
}
