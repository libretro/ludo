BUNDLENAME=Play\ Them\ All

go-playthemall:
	go build

Play\ Them\ All.app: go-playthemall
	mkdir -p ${BUNDLENAME}.app/Contents/MacOS
	mkdir -p ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset
	cp Info.plist ${BUNDLENAME}.app/Contents/
	echo "APPL????" > ${BUNDLENAME}.app/Contents/PkgInfo
	cp -r assets ${BUNDLENAME}.app/Contents/Resources
	sips -z 16 16     assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_16x16.png
	sips -z 32 32     assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_16x16@2x.png
	sips -z 32 32     assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_32x32.png
	sips -z 64 64     assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_32x32@2x.png
	sips -z 128 128   assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_128x128.png
	sips -z 256 256   assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_128x128@2x.png
	sips -z 256 256   assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_256x256.png
	sips -z 512 512   assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_256x256@2x.png
	sips -z 512 512   assets/icon.png --out ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_512x512.png
	cp go-playthemall ${BUNDLENAME}.app/Contents/MacOS
	iconutil -c icns -o ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.icns ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset
	rm -rf ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset

clean:
	rm -rf ${BUNDLENAME}.app go-playthemall