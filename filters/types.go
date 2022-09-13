package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cockscomb/cel2sql"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

const (
	ExistsEquals     = "existsEquals"
	ExistsEqualsCI   = "existsEqualsCI"
	ExistsStarts     = "existsStarts"
	ExistsStartsCI   = "existsStartsCI"
	ExistsEnds       = "existsEnds"
	ExistsEndsCI     = "existsEndsCI"
	ExistsContains   = "existsContains"
	ExistsContainsCI = "existsContainsCI"
	ExistsRegexp     = "existsRegexp"
	ExistsRegexpCI   = "existsRegexpCI"
)

var ciFuncs = map[string]struct{}{
	ExistsEqualsCI:   {},
	ExistsStartsCI:   {},
	ExistsEndsCI:     {},
	ExistsContainsCI: {},
	ExistsRegexpCI:   {},
}

var Declarations = cel.Declarations(
	decls.NewFunction(ExistsEquals,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsEqualsCI,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsStarts,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsStartsCI,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsEnds,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsEndsCI,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsContains,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsContainsCI,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsRegexp,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
	decls.NewFunction(ExistsRegexpCI,
		decls.NewInstanceOverload("string_to_string", []*expr.Type{decls.String, decls.String}, decls.Bool),
		decls.NewInstanceOverload("string_to_list", []*expr.Type{decls.String, decls.NewListType(decls.String)}, decls.Bool),
		decls.NewInstanceOverload("list_to_string", []*expr.Type{decls.NewListType(decls.String), decls.String}, decls.Bool),
		decls.NewInstanceOverload("list_to_list", []*expr.Type{decls.NewListType(decls.String), decls.NewListType(decls.String)}, decls.Bool),
	),
)

type Extension struct{}

func (ext *Extension) ImplementsFunction(fun string) bool {
	switch fun {
	case ExistsEquals, ExistsEqualsCI, ExistsStarts, ExistsStartsCI, ExistsEnds, ExistsEndsCI, ExistsContains, ExistsContainsCI, ExistsRegexp, ExistsRegexpCI:
		return true
	}
	return false
}

func (ext *Extension) CallFunction(con *cel2sql.Converter, function string, target *expr.Expr, args []*expr.Expr) error {
	tgtType := con.GetType(target)
	argType := con.GetType(args[0])
	switch function {
	case ExistsEquals, ExistsEqualsCI:
		switch {
		case tgtType.GetPrimitive() == expr.Type_STRING:
			if err := writeTarget(con, function, target); err != nil {
				return err
			}
			switch {
			case argType.GetPrimitive() == expr.Type_STRING:
				con.WriteString(" = ")
				return con.Visit(args[0])
			case cel2sql.IsListType(argType):
				con.WriteString(" IN UNNEST(")
				if err := con.Visit(args[0]); err != nil {
					return err
				}
				con.WriteString(")")
				return nil
			}
		case cel2sql.IsListType(tgtType):
			switch {
			case argType.GetPrimitive() == expr.Type_STRING:
				return ext.CallFunction(con, function, args[0], []*expr.Expr{target})
			case cel2sql.IsListType(argType):
				return ext.callRegexp(con, target, args, regexpOptions{caseInsensitive: function == ExistsEqualsCI, start: true, end: true, regexEscape: true})
			}
		}
	case ExistsStarts, ExistsStartsCI:
		if tgtType.GetPrimitive() == expr.Type_STRING && argType.GetPrimitive() == expr.Type_STRING {
			if err := writeSimpleCall("STARTS_WITH", con, function, target, args[0]); err != nil {
				return err
			}
			return nil
		}
		return ext.callRegexp(con, target, args, regexpOptions{caseInsensitive: function == ExistsStartsCI, start: true, regexEscape: true})
	case ExistsEnds, ExistsEndsCI:
		if tgtType.GetPrimitive() == expr.Type_STRING && argType.GetPrimitive() == expr.Type_STRING {
			if err := writeSimpleCall("ENDS_WITH", con, function, target, args[0]); err != nil {
				return err
			}
			return nil
		}
		return ext.callRegexp(con, target, args, regexpOptions{caseInsensitive: function == ExistsEndsCI, end: true, regexEscape: true})
	case ExistsContains, ExistsContainsCI:
		if tgtType.GetPrimitive() == expr.Type_STRING && argType.GetPrimitive() == expr.Type_STRING {
			if err := writeSimpleCall("0 != INSTR", con, function, target, args[0]); err != nil {
				return err
			}
			return nil
		}
		return ext.callRegexp(con, target, args, regexpOptions{caseInsensitive: function == ExistsContainsCI, regexEscape: true})
	case ExistsRegexp, ExistsRegexpCI:
		return ext.callRegexp(con, target, args, regexpOptions{caseInsensitive: function == ExistsRegexpCI, start: true, end: true})
	default:
		return fmt.Errorf("unsupported filter: %v", function)
	}
	return fmt.Errorf("unsupported types: %v.(%v)", tgtType, argType)
}

type regexpOptions struct {
	caseInsensitive bool
	start           bool
	end             bool
	regexEscape     bool
}

func writeTarget(con *cel2sql.Converter, function string, target *expr.Expr) error {
	if _, has := ciFuncs[function]; has {
		con.WriteString("COLLATE(")
	}
	if err := con.Visit(target); err != nil {
		return err
	}
	if _, has := ciFuncs[function]; has {
		con.WriteString(", \"und:ci\")")
	}
	return nil
}

func writeSimpleCall(sqlFunc string, con *cel2sql.Converter, function string, target, arg *expr.Expr) error {
	con.WriteString(sqlFunc + "(")
	if err := writeTarget(con, function, target); err != nil {
		return err
	}
	con.WriteString(", ")
	if err := con.Visit(arg); err != nil {
		return err
	}
	con.WriteString(")")
	return nil
}

// REGEXP_CONTAINS("\x00" || ARRAY_TO_STRING(target, "\x00") || "\x00", r"\x00(arg1|arg2|arg3)\x00")
func (ext *Extension) callRegexp(con *cel2sql.Converter, target *expr.Expr, args []*expr.Expr, opts regexpOptions) error {
	// Special case for one value and non-slice fields.
	tgtType := con.GetType(target)
	argType := con.GetType(args[0])
	if tgtType.GetPrimitive() == expr.Type_STRING && argType.GetPrimitive() == expr.Type_STRING {
		arg, err := cel2sql.GetConstValue(args[0])
		if err != nil {
			return fmt.Errorf("failed to get const value of regexp: %w", err)
		}
		re, ok := arg.(string)
		if !ok {
			return fmt.Errorf("regexp's value is %T, want a string", arg)
		}
		re = "^(" + re + ")$"
		if opts.caseInsensitive {
			re = "(?i)" + re
		}
		con.WriteString("REGEXP_CONTAINS(")
		if err := con.Visit(target); err != nil {
			return err
		}
		con.WriteString(", ")
		con.WriteValue(re)
		con.WriteString(")")
		return nil
	}

	con.WriteString("REGEXP_CONTAINS(\"\\x00\" || ")
	switch {
	case tgtType.GetPrimitive() == expr.Type_STRING:
		if err := con.Visit(target); err != nil {
			return err
		}
	case cel2sql.IsListType(tgtType):
		con.WriteString("ARRAY_TO_STRING(")
		if err := con.Visit(target); err != nil {
			return err
		}
		con.WriteString(", \"\\x00\")")
	}
	con.WriteString(" || \"\\x00\", ")
	regexp, err := buildRegex(args[0], opts)
	if err != nil {
		return err
	}
	//replace con.WriteValue with this if params don't work for some reason
	//con.WriteString(fmt.Sprintf("%q", regexp))
	con.WriteValue(regexp)
	con.WriteString(")")
	return nil
}

func buildRegex(expression *expr.Expr, opts regexpOptions) (string, error) {
	builder := strings.Builder{}
	if opts.caseInsensitive {
		builder.WriteString("(?i)")
	}
	if opts.start {
		builder.WriteString("\x00")
	}
	builder.WriteString("(")

	arg, err := cel2sql.GetConstValue(expression)
	if err != nil {
		return "", err
	}
	switch value := arg.(type) {
	case string:
		builder.WriteString(joinRegexps([]string{value}, opts.regexEscape))
	case []interface{}:
		patterns := make([]string, 0, len(value))
		for _, val := range value {
			if pattern, ok := val.(string); ok {
				patterns = append(patterns, pattern)
			} else {
				return "", fmt.Errorf("wrong const value: %v", pattern)
			}
		}
		builder.WriteString(joinRegexps(patterns, opts.regexEscape))
	default:
		return "", fmt.Errorf("wrong const value: %v", value)
	}
	builder.WriteString(")")
	if opts.end {
		builder.WriteString("\x00")
	}
	return builder.String(), nil
}

func joinRegexps(patterns []string, escapeItems bool) string {
	parts := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if escapeItems {
			p = regexp.QuoteMeta(p)
		} else {
			p = fmt.Sprintf("(%s)", p)
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, "|")
}
