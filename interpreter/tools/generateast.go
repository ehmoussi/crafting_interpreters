package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func defineAst(outputDir string, interfaceName string, types []string) {
	log.Println("defineAst: start")
	path := filepath.Join(outputDir, strings.ToLower(interfaceName)+".go")
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	defer file.Close()
	file.WriteString("package golox\n")
	file.WriteString("\n")
	file.WriteString("type " + interfaceName + "[T any] interface {\n")
	file.WriteString("    accept(visitor " + interfaceName + "Visitor[T]) (T, error)\n")
	file.WriteString("}\n")
	file.WriteString("\n")
	structNames := make([]string, len(types))
	for i, classType := range types {
		structDefs := strings.Split(classType, ":")
		structName := strings.TrimSpace(structDefs[0])
		structNames[i] = structName
		// structName += "[T]"
		log.Printf("defineAst: Add %s\n", structName)
		fields := strings.TrimSpace(structDefs[1])
		defineType(file, interfaceName, structName, fields)
	}
	defineVisitor(file, interfaceName, structNames)
	log.Println("defineAst: end")
}

func defineVisitor(file *os.File, interfaceName string, structNames []string) {
	file.WriteString("type " + interfaceName + "Visitor[T any] interface {\n")
	for _, structName := range structNames {
		file.WriteString("    visit" + structName + interfaceName + "(" + strings.ToLower(interfaceName) + " *" + structName + "[T]) (T, error)\n")
	}
	file.WriteString("}\n")
	file.WriteString("\n")
}

func defineType(file *os.File, interfaceName string, structName string, fields string) {
	fieldList := strings.Split(fields, ", ")
	file.WriteString("type " + structName + "[T any] struct {\n")
	parameterList := []string{}
	fieldNames := []string{}
	for _, fieldDef := range fieldList {
		fieldDefList := strings.Split(fieldDef, " ")
		fieldName := fieldDefList[1]
		fieldNames = append(fieldNames, fieldName)
		fieldType := fieldDefList[0]
		if fieldType == "Object" {
			fieldType = "any"
		} else if fieldType == "Expr" || fieldType == "Stmt" {
			fieldType = fieldType + "[T]"
		} else if len(fieldType) > 5 && fieldType[0:5] == "List<" && fieldType[len(fieldType)-1] == '>' {
			fieldType = "[]" + fieldType[5:len(fieldType)-1] + "[T]"
		} else {
			fieldType = "*" + fieldType
		}
		parameterList = append(parameterList, fieldName+" "+fieldType)
		file.WriteString("    " + fieldName + " " + fieldType + "\n")
	}
	file.WriteString("}\n")
	file.WriteString("\n")
	// Create constructor function
	parameters := strings.Join(parameterList, ", ")
	file.WriteString("func New" + structName + "[T any](" + parameters + ") *" + structName + "[T] {\n")
	file.WriteString("    return &" + structName + "[T]{\n")
	for _, fieldName := range fieldNames {
		file.WriteString("        " + fieldName + ": " + fieldName + ",\n")
	}
	file.WriteString("    }\n")
	file.WriteString("}\n")
	file.WriteString("\n")
	// Implement interfaceName interface methods
	file.WriteString("func (e *" + structName + "[T]) accept(visitor " + interfaceName + "Visitor[T]) (T, error){\n")
	file.WriteString("    return visitor.visit" + structName + interfaceName + "(e)\n")
	file.WriteString("}\n")
	file.WriteString("\n")
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: generate_ast <output directory>")
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Assign   : Token name, Expr value",
		"Binary   : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : Object value",
		"Logical  : Expr left, Token operator, Expr right",
		"Unary    : Token operator, Expr right",
		"Variable : Token name",
	})
	defineAst(outputDir, "Stmt", []string{
		"Block      : List<Stmt> statements",
		"Expression : Expr expression",
		"If         : Expr condition, Stmt thenBranch, Stmt elseBranch",
		"Print      : Expr expression",
		"Var        : Token name, Expr initializer",
		"While      : Expr condition, Stmt body",
		// "For        : Stmt initializer, Expr condition, Expr increment, Stmt body",
	})
}
