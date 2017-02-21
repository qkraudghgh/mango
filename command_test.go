package main

import (
	"os"
	"testing"

	"github.com/qkraudghgh/mango/utils"
)

func TestCreateBucketFunc(t *testing.T) {
	t.Log("Create Mango Bucket")

	defer TestClearFunc(t)

	if err := mangoUtils.CheckBucketAndMake(); err != nil {
		t.Errorf("Error occur when make bucket: %v", err)
	}
}

func TestClearFunc(t *testing.T) {
	t.Log("Clear DB testcase")

	os.Remove(mangoUtils.GetDbPath())
}

func TestAddFunc(t *testing.T) {
	t.Log("Add Todo testcase")

	TestCreateBucketFunc(t)
	defer TestClearFunc(t)

	var args []string

	args = []string{"This is Test"}

	if err := addFunc(args); err != nil {
		t.Errorf("Error occur when run add command: %v", err)
	}

}

func TestDeleteFunc(t *testing.T) {
	t.Log("Delete Todo testcase")

	TestCreateBucketFunc(t)
	defer TestClearFunc(t)

	var args []string

	args = []string{"Delete Soon"}

	addFunc(args)

	args = []string{"2"}

	if err := deleteFunc(args); err == nil {
		t.Error("Expect error")
	}

	args = []string{"a"}

	if err := deleteFunc(args); err == nil {
		t.Error("Expect error")
	}

	args = []string{"1"}

	if err := deleteFunc(args); err != nil {
		t.Errorf("Error occur when run delete command: %v", err)
	}
}

func TestDoneFunc(t *testing.T) {
	t.Log("Done Todo testcase")

	TestCreateBucketFunc(t)
	defer TestClearFunc(t)

	var args []string

	args = []string{"Done Soon"}

	addFunc(args)

	args = []string{"2"}

	if err := doneFunc(args); err == nil {
		t.Error("Expect error")
	}

	args = []string{"a"}

	if err := doneFunc(args); err == nil {
		t.Error("Expect error")
	}

	args = []string{"1"}

	if err := doneFunc(args); err != nil {
		t.Errorf("Error occur when run done command: %v", err)
	}
}

func TestUnDoneFunc(t *testing.T) {
	t.Log("Done Todo testcase")

	TestCreateBucketFunc(t)
	defer TestClearFunc(t)

	var args []string

	args = []string{"unDone Soon"}

	addFunc(args)

	args = []string{"1"}

	doneFunc(args)

	args = []string{"2"}

	if err := unDoneFunc(args); err == nil {
		t.Error("Expect error")
	}

	args = []string{"a"}

	if err := unDoneFunc(args); err == nil {
		t.Error("Expect error")
	}

	args = []string{"1"}

	if err := unDoneFunc(args); err != nil {
		t.Errorf("Error occur when run done command: %v", err)
	}
}
