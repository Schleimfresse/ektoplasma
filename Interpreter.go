package main

import (
	"fmt"
	"log"
	"reflect"
)

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) visit(node Node, context *Context) *RTResult {
	switch n := node.(type) {
	case *UnaryOpNode:
		return i.visitUnaryOpNode(*n, context)
	case *BinOpNode:
		return i.visitBinOpNode(*n, context)
	case *NumberNode:
		return i.visitNumberNode(*n, context)
	case *VarAccessNode:
		return i.visitVarAccessNode(*n, context)
	case *VarAssignNode:
		return i.visitVarAssignNode(*n, context)
	case *IfNode:
		return i.visitIfNode(*n, context)
	case *ForNode:
		return i.visitForNode(*n, context)
	case *WhileNode:
		return i.visitWhileNode(*n, context)
	case *CallNode:
		return i.visitCallNode(*n, context)
	case *FuncDefNode:
		return i.visitFuncDefNode(*n, context)
	case *StringNode:
		return i.visitStringNode(*n, context)
	case *ArrayNode:
		return i.visitArrayNode(*n, context)
	case *IndexNode:
		return i.visitIndexNode(*n, context)
	case *ReturnNode:
		return i.visitReturnNode(*n, context)
	case *ContinueNode:
		return i.visitContinueNode()
	case *BreakNode:
		return i.visitBreakNode()
	case *ImportNode:
		return i.visitImportNode(*n, context)
	case *PackageMethod:
		return i.visitPackageMethodNode(*n, context)
	default:
		// Handle unknown node types
		return NewRTResult().Failure(NewRTError(node.PosStart(), node.PosEnd(), fmt.Sprintf("No visit method defined for node type %T", node), context))
	}
}

// NoVisitMethod raises an exception for non-existing visit methods.
func (i *Interpreter) NoVisitMethod(node Node, context *Context) *RTResult {
	errMsg := fmt.Sprintf("No visit_%T method defined", reflect.TypeOf(node))
	return NewRTResult().Failure(NewRTError(node.PosStart(), node.PosEnd(), errMsg, context))
}

func (i *Interpreter) visitNumberNode(node NumberNode, context *Context) *RTResult {
	value := NewNumber(node.Value)
	value.SetContext(context).SetPos(node.PosStart(), node.PosEnd())
	return NewRTResult().Success(value)
}

func (i *Interpreter) visitStringNode(node StringNode, context *Context) *RTResult {
	if TokenValue, ok := node.Value.(string); ok {
		value := NewString(TokenValue)
		value.SetContext(context).SetPos(node.PosStart(), node.PosEnd())
		return NewRTResult().Success(value)
	}
	return NewRTResult().Failure(NewRTError(node.PosStart(), node.PosEnd(), fmt.Sprintf("%s is not of type string", reflect.TypeOf(node.Value)), context))
}

func (i *Interpreter) visitImportNode(node ImportNode, context *Context) *RTResult {
	res := NewRTResult()

	if node.ImportNames != nil {
		packageName := node.PackageNames[0]
		importNames := node.ImportNames

		// Check if it's a native package
		if pkg, exists := nativePackages[packageName.Value.(string)]; exists {
			for _, importName := range importNames {
				name := importName.Value.(string)
				if function, exists := pkg.Methods[name]; exists {
					context.SymbolTable.Set(name, function, false)
				} else {
					return res.Failure(NewRTError(
						node.PositionStart, node.PositionEnd,
						fmt.Sprintf("'%s' not found in native package '%s'", name, packageName.Value),
						context,
					))
				}
			}
			return res.Success(NewEmptyValue())
		}

		packageContent, err := LoadPackage(packageName.Value.(string))
		if err != nil {
			return res.Failure(NewRTError(
				node.PositionStart, node.PositionEnd,
				fmt.Sprintf("Failed to load package, no package named '%s', ", packageName.Value),
				context,
			))
		}

		packageContext := NewContext(fmt.Sprintf("<package %v>", packageName.Value.(string)), context, node.PositionStart)
		packageContext.SymbolTable = NewSymbolTable(nil)

		tokens, Err := NewLexer(packageName.Value.(string)+".ecp", packageContent).MakeTokens()
		if Err != nil {
			return res.Failure(NewRTError(node.PositionStart, node.PositionEnd, fmt.Sprintf("Error tokenizing imported package %v", packageName), context))
		}

		parser := NewParser(tokens)
		packageAst := parser.Parse()
		if packageAst.Error != nil {
			return res.Failure(NewRTError(packageAst.Error.PosStart, packageAst.Error.PosEnd, packageAst.Error.Details, packageContext))
		}

		interpreter := NewInterpreter()
		res.Register(interpreter.visit(packageAst.Node, packageContext))
		if res.Error != nil {
			return res
		}

		// Retrieve the function from the package's symbol table
		for _, importName := range importNames {
			value, exists, isConst := packageContext.SymbolTable.Get(importName.Value.(string))
			if !exists {
				return res.Failure(NewRTError(
					node.PositionStart, node.PositionEnd,
					fmt.Sprintf("Function '%s' not declared in package '%s'", importName.Value.(string), packageName.Value.(string)),
					context,
				))
			}

			// Register the function in the current context
			context.SymbolTable.Set(importName.Value.(string), value, isConst)
		}
	} else {
		packageNames := node.PackageNames

		// add native support

		for _, pkg := range packageNames {
			// Check if it's a native package
			if _pkg, exists := nativePackages[pkg.Value.(string)]; exists {
				context.SymbolTable.SetPackage(pkg.Value.(string), _pkg)
				return res.Success(NewEmptyValue())
			}

			packageName := pkg.Value.(string)
			packageContent, err := LoadPackage(packageName)
			if err != nil {
				return res.Failure(NewRTError(
					node.PositionStart, node.PositionEnd,
					fmt.Sprintf("Failed to load package, no package named '%s', ", packageName),
					context,
				))
			}

			packageContext := NewContext(fmt.Sprintf("<package %v>", packageName), context, node.PositionStart)
			packageContext.SymbolTable = NewSymbolTable(nil)

			tokens, Err := NewLexer(packageName+".ecp", packageContent).MakeTokens()
			if Err != nil {
				return res.Failure(NewRTError(node.PositionStart, node.PositionEnd, fmt.Sprintf("Error tokenizing imported package %v", packageName), context))
			}

			parser := NewParser(tokens)
			packageAst := parser.Parse()
			if packageAst.Error != nil {
				return res.Failure(NewRTError(packageAst.Error.PosStart, packageAst.Error.PosEnd, packageAst.Error.Details, packageContext))
			}

			interpreter := NewInterpreter()
			res.Register(interpreter.visit(packageAst.Node, packageContext))
			if res.Error != nil {
				return res
			}

			// Export all functions and constants from packageContext to the global scope
			for name, value := range packageContext.SymbolTable.symbols {
				context.SymbolTable.Set(name, value, false)
			}
			for name, value := range packageContext.SymbolTable.constants {
				context.SymbolTable.Set(name, value, true)
			}
		}
	}

	return res.Success(NewEmptyValue())
}

func (i *Interpreter) visitArrayNode(node ArrayNode, context *Context) *RTResult {
	res := NewRTResult()
	var elements []*Value

	for _, elementNode := range node.ElementNodes {
		result := i.visit(elementNode, context)
		if result.Error != nil {
			return result
		}
		value := res.Register(result)
		if !value.IsEmpty() {
			elements = append(elements, value)
		}
		if res.ShouldReturn() {
			return res
		}
	}

	newArray := NewArray(elements).SetContext(context).SetPos(node.PosStart(), node.PosEnd())

	return res.Success(newArray)
}

func (i *Interpreter) visitBinOpNode(node BinOpNode, context *Context) *RTResult {
	res := NewRTResult()

	leftRTValue := i.visit(node.LeftNode, context)
	res.Register(leftRTValue)
	if res.ShouldReturn() {
		return res
	}

	rightRTValue := i.visit(node.RightNode, context)
	res.Register(rightRTValue)
	if res.ShouldReturn() {
		return res
	}

	left := leftRTValue.Value
	right := rightRTValue.Value

	var result *Value
	var err *RuntimeError

	// get the operation type and use the left and the right.Number node from the operation symbol as values
	switch node.OpTok.Type {
	case TT_PLUS:
		if left.String != nil && right.String != nil {
			result, err = left.String.AddedTo(right.String)
		} else if left.Number != nil && right.Number != nil {
			result, err = left.Number.AddedTo(right.Number)
		} else if left.Array != nil && right.Number != nil {
			result, err = left.Array.AddedTo(right)
		} else if right.Array != nil && left.Array != nil {
			result, err = left.Array.AddedTo(right)
		} else if right.ByteArray != nil && left.ByteArray != nil {
			result, err = left.ByteArray.AddedTo(right.ByteArray)
		} else {
			return res.Failure(NewRTError(left.GetPosStart(), right.GetPosEnd(), fmt.Sprintf("Cannot add values of types %s and %s together ", left.Type(), right.Type()), context))
		}
	case TT_MINUS:
		if left.Array != nil && right.Number != nil {
			result, err = left.Array.SubtractedBy(right)
		} else {
			result, err = left.Number.SubtractedBy(right.Number)
		}
	case TT_MUL:
		if left.String != nil && right.Number != nil {
			result, err = left.String.MultipliedBy(right.Number)
		} else if left.Number != nil && right.Number != nil {
			result, err = left.Number.MultipliedBy(right.Number)
		} else if left.Array != nil && right.Array != nil {
			result, err = left.Array.MultipliedBy(right)
		} else if right.Number != nil && left.ByteArray != nil {
			result, err = left.ByteArray.MultipliedBy(right.Number)
		} else {
			return res.Failure(NewRTError(left.GetPosStart(), right.GetPosEnd(), fmt.Sprintf("Cannot multiply values of types %s and %s", left.Type(), right.Type()), context))
		}
	case TT_DIV:
		result, err = left.Number.DividedBy(right.Number)
	case TT_POW:
		result, err = left.Number.PowedBy(right.Number)
	case TT_EE:
		if left.Number != nil {
			result, err = left.Number.GetComparisonEq(right)
		} else if left.Null != nil {
			result, err = left.Null.GetComparisonEq(right)
		} else if right.Null != nil {
			result, err = right.Null.GetComparisonEq(left)
		} else if right.Boolean != nil && left.Boolean != nil {
			result, err = left.Boolean.GetComparisonEq(right.Boolean)
		} else if left.String != nil {
			result, err = left.String.GetComparisonEq(left)
		} else if right.ByteArray != nil && left.ByteArray != nil {
			result, err = left.ByteArray.GetComparisonEq(right)
		} else {
			return res.Failure(NewRTError(left.GetPosStart(), right.GetPosEnd(), fmt.Sprintf("Cannot compaire values of types %s and %s ", left.Type(), right.Type()), context))
		}
	case TT_NE:
		if left.Number != nil && right.Number != nil {
			result, err = left.Number.GetComparisonNe(right.Number)
		} else if right.Boolean != nil && left.Boolean != nil {
			result, err = left.Boolean.GetComparisonNe(right.Boolean)
		} else if left.Null != nil {
			result, err = left.Null.GetComparisonNe(right)
		} else if right.Null != nil {
			result, err = right.Null.GetComparisonNe(left)
		}
	case TT_LT:
		result, err = left.Number.GetComparisonLt(right.Number)
	case TT_GT:
		if left.Number != nil && right.Number != nil {
			result, err = left.Number.GetComparisonGt(right.Number)
		} else {
			return res.Failure(NewRTError(left.GetPosStart(), right.GetPosEnd(), fmt.Sprintf("Cannot compaire values of types %s and %s ", left.Type(), right.Type()), context))
		}
	case TT_LTE:
		result, err = left.Number.GetComparisonLte(right.Number)
	case TT_GTE:
		result, err = left.Number.GetComparisonGte(right.Number)
	case TT_KEYWORD:
		if node.OpTok.Value == "AND" {
			result, err = left.Number.AndedBy(right.Number)
		} else if node.OpTok.Value == "OR" {
			result, err = left.Number.OredBy(right.Number)
		}
	default:
		return res.Failure(NewRTError(node.OpTok.PosStart, node.OpTok.PosEnd, "Invalid operation", context))
	}
	if err != nil {
		return res.Failure(err)
	}

	result.SetContext(context).SetPos(node.PosStart(), node.PosEnd())

	return res.Success(result)
}

func (i *Interpreter) visitUnaryOpNode(node UnaryOpNode, context *Context) *RTResult {
	res := NewRTResult()
	numValue := res.Register(i.visit(node.Node, context))
	if res.ShouldReturn() {
		return res
	}

	var result *Value
	var err *RuntimeError

	num := numValue.Number
	if num == nil {
		return res.Failure(NewRTError(node.Node.PosStart(), node.Node.PosEnd(), "Expected a number", context))
	}

	// else if for some reason required, when not expressions like +1 won't work because the context is not set
	if node.OpTok.Type == TT_MINUS {
		result, err = num.MultipliedBy(NewNumber(-1).Number)
		if err != nil {
			return res.Failure(err)
		}
	} else if node.OpTok.Type == TT_PLUS {
		result, err = num.MultipliedBy(NewNumber(1).Number)
		if err != nil {
			return res.Failure(err)
		}
	} else if node.OpTok.Matches(TT_KEYWORD, "not") {
		result, err = num.Notted()
	}

	if err != nil {
		return res.Failure(err)
	} else {
		result.SetContext(context).SetPos(node.PosStart(), node.PosEnd())
		return res.Success(result)
	}
}

func (i *Interpreter) visitPackageMethodNode(node PackageMethod, context *Context) *RTResult {
	res := NewRTResult()
	var result *Value
	packageMethod, exists := context.SymbolTable.GetPackageMethod(node.PackageName, node.MethodName)
	if !exists {
		return res.Failure(NewRTError(
			node.PosStart(), node.PosEnd(),
			fmt.Sprintf("Unresolved package reference '%s'", node.PackageName),
			context))
	}

	packageMethod.SetPos(node.PosStart(), node.PosEnd()).SetContext(context)

	callNode := *node.CallNode
	if _, ok := callNode.(*CallNode); ok {
		var args []*Value
		for _, argNode := range callNode.(*CallNode).ArgNodes {
			arg := res.Register(i.visit(argNode, context))
			args = append(args, arg)
		}

		result = res.Register(packageMethod.StdLibFunction.Function(args))
		if res.Error != nil {
			res.Error.SetLocation(packageMethod.StdLibFunction.Base)
			return res
		}
	} else if _, ok := callNode.(*VarAccessNode); ok {
		result = packageMethod
	}

	return res.Success(result)
}

// visitVarAccessNode visits a VarAccessNode and retrieves its value from the symbol table.
func (i *Interpreter) visitVarAccessNode(node VarAccessNode, context *Context) *RTResult {
	res := NewRTResult()
	varName := node.VarNameTok.Value

	value, exists, _ := context.SymbolTable.Get(varName.(string))
	if !exists {
		log.Println("VAR NAME:", varName, context.SymbolTable.packages)
		if _, exists := context.SymbolTable.packages[varName.(string)]; exists {
			// TODO error for package with dot but no func -> parser917
			return res.Failure(NewRTError(
				node.PosStart(), node.PosEnd(),
				fmt.Sprintf("Use of the package '%s' without a selector", varName),
				context))
		}
		return res.Failure(NewRTError(
			node.PosStart(), node.PosEnd(),
			fmt.Sprintf("Unresolved reference '%s'", varName),
			context))
	}

	// value.SetPos(value.GetPosStart(), value.GetPosEnd()).SetContext(context)
	value.SetContext(context)

	return res.Success(value)
}

// visitVarAssignNode visits a VarAssignNode and assigns a value to the variable in the symbol table.
func (i *Interpreter) visitVarAssignNode(node VarAssignNode, context *Context) *RTResult {
	res := NewRTResult()
	varName := node.VarNameTok.Value
	var value *Value

	if context.SymbolTable.Contains(varName.(string)) && node.declaration && GlobalSymbolTable == context.SymbolTable {
		return res.Failure(NewRTError(node.PosStart(), node.PosEnd(), fmt.Sprintf("Variable '%s' redeclared in scope", varName), context))
	} else if !node.declaration && !context.SymbolTable.Contains(varName.(string)) {
		return res.Failure(NewRTError(
			node.PosStart(), node.PosEnd(),
			fmt.Sprintf("Unresolved reference '%s'", varName),
			context))
	}
	if node.ValueNode != nil {
		value = res.Register(i.visit(*node.ValueNode, context))
	} else {
		value = NewNull()
	}
	if res.ShouldReturn() {
		return res
	}

	err := context.SymbolTable.Set(varName.(string), value, node.isConst)
	if err != nil {
		return res.Failure(err)
	}

	return res.Success(NewEmptyValue())
}

func (i *Interpreter) visitIfNode(node IfNode, context *Context) *RTResult {
	res := NewRTResult()

	for _, ifcase := range node.Cases {

		value := res.Register(i.visit(ifcase.Condition, context))
		if res.ShouldReturn() {
			return res
		}

		conditionValue := value.Boolean
		if conditionValue.IsTrue() {
			exprValue := res.Register(i.visit(ifcase.Expr, context))
			if res.ShouldReturn() {
				return res
			}
			if ifcase.Flag {
				return res.Success(NewNull())
			}
			return res.Success(exprValue)
		}
	}

	if node.ElseCase != nil {
		elseValue := res.Register(i.visit(node.ElseCase.Expr, context))
		if res.ShouldReturn() {
			return res
		}
		if node.ElseCase.Flag {
			return res.Success(NewNull())
		}
		return res.Success(elseValue)
	}

	return res.Success(NewNull())
}

func (i *Interpreter) visitForNode(node ForNode, context *Context) *RTResult {
	res := NewRTResult()
	var elements []*Value

	start := res.Register(i.visit(node.StartValueNode, context))
	if res.ShouldReturn() {
		return res
	}

	end := res.Register(i.visit(node.EndValueNode, context))
	if res.ShouldReturn() {
		return res
	}

	var stepValue *Number
	if node.StepValueNode != nil {
		stepValue = res.Register(i.visit(node.StepValueNode, context)).Number
		if res.ShouldReturn() {
			return res
		}
	} else {
		stepValue = NewNumber(1).Number
	}

	var iVal int
	var endValue int

	if IsInt(start.Number.ValueField) {
		iVal = start.Number.ValueField.(int)
	} else {
		startNum := start.Number
		return res.Failure(NewRTError(startNum.PosStart(), startNum.PosEnd(), fmt.Sprintf("Can not use type %s as type int", reflect.TypeOf(startNum.ValueField)), startNum.Context))
	}
	if IsInt(end.Number.ValueField) {
		endValue = end.Number.ValueField.(int)
	} else {
		endNum := end.Number
		return res.Failure(NewRTError(endNum.PosStart(), endNum.PosEnd(), fmt.Sprintf("Can not use type %s as type int", reflect.TypeOf(endNum.ValueField)), endNum.Context))
	}

	var condition func() bool
	if stepValue.ValueField.(int) >= 0 {
		condition = func() bool { return iVal < endValue }
	} else {
		condition = func() bool { return iVal > endValue }
	}

	for condition() {
		context.SymbolTable.Set(node.VarNameTok.Value.(string), NewNumber(iVal), false)

		iVal += stepValue.ValueField.(int)

		value := res.Register(i.visit(node.BodyNode, context))

		if res.ShouldReturn() && res.LoopShouldContinue == false && res.LoopShouldBreak == false {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}
		if res.LoopShouldBreak {
			break
		}

		elements = append(elements, value)
	}

	if node.Flag {
		return res.Success(NewEmptyValue())
	}
	return res.Success(NewArray(elements).SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visitWhileNode(node WhileNode, context *Context) *RTResult {
	res := NewRTResult()
	var elements []*Value

	for {
		condition := res.Register(i.visit(node.ConditionNode, context))
		if res.ShouldReturn() {
			return res
		}

		if !condition.Boolean.IsTrue() {
			break
		}

		value := res.Register(i.visit(node.BodyNode, context))

		if res.ShouldReturn() && res.LoopShouldContinue == false && res.LoopShouldBreak == false {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}
		if res.LoopShouldBreak {
			break
		}

		elements = append(elements, value)
	}

	if node.Flag {
		return res.Success(NewEmptyValue())
	}
	return res.Success(NewArray(elements).SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visitFuncDefNode(node FuncDefNode, context *Context) *RTResult {
	res := NewRTResult()

	var funcName *string
	if node.VarNameTok != nil {
		value := node.VarNameTok.Value.(string)
		funcName = &value
	} else {
		funcName = nil
	}

	argNames := make([]string, len(node.ArgNameToks))
	for idx, argName := range node.ArgNameToks {
		argNames[idx] = argName.Value.(string)
	}

	value := NewFunction(funcName, &node.BodyNode, argNames, node.Flag)
	value.SetContext(context).SetPos(node.PosStart(), node.PosEnd())

	if node.VarNameTok != nil {
		context.SymbolTable.Set(*funcName, value, false)
	}

	return res.Success(NewEmptyValue())
}

func (i *Interpreter) visitCallNode(node CallNode, context *Context) *RTResult {
	res := NewRTResult()
	var args []*Value

	Call := res.Register(i.visit(node.NodeToCall, context))
	if res.ShouldReturn() {
		return res
	}
	valueToCall := Call.SetPos(node.PosStart(), node.PosEnd()).SetContext(context)

	for _, argNode := range node.ArgNodes {
		args = append(args, res.Register(i.visit(argNode, context)))
		if res.Error != nil {
			return res
		}
		if res.ShouldReturn() {
			return res
		}
	}

	var returnValue *Value
	if valueToCall.Function != nil {
		returnValue = res.Register(valueToCall.Function.Execute(args))
	} else if valueToCall.BuildInFunction != nil {
		returnValue = res.Register(valueToCall.BuildInFunction.Execute(args...))
	} else if valueToCall.StdLibFunction != nil {
		returnValue = res.Register(valueToCall.StdLibFunction.Function(args))
		if res.Error != nil {
			res.Error.SetLocation(valueToCall.StdLibFunction.Base)
		}
	}
	if res.ShouldReturn() {
		return res
	}
	returnValue = returnValue.Copy().SetPos(node.PosStart(), node.PosEnd()).SetContext(context)
	return res.Success(returnValue)
}

func (i *Interpreter) visitIndexNode(node IndexNode, context *Context) *RTResult {
	res := NewRTResult()
	array := i.visit(node.VarAccessNode, context)
	index := i.visit(node.IndexNode, context)

	if array.Value.Array != nil {
		value, err := array.Value.Array.GetIndex(index.Value.Number)

		res.Failure(err)
		if err != nil {
			return res
		}
		return res.Success(value)
	} else {
		return res.Failure(NewRTError(node.PosStart(), node.PosEnd(), fmt.Sprintf("Element at index %v could not be retrieved, index is out of bounds with length 0", node.IndexNode.Value), context))
	}
}

func (i *Interpreter) visitReturnNode(node ReturnNode, context *Context) *RTResult {
	res := NewRTResult()

	var value *Value
	if node.NodeToReturn != nil {
		value = res.Register(i.visit(node.NodeToReturn, context))
		if res.ShouldReturn() {
			return res
		}
	} else {
		value = NewNull()
	}
	return res.SuccessReturn(value)
}

func (i *Interpreter) visitContinueNode() *RTResult {
	return NewRTResult().SuccessContinue()
}

func (i *Interpreter) visitBreakNode() *RTResult {
	return NewRTResult().SuccessBreak()
}

// NewRTResult creates a new RTResult instance.
func NewRTResult() *RTResult {
	return &RTResult{}
}

// Register registers the result of a runtime operation.
func (r *RTResult) Register(res *RTResult) *Value {
	r.Error = res.Error
	r.FuncReturnValue = res.FuncReturnValue
	r.LoopShouldContinue = res.LoopShouldContinue
	r.LoopShouldBreak = res.LoopShouldBreak
	return res.Value
}

func (r *RTResult) Reset() {
	r.Value = nil
	r.Error = nil
	r.FuncReturnValue = nil
	r.LoopShouldContinue = false
	r.LoopShouldBreak = false
}

// Success indicates a successful runtime operation.
func (r *RTResult) Success(value *Value) *RTResult {
	r.Reset()
	r.Value = value
	return r
}

func (r *RTResult) SuccessReturn(value *Value) *RTResult {
	r.Reset()
	r.FuncReturnValue = value
	return r
}

func (r *RTResult) SuccessContinue() *RTResult {
	r.Reset()
	r.LoopShouldContinue = true
	return r
}

func (r *RTResult) SuccessBreak() *RTResult {
	r.Reset()
	r.LoopShouldBreak = true
	return r
}

// Failure indicates a failed runtime operation.
func (r *RTResult) Failure(error *RuntimeError) *RTResult {
	r.Reset()
	r.Error = error
	return r
}

func (r *RTResult) ShouldReturn() bool {
	return r.Error != nil || r.FuncReturnValue != nil || r.LoopShouldContinue || r.LoopShouldBreak
}

// NewContext creates a new context with the given display name, parent, and parent entry position.
func NewContext(displayName string, parent *Context, parentEntryPos *Position) *Context {
	return &Context{
		DisplayName:    displayName,
		Parent:         parent,
		ParentEntryPos: parentEntryPos,
		SymbolTable:    &SymbolTable{},
	}
}

// NewSymbolTable creates a new SymbolTable instance.
func NewSymbolTable(symboltable *SymbolTable) *SymbolTable {
	if symboltable == nil {

		return &SymbolTable{
			symbols:   make(map[string]*Value),
			constants: make(map[string]*Value),
			buildIn:   make(map[string]*Value),
			packages:  make(map[string]*Package),
			parent:    nil,
		}
	} else {
		return &SymbolTable{parent: symboltable, symbols: symboltable.symbols}
	}
}

// Get retrieves the value associated with the name from the symbol table. 1 arg Value, 2 exists, 3 isConstant
func (st *SymbolTable) Get(name string) (*Value, bool, bool) {
	if value, exists := st.symbols[name]; exists {
		return value, exists, false
	}
	if value, exists := st.constants[name]; exists {
		return value, exists, true
	}
	if value, exists := st.buildIn[name]; exists {
		return value, exists, false
	}
	if st.parent != nil {
		return st.parent.Get(name)
	}
	return nil, false, false
}

// GetPackageMethod retrieves the package associated with the name from the symbol table.
func (st *SymbolTable) GetPackageMethod(packageName string, methodName string) (*Value, bool) {
	if value, exists := st.packages[packageName].Methods[methodName]; exists {
		return value, exists
	}
	if st.parent != nil {
		return st.parent.GetPackageMethod(packageName, methodName)
	}
	return nil, false
}

// Set sets the value associated with the name in the symbol table.
func (st *SymbolTable) Set(name string, value *Value, isConst bool) *RuntimeError {
	if isConst {
		if _, exists := st.constants[name]; exists {
			return NewRTError(value.GetPosStart(), value.GetPosEnd(), fmt.Sprintf("Cannot reassign constant '%v'", name), value.GetContext())
		}
		st.constants[name] = value
	} else {
		st.symbols[name] = value
	}
	return nil
}

// SetBuildIn sets a build-in function in the symbol table.
func (st *SymbolTable) SetBuildIn(name string, value *Value) *RuntimeError {
	st.buildIn[name] = value
	return nil
}

// SetPackage sets a package into the symbol table.
func (st *SymbolTable) SetPackage(name string, value *Package) *RuntimeError {
	st.packages[name] = value
	return nil
}

// Remove removes the entry associated with the name from the symbol table.
func (st *SymbolTable) Remove(name string) {
	delete(st.symbols, name)
}

// Set sets the value associated with the name in the symbol table.
func (st *SymbolTable) Contains(name string) bool {
	_, exists, _ := st.Get(name)
	return exists
}

func IsInt(value interface{}) bool {
	switch value.(type) {
	case int:
		return true
	default:
		return false
	}
}

func NewEmptyValue() *Value {
	return &Value{}
}
