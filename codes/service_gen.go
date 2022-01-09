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

	code, codeErr = svc.generateFileServiceNamespace()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceInternal()
	if codeErr != nil {
		err = fmt.Errorf("%s service generate failed, %v", svc.Name(), codeErr)
		return
	}
	file.AddCode(code)

	code, codeErr = svc.generateFileServiceBuild()
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

	code, codeErr = svc.generateFileServiceShutdown()
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
	v.Add("namespace", svc.Name())
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
		proxy.AddParam("ctx", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns"), "Context"))
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
			proxy.AddParam("param", typ)
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
			proxy.AddResult("v", typ)
		}
		proxy.AddResult("err", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/errors"), "CodeError"))
		// body
		body := gcg.Statements()
		body.Tab().Ident("proxy").Symbol(",").Space().Ident("proxyErr").Space().ColonEqual().Space().Token("ctx.App().ServiceProxy(ctx, namespace)").Line()
		body.Tab().Token("if proxyErr != nil {").Line()
		body.Tab().Tab().Token("err = errors.Warning(fmt.Sprintf(\"get %s proxy failed\", namespace)).WithCause(proxyErr)").Line()
		body.Tab().Tab().Return().Line()
		body.Tab().Token("}").Line()
		body.Tab().Token("arg, argErr := fns.NewArgument(param)").Line()
		body.Tab().Token("if argErr != nil {").Line()
		body.Tab().Tab().Token("err = errors.Warning(fmt.Sprintf(\"get %s proxy failed\", namespace)).WithCause(argErr)").Line()
		body.Tab().Tab().Return().Line()
		body.Tab().Token("}").Line()
		body.Tab().Token(fmt.Sprintf("r := proxy.Request(ctx, %s, arg)", fnName)).Line()
		if fn.HasResult() {
			typeToken := fn.Result.Type.CodeString()
			if fn.Result.Type.IsBuiltin() {
				body.Tab().Token(fmt.Sprintf("var x %s", typeToken)).Line()
			} else if fn.Result.Type.IsArray() {
				body.Tab().Token(fmt.Sprintf("x := make(%s, 0, 1)", typeToken)).Line()
			} else {
				if strings.Index(typeToken, "*") == 0 {
					typeToken = typeToken[1:]
				}
				body.Tab().Token(fmt.Sprintf("x := %s{}", typeToken)).Line()
			}
		} else {
			body.Tab().Token(fmt.Sprintf("x := %s{}", "json.RawMessage"), gcg.NewPackage("github.com/aacfactory/json")).Line()
		}

		body.Tab().Token("err = r.Get(ctx, &x)").Line()
		if fn.HasResult() {
			body.Tab().Token("if err == nil {").Line()
			if fn.Result.Type.IsBuiltin() {
				body.Tab().Tab().Token("v = x").Line()
			} else if fn.Result.Type.IsArray() {
				body.Tab().Tab().Token("v = x").Line()
			} else {
				body.Tab().Tab().Token("v = &x").Line()
			}
			body.Tab().Token("}").Line()
		}
		body.Tab().Return()
		proxy.Body(body)
		v.Add(proxy.Build(), gcg.Line())
	}
	code = v
	return
}

func (svc *Service) generateFileService() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Service")
	v.AddResult("v", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns"), "Service"))
	body := gcg.Statements()
	body.Tab().Ident("v").Space().Equal().Space().Token("&service{}").Line()
	body.Tab().Return()
	v.Body(body)
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceStruct() (code gcg.Code, err error) {
	code = gcg.Type("service", gcg.Struct().Build())
	return
}

func (svc *Service) generateFileServiceNamespace() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Namespace")
	v.Receiver("s", gcg.Star().Ident("service"))
	v.AddResult("v", gcg.String())
	body := gcg.Statements()
	body.Tab().Ident("v").Space().Equal().Space().Ident("namespace").Line()
	body.Tab().Return()
	v.Body(body)
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceInternal() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Internal")
	v.Receiver("s", gcg.Star().Ident("service"))
	v.AddResult("v", gcg.Ident("bool"))
	body := gcg.Statements()
	body.Tab().Ident("v").Space().Equal().Space().Literal(svc.Internal()).Line()
	body.Tab().Return()
	v.Body(body)
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceBuild() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Build")
	v.Receiver("s", gcg.Star().Ident("service"))
	v.AddParam("_", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/configuares"), "Config"))
	v.AddResult("err", gcg.Error())
	v.Body(gcg.Return())
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceHandle(fns []*Fn) (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Handle")
	v.Receiver("s", gcg.Star().Ident("service"))
	v.AddParam("ctx", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns"), "Context"))
	v.AddParam("fn", gcg.Ident("string"))
	v.AddParam("argument", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns"), "Argument"))
	v.AddResult("v", gcg.Token("interface{}"))
	v.AddParam("err", gcg.QualifiedIdent(gcg.NewPackage("github.com/aacfactory/errors"), "CodeError"))
	body := gcg.Statements()
	body.Tab().Token("switch fn {").Line()
	for _, fn := range fns {
		key, keyErr := fn.NameToConstName()
		if keyErr != nil {
			err = keyErr
			return
		}
		body.Tab().Token(fmt.Sprintf("case %s:", key)).Line()
		body.Tab().Tab().Token(fmt.Sprintf("ctx = fns.WithFn(ctx, %s)", key)).Line()
		// authorization
		if fn.HasAuthorization() {
			body.Tab().Tab().Token("// authorization").Line()
			body.Tab().Tab().Token("authorization, hasAuthorization := ctx.User().Authorization()").Line()
			body.Tab().Tab().Token("if !hasAuthorization {").Line()
			body.Tab().Tab().Tab().Token("err = errors.Unauthorized(\"authorization was not found in head\")").Line()
			body.Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
			body.Tab().Tab().Token("authorizationErr := ctx.App().Authorizations().Decode(ctx, authorization)").Line()
			body.Tab().Tab().Token("if authorizationErr != nil {").Line()
			body.Tab().Tab().Tab().Token("err = errors.Unauthorized(\"authorization is invalid\").WithCause(authorizationErr)").Line()
			body.Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
			// permission
			body.Tab().Tab().Token("// permission").Line()
			body.Tab().Tab().Token("permissionErr := ctx.App().Permissions().Validate(ctx, namespace, fn)").Line()
			body.Tab().Tab().Token("if permissionErr != nil {").Line()
			body.Tab().Tab().Tab().Token("err = errors.Forbidden(\"forbidden\").WithCause(permissionErr)").Line()
			body.Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
		}

		// param
		if fn.HasParam() {
			body.Tab().Tab().Token("// param").Line()
			if fn.Param.Type.IsBuiltin() {
				body.Tab().Tab().Token(fmt.Sprintf("var param %s", fn.Param.Type.Indent)).Line()
			} else if fn.Param.Type.IsStruct() {
				if fn.Param.InFile {
					body.Tab().Tab().Token(fmt.Sprintf("param := %s{}", fn.Param.Type.Struct.Name)).Line()
				} else {
					if fn.Param.Type.Import.Alias != "" {
						body.Tab().Tab().Token(fmt.Sprintf("param := %s.%s{}", fn.Param.Type.Import.Alias, fn.Param.Type.Struct.Name), gcg.NewPackageWithAlias(fn.Param.Type.Import.Path, fn.Param.Type.Import.Alias)).Line()
					} else {
						body.Tab().Tab().Token(fmt.Sprintf("param := %s.%s{}", fn.Param.Type.Import.Name, fn.Param.Type.Struct.Name), gcg.NewPackage(fn.Param.Type.Import.Path)).Line()
					}
				}
			}
			body.Tab().Tab().Token("scanErr := argument.As(&param)").Line()
			body.Tab().Tab().Token("if scanErr != nil {").Line()
			body.Tab().Tab().Tab().Token("err = errors.BadRequest(\"parse request body failed\").WithCause(scanErr).WithMeta(\"_key\", \"parse_body_failed\")").Line()
			body.Tab().Tab().Token("return").Line()
			body.Tab().Tab().Token("}").Line()
			if fn.HasValidate() {
				body.Tab().Tab().Token("validateErr := ctx.App().Validate(param)").Line()
				body.Tab().Tab().Token("if validateErr != nil {").Line()
				body.Tab().Tab().Tab().Token("err = validateErr").Line()
				body.Tab().Tab().Token("return").Line()
				body.Tab().Tab().Token("}").Line()
			}
		}

		// todo: cache

		// tx
		txKind, txOpts, hasTx := fn.HasTx()
		if hasTx {
			switch txKind {
			case "sql":
				body.Tab().Tab().Token("// tx").Line()
				if txOpts != nil && len(txOpts) > 0 {
					sqlTimeout := strings.TrimSpace(txOpts[0])
					sqlISO := "0"
					if len(txOpts) > 1 {
						sqlISO = strings.TrimSpace(txOpts[1])
					}
					body.Tab().Tab().Token(fmt.Sprintf("txBegErr := sql.TxBeginWithOption(ctx, sql.TxOption(\"%s\", %s)", sqlTimeout, sqlISO), gcg.NewPackage("github.com/aacfactory/fns-contrib/databases/sql")).Line()
				} else {
					body.Tab().Tab().Token("txBegErr := sql.TxBegin(ctx)", gcg.NewPackage("github.com/aacfactory/fns-contrib/databases/sql")).Line()
				}
				body.Tab().Tab().Token("if txBegErr != nil {").Line()
				body.Tab().Tab().Tab().Token("err = errors.Warning(\"begin sql tx failed\").WithCause(txBegErr)").Line()
				body.Tab().Tab().Tab().Token("return").Line()
				body.Tab().Tab().Token("}").Line()
			default:
				err = fmt.Errorf("tx kind is not supported, kind is %s", txKind)
				return
			}
		}
		// do
		body.Tab().Tab().Token("// handle").Line()
		if fn.HasResult() {
			if fn.HasParam() {
				body.Tab().Tab().Token(fmt.Sprintf("v, err = %s(ctx, param)", fn.FuncName)).Line()
			} else {
				body.Tab().Tab().Token(fmt.Sprintf("v, err = %s(ctx)", fn.FuncName)).Line()
			}
		} else {
			if fn.HasParam() {
				body.Tab().Tab().Token(fmt.Sprintf("err = %s(ctx, param)", fn.FuncName)).Line()
			} else {
				body.Tab().Tab().Token(fmt.Sprintf("err = %s(ctx)", fn.FuncName)).Line()
			}
		}
		// todo: cache

		// tx
		if hasTx {
			body.Tab().Tab().Token("if err != nil {").Line()
			body.Tab().Tab().Tab().Token("_ = sql.TxRollback(ctx)").Line()
			body.Tab().Tab().Token("} else {").Line()
			body.Tab().Tab().Tab().Token("txCmtErr := sql.TxCommit(ctx)").Line()
			body.Tab().Tab().Tab().Token("if txCmtErr != nil {").Line()
			body.Tab().Tab().Tab().Tab().Token("err = errors.Warning(\"commit sql tx failed\").WithCause(txCmtErr)").Line()
			body.Tab().Tab().Tab().Tab().Token("return").Line()
			body.Tab().Tab().Tab().Token("}").Line()
			body.Tab().Tab().Token("}").Line()
		}

	}
	body.Tab().Token("default:").Line()
	body.Tab().Tab().Token("err = errors.NotFound(fmt.Sprintf(\"%s/%s was not found\", namespace, fn))").Line()
	body.Tab().Token("}").Line()
	body.Tab().Token("return").Line()
	v.Body(body)
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceShutdown() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Shutdown")
	v.Receiver("s", gcg.Star().Ident("service"))
	v.AddResult("err", gcg.Error())
	v.Body(gcg.Return())
	code = v.Build()
	return
}

func (svc *Service) generateFileServiceDocument() (code gcg.Code, err error) {
	v := gcg.Func()
	v.Name("Document")
	v.Receiver("s", gcg.Star().Ident("service"))
	v.AddResult("doc", gcg.Star().QualifiedIdent(gcg.NewPackage("github.com/aacfactory/fns"), "ServiceDocument"))
	// todo
	body := gcg.Statements()
	if len(svc.fns) > 0 {
		objects := make(map[string]string)
		body.Ident("doc").Space().Equal().Space().Token(fmt.Sprintf("fns.NewServiceDocument(namespace, \"%s\")", svc.Description())).Line()
		for _, fn := range svc.fns {
			fnConstName, _ := fn.NameToConstName()
			body.Comments("// " + fnConstName)
			body.Token(fmt.Sprintf("%sDoc := fns.NewFnDocument(\"%s\", \"%s\", \"%s\", %v, %v)", fnConstName, fn.Name(), fn.Title(), fn.Description(), fn.HasAuthorization(), fn.HasDeprecated())).Line()
			if fn.HasParam() {
				if fn.Param.Type.IsBuiltin() {
					doc, docErr := fn.Param.MapToDocCode(fmt.Sprintf("%sArgDoc", fnConstName))
					if docErr != nil {
						err = docErr
						return
					}
					body.Add(doc)
					body.Token(fmt.Sprintf("%sDoc.Argument = %sArgDoc", fnConstName, fnConstName)).Line()
				} else {
					// todo: use struct map to code, delete type map to code
					objectKey := fn.Param.Type.ObjectKey()
					typeCode, typeCodeErr := fn.Param.Type.MapToDocCode(objects)
					if typeCodeErr != nil {
						err = typeCodeErr
						return
					}
					body.Add(typeCode)
					body.Token(fmt.Sprintf("%sDoc.Argument = %s", fnConstName, objectKey)).Line()
				}

			}
			if fn.HasResult() {
				if fn.Result.Type.IsBuiltin() {
					doc, docErr := fn.Result.MapToDocCode(fmt.Sprintf("%sResultDoc", fnConstName))
					if docErr != nil {
						err = docErr
						return
					}
					body.Add(doc)
					body.Token(fmt.Sprintf("%sDoc.Result = %sResultDoc", fnConstName, fnConstName)).Line()
				} else {
					objectKey := fn.Result.Type.ObjectKey()
					typeCode, typeCodeErr := fn.Result.Type.MapToDocCode(objects)
					if typeCodeErr != nil {
						err = typeCodeErr
						return
					}
					body.Add(typeCode)
					body.Token(fmt.Sprintf("%sDoc.Result = %s", fnConstName, objectKey)).Line()
				}

			}
			body.Token("doc.AddFn(*fnDoc)").Line()
		}
	}
	body.Return()
	v.Body(body)
	code = v.Build()
	return
}
