package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/documents/pkg/errors"
)

func TestNewRuleNotSatisfied(t *testing.T) {
	var RuleNotSatisfiedCode errors.Code
	RuleNotSatisfiedCode = "RULE_NOT_SATISFIED"

	RuleNotSatisfied := errors.NewRuleNotSatisfied(RuleNotSatisfiedCode)
	assert.Equal(t, "<RULE_NOT_SATISFIED> internal rule was not satisfied", RuleNotSatisfied.Error())
	assert.Equal(t, RuleNotSatisfiedCode, RuleNotSatisfied.GetCode())
	assert.Equal(t, "", RuleNotSatisfied.GetMessage())

	RuleNotSatisfied = errors.NewRuleNotSatisfied(RuleNotSatisfiedCode).WithArg("id", "123abc123abc")
	assert.Equal(t, "<RULE_NOT_SATISFIED> internal rule was not satisfied (id: 123abc123abc)", RuleNotSatisfied.Error())

	RuleNotSatisfied = errors.NewRuleNotSatisfied(RuleNotSatisfiedCode).WithMessage("custom message rule not satisfied").WithArg("id", "123abc123abc")
	assert.Equal(t, "<RULE_NOT_SATISFIED> custom message rule not satisfied (id: 123abc123abc)", RuleNotSatisfied.Error())
	assert.Equal(t, "custom message rule not satisfied", RuleNotSatisfied.GetMessage())
}
