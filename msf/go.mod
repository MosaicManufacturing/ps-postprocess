module mosaicmfg.com/ps-postprocess/msf

go 1.16

require (
	mosaicmfg.com/ps-postprocess/gcode v0.0.0
	mosaicmfg.com/ps-postprocess/printerscript v0.0.0
	mosaicmfg.com/ps-postprocess/sequences v0.0.0
)

replace (
	mosaicmfg.com/ps-postprocess/gcode => ../gcode
	mosaicmfg.com/ps-postprocess/printerscript => ../printerscript
	mosaicmfg.com/ps-postprocess/printerscript/interpreter => ../printerscript/interpreter
	mosaicmfg.com/ps-postprocess/sequences => ../sequences
)
