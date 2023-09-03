package main

import (
	"strings"
	"testing"
)

func TestGenerateFileName(t *testing.T) {
	fileName := GenerateName("i test this function")

	isOkExtension := strings.Contains(fileName, ".mp3")
	isOkFormat := strings.Contains(fileName, "irp")

	if !isOkExtension {
		t.Fail()
		t.Log("filename should contain .mp3 suffix")
	}

	if !isOkFormat {
		t.Fail()
		t.Log("filename should contain irp preffix")
	}
}
