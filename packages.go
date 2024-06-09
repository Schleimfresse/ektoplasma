package main

var nativePackages = map[string]*Package{"os": packageOs}

func (s *StdLibFunction) Copy() *Value {
	f := &Value{StdLibFunction: &StdLibFunction{Base: s.Base, Function: s.Function}}
	return f.SetPos(s.Base.PosStart(), s.Base.PosEnd()).SetContext(s.Base.Context)
}

func (s *StdLibFunction) String() string {
	return "<function " + s.Base.Name + " from " + s.PackageName + ">"
}

func NewStdLibFunction(funcMethod func(args []*Value) *RTResult, base *BaseFunction, packageName string) *Value {
	return &Value{StdLibFunction: &StdLibFunction{base, packageName, funcMethod}}
}

var packageOs = &Package{Methods: map[string]*Value{
	"CreateFile":          NewStdLibFunction(CreateFile, &BaseFunction{Name: "CreateFile"}, "os"),
	"OpenFile":            NewStdLibFunction(OpenFile, &BaseFunction{Name: "OpenFile"}, "os"),
	"WriteFile":           NewStdLibFunction(WriteFile, &BaseFunction{Name: "WriteFile"}, "os"),
	"ReadFile":            NewStdLibFunction(ReadFile, &BaseFunction{Name: "ReadFile"}, "os"),
	"DeleteFile":          NewStdLibFunction(DeleteFile, &BaseFunction{Name: "DeleteFile"}, "os"),
	"CloseFile":           NewStdLibFunction(CloseFile, &BaseFunction{Name: "CloseFile"}, "os"),
	"CopyFile":            NewStdLibFunction(CopyFile, &BaseFunction{Name: "CopyFile"}, "os"),
	"MoveFile":            NewStdLibFunction(MoveFile, &BaseFunction{Name: "MoveFile"}, "os"),
	"RenameFile":          NewStdLibFunction(RenameFile, &BaseFunction{Name: "RenameFile"}, "os"),
	"FileExists":          NewStdLibFunction(FileExists, &BaseFunction{Name: "FileExists"}, "os"),
	"CreateDirectory":     NewStdLibFunction(CreateDirectory, &BaseFunction{Name: "CreateDirectory"}, "os"),
	"DeleteDirectory":     NewStdLibFunction(DeleteDirectory, &BaseFunction{Name: "DeleteDirectory"}, "os"),
	"ReadDirectory":       NewStdLibFunction(ReadDirectory, &BaseFunction{Name: "ReadDirectory"}, "os"),
	"GetCurrentDirectory": NewStdLibFunction(GetCurrentDirectory, &BaseFunction{Name: "GetCurrentDirectory"}, "os"),
	"ChangeDirectory":     NewStdLibFunction(ChangeDirectory, &BaseFunction{Name: "ChangeDirectory"}, "os"),
	"DirectoryExists":     NewStdLibFunction(DirectoryExists, &BaseFunction{Name: "DirectoryExists"}, "os"),
	"CopyDirectory":       NewStdLibFunction(CopyDirectory, &BaseFunction{Name: "CopyDirectory"}, "os"),
	"IsDirectory":         NewStdLibFunction(IsDirectory, &BaseFunction{Name: "IsDirectory"}, "os"),
	"MoveDirectory":       NewStdLibFunction(MoveDirectory, &BaseFunction{Name: "MoveDirectory"}, "os"),
	"GetFileSize":         NewStdLibFunction(GetFileSize, &BaseFunction{Name: "GetFileSize"}, "os"),
	"GetFilePermissions":  NewStdLibFunction(GetFilePermissions, &BaseFunction{Name: "GetFilePermissions"}, "os"),
	"SetFilePermissions":  NewStdLibFunction(SetFilePermissions, &BaseFunction{Name: "SetFilePermissions"}, "os"),
	"GetFileOwner":        NewStdLibFunction(GetFileOwner, &BaseFunction{Name: "GetFileOwner"}, "os"),
	"SetFileOwner":        NewStdLibFunction(SetFileOwner, &BaseFunction{Name: "SetFileOwner"}, "os"),
	"StartProcess":        NewStdLibFunction(StartProcess, &BaseFunction{Name: "StartProcess"}, "os"),
	"KillProcess":         NewStdLibFunction(KillProcess, &BaseFunction{Name: "KillProcess"}, "os"),
	"GetEnv":              NewStdLibFunction(GetEnv, &BaseFunction{Name: "GetEnv"}, "os"),
	"SetEnv":              NewStdLibFunction(SetEnv, &BaseFunction{Name: "SetEnv"}, "os"),
	"ListEnv":             NewStdLibFunction(ListEnv, &BaseFunction{Name: "ListEnv"}, "os"),
	"ExecCommand":         NewStdLibFunction(ExecCommand, &BaseFunction{Name: "ExecCommand"}, "os"),
}}
