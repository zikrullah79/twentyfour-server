package services

import (
	"errors"
	"fmt"

	"github.com/Knetic/govaluate"
)

func EvaluateFormula(formula string) error {
	res, err := govaluate.NewEvaluableExpression(formula)
	if err != nil {
		return err
	}
	if count := countModifier(res.Tokens()); count != 3 {
		return errors.New("wrong Formula : modifier != 3")
	}
	re, err := res.Evaluate(nil)
	if err != nil {
		return err
	}

	fmt.Printf("(%v, %T)\n", re, re)
	if re != 24.0 {
		return errors.New("wrong answer")
	}
	return err

}

func countModifier(formula []govaluate.ExpressionToken) int {
	total := 0
	for _, v := range formula {
		if v.Kind == govaluate.MODIFIER {
			total += 1
		}
	}
	return total
}
