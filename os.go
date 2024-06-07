package main

import (
	"fmt"
	"syscall"
)

/* ***********************************************************
	File and Directory Operations
   *********************************************************** */

func CreateFile(args []*Value) *RTResult {
	res := NewRTResult()
	filePathPtr, err := syscall.UTF16PtrFromString(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to convert file path to UTF-16: %s", err), nil))
	}

	handle, err := syscall.CreateFile(filePathPtr, syscall.GENERIC_WRITE, 0, nil, syscall.CREATE_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to create file: %s", err), nil))
	}
	return res.Success(NewNumber(int(handle)))
}

func OpenFile(args []*Value) *RTResult {
	res := NewRTResult()
	handle, err := syscall.Open(args[0].Value().(string), syscall.O_RDWR|syscall.O_CLOEXEC, 0)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to open file: %s", err), nil))
	}
	return res.Success(NewNumber(int(handle)))
}

func WriteFile(args []*Value) *RTResult {
	res := NewRTResult()
	handle := syscall.Handle(args[0].Value().(int))

	write, err := syscall.Write(handle, interfaceToBytes(args[1].Value()))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to write to file: %s", err), nil))
	}

	// reset the pointer in the file
	_, err = syscall.Seek(handle, 0, 0)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to reset file pointer: %s", err), nil))
	}
	return res.Success(NewNumber(write))
}

// TODO fix append and write rethink openfile

func AppendFile(args []*Value) *RTResult {
	res := NewRTResult()
	buf := make([]byte, 1024)
	n, err := syscall.Read(syscall.Handle(args[0].Value().(int)), buf)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewByteArray(buf[:n]))
}

func ReadFile(args []*Value) *RTResult {
	res := NewRTResult()
	buf := make([]byte, 1024)
	n, err := syscall.Read(syscall.Handle(args[0].Value().(int)), buf)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to read file: %s", err), nil))
	}
	return res.Success(NewByteArray(buf[:n]))
}

func DeleteFile(args []*Value) *RTResult {
	res := NewRTResult()

	filePathPtr, err := syscall.UTF16PtrFromString(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	err = syscall.DeleteFile(filePathPtr)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewNull())
}

func CloseFile(args []*Value) *RTResult {
	res := NewRTResult()

	err := syscall.Close(syscall.Handle(args[0].Value().(int)))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewNull())
}
