package main

import (
	files "ektoplasma/sysops/files"
	metadata "ektoplasma/sysops/metadata"
	sysops "ektoplasma/sysops/process"
	"fmt"
	"log"
	"syscall"
)

/* ***********************************************************
	File and Directory Operations
   *********************************************************** */

func CreateFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}
	handle, err := files.Create_file(args[0].Value().(string), syscall.O_RDWR, 0666)

	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to create file: %s", err), nil))
	}
	return res.Success(NewNumber(handle))
}

func OpenFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}
	handle, err := files.Open_file(args[0].Value().(string), syscall.O_RDWR|syscall.O_CLOEXEC, 0)
	if err != nil || handle == -1 {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to open file: %s", err), nil))
	}
	return res.Success(NewNumber(int(handle)))
}

func WriteFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"Number", "Any"}); err != nil {
		return res.Failure(err)
	}
	handle := syscall.Handle(args[0].Value().(int))

	var write int
	var err error

	if args[1].Type() == "ByteArray" {
		write, err = syscall.Write(handle, args[1].Value().([]byte))
	} else {
		write, err = syscall.Write(handle, interfaceToBytes(args[1].Value()))
	}
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
	if err := checkArgumentTypes(args, []string{"Number", "Any"}); err != nil {
		return res.Failure(err)
	}
	buf := make([]byte, 1024)
	n, err := syscall.Read(syscall.Handle(args[0].Value().(int)), buf)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewByteArray(buf[:n]))
}

func ReadFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"Number"}); err != nil {
		return res.Failure(err)
	}
	handle := syscall.Handle(args[0].Value().(int))
	// TODO fix style to files
	var fileInfo syscall.ByHandleFileInformation
	err := syscall.GetFileInformationByHandle(handle, &fileInfo)
	fileSize := int64(fileInfo.FileSizeHigh)<<32 | int64(fileInfo.FileSizeLow)
	buffer := make([]byte, fileSize)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to get fileinfo: %s", err), nil))
	}

	n, err := syscall.Read(handle, buffer)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to read file: %s", err), nil))
	}

	return res.Success(NewByteArray(buffer[:n]))
}

func DeleteFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

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
	if err := checkArgumentTypes(args, []string{"Number"}); err != nil {
		return res.Failure(err)
	}

	err := syscall.Close(syscall.Handle(args[0].Value().(int)))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewNull())
}

func CopyFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}

	exists, err := files.File_exists(args[1].Value().(string))
	if exists || err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error with target file: %s", syscall.ERROR_ALREADY_EXISTS), nil))
	}

	sourceHandle, err := syscall.Open(args[0].Value().(string), syscall.O_RDONLY, 0)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to open file: %s", err), nil))
	}

	defer syscall.Close(sourceHandle)

	filePathPtr, err := syscall.UTF16PtrFromString(args[1].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Failed to convert file path to UTF-16: %s", err), nil))
	}

	buffer := res.Register(ReadFile([]*Value{{Number: &Number{ValueField: int(sourceHandle)}}}))
	if res.Error != nil {
		return res
	}

	copyHandle, err := syscall.CreateFile(filePathPtr, syscall.GENERIC_WRITE, 0, nil, syscall.CREATE_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	defer func() *RTResult {
		err := syscall.Close(sourceHandle)
		if err != nil {
			return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
		}
		err = syscall.Close(copyHandle)
		if err != nil {
			return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
		}
		return nil
	}()

	writeRes := WriteFile([]*Value{{Number: &Number{ValueField: int(copyHandle)}}, buffer})
	res.Register(writeRes)

	fmt.Println("Copied file:", args[0].Value().(string))

	if res.Error != nil {
		return res
	}
	return res.Success(NewNumber(writeRes.Value.Number.ValueField))
}

func MoveFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}

	bytesWritten := res.Register(CopyFile(args))
	if res.Error != nil {
		return res
	}
	res.Register(DeleteFile([]*Value{args[0]}))
	if res.Error != nil {
		return res
	}

	return res.Success(bytesWritten)
}

func RenameFile(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}
	err := syscall.Rename(args[0].Value().(string), args[1].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewNull())
}

func FileExists(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	exists, err := files.File_exists(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error accessing file: %s", syscall.ERROR_FILE_EXISTS), nil))
	}
	return res.Success(NewBoolean(ConvertBoolToInt(exists)))
}

func CreateDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	err := files.MkDir(args[0].Value().(string), 0755)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error creating directory: %s", err), nil))
	}
	return res.Success(NewNull())
}

// DeleteDirectory deletes a directory and all its contents recursively.
func DeleteDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	err := files.RmDir(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error deleting directory: %s", err), nil))
	}
	return res.Success(NewNull())
}

func ReadDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	files, err := files.ReadDir(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error reading directory: %s", err), nil))
	}

	var FileArray []*Value
	for _, file := range files {
		FileArray = append(FileArray, NewString(file))
	}

	return res.Success(NewArray(FileArray))
}

func ChangeDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	err := files.Cd(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error changing directory: %s", err), nil))
	}

	return GetCurrentDirectory([]*Value{})
}

func GetCurrentDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{}); err != nil {
		return res.Failure(err)
	}

	cwd, err := files.Cwd()
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewString(cwd))
}

func DirectoryExists(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	exists := files.Dir_exists(args[0].Value().(string))
	return res.Success(NewBoolean(ConvertBoolToInt(exists)))
}

func CopyDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}
	var dir = args[0].Value().(string)

	var fileLen int
	readDir, err := files.ReadDir(dir)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error reading directory: %s", err), nil))
	}
	fileLen = len(readDir)

	if !files.Dir_exists(args[1].Value().(string)) {
		res.Register(CreateDirectory([]*Value{args[1]}))
		if res.Error != nil {
			return res
		}
	}

	for _, file := range readDir {
		if exists, err := files.Is_dir(dir + "/" + file); exists && err == nil {
			fileLen--
			files := res.Register(CopyDirectory([]*Value{NewString(dir + "/" + file + "/"), NewString(args[1].Value().(string) + "/" + file + "/")}))
			if res.Error != nil {
				return res
			}
			fileLen += files.Value().(int)
		} else if err != nil {
			return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
		} else {
			result := CopyFile([]*Value{NewString(dir + "/" + file), NewString(args[1].Value().(string) + "/" + file)})
			res.Register(result)
			if res.Error != nil {
				return res
			}
		}

	}
	return res.Success(NewNumber(fileLen))
}

func IsDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	isDir, err := files.Is_dir(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewBoolean(ConvertBoolToInt(isDir)))
}

func MoveDirectory(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}
	var dir = args[0].Value().(string)
	var fileLen int
	readDir, err := files.ReadDir(dir)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("Error reading directory: %s", err), nil))
	}
	fileLen = len(readDir)

	if !files.Dir_exists(args[1].Value().(string)) {
		res.Register(CreateDirectory([]*Value{args[1]}))
		if res.Error != nil {
			return res
		}
	}

	for _, file := range readDir {
		exists, err := files.Is_dir(dir + "/" + file)
		log.Println(exists, err)
		if exists, err := files.Is_dir(dir + "/" + file); exists && err == nil {
			fileLen--
			files := res.Register(MoveDirectory([]*Value{NewString(dir + "/" + file + "/"), NewString(args[1].Value().(string) + "/" + file + "/")}))
			if res.Error != nil {
				return res
			}
			log.Println(files.Value())
			fileLen += files.Value().(int)
		} else if err != nil {
			return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
		} else {
			result := MoveFile([]*Value{NewString(dir + "/" + file), NewString(args[1].Value().(string) + "/" + file)})
			res.Register(result)
			if res.Error != nil {
				return res
			}
		}

	}

	res.Register(DeleteDirectory([]*Value{args[0]}))
	if res.Error != nil {
		return res
	}

	return res.Success(NewNumber(fileLen))
}

/* ***********************************************************
	File Metadata Operations
   *********************************************************** */

func GetFileSize(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	filesize, err := metadata.Get_file_size(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewNumber(filesize))
}

func GetFilePermissions(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	perms, err := metadata.Get_file_perms(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewNumber(perms))
}

func SetFilePermissions(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "Number"}); err != nil {
		return res.Failure(err)
	}

	err := metadata.Set_file_perm(args[0].Value().(string), uint32(args[1].Value().(int)))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewNull())
}

func GetFileOwner(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	owner, err := metadata.Get_file_owner(args[0].Value().(string))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}

	return res.Success(NewString(owner))
}

func SetFileOwner(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}

	err := metadata.Set_file_owner(args[0].String.ValueField, args[1].String.ValueField)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	log.Println("ERR:", err)
	return res.Success(NewNull())
}

/* ***********************************************************
	Process Management
   *********************************************************** */

func StartProcess(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "Array"}); err != nil { // TODO array should not be required
		return res.Failure(err)
	}

	flags := make([]string, len(args[1].Array.Elements))

	for i, element := range args[1].Array.Elements {
		flags[i] = element.Value().(string)
	}

	pid, err := sysops.Start_proc(args[0].String.ValueField, flags)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewNumber(pid))
}

func KillProcess(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"Number"}); err != nil {
		return res.Failure(err)
	}

	err := sysops.Kill_process(args[0].Number.ValueField.(int))
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewNull())
}

func GetEnv(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String"}); err != nil {
		return res.Failure(err)
	}

	env, err := sysops.Get_env(args[0].String.ValueField)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewString(env))
}

func SetEnv(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "String"}); err != nil {
		return res.Failure(err)
	}

	err := sysops.Set_env(args[0].String.ValueField, args[1].String.ValueField)
	if err != nil {
		return res.Failure(NewRTError(nil, nil, fmt.Sprintf("%s", err), nil))
	}
	return res.Success(NewNull())
}

func ListEnv(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{}); err != nil {
		return res.Failure(err)
	}

	env := sysops.List_env()
	var arr = make([]*Value, len(env))
	for k, v := range env {
		arr[k] = NewString(v)
	}

	return res.Success(NewArray(arr))
}

func ExecCommand(args []*Value) *RTResult {
	res := NewRTResult()
	if err := checkArgumentTypes(args, []string{"String", "Array"}); err != nil {
		return res.Failure(err)
	}

	var argv = make([]string, len(args[1].Array.Elements))
	for i, arg := range args[1].Array.Elements {
		argv[i] = arg.Value().(string)
	}

	err := sysops.Exec_cmd(args[0].String.ValueField, argv)
	if err != nil {
		return nil
	}
	return res.Success(NewNull())
}
