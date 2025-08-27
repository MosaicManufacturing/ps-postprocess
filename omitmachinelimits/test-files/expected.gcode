; external perimeters extrusion width = 0.40mm
; perimeters extrusion width = 0.40mm
; infill extrusion width = 0.40mm
; solid infill extrusion width = 0.40mm
; top infill extrusion width = 0.40mm
; support material extrusion width = 0.40mm
; first layer extrusion width = 0.40mm


; printing object Group id:0 copy 0
; stop printing object Group id:0 copy 0

;TYPE:Custom
M115 U3.1.0 ; use the latest firmware version
G28 W ; home all axes without mesh bed leveling
M140 S60 ; start heating the bed
M190 S60 ; wait until bed heated
M104 S220 ; start heating extruder
G80 ; run mesh bed leveling routine
M109 S220 ; wait for extruder temperature
G1 Y-3.0 F1000.0 ; prepare to prime
G92 E0 ; reset extrusion distance
G1 X60.0 E9.0  F1000.0 ; priming
G1 X100.0 E12.5  F1000.0 ; priming
M201 X5000 Y5000 Z50 E1000 ; should NOT be removed - after TYPE:Custom
M203 X150 Y150 Z12 E50 ; should NOT be removed - after TYPE:Custom  
M204 P500 R750 T1500 ; should NOT be removed - after TYPE:Custom
M205 X5.00 Y5.00 Z0.15 E2.50 ; should NOT be removed - after TYPE:Custom
M900 K30; K45 for PET, K30 for PLA/ABS
G92 E0 ; reset extrusion distance
;START_OF_PRINT
G21 ; set units to millimeters
G90 ; use absolute coordinates
M82 ; use absolute distances for extrusion
M201 X2000 Y2000 Z20 E500 ; should NOT be removed - after TYPE:Custom
G92 E0 ; reset extrusion distance
M203 X100 Y100 Z8 E25 ; should NOT be removed - after TYPE:Custom
M204 P250 R400 T800 ; should NOT be removed - after TYPE:Custom  
M205 X2.50 Y2.50 Z0.10 E1.25 ; should NOT be removed - after TYPE:Custom
;
