// Code generated from SequenceLexer.g4 by ANTLR 4.10.1. DO NOT EDIT.

package interpreter

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type SequenceLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var sequencelexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sequencelexerLexerInit() {
	staticData := &sequencelexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE", "LINE_COMMENT", "BLOCK_COMMENT", "GCODE_SQ", "GCODE_DQ",
	}
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
		"DIGIT", "IF", "ELSE", "WHILE", "TRUE", "FALSE", "L_PAREN", "R_PAREN",
		"LL_BRACE", "L_BRACE", "R_BRACE", "AND", "OR", "EQ", "N_EQ", "LT_EQ",
		"GT_EQ", "ASSIGN", "NOT", "MINUS", "PLUS", "TIMES", "DIV", "MOD", "LT",
		"GT", "COMMA", "GCODE_ESCAPE", "OPEN_GCODE_SQ", "OPEN_GCODE_DQ", "EXIT_EXPR",
		"IDENTIFIER", "INT", "FLOAT", "WS", "OPEN_LINE_COMMENT", "OPEN_BLOCK_COMMENT",
		"CLOSE_LINE_COMMENT", "LINE_COMMENT_TEXT", "CLOSE_BLOCK_COMMENT", "ESCAPED_CLOSE_BLOCK_COMMENT",
		"BLOCK_COMMENT_TEXT", "CLOSE_GCODE_SQ", "ENTER_EXPR_SQ", "ESCAPE_SEQUENCE_SQ",
		"TEXT_SQ", "CLOSE_GCODE_DQ", "ENTER_EXPR_DQ", "ESCAPE_SEQUENCE_DQ",
		"TEXT_DQ",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 49, 314, 6, -1, 6, -1, 6, -1, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7,
		1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7,
		7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2,
		13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18,
		7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7,
		23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28,
		2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2,
		34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39,
		7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7,
		44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47, 7, 47, 2, 48, 7, 48, 2, 49, 7, 49,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5,
		1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9,
		1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1,
		13, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17,
		1, 17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 20, 1, 20, 1, 21, 1, 21, 1, 22, 1,
		22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27,
		1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 3, 27, 191, 8, 27, 1,
		28, 1, 28, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 30, 1, 30, 1, 30,
		1, 30, 1, 30, 1, 31, 1, 31, 5, 31, 208, 8, 31, 10, 31, 12, 31, 211, 9,
		31, 1, 32, 4, 32, 214, 8, 32, 11, 32, 12, 32, 215, 1, 33, 5, 33, 219, 8,
		33, 10, 33, 12, 33, 222, 9, 33, 1, 33, 1, 33, 4, 33, 226, 8, 33, 11, 33,
		12, 33, 227, 1, 34, 4, 34, 231, 8, 34, 11, 34, 12, 34, 232, 1, 34, 1, 34,
		1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36, 1, 36, 1, 36, 1, 36, 1,
		36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37, 1, 37, 1, 38, 4, 38, 255, 8, 38,
		11, 38, 12, 38, 256, 1, 38, 1, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1,
		39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 41, 4, 41, 273, 8, 41, 11, 41,
		12, 41, 274, 1, 41, 1, 41, 1, 42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1,
		43, 1, 43, 1, 44, 1, 44, 1, 45, 4, 45, 290, 8, 45, 11, 45, 12, 45, 291,
		1, 45, 3, 45, 295, 8, 45, 1, 46, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1,
		47, 1, 47, 1, 48, 1, 48, 1, 49, 4, 49, 308, 8, 49, 11, 49, 12, 49, 309,
		1, 49, 3, 49, 313, 8, 49, 0, 0, 50, 5, 0, 7, 1, 9, 2, 11, 3, 13, 4, 15,
		5, 17, 6, 19, 7, 21, 8, 23, 9, 25, 10, 27, 11, 29, 12, 31, 13, 33, 14,
		35, 15, 37, 16, 39, 17, 41, 18, 43, 19, 45, 20, 47, 21, 49, 22, 51, 23,
		53, 24, 55, 25, 57, 26, 59, 27, 61, 28, 63, 29, 65, 30, 67, 31, 69, 32,
		71, 33, 73, 34, 75, 35, 77, 36, 79, 37, 81, 38, 83, 39, 85, 40, 87, 41,
		89, 42, 91, 43, 93, 44, 95, 45, 97, 46, 99, 47, 101, 48, 103, 49, 5, 0,
		1, 2, 3, 4, 9, 1, 0, 48, 57, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57,
		65, 90, 95, 95, 97, 122, 2, 0, 9, 10, 32, 32, 1, 0, 10, 10, 1, 0, 47, 47,
		1, 0, 42, 42, 4, 0, 10, 10, 39, 39, 92, 92, 123, 123, 4, 0, 10, 10, 34,
		34, 92, 92, 123, 123, 322, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1,
		0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19,
		1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0,
		27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0,
		0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0,
		0, 0, 43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0,
		0, 0, 0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1,
		0, 0, 0, 0, 59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 65,
		1, 0, 0, 0, 0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 0, 0,
		73, 1, 0, 0, 0, 0, 75, 1, 0, 0, 0, 0, 77, 1, 0, 0, 0, 1, 79, 1, 0, 0, 0,
		1, 81, 1, 0, 0, 0, 2, 83, 1, 0, 0, 0, 2, 85, 1, 0, 0, 0, 2, 87, 1, 0, 0,
		0, 3, 89, 1, 0, 0, 0, 3, 91, 1, 0, 0, 0, 3, 93, 1, 0, 0, 0, 3, 95, 1, 0,
		0, 0, 4, 97, 1, 0, 0, 0, 4, 99, 1, 0, 0, 0, 4, 101, 1, 0, 0, 0, 4, 103,
		1, 0, 0, 0, 5, 105, 1, 0, 0, 0, 7, 107, 1, 0, 0, 0, 9, 110, 1, 0, 0, 0,
		11, 115, 1, 0, 0, 0, 13, 121, 1, 0, 0, 0, 15, 126, 1, 0, 0, 0, 17, 132,
		1, 0, 0, 0, 19, 134, 1, 0, 0, 0, 21, 136, 1, 0, 0, 0, 23, 139, 1, 0, 0,
		0, 25, 141, 1, 0, 0, 0, 27, 143, 1, 0, 0, 0, 29, 146, 1, 0, 0, 0, 31, 149,
		1, 0, 0, 0, 33, 152, 1, 0, 0, 0, 35, 155, 1, 0, 0, 0, 37, 158, 1, 0, 0,
		0, 39, 161, 1, 0, 0, 0, 41, 163, 1, 0, 0, 0, 43, 165, 1, 0, 0, 0, 45, 167,
		1, 0, 0, 0, 47, 169, 1, 0, 0, 0, 49, 171, 1, 0, 0, 0, 51, 173, 1, 0, 0,
		0, 53, 175, 1, 0, 0, 0, 55, 177, 1, 0, 0, 0, 57, 179, 1, 0, 0, 0, 59, 190,
		1, 0, 0, 0, 61, 192, 1, 0, 0, 0, 63, 196, 1, 0, 0, 0, 65, 200, 1, 0, 0,
		0, 67, 205, 1, 0, 0, 0, 69, 213, 1, 0, 0, 0, 71, 220, 1, 0, 0, 0, 73, 230,
		1, 0, 0, 0, 75, 236, 1, 0, 0, 0, 77, 242, 1, 0, 0, 0, 79, 248, 1, 0, 0,
		0, 81, 254, 1, 0, 0, 0, 83, 260, 1, 0, 0, 0, 85, 266, 1, 0, 0, 0, 87, 272,
		1, 0, 0, 0, 89, 278, 1, 0, 0, 0, 91, 282, 1, 0, 0, 0, 93, 286, 1, 0, 0,
		0, 95, 294, 1, 0, 0, 0, 97, 296, 1, 0, 0, 0, 99, 300, 1, 0, 0, 0, 101,
		304, 1, 0, 0, 0, 103, 312, 1, 0, 0, 0, 105, 106, 7, 0, 0, 0, 106, 6, 1,
		0, 0, 0, 107, 108, 5, 105, 0, 0, 108, 109, 5, 102, 0, 0, 109, 8, 1, 0,
		0, 0, 110, 111, 5, 101, 0, 0, 111, 112, 5, 108, 0, 0, 112, 113, 5, 115,
		0, 0, 113, 114, 5, 101, 0, 0, 114, 10, 1, 0, 0, 0, 115, 116, 5, 119, 0,
		0, 116, 117, 5, 104, 0, 0, 117, 118, 5, 105, 0, 0, 118, 119, 5, 108, 0,
		0, 119, 120, 5, 101, 0, 0, 120, 12, 1, 0, 0, 0, 121, 122, 5, 116, 0, 0,
		122, 123, 5, 114, 0, 0, 123, 124, 5, 117, 0, 0, 124, 125, 5, 101, 0, 0,
		125, 14, 1, 0, 0, 0, 126, 127, 5, 102, 0, 0, 127, 128, 5, 97, 0, 0, 128,
		129, 5, 108, 0, 0, 129, 130, 5, 115, 0, 0, 130, 131, 5, 101, 0, 0, 131,
		16, 1, 0, 0, 0, 132, 133, 5, 40, 0, 0, 133, 18, 1, 0, 0, 0, 134, 135, 5,
		41, 0, 0, 135, 20, 1, 0, 0, 0, 136, 137, 5, 123, 0, 0, 137, 138, 5, 123,
		0, 0, 138, 22, 1, 0, 0, 0, 139, 140, 5, 123, 0, 0, 140, 24, 1, 0, 0, 0,
		141, 142, 5, 125, 0, 0, 142, 26, 1, 0, 0, 0, 143, 144, 5, 38, 0, 0, 144,
		145, 5, 38, 0, 0, 145, 28, 1, 0, 0, 0, 146, 147, 5, 124, 0, 0, 147, 148,
		5, 124, 0, 0, 148, 30, 1, 0, 0, 0, 149, 150, 5, 61, 0, 0, 150, 151, 5,
		61, 0, 0, 151, 32, 1, 0, 0, 0, 152, 153, 5, 33, 0, 0, 153, 154, 5, 61,
		0, 0, 154, 34, 1, 0, 0, 0, 155, 156, 5, 60, 0, 0, 156, 157, 5, 61, 0, 0,
		157, 36, 1, 0, 0, 0, 158, 159, 5, 62, 0, 0, 159, 160, 5, 61, 0, 0, 160,
		38, 1, 0, 0, 0, 161, 162, 5, 61, 0, 0, 162, 40, 1, 0, 0, 0, 163, 164, 5,
		33, 0, 0, 164, 42, 1, 0, 0, 0, 165, 166, 5, 45, 0, 0, 166, 44, 1, 0, 0,
		0, 167, 168, 5, 43, 0, 0, 168, 46, 1, 0, 0, 0, 169, 170, 5, 42, 0, 0, 170,
		48, 1, 0, 0, 0, 171, 172, 5, 47, 0, 0, 172, 50, 1, 0, 0, 0, 173, 174, 5,
		37, 0, 0, 174, 52, 1, 0, 0, 0, 175, 176, 5, 60, 0, 0, 176, 54, 1, 0, 0,
		0, 177, 178, 5, 62, 0, 0, 178, 56, 1, 0, 0, 0, 179, 180, 5, 44, 0, 0, 180,
		58, 1, 0, 0, 0, 181, 182, 5, 92, 0, 0, 182, 183, 5, 123, 0, 0, 183, 191,
		5, 123, 0, 0, 184, 185, 5, 92, 0, 0, 185, 191, 5, 39, 0, 0, 186, 187, 5,
		92, 0, 0, 187, 191, 5, 34, 0, 0, 188, 189, 5, 92, 0, 0, 189, 191, 5, 92,
		0, 0, 190, 181, 1, 0, 0, 0, 190, 184, 1, 0, 0, 0, 190, 186, 1, 0, 0, 0,
		190, 188, 1, 0, 0, 0, 191, 60, 1, 0, 0, 0, 192, 193, 5, 39, 0, 0, 193,
		194, 1, 0, 0, 0, 194, 195, 6, 28, 0, 0, 195, 62, 1, 0, 0, 0, 196, 197,
		5, 34, 0, 0, 197, 198, 1, 0, 0, 0, 198, 199, 6, 29, 1, 0, 199, 64, 1, 0,
		0, 0, 200, 201, 5, 125, 0, 0, 201, 202, 5, 125, 0, 0, 202, 203, 1, 0, 0,
		0, 203, 204, 6, 30, 2, 0, 204, 66, 1, 0, 0, 0, 205, 209, 7, 1, 0, 0, 206,
		208, 7, 2, 0, 0, 207, 206, 1, 0, 0, 0, 208, 211, 1, 0, 0, 0, 209, 207,
		1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 68, 1, 0, 0, 0, 211, 209, 1, 0,
		0, 0, 212, 214, 3, 5, 0, 0, 213, 212, 1, 0, 0, 0, 214, 215, 1, 0, 0, 0,
		215, 213, 1, 0, 0, 0, 215, 216, 1, 0, 0, 0, 216, 70, 1, 0, 0, 0, 217, 219,
		3, 5, 0, 0, 218, 217, 1, 0, 0, 0, 219, 222, 1, 0, 0, 0, 220, 218, 1, 0,
		0, 0, 220, 221, 1, 0, 0, 0, 221, 223, 1, 0, 0, 0, 222, 220, 1, 0, 0, 0,
		223, 225, 5, 46, 0, 0, 224, 226, 3, 5, 0, 0, 225, 224, 1, 0, 0, 0, 226,
		227, 1, 0, 0, 0, 227, 225, 1, 0, 0, 0, 227, 228, 1, 0, 0, 0, 228, 72, 1,
		0, 0, 0, 229, 231, 7, 3, 0, 0, 230, 229, 1, 0, 0, 0, 231, 232, 1, 0, 0,
		0, 232, 230, 1, 0, 0, 0, 232, 233, 1, 0, 0, 0, 233, 234, 1, 0, 0, 0, 234,
		235, 6, 34, 3, 0, 235, 74, 1, 0, 0, 0, 236, 237, 5, 47, 0, 0, 237, 238,
		5, 47, 0, 0, 238, 239, 1, 0, 0, 0, 239, 240, 6, 35, 4, 0, 240, 241, 6,
		35, 5, 0, 241, 76, 1, 0, 0, 0, 242, 243, 5, 47, 0, 0, 243, 244, 5, 42,
		0, 0, 244, 245, 1, 0, 0, 0, 245, 246, 6, 36, 6, 0, 246, 247, 6, 36, 5,
		0, 247, 78, 1, 0, 0, 0, 248, 249, 5, 10, 0, 0, 249, 250, 1, 0, 0, 0, 250,
		251, 6, 37, 2, 0, 251, 252, 6, 37, 5, 0, 252, 80, 1, 0, 0, 0, 253, 255,
		8, 4, 0, 0, 254, 253, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 254, 1, 0,
		0, 0, 256, 257, 1, 0, 0, 0, 257, 258, 1, 0, 0, 0, 258, 259, 6, 38, 5, 0,
		259, 82, 1, 0, 0, 0, 260, 261, 5, 42, 0, 0, 261, 262, 5, 47, 0, 0, 262,
		263, 1, 0, 0, 0, 263, 264, 6, 39, 2, 0, 264, 265, 6, 39, 5, 0, 265, 84,
		1, 0, 0, 0, 266, 267, 5, 42, 0, 0, 267, 268, 8, 5, 0, 0, 268, 269, 1, 0,
		0, 0, 269, 270, 6, 40, 5, 0, 270, 86, 1, 0, 0, 0, 271, 273, 8, 6, 0, 0,
		272, 271, 1, 0, 0, 0, 273, 274, 1, 0, 0, 0, 274, 272, 1, 0, 0, 0, 274,
		275, 1, 0, 0, 0, 275, 276, 1, 0, 0, 0, 276, 277, 6, 41, 5, 0, 277, 88,
		1, 0, 0, 0, 278, 279, 5, 39, 0, 0, 279, 280, 1, 0, 0, 0, 280, 281, 6, 42,
		2, 0, 281, 90, 1, 0, 0, 0, 282, 283, 3, 21, 8, 0, 283, 284, 1, 0, 0, 0,
		284, 285, 6, 43, 7, 0, 285, 92, 1, 0, 0, 0, 286, 287, 3, 59, 27, 0, 287,
		94, 1, 0, 0, 0, 288, 290, 8, 7, 0, 0, 289, 288, 1, 0, 0, 0, 290, 291, 1,
		0, 0, 0, 291, 289, 1, 0, 0, 0, 291, 292, 1, 0, 0, 0, 292, 295, 1, 0, 0,
		0, 293, 295, 3, 23, 9, 0, 294, 289, 1, 0, 0, 0, 294, 293, 1, 0, 0, 0, 295,
		96, 1, 0, 0, 0, 296, 297, 5, 34, 0, 0, 297, 298, 1, 0, 0, 0, 298, 299,
		6, 46, 2, 0, 299, 98, 1, 0, 0, 0, 300, 301, 3, 21, 8, 0, 301, 302, 1, 0,
		0, 0, 302, 303, 6, 47, 7, 0, 303, 100, 1, 0, 0, 0, 304, 305, 3, 59, 27,
		0, 305, 102, 1, 0, 0, 0, 306, 308, 8, 8, 0, 0, 307, 306, 1, 0, 0, 0, 308,
		309, 1, 0, 0, 0, 309, 307, 1, 0, 0, 0, 309, 310, 1, 0, 0, 0, 310, 313,
		1, 0, 0, 0, 311, 313, 3, 23, 9, 0, 312, 307, 1, 0, 0, 0, 312, 311, 1, 0,
		0, 0, 313, 104, 1, 0, 0, 0, 17, 0, 1, 2, 3, 4, 190, 209, 215, 220, 227,
		232, 256, 274, 291, 294, 309, 312, 8, 5, 3, 0, 5, 4, 0, 4, 0, 0, 0, 1,
		0, 5, 1, 0, 6, 0, 0, 5, 2, 0, 5, 0, 0,
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

// SequenceLexerInit initializes any static state used to implement SequenceLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewSequenceLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SequenceLexerInit() {
	staticData := &sequencelexerLexerStaticData
	staticData.once.Do(sequencelexerLexerInit)
}

// NewSequenceLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewSequenceLexer(input antlr.CharStream) *SequenceLexer {
	SequenceLexerInit()
	l := new(SequenceLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &sequencelexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "SequenceLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// SequenceLexer tokens.
const (
	SequenceLexerIF                          = 1
	SequenceLexerELSE                        = 2
	SequenceLexerWHILE                       = 3
	SequenceLexerTRUE                        = 4
	SequenceLexerFALSE                       = 5
	SequenceLexerL_PAREN                     = 6
	SequenceLexerR_PAREN                     = 7
	SequenceLexerLL_BRACE                    = 8
	SequenceLexerL_BRACE                     = 9
	SequenceLexerR_BRACE                     = 10
	SequenceLexerAND                         = 11
	SequenceLexerOR                          = 12
	SequenceLexerEQ                          = 13
	SequenceLexerN_EQ                        = 14
	SequenceLexerLT_EQ                       = 15
	SequenceLexerGT_EQ                       = 16
	SequenceLexerASSIGN                      = 17
	SequenceLexerNOT                         = 18
	SequenceLexerMINUS                       = 19
	SequenceLexerPLUS                        = 20
	SequenceLexerTIMES                       = 21
	SequenceLexerDIV                         = 22
	SequenceLexerMOD                         = 23
	SequenceLexerLT                          = 24
	SequenceLexerGT                          = 25
	SequenceLexerCOMMA                       = 26
	SequenceLexerGCODE_ESCAPE                = 27
	SequenceLexerOPEN_GCODE_SQ               = 28
	SequenceLexerOPEN_GCODE_DQ               = 29
	SequenceLexerEXIT_EXPR                   = 30
	SequenceLexerIDENTIFIER                  = 31
	SequenceLexerINT                         = 32
	SequenceLexerFLOAT                       = 33
	SequenceLexerWS                          = 34
	SequenceLexerOPEN_LINE_COMMENT           = 35
	SequenceLexerOPEN_BLOCK_COMMENT          = 36
	SequenceLexerCLOSE_LINE_COMMENT          = 37
	SequenceLexerLINE_COMMENT_TEXT           = 38
	SequenceLexerCLOSE_BLOCK_COMMENT         = 39
	SequenceLexerESCAPED_CLOSE_BLOCK_COMMENT = 40
	SequenceLexerBLOCK_COMMENT_TEXT          = 41
	SequenceLexerCLOSE_GCODE_SQ              = 42
	SequenceLexerENTER_EXPR_SQ               = 43
	SequenceLexerESCAPE_SEQUENCE_SQ          = 44
	SequenceLexerTEXT_SQ                     = 45
	SequenceLexerCLOSE_GCODE_DQ              = 46
	SequenceLexerENTER_EXPR_DQ               = 47
	SequenceLexerESCAPE_SEQUENCE_DQ          = 48
	SequenceLexerTEXT_DQ                     = 49
)

// SequenceLexer modes.
const (
	SequenceLexerLINE_COMMENT = iota + 1
	SequenceLexerBLOCK_COMMENT
	SequenceLexerGCODE_SQ
	SequenceLexerGCODE_DQ
)
