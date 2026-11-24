package generic

//go:generate genny -in=generic.tmpl -pkg=losertree -out=../generic_int.go gen "GTYPE=int"
//go:generate genny -in=generic_int.tmpl -pkg=losertree -out=../generic_int_spec.go gen "GTYPE=int"
//go:generate genny -in=generic_unsafe.tmpl -pkg=losertree -out=../generic_int_unsafe.go gen "GTYPE=int"
