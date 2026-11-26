// Package generic holds the genny templates used to generate concrete
// losertree variants. GTYPE is replaced with the key type by the directives
// below; edit the .tmpl files and run `go generate ./genericTemplate`.
//
// A new key type needs a type-specific template for deadValGTYPE and the
// battle functions (see generic_int.tmpl); generic.tmpl and
// generic_unsafe.tmpl are reused.
package generic

//go:generate genny -in=generic.tmpl -pkg=losertree -out=../generic_int.go gen "GTYPE=int"
//go:generate genny -in=generic_int.tmpl -pkg=losertree -out=../generic_int_spec.go gen "GTYPE=int"
//go:generate genny -in=generic_unsafe.tmpl -pkg=losertree -out=../generic_int_unsafe.go gen "GTYPE=int"
