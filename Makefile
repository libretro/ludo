BUNDLENAME=Play\ Them\ All

go-playthemall:
	go build

Play\ Them\ All.app: go-playthemall
	mkdir -p ${BUNDLENAME}.app/Contents/MacOS
	mkdir -p ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset
	cp Info.plist ${BUNDLENAME}.app/Contents/
	echo "APPL????" > ${BUNDLENAME}.app/Contents/PkgInfo
	cp -r assets ${BUNDLENAME}.app/Contents/Resources
	cp assets/icon.png ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset/icon_256x256.png
	cp go-playthemall ${BUNDLENAME}.app/Contents/MacOS
	iconutil -c icns -o ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.icns ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset
	rm -rf ${BUNDLENAME}.app/Contents/Resources/${BUNDLENAME}.iconset

clean:
	rm -rf ${BUNDLENAME}.app go-playthemall