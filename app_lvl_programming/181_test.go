package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)


func TestMultipartPost(t *testing.T) {
	reqBody := new(bytes.Buffer)

	w := multipart.NewWriter(reqBody)

	for k, v := range map[string]string {
	"date":		time.Now().Format(time.RFC3339),
	"description":	"Form values with attached files",
	} {
		err := w.WriteField(k, v)
		if err != nil {
			t.Fatal(err)
		}

	}

	for i, file := range []string {
		"./files/hello.txt",
		"./files/goodbye.txt",
	} {
		filePart, err := w.CreateFormFile(fmt.Sprintf("field%d", i+1),
			filepath.Base(file))
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(filePart, f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,http.MethodPost,
		"https://httpbin.org/post", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {_ = resp.Body.Close()}()

	b, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		t.Fatal(err)
	}


	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; acutal status %d",
			http.StatusOK, resp.StatusCode)
	}

	t.Logf("\n%s", b)

}


