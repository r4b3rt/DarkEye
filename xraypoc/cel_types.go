package xraypoc

import (
	"github.com/google/cel-go/checker/decls"
	"google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"strings"
)

//ConvertDeclType add comment
func ConvertDeclType(value string) *expr.Type {
	if strings.HasPrefix(value, "randomInt") {
		return decls.Int
	} else if strings.HasPrefix(value, "newReverse") {
		return decls.NewObjectType("xraypoc.celtypes.Reverse")
	} else {
		return decls.String
	}
}
