// Code generated by "stringer -type=ArithToken"; DO NOT EDIT

package arith

import "fmt"

const _ArithToken_name = "ArithErrorArithAssignmentArithNotArithAndArithOrArithNumberArithVariableArithLessEqualArithGreaterEqualArithLessThanArithGreaterThanArithEqualArithNotEqualArithBinaryAndArithBinaryOrArithBinaryXorArithLeftShiftArithRightShiftArithRemainderArithMultiplyArithDivideArithSubtractArithAddArithAssignBinaryAndArithAssignBinaryOrArithAssignBinaryXorArithAssignLeftShiftArithAssignRightShiftArithAssignRemainderArithAssignMultiplyArithAssignDivideArithAssignSubtractArithAssignAddArithLeftParenArithRightParenArithBinaryNotArithQuestionMarkArithColonArithEOF"

var _ArithToken_index = [...]uint16{0, 10, 25, 33, 41, 48, 59, 72, 86, 103, 116, 132, 142, 155, 169, 182, 196, 210, 225, 239, 252, 263, 276, 284, 304, 323, 343, 363, 384, 404, 423, 440, 459, 473, 487, 502, 516, 533, 543, 551}

func (i ArithToken) String() string {
	i -= 1
	if i < 0 || i >= ArithToken(len(_ArithToken_index)-1) {
		return fmt.Sprintf("ArithToken(%d)", i+1)
	}
	return _ArithToken_name[_ArithToken_index[i]:_ArithToken_index[i+1]]
}
