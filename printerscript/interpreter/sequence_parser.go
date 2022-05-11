// Code generated from SequenceParser.g4 by ANTLR 4.10.1. DO NOT EDIT.

package interpreter // SequenceParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type SequenceParser struct {
	*antlr.BaseParser
}

var sequenceparserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sequenceparserParserInit() {
	staticData := &sequenceparserParserStaticData
	staticData.literalNames = []string{
		"", "'if'", "'else'", "'while'", "'true'", "'false'", "'('", "')'",
		"'{{'", "'{'", "'}'", "'&&'", "'||'", "'=='", "'!='", "'<='", "'>='",
		"'='", "'!'", "'-'", "'+'", "'*'", "'/'", "'%'", "'<'", "'>'", "','",
		"", "", "", "'}}'", "", "", "", "", "'//'", "'/*'", "'\\n'", "", "'*/'",
	}
	staticData.symbolicNames = []string{
		"", "IF", "ELSE", "WHILE", "TRUE", "FALSE", "L_PAREN", "R_PAREN", "LL_BRACE",
		"L_BRACE", "R_BRACE", "AND", "OR", "EQ", "N_EQ", "LT_EQ", "GT_EQ", "ASSIGN",
		"NOT", "MINUS", "PLUS", "TIMES", "DIV", "MOD", "LT", "GT", "COMMA",
		"GCODE_ESCAPE", "OPEN_GCODE_SQ", "OPEN_GCODE_DQ", "EXIT_EXPR", "IDENTIFIER",
		"INT", "FLOAT", "WS", "OPEN_LINE_COMMENT", "OPEN_BLOCK_COMMENT", "CLOSE_LINE_COMMENT",
		"LINE_COMMENT_TEXT", "CLOSE_BLOCK_COMMENT", "ESCAPED_CLOSE_BLOCK_COMMENT",
		"BLOCK_COMMENT_TEXT", "CLOSE_GCODE_SQ", "ENTER_EXPR_SQ", "ESCAPE_SEQUENCE_SQ",
		"TEXT_SQ", "CLOSE_GCODE_DQ", "ENTER_EXPR_DQ", "ESCAPE_SEQUENCE_DQ",
		"TEXT_DQ",
	}
	staticData.ruleNames = []string{
		"sequence", "statements", "ifBlock", "optionalElseBlock", "whileBlock",
		"statement", "assignment", "gCode", "gCodePart", "expression",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 49, 143, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
		1, 34, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 53, 8, 3, 1, 4, 1, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 3, 5, 65, 8, 5, 1, 6, 1, 6,
		1, 6, 1, 6, 1, 7, 1, 7, 5, 7, 73, 8, 7, 10, 7, 12, 7, 76, 9, 7, 1, 7, 1,
		7, 1, 7, 5, 7, 81, 8, 7, 10, 7, 12, 7, 84, 9, 7, 1, 7, 3, 7, 87, 8, 7,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 95, 8, 8, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 5, 9, 103, 8, 9, 10, 9, 12, 9, 106, 9, 9, 3, 9, 108, 8,
		9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 3, 9, 124, 8, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 5, 9, 138, 8, 9, 10, 9, 12, 9, 141, 9,
		9, 1, 9, 0, 1, 18, 10, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 0, 7, 2, 0, 45,
		45, 49, 49, 2, 0, 44, 44, 48, 48, 2, 0, 43, 43, 47, 47, 1, 0, 21, 23, 1,
		0, 19, 20, 2, 0, 13, 16, 24, 25, 1, 0, 11, 12, 157, 0, 20, 1, 0, 0, 0,
		2, 33, 1, 0, 0, 0, 4, 35, 1, 0, 0, 0, 6, 52, 1, 0, 0, 0, 8, 54, 1, 0, 0,
		0, 10, 64, 1, 0, 0, 0, 12, 66, 1, 0, 0, 0, 14, 86, 1, 0, 0, 0, 16, 94,
		1, 0, 0, 0, 18, 123, 1, 0, 0, 0, 20, 21, 3, 2, 1, 0, 21, 22, 5, 0, 0, 1,
		22, 1, 1, 0, 0, 0, 23, 34, 1, 0, 0, 0, 24, 25, 3, 4, 2, 0, 25, 26, 3, 2,
		1, 0, 26, 34, 1, 0, 0, 0, 27, 28, 3, 8, 4, 0, 28, 29, 3, 2, 1, 0, 29, 34,
		1, 0, 0, 0, 30, 31, 3, 10, 5, 0, 31, 32, 3, 2, 1, 0, 32, 34, 1, 0, 0, 0,
		33, 23, 1, 0, 0, 0, 33, 24, 1, 0, 0, 0, 33, 27, 1, 0, 0, 0, 33, 30, 1,
		0, 0, 0, 34, 3, 1, 0, 0, 0, 35, 36, 5, 1, 0, 0, 36, 37, 5, 6, 0, 0, 37,
		38, 3, 18, 9, 0, 38, 39, 5, 7, 0, 0, 39, 40, 5, 9, 0, 0, 40, 41, 3, 2,
		1, 0, 41, 42, 5, 10, 0, 0, 42, 43, 3, 6, 3, 0, 43, 5, 1, 0, 0, 0, 44, 53,
		1, 0, 0, 0, 45, 46, 5, 2, 0, 0, 46, 53, 3, 4, 2, 0, 47, 48, 5, 2, 0, 0,
		48, 49, 5, 9, 0, 0, 49, 50, 3, 2, 1, 0, 50, 51, 5, 10, 0, 0, 51, 53, 1,
		0, 0, 0, 52, 44, 1, 0, 0, 0, 52, 45, 1, 0, 0, 0, 52, 47, 1, 0, 0, 0, 53,
		7, 1, 0, 0, 0, 54, 55, 5, 3, 0, 0, 55, 56, 5, 6, 0, 0, 56, 57, 3, 18, 9,
		0, 57, 58, 5, 7, 0, 0, 58, 59, 5, 9, 0, 0, 59, 60, 3, 2, 1, 0, 60, 61,
		5, 10, 0, 0, 61, 9, 1, 0, 0, 0, 62, 65, 3, 12, 6, 0, 63, 65, 3, 14, 7,
		0, 64, 62, 1, 0, 0, 0, 64, 63, 1, 0, 0, 0, 65, 11, 1, 0, 0, 0, 66, 67,
		5, 31, 0, 0, 67, 68, 5, 17, 0, 0, 68, 69, 3, 18, 9, 0, 69, 13, 1, 0, 0,
		0, 70, 74, 5, 28, 0, 0, 71, 73, 3, 16, 8, 0, 72, 71, 1, 0, 0, 0, 73, 76,
		1, 0, 0, 0, 74, 72, 1, 0, 0, 0, 74, 75, 1, 0, 0, 0, 75, 77, 1, 0, 0, 0,
		76, 74, 1, 0, 0, 0, 77, 87, 5, 42, 0, 0, 78, 82, 5, 29, 0, 0, 79, 81, 3,
		16, 8, 0, 80, 79, 1, 0, 0, 0, 81, 84, 1, 0, 0, 0, 82, 80, 1, 0, 0, 0, 82,
		83, 1, 0, 0, 0, 83, 85, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0, 85, 87, 5, 46,
		0, 0, 86, 70, 1, 0, 0, 0, 86, 78, 1, 0, 0, 0, 87, 15, 1, 0, 0, 0, 88, 95,
		7, 0, 0, 0, 89, 95, 7, 1, 0, 0, 90, 91, 7, 2, 0, 0, 91, 92, 3, 18, 9, 0,
		92, 93, 5, 30, 0, 0, 93, 95, 1, 0, 0, 0, 94, 88, 1, 0, 0, 0, 94, 89, 1,
		0, 0, 0, 94, 90, 1, 0, 0, 0, 95, 17, 1, 0, 0, 0, 96, 97, 6, 9, -1, 0, 97,
		98, 5, 31, 0, 0, 98, 107, 5, 6, 0, 0, 99, 104, 3, 18, 9, 0, 100, 101, 5,
		26, 0, 0, 101, 103, 3, 18, 9, 0, 102, 100, 1, 0, 0, 0, 103, 106, 1, 0,
		0, 0, 104, 102, 1, 0, 0, 0, 104, 105, 1, 0, 0, 0, 105, 108, 1, 0, 0, 0,
		106, 104, 1, 0, 0, 0, 107, 99, 1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 109,
		1, 0, 0, 0, 109, 124, 5, 7, 0, 0, 110, 124, 5, 31, 0, 0, 111, 124, 5, 32,
		0, 0, 112, 124, 5, 33, 0, 0, 113, 124, 5, 4, 0, 0, 114, 124, 5, 5, 0, 0,
		115, 116, 5, 6, 0, 0, 116, 117, 3, 18, 9, 0, 117, 118, 5, 7, 0, 0, 118,
		124, 1, 0, 0, 0, 119, 120, 5, 18, 0, 0, 120, 124, 3, 18, 9, 6, 121, 122,
		5, 19, 0, 0, 122, 124, 3, 18, 9, 5, 123, 96, 1, 0, 0, 0, 123, 110, 1, 0,
		0, 0, 123, 111, 1, 0, 0, 0, 123, 112, 1, 0, 0, 0, 123, 113, 1, 0, 0, 0,
		123, 114, 1, 0, 0, 0, 123, 115, 1, 0, 0, 0, 123, 119, 1, 0, 0, 0, 123,
		121, 1, 0, 0, 0, 124, 139, 1, 0, 0, 0, 125, 126, 10, 4, 0, 0, 126, 127,
		7, 3, 0, 0, 127, 138, 3, 18, 9, 5, 128, 129, 10, 3, 0, 0, 129, 130, 7,
		4, 0, 0, 130, 138, 3, 18, 9, 4, 131, 132, 10, 2, 0, 0, 132, 133, 7, 5,
		0, 0, 133, 138, 3, 18, 9, 3, 134, 135, 10, 1, 0, 0, 135, 136, 7, 6, 0,
		0, 136, 138, 3, 18, 9, 2, 137, 125, 1, 0, 0, 0, 137, 128, 1, 0, 0, 0, 137,
		131, 1, 0, 0, 0, 137, 134, 1, 0, 0, 0, 138, 141, 1, 0, 0, 0, 139, 137,
		1, 0, 0, 0, 139, 140, 1, 0, 0, 0, 140, 19, 1, 0, 0, 0, 141, 139, 1, 0,
		0, 0, 12, 33, 52, 64, 74, 82, 86, 94, 104, 107, 123, 137, 139,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// SequenceParserInit initializes any static state used to implement SequenceParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewSequenceParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SequenceParserInit() {
	staticData := &sequenceparserParserStaticData
	staticData.once.Do(sequenceparserParserInit)
}

// NewSequenceParser produces a new parser instance for the optional input antlr.TokenStream.
func NewSequenceParser(input antlr.TokenStream) *SequenceParser {
	SequenceParserInit()
	this := new(SequenceParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &sequenceparserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "SequenceParser.g4"

	return this
}

// SequenceParser tokens.
const (
	SequenceParserEOF                         = antlr.TokenEOF
	SequenceParserIF                          = 1
	SequenceParserELSE                        = 2
	SequenceParserWHILE                       = 3
	SequenceParserTRUE                        = 4
	SequenceParserFALSE                       = 5
	SequenceParserL_PAREN                     = 6
	SequenceParserR_PAREN                     = 7
	SequenceParserLL_BRACE                    = 8
	SequenceParserL_BRACE                     = 9
	SequenceParserR_BRACE                     = 10
	SequenceParserAND                         = 11
	SequenceParserOR                          = 12
	SequenceParserEQ                          = 13
	SequenceParserN_EQ                        = 14
	SequenceParserLT_EQ                       = 15
	SequenceParserGT_EQ                       = 16
	SequenceParserASSIGN                      = 17
	SequenceParserNOT                         = 18
	SequenceParserMINUS                       = 19
	SequenceParserPLUS                        = 20
	SequenceParserTIMES                       = 21
	SequenceParserDIV                         = 22
	SequenceParserMOD                         = 23
	SequenceParserLT                          = 24
	SequenceParserGT                          = 25
	SequenceParserCOMMA                       = 26
	SequenceParserGCODE_ESCAPE                = 27
	SequenceParserOPEN_GCODE_SQ               = 28
	SequenceParserOPEN_GCODE_DQ               = 29
	SequenceParserEXIT_EXPR                   = 30
	SequenceParserIDENTIFIER                  = 31
	SequenceParserINT                         = 32
	SequenceParserFLOAT                       = 33
	SequenceParserWS                          = 34
	SequenceParserOPEN_LINE_COMMENT           = 35
	SequenceParserOPEN_BLOCK_COMMENT          = 36
	SequenceParserCLOSE_LINE_COMMENT          = 37
	SequenceParserLINE_COMMENT_TEXT           = 38
	SequenceParserCLOSE_BLOCK_COMMENT         = 39
	SequenceParserESCAPED_CLOSE_BLOCK_COMMENT = 40
	SequenceParserBLOCK_COMMENT_TEXT          = 41
	SequenceParserCLOSE_GCODE_SQ              = 42
	SequenceParserENTER_EXPR_SQ               = 43
	SequenceParserESCAPE_SEQUENCE_SQ          = 44
	SequenceParserTEXT_SQ                     = 45
	SequenceParserCLOSE_GCODE_DQ              = 46
	SequenceParserENTER_EXPR_DQ               = 47
	SequenceParserESCAPE_SEQUENCE_DQ          = 48
	SequenceParserTEXT_DQ                     = 49
)

// SequenceParser rules.
const (
	SequenceParserRULE_sequence          = 0
	SequenceParserRULE_statements        = 1
	SequenceParserRULE_ifBlock           = 2
	SequenceParserRULE_optionalElseBlock = 3
	SequenceParserRULE_whileBlock        = 4
	SequenceParserRULE_statement         = 5
	SequenceParserRULE_assignment        = 6
	SequenceParserRULE_gCode             = 7
	SequenceParserRULE_gCodePart         = 8
	SequenceParserRULE_expression        = 9
)

// ISequenceContext is an interface to support dynamic dispatch.
type ISequenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSequenceContext differentiates from other interfaces.
	IsSequenceContext()
}

type SequenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySequenceContext() *SequenceContext {
	var p = new(SequenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_sequence
	return p
}

func (*SequenceContext) IsSequenceContext() {}

func NewSequenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SequenceContext {
	var p = new(SequenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_sequence

	return p
}

func (s *SequenceContext) GetParser() antlr.Parser { return s.parser }

func (s *SequenceContext) Statements() IStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *SequenceContext) EOF() antlr.TerminalNode {
	return s.GetToken(SequenceParserEOF, 0)
}

func (s *SequenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SequenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SequenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterSequence(s)
	}
}

func (s *SequenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitSequence(s)
	}
}

func (s *SequenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitSequence(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) Sequence() (localctx ISequenceContext) {
	this := p
	_ = this

	localctx = NewSequenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SequenceParserRULE_sequence)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(20)
		p.Statements()
	}
	{
		p.SetState(21)
		p.Match(SequenceParserEOF)
	}

	return localctx
}

// IStatementsContext is an interface to support dynamic dispatch.
type IStatementsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementsContext differentiates from other interfaces.
	IsStatementsContext()
}

type StatementsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementsContext() *StatementsContext {
	var p = new(StatementsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_statements
	return p
}

func (*StatementsContext) IsStatementsContext() {}

func NewStatementsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementsContext {
	var p = new(StatementsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_statements

	return p
}

func (s *StatementsContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementsContext) IfBlock() IIfBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfBlockContext)
}

func (s *StatementsContext) Statements() IStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *StatementsContext) WhileBlock() IWhileBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhileBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhileBlockContext)
}

func (s *StatementsContext) Statement() IStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *StatementsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterStatements(s)
	}
}

func (s *StatementsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitStatements(s)
	}
}

func (s *StatementsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitStatements(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) Statements() (localctx IStatementsContext) {
	this := p
	_ = this

	localctx = NewStatementsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SequenceParserRULE_statements)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(33)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SequenceParserEOF, SequenceParserR_BRACE:
		p.EnterOuterAlt(localctx, 1)

	case SequenceParserIF:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(24)
			p.IfBlock()
		}
		{
			p.SetState(25)
			p.Statements()
		}

	case SequenceParserWHILE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(27)
			p.WhileBlock()
		}
		{
			p.SetState(28)
			p.Statements()
		}

	case SequenceParserOPEN_GCODE_SQ, SequenceParserOPEN_GCODE_DQ, SequenceParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(30)
			p.Statement()
		}
		{
			p.SetState(31)
			p.Statements()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IIfBlockContext is an interface to support dynamic dispatch.
type IIfBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIfBlockContext differentiates from other interfaces.
	IsIfBlockContext()
}

type IfBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIfBlockContext() *IfBlockContext {
	var p = new(IfBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_ifBlock
	return p
}

func (*IfBlockContext) IsIfBlockContext() {}

func NewIfBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfBlockContext {
	var p = new(IfBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_ifBlock

	return p
}

func (s *IfBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *IfBlockContext) IF() antlr.TerminalNode {
	return s.GetToken(SequenceParserIF, 0)
}

func (s *IfBlockContext) L_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_PAREN, 0)
}

func (s *IfBlockContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *IfBlockContext) R_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_PAREN, 0)
}

func (s *IfBlockContext) L_BRACE() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_BRACE, 0)
}

func (s *IfBlockContext) Statements() IStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *IfBlockContext) R_BRACE() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_BRACE, 0)
}

func (s *IfBlockContext) OptionalElseBlock() IOptionalElseBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOptionalElseBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOptionalElseBlockContext)
}

func (s *IfBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IfBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterIfBlock(s)
	}
}

func (s *IfBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitIfBlock(s)
	}
}

func (s *IfBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitIfBlock(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) IfBlock() (localctx IIfBlockContext) {
	this := p
	_ = this

	localctx = NewIfBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SequenceParserRULE_ifBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(35)
		p.Match(SequenceParserIF)
	}
	{
		p.SetState(36)
		p.Match(SequenceParserL_PAREN)
	}
	{
		p.SetState(37)
		p.expression(0)
	}
	{
		p.SetState(38)
		p.Match(SequenceParserR_PAREN)
	}
	{
		p.SetState(39)
		p.Match(SequenceParserL_BRACE)
	}
	{
		p.SetState(40)
		p.Statements()
	}
	{
		p.SetState(41)
		p.Match(SequenceParserR_BRACE)
	}
	{
		p.SetState(42)
		p.OptionalElseBlock()
	}

	return localctx
}

// IOptionalElseBlockContext is an interface to support dynamic dispatch.
type IOptionalElseBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOptionalElseBlockContext differentiates from other interfaces.
	IsOptionalElseBlockContext()
}

type OptionalElseBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOptionalElseBlockContext() *OptionalElseBlockContext {
	var p = new(OptionalElseBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_optionalElseBlock
	return p
}

func (*OptionalElseBlockContext) IsOptionalElseBlockContext() {}

func NewOptionalElseBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionalElseBlockContext {
	var p = new(OptionalElseBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_optionalElseBlock

	return p
}

func (s *OptionalElseBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *OptionalElseBlockContext) ELSE() antlr.TerminalNode {
	return s.GetToken(SequenceParserELSE, 0)
}

func (s *OptionalElseBlockContext) IfBlock() IIfBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfBlockContext)
}

func (s *OptionalElseBlockContext) L_BRACE() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_BRACE, 0)
}

func (s *OptionalElseBlockContext) Statements() IStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *OptionalElseBlockContext) R_BRACE() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_BRACE, 0)
}

func (s *OptionalElseBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OptionalElseBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OptionalElseBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterOptionalElseBlock(s)
	}
}

func (s *OptionalElseBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitOptionalElseBlock(s)
	}
}

func (s *OptionalElseBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitOptionalElseBlock(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) OptionalElseBlock() (localctx IOptionalElseBlockContext) {
	this := p
	_ = this

	localctx = NewOptionalElseBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SequenceParserRULE_optionalElseBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(45)
			p.Match(SequenceParserELSE)
		}
		{
			p.SetState(46)
			p.IfBlock()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(47)
			p.Match(SequenceParserELSE)
		}
		{
			p.SetState(48)
			p.Match(SequenceParserL_BRACE)
		}
		{
			p.SetState(49)
			p.Statements()
		}
		{
			p.SetState(50)
			p.Match(SequenceParserR_BRACE)
		}

	}

	return localctx
}

// IWhileBlockContext is an interface to support dynamic dispatch.
type IWhileBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhileBlockContext differentiates from other interfaces.
	IsWhileBlockContext()
}

type WhileBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhileBlockContext() *WhileBlockContext {
	var p = new(WhileBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_whileBlock
	return p
}

func (*WhileBlockContext) IsWhileBlockContext() {}

func NewWhileBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhileBlockContext {
	var p = new(WhileBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_whileBlock

	return p
}

func (s *WhileBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *WhileBlockContext) WHILE() antlr.TerminalNode {
	return s.GetToken(SequenceParserWHILE, 0)
}

func (s *WhileBlockContext) L_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_PAREN, 0)
}

func (s *WhileBlockContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *WhileBlockContext) R_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_PAREN, 0)
}

func (s *WhileBlockContext) L_BRACE() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_BRACE, 0)
}

func (s *WhileBlockContext) Statements() IStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *WhileBlockContext) R_BRACE() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_BRACE, 0)
}

func (s *WhileBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhileBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhileBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterWhileBlock(s)
	}
}

func (s *WhileBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitWhileBlock(s)
	}
}

func (s *WhileBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitWhileBlock(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) WhileBlock() (localctx IWhileBlockContext) {
	this := p
	_ = this

	localctx = NewWhileBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SequenceParserRULE_whileBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.Match(SequenceParserWHILE)
	}
	{
		p.SetState(55)
		p.Match(SequenceParserL_PAREN)
	}
	{
		p.SetState(56)
		p.expression(0)
	}
	{
		p.SetState(57)
		p.Match(SequenceParserR_PAREN)
	}
	{
		p.SetState(58)
		p.Match(SequenceParserL_BRACE)
	}
	{
		p.SetState(59)
		p.Statements()
	}
	{
		p.SetState(60)
		p.Match(SequenceParserR_BRACE)
	}

	return localctx
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) Assignment() IAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *StatementContext) GCode() IGCodeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGCodeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGCodeContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) Statement() (localctx IStatementContext) {
	this := p
	_ = this

	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SequenceParserRULE_statement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(64)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SequenceParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(62)
			p.Assignment()
		}

	case SequenceParserOPEN_GCODE_SQ, SequenceParserOPEN_GCODE_DQ:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(63)
			p.GCode()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAssignmentContext is an interface to support dynamic dispatch.
type IAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetIdent returns the ident token.
	GetIdent() antlr.Token

	// SetIdent sets the ident token.
	SetIdent(antlr.Token)

	// GetExpr returns the expr rule contexts.
	GetExpr() IExpressionContext

	// SetExpr sets the expr rule contexts.
	SetExpr(IExpressionContext)

	// IsAssignmentContext differentiates from other interfaces.
	IsAssignmentContext()
}

type AssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	ident  antlr.Token
	expr   IExpressionContext
}

func NewEmptyAssignmentContext() *AssignmentContext {
	var p = new(AssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_assignment
	return p
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) GetIdent() antlr.Token { return s.ident }

func (s *AssignmentContext) SetIdent(v antlr.Token) { s.ident = v }

func (s *AssignmentContext) GetExpr() IExpressionContext { return s.expr }

func (s *AssignmentContext) SetExpr(v IExpressionContext) { s.expr = v }

func (s *AssignmentContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(SequenceParserASSIGN, 0)
}

func (s *AssignmentContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SequenceParserIDENTIFIER, 0)
}

func (s *AssignmentContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (s *AssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) Assignment() (localctx IAssignmentContext) {
	this := p
	_ = this

	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SequenceParserRULE_assignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)

		var _m = p.Match(SequenceParserIDENTIFIER)

		localctx.(*AssignmentContext).ident = _m
	}
	{
		p.SetState(67)
		p.Match(SequenceParserASSIGN)
	}
	{
		p.SetState(68)

		var _x = p.expression(0)

		localctx.(*AssignmentContext).expr = _x
	}

	return localctx
}

// IGCodeContext is an interface to support dynamic dispatch.
type IGCodeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGCodeContext differentiates from other interfaces.
	IsGCodeContext()
}

type GCodeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGCodeContext() *GCodeContext {
	var p = new(GCodeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_gCode
	return p
}

func (*GCodeContext) IsGCodeContext() {}

func NewGCodeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GCodeContext {
	var p = new(GCodeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_gCode

	return p
}

func (s *GCodeContext) GetParser() antlr.Parser { return s.parser }

func (s *GCodeContext) OPEN_GCODE_SQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserOPEN_GCODE_SQ, 0)
}

func (s *GCodeContext) CLOSE_GCODE_SQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserCLOSE_GCODE_SQ, 0)
}

func (s *GCodeContext) AllGCodePart() []IGCodePartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IGCodePartContext); ok {
			len++
		}
	}

	tst := make([]IGCodePartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IGCodePartContext); ok {
			tst[i] = t.(IGCodePartContext)
			i++
		}
	}

	return tst
}

func (s *GCodeContext) GCodePart(i int) IGCodePartContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGCodePartContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGCodePartContext)
}

func (s *GCodeContext) OPEN_GCODE_DQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserOPEN_GCODE_DQ, 0)
}

func (s *GCodeContext) CLOSE_GCODE_DQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserCLOSE_GCODE_DQ, 0)
}

func (s *GCodeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GCodeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GCodeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterGCode(s)
	}
}

func (s *GCodeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitGCode(s)
	}
}

func (s *GCodeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitGCode(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) GCode() (localctx IGCodeContext) {
	this := p
	_ = this

	localctx = NewGCodeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SequenceParserRULE_gCode)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(86)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SequenceParserOPEN_GCODE_SQ:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(70)
			p.Match(SequenceParserOPEN_GCODE_SQ)
		}
		p.SetState(74)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-43)&-(0x1f+1)) == 0 && ((1<<uint((_la-43)))&((1<<(SequenceParserENTER_EXPR_SQ-43))|(1<<(SequenceParserESCAPE_SEQUENCE_SQ-43))|(1<<(SequenceParserTEXT_SQ-43))|(1<<(SequenceParserENTER_EXPR_DQ-43))|(1<<(SequenceParserESCAPE_SEQUENCE_DQ-43))|(1<<(SequenceParserTEXT_DQ-43)))) != 0 {
			{
				p.SetState(71)
				p.GCodePart()
			}

			p.SetState(76)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(77)
			p.Match(SequenceParserCLOSE_GCODE_SQ)
		}

	case SequenceParserOPEN_GCODE_DQ:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(78)
			p.Match(SequenceParserOPEN_GCODE_DQ)
		}
		p.SetState(82)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-43)&-(0x1f+1)) == 0 && ((1<<uint((_la-43)))&((1<<(SequenceParserENTER_EXPR_SQ-43))|(1<<(SequenceParserESCAPE_SEQUENCE_SQ-43))|(1<<(SequenceParserTEXT_SQ-43))|(1<<(SequenceParserENTER_EXPR_DQ-43))|(1<<(SequenceParserESCAPE_SEQUENCE_DQ-43))|(1<<(SequenceParserTEXT_DQ-43)))) != 0 {
			{
				p.SetState(79)
				p.GCodePart()
			}

			p.SetState(84)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(85)
			p.Match(SequenceParserCLOSE_GCODE_DQ)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IGCodePartContext is an interface to support dynamic dispatch.
type IGCodePartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGCodePartContext differentiates from other interfaces.
	IsGCodePartContext()
}

type GCodePartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGCodePartContext() *GCodePartContext {
	var p = new(GCodePartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_gCodePart
	return p
}

func (*GCodePartContext) IsGCodePartContext() {}

func NewGCodePartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GCodePartContext {
	var p = new(GCodePartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_gCodePart

	return p
}

func (s *GCodePartContext) GetParser() antlr.Parser { return s.parser }

func (s *GCodePartContext) CopyFrom(ctx *GCodePartContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *GCodePartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GCodePartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type GCodeSubExpressionContext struct {
	*GCodePartContext
}

func NewGCodeSubExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GCodeSubExpressionContext {
	var p = new(GCodeSubExpressionContext)

	p.GCodePartContext = NewEmptyGCodePartContext()
	p.parser = parser
	p.CopyFrom(ctx.(*GCodePartContext))

	return p
}

func (s *GCodeSubExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GCodeSubExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *GCodeSubExpressionContext) EXIT_EXPR() antlr.TerminalNode {
	return s.GetToken(SequenceParserEXIT_EXPR, 0)
}

func (s *GCodeSubExpressionContext) ENTER_EXPR_SQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserENTER_EXPR_SQ, 0)
}

func (s *GCodeSubExpressionContext) ENTER_EXPR_DQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserENTER_EXPR_DQ, 0)
}

func (s *GCodeSubExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterGCodeSubExpression(s)
	}
}

func (s *GCodeSubExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitGCodeSubExpression(s)
	}
}

func (s *GCodeSubExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitGCodeSubExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type GCodeTextContext struct {
	*GCodePartContext
}

func NewGCodeTextContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GCodeTextContext {
	var p = new(GCodeTextContext)

	p.GCodePartContext = NewEmptyGCodePartContext()
	p.parser = parser
	p.CopyFrom(ctx.(*GCodePartContext))

	return p
}

func (s *GCodeTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GCodeTextContext) TEXT_SQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserTEXT_SQ, 0)
}

func (s *GCodeTextContext) TEXT_DQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserTEXT_DQ, 0)
}

func (s *GCodeTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterGCodeText(s)
	}
}

func (s *GCodeTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitGCodeText(s)
	}
}

func (s *GCodeTextContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitGCodeText(s)

	default:
		return t.VisitChildren(s)
	}
}

type GCodeEscapedTextContext struct {
	*GCodePartContext
}

func NewGCodeEscapedTextContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GCodeEscapedTextContext {
	var p = new(GCodeEscapedTextContext)

	p.GCodePartContext = NewEmptyGCodePartContext()
	p.parser = parser
	p.CopyFrom(ctx.(*GCodePartContext))

	return p
}

func (s *GCodeEscapedTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GCodeEscapedTextContext) ESCAPE_SEQUENCE_SQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserESCAPE_SEQUENCE_SQ, 0)
}

func (s *GCodeEscapedTextContext) ESCAPE_SEQUENCE_DQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserESCAPE_SEQUENCE_DQ, 0)
}

func (s *GCodeEscapedTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterGCodeEscapedText(s)
	}
}

func (s *GCodeEscapedTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitGCodeEscapedText(s)
	}
}

func (s *GCodeEscapedTextContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitGCodeEscapedText(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) GCodePart() (localctx IGCodePartContext) {
	this := p
	_ = this

	localctx = NewGCodePartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SequenceParserRULE_gCodePart)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(94)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SequenceParserTEXT_SQ, SequenceParserTEXT_DQ:
		localctx = NewGCodeTextContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(88)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SequenceParserTEXT_SQ || _la == SequenceParserTEXT_DQ) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case SequenceParserESCAPE_SEQUENCE_SQ, SequenceParserESCAPE_SEQUENCE_DQ:
		localctx = NewGCodeEscapedTextContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(89)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SequenceParserESCAPE_SEQUENCE_SQ || _la == SequenceParserESCAPE_SEQUENCE_DQ) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case SequenceParserENTER_EXPR_SQ, SequenceParserENTER_EXPR_DQ:
		localctx = NewGCodeSubExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(90)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SequenceParserENTER_EXPR_SQ || _la == SequenceParserENTER_EXPR_DQ) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(91)
			p.expression(0)
		}
		{
			p.SetState(92)
			p.Match(SequenceParserEXIT_EXPR)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SequenceParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SequenceParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type IdentExprContext struct {
	*ExpressionContext
}

func NewIdentExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentExprContext {
	var p = new(IdentExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *IdentExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SequenceParserIDENTIFIER, 0)
}

func (s *IdentExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterIdentExpr(s)
	}
}

func (s *IdentExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitIdentExpr(s)
	}
}

func (s *IdentExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitIdentExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type FloatExprContext struct {
	*ExpressionContext
}

func NewFloatExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FloatExprContext {
	var p = new(FloatExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *FloatExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FloatExprContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(SequenceParserFLOAT, 0)
}

func (s *FloatExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterFloatExpr(s)
	}
}

func (s *FloatExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitFloatExpr(s)
	}
}

func (s *FloatExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitFloatExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type UnaryOpExprContext struct {
	*ExpressionContext
}

func NewUnaryOpExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnaryOpExprContext {
	var p = new(UnaryOpExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *UnaryOpExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryOpExprContext) NOT() antlr.TerminalNode {
	return s.GetToken(SequenceParserNOT, 0)
}

func (s *UnaryOpExprContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *UnaryOpExprContext) MINUS() antlr.TerminalNode {
	return s.GetToken(SequenceParserMINUS, 0)
}

func (s *UnaryOpExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterUnaryOpExpr(s)
	}
}

func (s *UnaryOpExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitUnaryOpExpr(s)
	}
}

func (s *UnaryOpExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitUnaryOpExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type IntExprContext struct {
	*ExpressionContext
}

func NewIntExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntExprContext {
	var p = new(IntExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *IntExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntExprContext) INT() antlr.TerminalNode {
	return s.GetToken(SequenceParserINT, 0)
}

func (s *IntExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterIntExpr(s)
	}
}

func (s *IntExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitIntExpr(s)
	}
}

func (s *IntExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitIntExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionCallContext struct {
	*ExpressionContext
}

func NewFunctionCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallContext {
	var p = new(FunctionCallContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SequenceParserIDENTIFIER, 0)
}

func (s *FunctionCallContext) L_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_PAREN, 0)
}

func (s *FunctionCallContext) R_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_PAREN, 0)
}

func (s *FunctionCallContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FunctionCallContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SequenceParserCOMMA)
}

func (s *FunctionCallContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SequenceParserCOMMA, i)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type BinaryOpExprContext struct {
	*ExpressionContext
}

func NewBinaryOpExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BinaryOpExprContext {
	var p = new(BinaryOpExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *BinaryOpExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryOpExprContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *BinaryOpExprContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *BinaryOpExprContext) TIMES() antlr.TerminalNode {
	return s.GetToken(SequenceParserTIMES, 0)
}

func (s *BinaryOpExprContext) DIV() antlr.TerminalNode {
	return s.GetToken(SequenceParserDIV, 0)
}

func (s *BinaryOpExprContext) MOD() antlr.TerminalNode {
	return s.GetToken(SequenceParserMOD, 0)
}

func (s *BinaryOpExprContext) PLUS() antlr.TerminalNode {
	return s.GetToken(SequenceParserPLUS, 0)
}

func (s *BinaryOpExprContext) MINUS() antlr.TerminalNode {
	return s.GetToken(SequenceParserMINUS, 0)
}

func (s *BinaryOpExprContext) EQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserEQ, 0)
}

func (s *BinaryOpExprContext) N_EQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserN_EQ, 0)
}

func (s *BinaryOpExprContext) LT() antlr.TerminalNode {
	return s.GetToken(SequenceParserLT, 0)
}

func (s *BinaryOpExprContext) LT_EQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserLT_EQ, 0)
}

func (s *BinaryOpExprContext) GT() antlr.TerminalNode {
	return s.GetToken(SequenceParserGT, 0)
}

func (s *BinaryOpExprContext) GT_EQ() antlr.TerminalNode {
	return s.GetToken(SequenceParserGT_EQ, 0)
}

func (s *BinaryOpExprContext) OR() antlr.TerminalNode {
	return s.GetToken(SequenceParserOR, 0)
}

func (s *BinaryOpExprContext) AND() antlr.TerminalNode {
	return s.GetToken(SequenceParserAND, 0)
}

func (s *BinaryOpExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterBinaryOpExpr(s)
	}
}

func (s *BinaryOpExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitBinaryOpExpr(s)
	}
}

func (s *BinaryOpExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitBinaryOpExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type BoolExprContext struct {
	*ExpressionContext
}

func NewBoolExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolExprContext {
	var p = new(BoolExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *BoolExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprContext) TRUE() antlr.TerminalNode {
	return s.GetToken(SequenceParserTRUE, 0)
}

func (s *BoolExprContext) FALSE() antlr.TerminalNode {
	return s.GetToken(SequenceParserFALSE, 0)
}

func (s *BoolExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterBoolExpr(s)
	}
}

func (s *BoolExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitBoolExpr(s)
	}
}

func (s *BoolExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitBoolExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenExprContext struct {
	*ExpressionContext
}

func NewParenExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenExprContext {
	var p = new(ParenExprContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ParenExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenExprContext) L_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserL_PAREN, 0)
}

func (s *ParenExprContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenExprContext) R_PAREN() antlr.TerminalNode {
	return s.GetToken(SequenceParserR_PAREN, 0)
}

func (s *ParenExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.EnterParenExpr(s)
	}
}

func (s *ParenExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SequenceParserListener); ok {
		listenerT.ExitParenExpr(s)
	}
}

func (s *ParenExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SequenceParserVisitor:
		return t.VisitParenExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SequenceParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *SequenceParser) expression(_p int) (localctx IExpressionContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 18
	p.EnterRecursionRule(localctx, 18, SequenceParserRULE_expression, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFunctionCallContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(97)
			p.Match(SequenceParserIDENTIFIER)
		}
		{
			p.SetState(98)
			p.Match(SequenceParserL_PAREN)
		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la-4)&-(0x1f+1)) == 0 && ((1<<uint((_la-4)))&((1<<(SequenceParserTRUE-4))|(1<<(SequenceParserFALSE-4))|(1<<(SequenceParserL_PAREN-4))|(1<<(SequenceParserNOT-4))|(1<<(SequenceParserMINUS-4))|(1<<(SequenceParserIDENTIFIER-4))|(1<<(SequenceParserINT-4))|(1<<(SequenceParserFLOAT-4)))) != 0 {
			{
				p.SetState(99)
				p.expression(0)
			}
			p.SetState(104)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for _la == SequenceParserCOMMA {
				{
					p.SetState(100)
					p.Match(SequenceParserCOMMA)
				}
				{
					p.SetState(101)
					p.expression(0)
				}

				p.SetState(106)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(109)
			p.Match(SequenceParserR_PAREN)
		}

	case 2:
		localctx = NewIdentExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(110)
			p.Match(SequenceParserIDENTIFIER)
		}

	case 3:
		localctx = NewIntExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(111)
			p.Match(SequenceParserINT)
		}

	case 4:
		localctx = NewFloatExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(112)
			p.Match(SequenceParserFLOAT)
		}

	case 5:
		localctx = NewBoolExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(113)
			p.Match(SequenceParserTRUE)
		}

	case 6:
		localctx = NewBoolExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(114)
			p.Match(SequenceParserFALSE)
		}

	case 7:
		localctx = NewParenExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(115)
			p.Match(SequenceParserL_PAREN)
		}
		{
			p.SetState(116)
			p.expression(0)
		}
		{
			p.SetState(117)
			p.Match(SequenceParserR_PAREN)
		}

	case 8:
		localctx = NewUnaryOpExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(119)
			p.Match(SequenceParserNOT)
		}
		{
			p.SetState(120)
			p.expression(6)
		}

	case 9:
		localctx = NewUnaryOpExprContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(121)
			p.Match(SequenceParserMINUS)
		}
		{
			p.SetState(122)
			p.expression(5)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(139)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(137)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
			case 1:
				localctx = NewBinaryOpExprContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SequenceParserRULE_expression)
				p.SetState(125)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(126)
					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SequenceParserTIMES)|(1<<SequenceParserDIV)|(1<<SequenceParserMOD))) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(127)
					p.expression(5)
				}

			case 2:
				localctx = NewBinaryOpExprContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SequenceParserRULE_expression)
				p.SetState(128)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(129)
					_la = p.GetTokenStream().LA(1)

					if !(_la == SequenceParserMINUS || _la == SequenceParserPLUS) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(130)
					p.expression(4)
				}

			case 3:
				localctx = NewBinaryOpExprContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SequenceParserRULE_expression)
				p.SetState(131)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(132)
					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SequenceParserEQ)|(1<<SequenceParserN_EQ)|(1<<SequenceParserLT_EQ)|(1<<SequenceParserGT_EQ)|(1<<SequenceParserLT)|(1<<SequenceParserGT))) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(133)
					p.expression(3)
				}

			case 4:
				localctx = NewBinaryOpExprContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SequenceParserRULE_expression)
				p.SetState(134)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(135)
					_la = p.GetTokenStream().LA(1)

					if !(_la == SequenceParserAND || _la == SequenceParserOR) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(136)
					p.expression(2)
				}

			}

		}
		p.SetState(141)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())
	}

	return localctx
}

func (p *SequenceParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 9:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *SequenceParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
