package xraypoc

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/xraypoc/celtypes"
	"google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

//Eval add comment
func Eval(env *cel.Env, expression string, params map[string]interface{}) (ref.Val, error) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	out, _, err := prg.Eval(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//CustomLib add comment
type CustomLib struct {
	envOptions []cel.EnvOption
}

//CompileOptions add comment
func (cus *CustomLib) CompileOptions() []cel.EnvOption {
	return cus.envOptions
}

//UpdateOptions add comment
func (cus *CustomLib) UpdateOptions(param string, declType *expr.Type) {
	//动态注入类型
	opts := cel.Declarations(
		decls.NewVar(param, declType),
	)
	cus.envOptions = append(cus.envOptions, opts)
}

//ProgramOptions add comment
func (cus *CustomLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "reverse_wait_int",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					reverse, ok := lhs.Value().(*xraypoc_celtypes.Reverse)
					if !ok {
						return types.ValOrErr(lhs, "值校验 '%v'.reverse_wait", lhs.Type())
					}
					timeout, ok := rhs.Value().(int64)
					if !ok {
						return types.ValOrErr(rhs, "值校验 reverse_wait('%v')", rhs.Type())
					}
					return types.Bool(myReverseCheck(reverse, timeout))
				},
			},
			&functions.Overload{
				Operator: "bytes_bcontains_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "值校验 '%v'.%s", lhs.Type(), "bytes_bcontains_bytes")
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "值校验 %s('%v')", "bytes_bcontains_bytes", rhs.Type())
					}
					return types.Bool(bytes.Contains(v1, v2))
				},
			},
			&functions.Overload{
				Operator: "string_bmatches_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "值校验 '%v'.%s", lhs.Type(), "string_bmatches_bytes")
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "值校验 %s('%v')", "string_bmatches_bytes", rhs.Type())
					}
					return types.Bool(bytes.Contains([]byte(v1), v2))
				},
			},
			&functions.Overload{
				Operator: "string_startsWith_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "值校验 '%v'.%s", lhs.Type(), "string_startsWith_string")
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "值校验 %s('%v')", "string_startsWith_string", rhs.Type())
					}
					return types.Bool(strings.HasPrefix(string(v1), string(v2)))
				},
			},
			&functions.Overload{
				Operator: "string_endsWith_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "值校验 '%v'.%s", lhs.Type(), "string_endsWith_string")
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "值校验 %s('%v')", "string_endsWith_string", rhs.Type())
					}
					return types.Bool(strings.HasPrefix(string(v1), string(v2)))
				},
			},
			&functions.Overload{
				Operator: "md5_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "md5_string", value.Type())
					}
					return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
				},
			},
			&functions.Overload{
				Operator: "randomInt_int_int",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					min, ok := lhs.(types.Int)
					if !ok {
						return types.ValOrErr(lhs, "值校验 randomInt('%v')", lhs.Type())
					}
					max, ok := rhs.(types.Int)
					if !ok {
						return types.ValOrErr(rhs, "值校验 randomInt('%v')", rhs.Type())
					}
					return types.Int(rand.Intn(int(max)-int(min)) + int(min))
				},
			},
			&functions.Overload{
				Operator: "randomLowercase_int",
				Unary: func(value ref.Val) ref.Val {
					n, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "randomLowercase_int", value.Type())
					}
					result := ""
					r := rand.New(rand.NewSource(time.Now().Unix()))
					i := 0
					for i < int(n) {
						result += string(common.LowCaseAlpha[r.Int()%len(common.LowCaseAlpha)])
						i++
					}
					return types.String(result)
				},
			},
			&functions.Overload{
				Operator: "base64_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "base64_string", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
				},
			},
			&functions.Overload{
				Operator: "base64_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "base64_bytes", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString(v))
				},
			},
			&functions.Overload{
				Operator: "base64Decode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "base64Decode_string", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "base64Decode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "base64Decode_bytes", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "urlencode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "urlencode_string", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urlencode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "urlencode_bytes", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urldecode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "urlencode_string", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "urldecode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "值校验 %s('%v')", "urldecode_bytes", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "substr_string_int_int",
				Function: func(values ...ref.Val) ref.Val {
					if len(values) != 3 {
						return types.NewErr("'substr' 参数应该是3个")
					}
					str, ok := values[0].(types.String)
					if !ok {
						return types.ValOrErr(values[0], "值校验 %s('%v')", "substr_string_int_int", values[0].Type())
					}
					start, ok := values[1].(types.Int)
					if !ok {
						return types.ValOrErr(values[1], "值校验 %s('%v')", "substr_string_int_int", values[1].Type())
					}
					length, ok := values[2].(types.Int)
					if !ok {
						return types.ValOrErr(values[2], "值校验 %s('%v')", "substr_string_int_int", values[2].Type())
					}
					bstr := []byte(str)
					if start < 0 || length < 0 || int(start+length) > len(str) {
						return types.NewErr("'substr' 参数值不合法")
					}
					return types.String(bstr[start : start+length])
				},
			},
			&functions.Overload{
				Operator: "string_icontains_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "值校验 %s('%v')", "string_icontains_string", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "值校验 %s('%v')", "string_icontains_string", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(strings.Contains(strings.ToLower(string(v1)), strings.ToLower(string(v2))))
				},
			},
		),
		//添加更多functions
	}
}

//NewCustomLib add comment
func NewCustomLib() *CustomLib {
	return &CustomLib{
		envOptions: []cel.EnvOption{
			cel.Container("xraypoc.celtypes"),
			//
			cel.Types(&xraypoc_celtypes.Response{}),
			cel.Types(&xraypoc_celtypes.Request{}),
			cel.Types(&xraypoc_celtypes.Reverse{}),
			cel.Types(&xraypoc_celtypes.UrlType{}),
			cel.Declarations(
				//注入类型
				decls.NewVar("response", decls.NewObjectType("xraypoc.celtypes.Response")),
				decls.NewVar("request", decls.NewObjectType("xraypoc.celtypes.Request")),
				//注入Reverse函数
				decls.NewFunction("wait",
					decls.NewInstanceOverload("reverse_wait_int",
						[]*expr.Type{decls.Any, decls.Int},
						decls.Bool)),
				//注入比较函数
				//Note: contains/matches 已经存在不需要注入
				decls.NewFunction("bcontains",
					decls.NewInstanceOverload("bytes_bcontains_bytes",
						[]*expr.Type{decls.Bytes, decls.Bytes},
						decls.Bool)),
				decls.NewFunction("bmatches",
					decls.NewInstanceOverload("string_bmatches_bytes",
						[]*expr.Type{decls.String, decls.Bytes},
						decls.Bool)),
				//功能函数
				decls.NewFunction("startsWith",
					decls.NewOverload("string_startsWith_string",
						[]*expr.Type{decls.String, decls.String},
						decls.Bool)),
				decls.NewFunction("endsWith",
					decls.NewOverload("string_endsWith_string",
						[]*expr.Type{decls.String, decls.String},
						decls.Bool)),
				//todo NewFunction("in")
				decls.NewFunction("md5",
					decls.NewOverload("md5_string",
						[]*expr.Type{decls.String},
						decls.String)),
				decls.NewFunction("randomInt",
					decls.NewOverload("randomInt_int_int",
						[]*expr.Type{decls.Int, decls.Int},
						decls.Int)),
				//随机n长度小写字母
				decls.NewFunction("randomLowercase",
					decls.NewOverload("randomLowercase_int",
						[]*expr.Type{decls.Int},
						decls.String)),
				decls.NewFunction("base64",
					decls.NewOverload("base64_string",
						[]*expr.Type{decls.String},
						decls.String)),
				decls.NewFunction("base64",
					decls.NewOverload("base64_bytes",
						[]*expr.Type{decls.Bytes},
						decls.String)),
				decls.NewFunction("base64Decode",
					decls.NewOverload("base64Decode_string",
						[]*expr.Type{decls.String},
						decls.String)),
				decls.NewFunction("base64Decode",
					decls.NewOverload("base64Decode_bytes",
						[]*expr.Type{decls.Bytes},
						decls.String)),
				decls.NewFunction("urlencode",
					decls.NewOverload("urlencode_string",
						[]*expr.Type{decls.String},
						decls.String)),
				decls.NewFunction("urlencode",
					decls.NewOverload("urlencode_bytes",
						[]*expr.Type{decls.Bytes},
						decls.String)),
				decls.NewFunction("urldecode",
					decls.NewOverload("urldecode_string",
						[]*expr.Type{decls.String},
						decls.String)),
				decls.NewFunction("urldecode",
					decls.NewOverload("urldecode_bytes",
						[]*expr.Type{decls.Bytes},
						decls.String)),
				decls.NewFunction("substr",
					decls.NewOverload("substr_string_int_int",
						[]*expr.Type{decls.String, decls.Int, decls.Int},
						decls.String)),
				decls.NewFunction("icontains",
					decls.NewInstanceOverload("string_icontains_string",
						[]*expr.Type{decls.String, decls.String},
						decls.Bool)),
			),
		},
	}
}
