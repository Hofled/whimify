package sort

import (
	"log"
	"testing"
)

func TestFileByDependencyOrder(t *testing.T) {
	filePath := ".\\tests\\functions.go"

	err := FileByDependencyOrder(filePath)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}
