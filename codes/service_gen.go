/*
 * Copyright 2021 Wang Min Xiang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package codes

import (
	"fmt"
	"github.com/aacfactory/cases"
	"github.com/aacfactory/gcg"
	"path/filepath"
	"strings"
)

func (svc *Service) generateFile() (err error) {

	fns := svc.Fns()
	if len(fns) == 0 {
		return
	}

	file := gcg.NewFile(svc.Package)

	// imports
	for _, _import := range svc.Imports {
		if _import.Alias != "" {
			file.AddImport(gcg.NewPackageWithAlias(_import.Path, _import.Alias))
		} else {
			file.AddImport(gcg.NewPackage(_import.Path))
		}
	}
	var code gcg.Code = nil
	var codeErr error = nil

	code, codeErr = svc.generateFileConst(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceProxy(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	if len(svc.Components) > 0 {
		code, codeErr = svc.generateFileServiceComponents()
		if codeErr != nil {
			err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
			return
		}
		file.AddCode(code)
	}

	code, codeErr = svc.generateFileService()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceStruct()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceHandle(fns)
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceDocument()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	// write
	filename := filepath.Join(svc.DirPath, "fns.go")
	w := gcg.FileRender(filename, true)
	wErr := file.Render(w)
	if wErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), wErr)
		return
	}
	closeFileErr := w.Close()
	if closeFileErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), closeFileErr)
		return
	}

	return
}

func (svc *Service) generateFileConst(fns []*Fn) (code gcg.Code, err error) {
	v := gcg.Constants()
	v.Add("_name", svc.Name())
	for _, fn := range fns {
		key, keyErr := fn.NameToConstName()
		if keyErr != nil {
			err = keyErr
			return
		}
		v.Add(key, fn.Name())
	}
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceProxy(fns []*Fn) (code gcg.Code, err error) {
	v := gcg.Statements()
	for _, fn := range fns {
		fnName, fnNameErr := fn.NameToConstName()
		if fnNameErr != nil {
			err = fnNameErr
			return
		}
		funcName, funcNameErr := fn.NameToProxyName()
		if funcNameErr != nil {
			err = funcNameErr
			return
		}
		proxy := gcg.Func()
		proxy.Name(funcName)
		proxy.AddParam("ctx", gcg.QualifiedIdent(gcg.NewPackage("context"), "Context"))
		if fn.HasParam() {
			var typ gcg.Code = nil
			if fn.Param.Type.IsBuiltin() {
				typ = gcg.Ident(fn.Param.Type.Indent)
			} else if fn.Param.Type.IsStruct() {
				if fn.Param.InFile {
					typ = gcg.Ident(fn.Param.Type.Struct.Name)
				} else {
					if fn.Param.Type.Import.Alias != "" {
						typ = gcg.QualifiedIdent(gcg.NewPackageWithAlias(fn.Param.Type.Import.Path, fn.Param.Type.Import.Alias), fn.Param.Type.Struct.Name)
					} else {
						typ = gcg.QualifiedIdent(gcg.NewPackage(fn.Param.Type.Import.Path), fn.Param.Type.Struct.Name)
					}
				}
			} else {
				err = fmt.Errorf("%s.%s param is invalid, type must be value object", svc.Name(), fn.FuncName)
				return
			}
			proxy.AddParam("argument", typ)
		}
		if fn.HasResult() {
			var typ gcg.Code
			typToken := fn.Result.Type.CodeString()
			if typToken == "" {
				err = fmt.Errorf("%s.%s result is invalid, type must be star object or array", svc.Name(), fn.FuncName)
				return
			}
			import_, hasImport := fn.Result.Type.GetImport()
			if hasImport {
				if import_.Alias == "" {
					typ = gcg.Token(typToken, gcg.NewPackage(import_.Path))
				} else {
					typ = gcg.Token(typToken, gcg.NewPackageWithAlias(import_.Path, import_.Alias))
				}
			} else {
				typ = gcg.Token(typToken)
			}
			proxy.AddResult("result", typ)
		}
		proxy.AddResult("err", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/errors"), "CodeError"))
		// body
		body := gcg.Statements()
		body.Tab().Ident("endpoint").Symbol(",").Space().Ident("hasEndpoint").Space().ColonEqual().Space().Token("service.GetEndpoint(ctx, _name)").Line()
		body.Tab().Token("if !hasEndpoint {").Line()
		body.Tab().Tab().Token("err = errors.NotFound(\"endpoint was not found\").WithMeta(\"name\", _name)").Line()
		body.Tab().Tab().Return().Line()
		body.Tab().Token("}").Line()
		body.Tab().Token(fmt.Sprintf("fr := endpoint.Request(ctx, %s, service.NewArgument(argument))", fnName)).Line()
		if fn.HasResult() {
			typeToken := fn.Result.Type.CodeString()
			if fn.Result.Type.IsBuiltin() {
				body.Tab().Token(fmt.Sprintf("var handled %s", typeToken)).Line()
			} else if fn.Result.Type.IsArray() {
				body.Tab().Token(fmt.Sprintf("handled := make(%s, 0, 1)", typeToken)).Line()
			} else {
				if strings.Index(typeToken, "*") == 0 {
					typeToken = typeToken[1:]
				}
				body.Tab().Token(fmt.Sprintf("handled := %s{}", typeToken)).Line()
			}
			body.Tab().Token("hasResult, handleErr := fr.Get(ctx, &handled)").Line()
			body.Tab().Token("if handleErr != nil {").Line()
			body.Tab().Tab().Token("err = handleErr").Line()
			body.Tab().Tab().Return().Line()
			body.Tab().Token("}").Line()
			body.Tab().Token("if hasResult {").Line()
			if fn.Result.Type.IsBuiltin() {
				body.Tab().Tab().Token("result = handled").Line()
			} else if fn.Result.Type.IsArray() {
				body.Tab().Tab().Token("result = handled").Line()
			} else {
				body.Tab().Tab().Token("result = &handled").Line()
			}
			body.Tab().Token("}").Line()
		} else {
			body.Tab().Token("_, handleErr := fr.Get(ctx, &json.RawMessage{})", gcg.NewPackage("github.com/aacfactory/json")).Line()
			body.Tab().Token("if handleErr != nil {").Line()
			body.Tab().Tab().Token("err = handleErr").Line()
			body.Tab().Tab().Return().Line()
			body.Tab().Token("}").Line()
		}
		body.Tab().Return()
		proxy.Body(body)
		v.Add(proxy.Build(), gcg.Line())
	}
	code = v
	return
}

func (svc *Service) generateFileServiceComponents() (code gcg.Code, err error) {
	stmt := gcg.Statements()
	for _, component := range svc.Components {
		v := gcg.Func()
		v.Name(cases.LowerCamel().Format([]string{"get", "component", component.Name}))
		v.AddParam("ctx", gcg.QualifiedIdent(gcg.NewPackage("context"), "Context"))
		structName := component.Struct
		var resultType gcg.Code
		var resultStruct string
		if strings.Contains(structName, ".") {
			sidx := strings.LastIndexByte(structName, '.')
			resultStructPkg := structName[0:sidx]
			resultStructName := structName[sidx+1:]
			resultType = gcg.Star().QualifiedIdent(gcg.NewPackage(resultStructPkg), resultStructName)
			fidx := strings.LastIndexByte(structName, '/')
			resultStruct = structName[fidx+1:]
		} else {
			resultType = gcg.Token(fmt.Sprintf("*%s", structName))
			resultStruct = structName
		}
		v.AddResult("v", resultType)
		body := gcg.Statements()
		body.Tab().Token(fmt.Sprintf("c, has := service.GetComponent(ctx, \"%s\")", component.Name)).Line()
		body.Tab().Token("if !has {").Line()
		body.Tab().Tab().Token("panic(fmt.Sprintf(\"%+v\"")
		body.Tab().Tab().Token(fmt.Sprintf(", errors.Warning(\"%s: get %s component failed cause not found in context\")))", svc.Name(), component.Name)).Line()
		body.Tab().Tab().Return().Line()
		body.Tab().Token("}").Line()
		body.Tab().Token("ok := false").Line()
		body.Tab().Token(fmt.Sprintf("v, ok = c.(*%s)", resultStruct)).Line()
		body.Tab().Token("if !ok {").Line()
		body.Tab().Tab().Token("panic(fmt.Sprintf(\"%+v\"")
		body.Tab().Tab().Token(fmt.Sprintf(", errors.Warning(\"%s: get %s component failed cause type is not matched\")))", svc.Name(), component.Name)).Line()
		body.Tab().Tab().Return().Line()
		body.Tab().Token("}").Line()
		body.Tab().Tab().Return()
		v.Body(body)
		stmt.Add(v.Build()).Line()
	}
	code = stmt
	return
}

func (svc *Service) generateFileService() (code gcg.Code, err error) {
	stmt := gcg.Statements()
	v := gcg.Func()
	v.Name("Service")
	v.AddResult("svc", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns/service"), "Service"))
	body := gcg.Statements()
	// todo components
	body.Tab().Token("components := []service.Component{", gcg.NewPackage("github.com/aacfactory/fns/service"))
	if len(svc.Components) > 0 {
		cidx := 0
		for _, component := range svc.Components {
			if cidx > 0 {
				body.Token(", ")
			}
			body.Token(fmt.Sprintf("%s()", component.Loader))
			cidx++
		}
	}
	body.Token("}").Line()
	body.Tab().Ident("svc").Space().Equal().Space().Token("&_service_{").Line()
	body.Tab().Tab().Token("Abstract:service.NewAbstract(").Line()
	body.Tab().Tab().Tab().Token("_name").Ident(",").Line()
	body.Tab().Tab().Tab().Token(fmt.Sprintf("%v", svc.Internal())).Ident(",").Line()
	body.Tab().Tab().Tab().Token("components...").Ident(",").Line()
	body.Tab().Tab().Token("),").Line()
	body.Tab().Token("}").Line()
	body.Tab().Return()
	v.Body(body)
	stmt.Add(v.Build()).Line()
	code = stmt
	return
}

func (svc *Service) generateFileServiceStruct() (code gcg.Code, err error) {
	v := gcg.Statements()
	v.Token("type _service_ struct {").Line()
	v.Tab().Token("service.Abstract").Line()
	v.Token("}").Line()
	code = v
	return
}

func (svc *Service) generateFileServiceHandle(fns []*Fn) (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Handle")
	v.Receiver("svc", gcg.Star().Ident("_service_"))
	v.AddParam("ctx", gcg.QualifiedIdent(gcg.NewPackage("context"), "Context"))
	v.AddParam("fn", gcg.Ident("string"))
	v.AddParam("argument", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns/service"), "Argument"))
	v.AddResult("v", gcg.Token("interface{}"))
	v.AddResult("err", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/errors"), "CodeError"))
	body := gcg.Statements()
	body.Tab().Token("switch fn {").Line()
	for _, fn := range fns {
		key, keyErr := fn.NameToConstName()
		if keyErr != nil {
			err = keyErr
			return
		}
		body.Tab().Token(fmt.Sprintf("case %s:", key)).Line()
		if fn.IsInternal() {
			body.Tab().Tab().Token("if !service.CanAccessInternal(ctx) {").Line()
			body.Tab().Tab().Tab().Token(fmt.Sprintf("err = errors.NotAcceptable(\"%s: %s cannot be accessed externally\")", svc.Name(), fn.Name())).Line()
			body.Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
		}
		// authorization
		if fn.HasAuthorization() {
			body.Tab().Tab().Token("// verify authorizations").Line()
			body.Tab().Tab().Token("verifyAuthorizationErr := authorizations.Verify(ctx)", gcg.NewPackage("github.com/aacfactory/fns/endpoints/authorizations")).Line()
			body.Tab().Tab().Token("if verifyAuthorizationErr != nil {").Line()
			body.Tab().Tab().Tab().Token(fmt.Sprintf("err = verifyAuthorizationErr.WithMeta(\"service\", _name).WithMeta(\"fn\", %s)", key)).Line()
			body.Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
		}
		// permission
		if fn.HasPermission() {
			allow := fn.Annotations["permission"]
			if allow == "" {
				err = fmt.Errorf("there is no roles in permission")
				return
			}
			allows := strings.Split(allow, ",")
			roles := make([]string, 0, 1)
			for _, s := range allows {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				roles = append(roles, fmt.Sprintf("\"%s\"", s))
			}
			if len(roles) == 0 {
				err = fmt.Errorf("there is no roles in permission")
				return
			}
			body.Tab().Tab().Token("// permission").Line()
			body.Tab().Tab().Token(fmt.Sprintf("verifyPermissionsErr := permissions.Verify(ctx, %s)", strings.Join(roles, ",")), gcg.NewPackage("github.com/aacfactory/fns/endpoints/permissions")).Line()
			body.Tab().Tab().Token("if verifyPermissionsErr != nil {").Line()
			body.Tab().Tab().Tab().Token(fmt.Sprintf("err = verifyPermissionsErr.WithMeta(\"service\", _name).WithMeta(\"fn\", %s)", key)).Line()
			body.Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
		}

		// param
		if fn.HasParam() {
			body.Tab().Tab().Token("// argument").Line()
			if fn.Param.Type.IsBuiltin() {
				body.Tab().Tab().Token(fmt.Sprintf("var arg %s", fn.Param.Type.Indent)).Line()
			} else if fn.Param.Type.IsStruct() {
				if fn.Param.InFile {
					body.Tab().Tab().Token(fmt.Sprintf("arg := %s{}", fn.Param.Type.Struct.Name)).Line()
				} else {
					if fn.Param.Type.Import.Alias != "" {
						body.Tab().Tab().Token(fmt.Sprintf("arg := %s.%s{}", fn.Param.Type.Import.Alias, fn.Param.Type.Struct.Name), gcg.NewPackageWithAlias(fn.Param.Type.Import.Path, fn.Param.Type.Import.Alias)).Line()
					} else {
						body.Tab().Tab().Token(fmt.Sprintf("arg := %s.%s{}", fn.Param.Type.Import.Name, fn.Param.Type.Struct.Name), gcg.NewPackage(fn.Param.Type.Import.Path)).Line()
					}
				}
			}
			body.Tab().Tab().Token("scanErr := argument.As(&arg)").Line()
			body.Tab().Tab().Token("if scanErr != nil {").Line()
			body.Tab().Tab().Tab().Token(fmt.Sprintf("err = errors.BadRequest(\"%s: scan request argument failed\").WithCause(scanErr).WithMeta(\"service\", _name).WithMeta(\"fn\", %s)", svc.Name(), key)).Line()
			body.Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
			if fn.HasValidate() {
				body.Tab().Tab().Token("validateErr := validators.Validate(arg)", gcg.NewPackage("github.com/aacfactory/fns/service/validators")).Line()
				body.Tab().Tab().Token("if validateErr != nil {").Line()
				body.Tab().Tab().Tab().Token(fmt.Sprintf("err = errors.BadRequest(\"%s: invalid request argument\").WithMeta(\"service\", _name).WithMeta(\"fn\", %s).WithCause(validateErr)", svc.Name(), key)).Line()
				body.Tab().Tab().Token("return").Line()
				body.Tab().Tab().Token("}").Line()
			}
		}

		// tx
		txKind, _, hasTx := fn.HasTx()
		if hasTx {
			switch txKind {
			case "sql":
				body.Tab().Tab().Token("// sql begin transaction").Line()
				body.Tab().Token("beginTransactionErr := sql.BeginTransaction(ctx)", gcg.NewPackage("github.com/aacfactory/fns-contrib/databases/sql")).Line()
				body.Tab().Tab().Token("if beginTransactionErr != nil {").Line()
				body.Tab().Tab().Tab().Token(fmt.Sprintf("err = errors.ServiceError(\"%s: begin sql transaction failed\").WithMeta(\"service\", _name).WithMeta(\"fn\", %s).WithCause(beginTransactionErr)", svc.Name(), key)).Line()
				body.Tab().Tab().Tab().Token("return").Line()
				body.Tab().Tab().Token("}").Line()
			default:
				err = fmt.Errorf("tx kind is not supported, kind is %s", txKind)
				return
			}
		}
		// do
		body.Tab().Tab().Token("// handle ").Line()
		if fn.HasResult() {
			if fn.HasParam() {
				body.Tab().Tab().Token(fmt.Sprintf("v, err = %s(ctx, arg)", fn.FuncName)).Line()
			} else {
				body.Tab().Tab().Token(fmt.Sprintf("v, err = %s(ctx)", fn.FuncName)).Line()
			}
		} else {
			if fn.HasParam() {
				body.Tab().Tab().Token(fmt.Sprintf("err = %s(ctx, arg)", fn.FuncName)).Line()
			} else {
				body.Tab().Tab().Token(fmt.Sprintf("err = %s(ctx)", fn.FuncName)).Line()
			}
		}
		// tx
		if hasTx {
			body.Tab().Tab().Token("// sql close transaction").Line()
			body.Tab().Tab().Token("if err == nil {").Line()
			body.Tab().Tab().Tab().Token("commitTransactionErr := sql.CommitTransaction(ctx)").Line()
			body.Tab().Tab().Tab().Token("if commitTransactionErr != nil {").Line()
			body.Tab().Tab().Tab().Tab().Token(fmt.Sprintf("err = errors.ServiceError(\"%s: commit sql transaction failed\").WithMeta(\"service\", _name).WithMeta(\"fn\", %s).WithCause(commitTransactionErr)", svc.Name(), key)).Line()
			body.Tab().Tab().Tab().Tab().Token("_ = sql.RollbackTransaction(ctx)").Line()
			body.Tab().Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Tab().Token("}").Line()
			body.Tab().Tab().Token("}").Line()
		}
		body.Tab().Tab().Token("break").Line()
	}
	body.Tab().Token("default:").Line()
	body.Tab().Tab().Token(fmt.Sprintf("err = errors.NotFound(\"%s: fn was not found\").WithMeta(\"service\", _name).WithMeta(\"fn\", fn)", svc.Name())).Line()
	body.Tab().Tab().Token("break").Line()
	body.Tab().Token("}").Line()
	body.Tab().Token("return").Line()
	v.Body(body)
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceDocument() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Document")
	v.Receiver("svc", gcg.Star().Ident("_service_"))
	v.AddResult("doc", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns/service"), "Document"))
	body := gcg.Statements()
	if !svc.Internal() {
		if len(svc.fns) > 0 {
			// service
			body.Token(fmt.Sprintf("sd := documents.NewService(_name, \"%s\")", svc.Description()), gcg.NewPackage("github.com/aacfactory/fns/service/documents")).Line()
			// fn
			for _, fn := range svc.fns {
				if fn.IsInternal() {
					continue
				}
				body.Token(fmt.Sprintf("sd.AddFn(")).Line()
				body.Token(fmt.Sprintf("\"%s\", \"%s\", \"%s\",%v, %v,", fn.Name(), fn.Title(), fn.Description(), fn.HasAuthorization(), fn.HasDeprecated())).Line()
				if fn.HasParam() {
					body.Add(fn.Param.generateObjectDocument()).Symbol(",").Line()
				}
				if fn.HasResult() {
					body.Add(fn.Result.generateObjectDocument()).Symbol(",").Line()
				}
				body.Token(")").Line()
			}
			body.Token("doc = sd").Line()
		}
	}
	body.Return()
	v.Body(body)
	code = v.Build()
	return
}
