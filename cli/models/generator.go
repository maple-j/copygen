package models

// Generator represents a code generator.
type Generator struct {
	GenFile    string     // The generated filepath.
	GenPackage string     // The generated package.
	Imports    []string   // The imports included in the generated file.
	Functions  []Function // The functions to generate.
}