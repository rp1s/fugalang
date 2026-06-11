//! DO NOT EDIT
package action

import (
	"fugu/pkg/reporter"
	. "fugu/pkg/token"
)

var Actions = []ActionStruct{
	Err(reporter.NoError), // 0
	Err(reporter.NoError), // 1
	Err(reporter.NoError), // 2
	Err(reporter.NoError), // 3
	Err(reporter.NoError), // 4
	Err(reporter.NoError), // 5
	Sh(2), // 6
	Err(reporter.NoError), // 7
	Err(reporter.NoError), // 8
	Err(reporter.NoError), // 9
	Err(reporter.NoError), // 10
	Err(reporter.NoError), // 11
	Err(reporter.NoError), // 12
	Err(reporter.NoError), // 13
	Err(reporter.NoError), // 14
	Err(reporter.NoError), // 15
	Err(reporter.NoError), // 16
	Err(reporter.NoError), // 17
	Err(reporter.NoError), // 18
	Err(reporter.NoError), // 19
	Err(reporter.NoError), // 20
	Err(reporter.NoError), // 21
	Err(reporter.NoError), // 22
	Err(reporter.NoError), // 23
	Err(reporter.NoError), // 24
	Err(reporter.NoError), // 25
	Err(reporter.NoError), // 26
	Err(reporter.NoError), // 27
	Err(reporter.NoError), // 28
	Err(reporter.NoError), // 29
	Err(reporter.NoError), // 30
	Err(reporter.NoError), // 31
	Err(reporter.NoError), // 32
	Err(reporter.NoError), // 33
	Err(reporter.NoError), // 34
	Err(reporter.NoError), // 35
	Err(reporter.NoError), // 36
	Err(reporter.NoError), // 37
	Err(reporter.NoError), // 38
	Err(reporter.NoError), // 39
	Err(reporter.NoError), // 40
	Err(reporter.NoError), // 41
	Err(reporter.NoError), // 42
	Err(reporter.NoError), // 43
	Err(reporter.NoError), // 44
	Err(reporter.NoError), // 45
	Err(reporter.NoError), // 46
	Err(reporter.NoError), // 47
	Err(reporter.NoError), // 48
	Err(reporter.NoError), // 49
	Err(reporter.NoError), // 50
	Err(reporter.NoError), // 51
	Err(reporter.NoError), // 52
	Err(reporter.NoError), // 53
	Err(reporter.NoError), // 54
	Err(reporter.NoError), // 55
	Err(reporter.NoError), // 56
	Sh(3), // 57
}

var Check = []int{
	-1, // 0
	-1, // 1
	-1, // 2
	-1, // 3
	-1, // 4
	-1, // 5
	0, // 6
	-1, // 7
	-1, // 8
	-1, // 9
	-1, // 10
	-1, // 11
	-1, // 12
	-1, // 13
	-1, // 14
	-1, // 15
	-1, // 16
	-1, // 17
	-1, // 18
	-1, // 19
	-1, // 20
	-1, // 21
	-1, // 22
	-1, // 23
	-1, // 24
	-1, // 25
	-1, // 26
	-1, // 27
	-1, // 28
	-1, // 29
	-1, // 30
	-1, // 31
	-1, // 32
	-1, // 33
	-1, // 34
	-1, // 35
	-1, // 36
	-1, // 37
	-1, // 38
	-1, // 39
	-1, // 40
	-1, // 41
	-1, // 42
	-1, // 43
	-1, // 44
	-1, // 45
	-1, // 46
	-1, // 47
	-1, // 48
	-1, // 49
	-1, // 50
	-1, // 51
	-1, // 52
	-1, // 53
	-1, // 54
	-1, // 55
	-1, // 56
	0, // 57
}

var Base = []int{
	0, // state 0
}

func Action(state int, tk TokenKind) ActionStruct {
	if state < 0 || state >= len(Base) {
		return Err(reporter.NoError)
	}
	b := Base[state]
	if b < 0 {
		return Err(reporter.NoError)
	}
	idx := b + int(tk)
	if idx >= 0 && idx < len(Actions) && Check[idx] == state {
		return Actions[idx]
	}
	return Err(reporter.NoError)
}
