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
	"CreateFile": NewStdLibFunction(CreateFile, &BaseFunction{Name: "CreateFile"}, "os"),
	"OpenFile":   NewStdLibFunction(OpenFile, &BaseFunction{Name: "OpenFile"}, "os"),
	"WriteFile":  NewStdLibFunction(WriteFile, &BaseFunction{Name: "WriteFile"}, "os"),
	"ReadFile":   NewStdLibFunction(ReadFile, &BaseFunction{Name: "ReadFile"}, "os"),
	"DeleteFile": NewStdLibFunction(DeleteFile, &BaseFunction{Name: "DeleteFile"}, "os"),
	"CloseFile":  NewStdLibFunction(CloseFile, &BaseFunction{Name: "CloseFile"}, "os"),
}}
