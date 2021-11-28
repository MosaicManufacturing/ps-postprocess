module mosaicmfg.com/ps-postprocess/printerscript

go 1.16

require (
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20211128173911-14938591d7dc
	mosaicmfg.com/ps-postprocess/gcode v0.0.0
	mosaicmfg.com/ps-postprocess/printerscript/interpreter v0.0.0
)

replace mosaicmfg.com/ps-postprocess/gcode => ../gcode
replace mosaicmfg.com/ps-postprocess/printerscript/interpreter => ./interpreter
