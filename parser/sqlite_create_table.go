package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenT int

const (
	// internals
	tokEOF tokenT = iota
	tokERROR
	tokIDENTIFIER
	tokCOMMENT

	// keywords
	tokCREATE
	tokTEMP
	tokTABLE
	tokIF
	tokNOT
	tokEXISTS
	tokAS
	tokWITHOUT
	tokROWID

	// separators
	tokDOT
	tokSEMICOLON
	tokCOMMA
	tokOPENparenthesis
	tokCLOSEDparenthesis

	// Constraints
	tokCONSTRAINT
	tokPRIMARY
	tokKEY
	tokUNIQUE
	tokCHECK
	tokFOREIGN
	tokON
	tokCONFLICT
	tokROLLBACK
	tokABORT
	tokFAIL
	tokIGNORE
	tokREPLACE
	tokCOLLATE
	tokASC
	tokDESC
	tokAUTOINCREMENT

	// foreign key clause
	tokREFERENCES
	tokDELETE
	tokUPDATE
	tokSET
	tokNULL
	tokDEFAULT
	tokCASCADE
	tokRESTRICT
	tokNO
	tokACTION
	tokMATCH
	tokDEFERRABLE
	tokINITIALLY
	tokDEFERRED
	tokIMMEDIATE
)

type ForeignKey struct {
	Table      string
	NumColumns int
	ColumnName []string
	OnDelete   FkAction
	OnUpdate   FkAction
	Match      string
	Deferrable FkDefType
}

type Column struct {
	Name                  string
	Type                  string
	Length                string
	ConstraintName        string
	IsPrimaryKey          bool
	IsAutoincrement       bool
	IsNotnull             bool
	IsUnique              bool
	PkOrder               OrderClause
	PkConflictClause      ConflictClause
	NotNullConflictClause ConflictClause
	UniqueConflictClause  ConflictClause
	CheckExpr             string
	DefaultExpr           string
	CollateName           string
	ForeignKeyClause      *ForeignKey
}

type TableConstraint struct {
	Type             ConstraintType
	Name             string
	NumIndexed       int
	IndexedColumns   []IdxColumn
	ConflictClause   ConflictClause
	CheckExpr        string
	ForeignKeyNum    int
	ForeignKeyName   []string
	ForeignKeyClause *ForeignKey
}

type Table struct {
	Name           string
	Schema         string
	IsTemporary    bool
	IsIfNotExists  bool
	IsWithoutRowid bool
	NumColumns     int
	Columns        []Column
	NumConstraint  int
	Constraints    []TableConstraint
}

type IdxColumn struct {
	Name        string
	CollateName string
	Order       OrderClause
}

type State struct {
	buffer     []rune
	size       int
	offset     int
	identifier string
	table      *Table
}

func isEOF(state *State) bool {
	return state.offset == state.size
}

func peek(state *State) rune {
	return state.buffer[state.offset]
}

func peek2(state *State) rune {
	return state.buffer[state.offset+1]
}

func next(state *State) rune {
	c := state.buffer[state.offset]
	state.offset++
	return c
}

func skip1(state *State) {
	state.offset++
}

func strNoCaseNcmp(s1, s2 string, n int) int {
	index := 0
	sl1 := strings.ToLower(s1)
	sl2 := strings.ToLower(s2)
	for n > 0 && index < len(sl2) && sl1[index] == sl2[index] {
		if sl1[index] == 0x00 {
			return 0
		}
		n--
		index++
	}
	if n == 0 {
		return 0
	}
	return 1
}

func symbolIsSpace(r rune) bool {
	switch r {
	case '\t', '\v', '\f', ' ':
		return true
	}
	return false
}

func symbolIsNewline(r rune) bool {
	switch r {
	case '\n', '\r':
		return true
	}
	return false
}

func symbolIsToSkip(r rune) bool {
	return symbolIsSpace(r) || symbolIsNewline(r)
}

func symbolIsComment(r rune, state *State) bool {
	if r == '-' && peek2(state) == '-' {
		return true
	}
	if r == '/' && peek2(state) == '*' {
		return true
	}
	return false
}

func symbolIsAlpha(r rune) bool {
	if r == '_' {
		return true
	}
	return unicode.IsLetter(r)
}

func symbolIsIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func symbolIsEscape(r rune) bool {
	return r == '`' || r == '\'' || r == '"' || r == '['
}

func symbolIsPunctuation(r rune) bool {
	return r == '.' || r == ',' || r == '(' || r == ')' || r == ';'
}

func tokenIsColumnConstraint(t tokenT) bool {
	return t == tokCONSTRAINT || t == tokPRIMARY || t == tokNOT || t == tokUNIQUE ||
		t == tokCHECK || t == tokDEFAULT || t == tokCOLLATE || t == tokREFERENCES
}

func tokenIsTableConstraint(t tokenT) bool {
	return t == tokCONSTRAINT || t == tokPRIMARY || t == tokUNIQUE ||
		t == tokCHECK || t == tokFOREIGN
}

func lexerKeyword(ptr string, length int) tokenT {
	switch length {
	case 2:
		if strNoCaseNcmp(ptr, "if", length) == 0 {
			return tokIF
		}
		if strNoCaseNcmp(ptr, "as", length) == 0 {
			return tokAS
		}
		if strNoCaseNcmp(ptr, "on", length) == 0 {
			return tokON
		}
		if strNoCaseNcmp(ptr, "no", length) == 0 {
			return tokNO
		}
	case 3:
		if strNoCaseNcmp(ptr, "not", length) == 0 {
			return tokNOT
		}
		if strNoCaseNcmp(ptr, "key", length) == 0 {
			return tokKEY
		}
		if strNoCaseNcmp(ptr, "asc", length) == 0 {
			return tokASC
		}
		if strNoCaseNcmp(ptr, "set", length) == 0 {
			return tokSET
		}
	case 4:
		if strNoCaseNcmp(ptr, "temp", length) == 0 {
			return tokTEMP
		}
		if strNoCaseNcmp(ptr, "desc", length) == 0 {
			return tokDESC
		}
		if strNoCaseNcmp(ptr, "null", length) == 0 {
			return tokNULL
		}
		if strNoCaseNcmp(ptr, "fail", length) == 0 {
			return tokFAIL
		}
	case 5:
		if strNoCaseNcmp(ptr, "table", length) == 0 {
			return tokTABLE
		}
		if strNoCaseNcmp(ptr, "rowid", length) == 0 {
			return tokROWID
		}
		if strNoCaseNcmp(ptr, "check", length) == 0 {
			return tokCHECK
		}
		if strNoCaseNcmp(ptr, "abort", length) == 0 {
			return tokABORT
		}
		if strNoCaseNcmp(ptr, "match", length) == 0 {
			return tokMATCH
		}
	case 6:
		if strNoCaseNcmp(ptr, "create", length) == 0 {
			return tokCREATE
		}
		if strNoCaseNcmp(ptr, "exists", length) == 0 {
			return tokEXISTS
		}
		if strNoCaseNcmp(ptr, "unique", length) == 0 {
			return tokUNIQUE
		}
		if strNoCaseNcmp(ptr, "ignore", length) == 0 {
			return tokIGNORE
		}
		if strNoCaseNcmp(ptr, "delete", length) == 0 {
			return tokDELETE
		}
		if strNoCaseNcmp(ptr, "update", length) == 0 {
			return tokUPDATE
		}
		if strNoCaseNcmp(ptr, "action", length) == 0 {
			return tokACTION
		}
	case 7:
		if strNoCaseNcmp(ptr, "without", length) == 0 {
			return tokWITHOUT
		}
		if strNoCaseNcmp(ptr, "primary", length) == 0 {
			return tokPRIMARY
		}
		if strNoCaseNcmp(ptr, "default", length) == 0 {
			return tokDEFAULT
		}
		if strNoCaseNcmp(ptr, "collate", length) == 0 {
			return tokCOLLATE
		}
		if strNoCaseNcmp(ptr, "replace", length) == 0 {
			return tokREPLACE
		}
		if strNoCaseNcmp(ptr, "cascade", length) == 0 {
			return tokCASCADE
		}
		if strNoCaseNcmp(ptr, "foreign", length) == 0 {
			return tokFOREIGN
		}
	case 8:
		if strNoCaseNcmp(ptr, "conflict", length) == 0 {
			return tokCONFLICT
		}
		if strNoCaseNcmp(ptr, "rollback", length) == 0 {
			return tokROLLBACK
		}
		if strNoCaseNcmp(ptr, "restrict", length) == 0 {
			return tokRESTRICT
		}
		if strNoCaseNcmp(ptr, "deferred", length) == 0 {
			return tokDEFERRED
		}
	case 9:
		if strNoCaseNcmp(ptr, "temporary", length) == 0 {
			return tokTEMP
		}
		if strNoCaseNcmp(ptr, "initially", length) == 0 {
			return tokINITIALLY
		}
		if strNoCaseNcmp(ptr, "immediate", length) == 0 {
			return tokIMMEDIATE
		}
	case 10:
		if strNoCaseNcmp(ptr, "constraint", length) == 0 {
			return tokCONSTRAINT
		}
		if strNoCaseNcmp(ptr, "references", length) == 0 {
			return tokREFERENCES
		}
		if strNoCaseNcmp(ptr, "deferrable", length) == 0 {
			return tokDEFERRABLE
		}
	case 13:
		if strNoCaseNcmp(ptr, "autoincrement", length) == 0 {
			return tokAUTOINCREMENT
		}
	}

	return tokIDENTIFIER
}

func lexerComment(state *State) tokenT {
	isCComment := next(state) == '/' && next(state) == '*'
	for {
		c1 := next(state)

		if c1 == 0x00 {
			if isCComment {
				return tokERROR
			}
			return tokCOMMENT
		}

		if !isCComment && symbolIsNewline(c1) {
			break
		}

		c2 := next(state)
		if isCComment && c1 == '*' && c2 == '/' {
			break
		}
	}
	return tokCOMMENT
}

func lexerPunctuation(state *State) tokenT {
	c := next(state)
	switch c {
	case ',':
		return tokCOMMA
	case '.':
		return tokDOT
	case '(':
		return tokOPENparenthesis
	case ')':
		return tokCLOSEDparenthesis
	case ';':
		return tokSEMICOLON
	}
	return tokERROR
}

func lexerAlpha(state *State) tokenT {
	offset := state.offset

	for symbolIsIdentifier(peek(state)) {
		skip1(state)
	}

	length := state.offset - offset
	ptr := string(state.buffer[offset:state.offset])

	t := lexerKeyword(ptr, length)

	if t != tokIDENTIFIER {
		return t
	}

	state.identifier = ptr

	return tokIDENTIFIER
}

func lexerEscape(state *State) tokenT {
	c := next(state)
	escaped := c
	if escaped == '[' {
		escaped = ']'
	}

	offset := state.offset
	c = next(state)
	for c != 0 && c != escaped {
		c = next(state)
	}

	ptr := string(state.buffer[offset : state.offset-1])

	if c != escaped {
		return tokERROR
	}

	state.identifier = ptr

	return tokIDENTIFIER
}

func lexerNext(state *State) tokenT {
	for {
		if isEOF(state) {
			return tokEOF
		}

		c := peek(state)
		if c == 0x00 {
			return tokEOF
		}

		if symbolIsToSkip(c) {
			skip1(state)
			continue
		}

		if symbolIsComment(c, state) {
			if lexerComment(state) != tokCOMMENT {
				return tokERROR
			}
			continue
		}

		if symbolIsPunctuation(c) {
			return lexerPunctuation(state)
		}

		if symbolIsAlpha(c) {
			return lexerAlpha(state)
		}

		if symbolIsEscape(c) {
			return lexerEscape(state)
		}

		return tokERROR
	}
}

func lexerPeek(state *State) tokenT {
	saved := state.offset
	token := lexerNext(state)
	state.offset = saved
	return token
}

func parseOptionalOrder(state *State, clause *OrderClause) ErrorCode {
	token := lexerPeek(state)
	*clause = ORDER_NONE

	if token == tokASC || token == tokDESC {
		lexerNext(state)
		if token == tokASC {
			*clause = ORDER_ASC
		} else {
			*clause = ORDER_DESC
		}
	}

	return ERROR_NONE
}

func parseOptionalConflictClause(state *State, conflict *ConflictClause) ErrorCode {
	token := lexerPeek(state)
	*conflict = CONFLICT_NONE

	if token == tokON {
		lexerNext(state)
		token = lexerNext(state)

		if token != tokCONFLICT {
			return ERROR_SYNTAX
		}

		token = lexerNext(state)
		if token == tokROLLBACK {
			*conflict = CONFLICT_ROOLBACK
		} else if token == tokABORT {
			*conflict = CONFLICT_ABORT
		} else if token == tokFAIL {
			*conflict = CONFLICT_FAIL
		} else if token == tokIGNORE {
			*conflict = CONFLICT_IGNORE
		} else if token == tokREPLACE {
			*conflict = CONFLICT_REPLACE
		} else {
			return ERROR_SYNTAX
		}
	}
	return ERROR_NONE
}

func parseForeignKeyClause(state *State) *ForeignKey {
	var fk ForeignKey

	token := lexerNext(state)
	if token != tokIDENTIFIER {
		fmt.Println("parseForeignKeyClause error")
		return nil
	}

	fk.Table = state.identifier

	if lexerPeek(state) == tokOPENparenthesis {
		lexerNext(state)

		token = lexerNext(state)
		if token != tokIDENTIFIER {
			fmt.Println("parseForeignKeyClause error")
			return nil
		}
		fk.ColumnName = []string{state.identifier}

		token = lexerPeek(state)
		if token == tokCOMMA {
			lexerNext(state)
		}
		for token == tokCOMMA {
			token = lexerNext(state)
			if token != tokIDENTIFIER {
				fmt.Println("parseForeignKeyClause error")
				return nil
			}

			fk.NumColumns++
			fk.ColumnName = append(fk.ColumnName, state.identifier)

			token = lexerPeek(state)
			if token == tokCOMMA {
				lexerNext(state)
			}
		}
		if lexerNext(state) != tokCLOSEDparenthesis {
			fmt.Println("parseForeignKeyClause error")
			return nil
		}
	}

	for {
		token = lexerPeek(state)
		if token == tokON || token == tokMATCH || token == tokNOT || token == tokDEFERRABLE {
			lexerNext(state)

			if token == tokMATCH {
				token = lexerNext(state)
				if token != tokIDENTIFIER {
					fmt.Println("parseForeignKeyClause error")
					return nil
				}
				fk.Match = state.identifier
				continue
			}

			if token == tokON {
				token = lexerNext(state)
				if token != tokDELETE && token != tokUPDATE {
					fmt.Println("parseForeignKeyClause error")
					return nil
				}
				isUpdate := token == tokUPDATE

				token = lexerNext(state)
				if token == tokCASCADE {
					if isUpdate {
						fk.OnUpdate = FKACTION_CASCADE
					} else {
						fk.OnDelete = FKACTION_CASCADE
					}
				} else if token == tokRESTRICT {
					if isUpdate {
						fk.OnUpdate = FKACTION_RESTRICT
					} else {
						fk.OnDelete = FKACTION_RESTRICT
					}
				} else if token == tokSET {
					token = lexerNext(state)
					if token != tokNULL && token != tokDEFAULT {
						fmt.Println("parseForeignKeyClause error")
						return nil
					}
					if token == tokNULL {
						if isUpdate {
							fk.OnUpdate = FKACTION_SETNULL
						} else {
							fk.OnDelete = FKACTION_SETNULL
						}
					} else {
						if isUpdate {
							fk.OnUpdate = FKACTION_SETDEFAULT
						} else {
							fk.OnDelete = FKACTION_SETDEFAULT
						}
					}
				} else if token == tokNO {
					if lexerNext(state) != tokACTION {
						fmt.Println("parseForeignKeyClause error")
						return nil
					}
					if isUpdate {
						fk.OnUpdate = FKACTION_NOACTION
					} else {
						fk.OnDelete = FKACTION_NOACTION
					}
				}
				continue
			}
			isNot := false
			if token == tokNOT {
				token = lexerNext(state)
				isNot = true
			}

			if token == tokDEFERRABLE {
				if isNot {
					fk.Deferrable = DEFTYPE_NOTDEFERRABLE
				} else {
					fk.Deferrable = DEFTYPE_DEFERRABLE
				}

				if lexerPeek(state) == tokINITIALLY {
					lexerNext(state)
					token = lexerNext(state)
					if token == tokDEFERRED {
						if isNot {
							fk.Deferrable = DEFTYPE_NOTDEFERRABLE_INITIALLY_DEFERRED
						} else {
							fk.Deferrable = DEFTYPE_DEFERRABLE_INITIALLY_DEFERRED
						}
					} else if token == tokIMMEDIATE {
						if isNot {
							fk.Deferrable = DEFTYPE_NOTDEFERRABLE_INITIALLY_IMMEDIATE
						} else {
							fk.Deferrable = DEFTYPE_DEFERRABLE_INITIALLY_IMMEDIATE
						}
					} else {
						fmt.Println("parseForeignKeyClause error")
						return nil
					}
				}
				continue
			}
			fmt.Println("parseForeignKeyClause error")
			return nil
		}
		return &fk
	}
}

func parseTableConstraint(state *State) *TableConstraint {
	token := lexerPeek(state)
	var constraint TableConstraint

	if token == tokCONSTRAINT {
		lexerNext(state)
		token = lexerNext(state)
		if token != tokIDENTIFIER {
			fmt.Println("parseTableConstraint error")
			return nil
		}
		constraint.Name = state.identifier

		token = lexerPeek(state)

		if token != tokCHECK && token != tokPRIMARY && token != tokUNIQUE && token != tokFOREIGN {
			fmt.Println("parseTableConstraint error")
			return nil
		}
	}

	if token == tokCHECK {
		token = lexerNext(state)
		constraint.Type = TABLECONSTRAINT_CHECK

		fmt.Println("parseTableConstraint error")
		return nil
	} else if token == tokPRIMARY || token == tokUNIQUE {
		token = lexerNext(state)
		if token == tokPRIMARY {
			if lexerNext(state) != tokKEY {
				fmt.Println("parseTableConstraint error")
				return nil
			}
			constraint.Type = TABLECONSTRAINT_PRIMARYKEY
		} else {
			constraint.Type = TABLECONSTRAINT_UNIQUE
		}

		if lexerNext(state) != tokOPENparenthesis {
			fmt.Println("parseTableConstraint error")
			return nil
		}

		//do
		var column IdxColumn

		token = lexerNext(state)
		if token != tokIDENTIFIER {
			fmt.Println("parseTableConstraint error")
			return nil
		}
		column.Name = state.identifier

		if lexerPeek(state) == tokCOLLATE {
			lexerNext(state)

			token := lexerNext(state)
			if token != tokIDENTIFIER {
				fmt.Println("parseTableConstraint error")
				return nil
			}
			column.CollateName = state.identifier
		}

		if parseOptionalOrder(state, &column.Order) != ERROR_NONE {
			fmt.Println("parseTableConstraint error")
			return nil
		}

		constraint.NumIndexed++
		constraint.IndexedColumns = []IdxColumn{column}

		token = lexerPeek(state)
		if token == tokCOMMA {
			lexerNext(state)
		}
		// while
		for token == tokCOMMA {
			var column IdxColumn

			token = lexerNext(state)
			if token != tokIDENTIFIER {
				fmt.Println("parseTableConstraint error")
				return nil
			}
			column.Name = state.identifier

			if lexerPeek(state) == tokCOLLATE {
				lexerNext(state)

				token := lexerNext(state)
				if token != tokIDENTIFIER {
					fmt.Println("parseTableConstraint error")
					return nil
				}
				column.CollateName = state.identifier
			}

			if parseOptionalOrder(state, &column.Order) != ERROR_NONE {
				fmt.Println("parseTableConstraint error")
				return nil
			}

			constraint.NumIndexed++
			constraint.IndexedColumns = append(constraint.IndexedColumns, column)

			token = lexerPeek(state)
			if token == tokCOMMA {
				lexerNext(state)
			}
		}
		if lexerNext(state) != tokCLOSEDparenthesis {
			fmt.Println("parseTableConstraint error")
			return nil
		}
		if parseOptionalConflictClause(state, &constraint.ConflictClause) != ERROR_NONE {
			fmt.Println("parseTableConstraint error")
			return nil
		}
	} else if token == tokFOREIGN {
		lexerNext(state)
		if lexerNext(state) != tokKEY {
			fmt.Println("parseTableConstraint error")
			return nil
		}
		if lexerNext(state) != tokOPENparenthesis {
			fmt.Println("parseTableConstraint error")
			return nil
		}

		constraint.Type = TABLECONSTRAINT_FOREIGNKEY
		//do
		token = lexerNext(state)
		if token != tokIDENTIFIER {
			fmt.Println("parseTableConstraint error")
			return nil
		}
		constraint.ForeignKeyNum++
		constraint.ForeignKeyName = []string{state.identifier}

		token = lexerPeek(state)
		if token == tokCOMMA {
			lexerNext(state)
		}
		//while
		for token == tokCOMMA {
			token = lexerNext(state)
			if token != tokIDENTIFIER {
				fmt.Println("parseTableConstraint error")
				return nil
			}
			constraint.ForeignKeyNum++
			constraint.ForeignKeyName = append(constraint.ForeignKeyName, state.identifier)

			token = lexerPeek(state)
			if token == tokCOMMA {
				lexerNext(state)
			}
		}

		if lexerNext(state) != tokCLOSEDparenthesis {
			fmt.Println("parseTableConstraint error")
			return nil
		}

		if lexerNext(state) != tokREFERENCES {
			fmt.Println("parseTableConstraint error")
			return nil
		}

		fk := parseForeignKeyClause(state)
		if fk == nil {
			fmt.Println("parseTableConstraint error")
			return nil
		}
		constraint.ForeignKeyClause = fk
	}

	return &constraint
}

func parseLiteral(state *State) ErrorCode {
	if lexerNext(state) == tokIDENTIFIER {
		return ERROR_NONE
	} else {
		return ERROR_SYNTAX
	}
}

func parseColumnType(state *State, column *Column) ErrorCode {
	offset := 0
	for lexerPeek(state) == tokIDENTIFIER {
		// consume identifier
		lexerNext(state)

		if offset == 0 {
			offset = state.offset - len(state.identifier)
		}
	}
	ptr := string(state.buffer[offset:state.offset])

	column.Type = ptr

	if lexerPeek(state) == tokOPENparenthesis {
		lexerNext(state)

		offset = state.offset
		var c rune
		c = next(state)
		for c != 0x00 && c != ')' {
			c = next(state)
		}

		if c != ')' {
			return ERROR_SYNTAX
		}

		ptr := string(state.buffer[offset:state.offset])

		column.Length = ptr
	}

	return ERROR_NONE
}

func parseColumnConstraints(state *State, column *Column) ErrorCode {
	for tokenIsColumnConstraint(lexerPeek(state)) {
		token := lexerNext(state)

		if token == tokCONSTRAINT {
			token = lexerNext(state)
			if token != tokIDENTIFIER {
				return ERROR_SYNTAX
			}
			column.ConstraintName = state.identifier
			token = lexerNext(state)
		}

		switch token {
		case tokPRIMARY:
			token = lexerNext(state)
			if token != tokKEY {
				return ERROR_SYNTAX
			}
			column.IsPrimaryKey = true
			if parseOptionalOrder(state, &column.PkOrder) != ERROR_NONE {
				return ERROR_SYNTAX
			}
			if parseOptionalConflictClause(state, &column.PkConflictClause) != ERROR_NONE {
				return ERROR_SYNTAX
			}
			if lexerPeek(state) == tokAUTOINCREMENT {
				lexerNext(state)
				column.IsAutoincrement = true
			}
		case tokNOT:
			token = lexerNext(state)
			if token != tokNULL {
				fmt.Println(string(state.buffer[state.offset:state.offset+4]), token)
				return ERROR_SYNTAX
			}
			column.IsNotnull = true
			if parseOptionalConflictClause(state, &column.NotNullConflictClause) != ERROR_NONE {
				return ERROR_SYNTAX
			}
		case tokUNIQUE:
			column.IsUnique = true
			if parseOptionalConflictClause(state, &column.UniqueConflictClause) != ERROR_NONE {
				return ERROR_SYNTAX
			}
		case tokCHECK:
			return ERROR_UNSUPPORTEDSQL
		case tokDEFAULT:
			if lexerPeek(state) == tokOPENparenthesis {
				return ERROR_UNSUPPORTEDSQL
			}
			if parseLiteral(state) != ERROR_NONE {
				return ERROR_SYNTAX
			}
			column.DefaultExpr = state.identifier
		case tokCOLLATE:
			token = lexerNext(state)
			if token != tokIDENTIFIER {
				return ERROR_SYNTAX
			}
			column.CollateName = state.identifier
		case tokREFERENCES:
			fk := parseForeignKeyClause(state)
			if fk == nil {
				return ERROR_SYNTAX
			}
			column.ForeignKeyClause = fk
		default:
			return ERROR_SYNTAX
		}
	}
	return ERROR_NONE
}

func parseColumn(state *State) *Column {
	var column Column

	token := lexerNext(state)

	if token != tokIDENTIFIER {
		fmt.Println("parseColumn error")
		return nil
	}

	column.Name = state.identifier

	if lexerPeek(state) == tokIDENTIFIER {
		if parseColumnType(state, &column) != ERROR_NONE {
			fmt.Println("parseColumn error")
			return nil
		}
	}

	if tokenIsColumnConstraint(lexerPeek(state)) {
		if parseColumnConstraints(state, &column) != ERROR_NONE {
			fmt.Println("parseColumn error")
			return nil
		}
	}

	return &column
}

func parse(state *State) ErrorCode {
	token := lexerNext(state)

	if token != tokCREATE {
		return ERROR_UNSUPPORTEDSQL
	}

	table := state.table
	token = lexerNext(state)
	if token == tokTEMP {
		table.IsTemporary = true

		token = lexerNext(state)
	}

	if token != tokTABLE {
		return ERROR_UNSUPPORTEDSQL
	}
	if lexerPeek(state) == tokIF {
		lexerNext(state)

		if lexerNext(state) != tokNOT {
			return ERROR_SYNTAX
		}

		if lexerNext(state) != tokEXISTS {
			return ERROR_SYNTAX
		}

		table.IsIfNotExists = true
	}

	if lexerNext(state) != tokIDENTIFIER {
		return ERROR_SYNTAX
	}

	identifier := state.identifier
	if identifier == "" {
		return ERROR_SYNTAX
	}

	table.Name = state.identifier

	if lexerPeek(state) == tokAS {
		return ERROR_UNSUPPORTEDSQL
	}

	token = lexerNext(state)
	if token != tokOPENparenthesis {
		return ERROR_SYNTAX
	}

	// parse column def
	for {
		token = lexerPeek(state)

		if token != tokIDENTIFIER {
			return ERROR_SYNTAX
		}

		column := parseColumn(state)
		if column == nil {
			return ERROR_SYNTAX
		}

		table.NumColumns++
		if table.Columns == nil {
			table.Columns = []Column{}
		}
		table.Columns = append(table.Columns, *column)

		token = lexerPeek(state)
		if token == tokCOMMA {
			lexerNext(state)
			token = lexerPeek(state)

			if tokenIsTableConstraint(token) {
				break
			} else {
				continue
			}
		}
		if token == tokCLOSEDparenthesis {
			break
		}

		return ERROR_SYNTAX
	}

	for tokenIsTableConstraint(token) {
		constraint := parseTableConstraint(state)
		if constraint == nil {

			return ERROR_SYNTAX
		}

		table.NumConstraint++
		if table.Constraints == nil {
			table.Constraints = []TableConstraint{}
		}
		table.Constraints = append(table.Constraints, *constraint)

		if lexerPeek(state) == tokCOMMA {
			lexerNext(state)
			token = lexerPeek(state)
			continue
		}

		if lexerPeek(state) == tokCLOSEDparenthesis {
			break
		}

		return ERROR_SYNTAX
	}

	token = lexerNext(state)

	if token != tokCLOSEDparenthesis {
		return ERROR_SYNTAX
	}

	if lexerPeek(state) == tokWITHOUT {
		lexerNext(state)

		if lexerNext(state) != tokROWID {
			return ERROR_SYNTAX
		}

		table.IsWithoutRowid = true
	}

	token = lexerPeek(state)
	if token == tokSEMICOLON {
		lexerNext(state)
	}
	return ERROR_NONE
}

func ParseTable(sql string, length int) (*Table, ErrorCode) {
	if sql == "" {
		return nil, ERROR_NONE
	}
	if length == 0 {
		length = len(sql)
	}
	if length == 0 {
		return nil, ERROR_NONE
	}

	var table Table

	state := State{
		buffer: []rune(sql),
		size:   length,
		table:  &table,
	}

	err := parse(&state)
	return &table, err
}
