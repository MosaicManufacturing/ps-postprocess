go install -v github.com/MosaicManufacturing/licensebot-client-go@0.1.0 || exit
echo "$PATH"
licensebot-client-go check
