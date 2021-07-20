package goparser

func IsExported(name string) bool {
	return name[0] >= 'A' && name[0] <= 'Z'
}


func ParseTypeFromDecl(decl interface{}) (typ Type) {

	return
}