package rules

import (
	"fmt"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
	"slices"
	"strings"

	"github.com/ci-space/protostyle/internal/utils"
)

type FieldWithOptionRule struct {
	ruleName  string
	optName   string
	constants []string
}

type FieldWithOptionVisitor struct {
	*visitor.BaseAddVisitor

	optName   string
	constants []string
}

func NewFieldWithOptionRule(
	ruleName string,
	optName string,
	constants []string,
) *FieldWithOptionRule {
	return &FieldWithOptionRule{
		ruleName:  ruleName,
		optName:   optName,
		constants: constants,
	}
}

func (r FieldWithOptionRule) ID() string {
	return r.ruleName
}

func (r FieldWithOptionRule) Purpose() string {
	return fmt.Sprintf("Field must have option %q", r.optName)
}

func (r FieldWithOptionRule) IsOfficial() bool {
	return false
}

func (r FieldWithOptionRule) Severity() rule.Severity {
	return rule.SeverityError
}

func (r FieldWithOptionRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &FieldWithOptionVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
		optName:        r.optName,
		constants:      r.constants,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

func (v *FieldWithOptionVisitor) VisitField(field *parser.Field) (next bool) {
	opt, ok := utils.GetOptionFromField(field, v.optName)
	if ok {
		if len(v.constants) == 0 {
			return true
		}

		if !slices.Contains(v.constants, opt.Constant) {
			v.AddFailuref(
				field.Meta.Pos,
				"Field %q have invalid value for option %q. Expected values: [%s]",
				field.FieldName,
				v.optName,
				strings.Join(v.constants, ", "),
			)
		}
	} else {
		v.AddFailuref(field.Meta.Pos, "Field %q must have option %q", field.FieldName, v.optName)
	}

	return true
}
