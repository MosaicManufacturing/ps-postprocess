module mosaicmfg.com/ps-postprocess

go 1.16

require (
	github.com/google/go-licenses v0.0.0-20211006200916-ceb292363ec8 // indirect
	mosaicmfg.com/ps-postprocess/comments v0.0.0
	mosaicmfg.com/ps-postprocess/flashforge v0.0.0
	mosaicmfg.com/ps-postprocess/msf v0.0.0
	mosaicmfg.com/ps-postprocess/ptp v0.0.0
	mosaicmfg.com/ps-postprocess/sequences v0.0.0
	mosaicmfg.com/ps-postprocess/ultimaker v0.0.0
)

replace (
	mosaicmfg.com/ps-postprocess/comments => ./comments
	mosaicmfg.com/ps-postprocess/flashforge => ./flashforge
	mosaicmfg.com/ps-postprocess/gcode => ./gcode
	mosaicmfg.com/ps-postprocess/msf => ./msf
	mosaicmfg.com/ps-postprocess/printerscript => ./printerscript
	mosaicmfg.com/ps-postprocess/printerscript/interpreter => ./printerscript/interpreter
	mosaicmfg.com/ps-postprocess/ptp => ./ptp
	mosaicmfg.com/ps-postprocess/sequences => ./sequences
	mosaicmfg.com/ps-postprocess/ultimaker => ./ultimaker
)
