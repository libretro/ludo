APP = Ludo
BUNDLENAME = $(APP)-$(OS)-$(VERSION)

ludo:
	go build

$(APP).app: ludo
	mkdir -p $(APP).app/Contents/MacOS
	mkdir -p $(APP).app/Contents/Resources/$(APP).iconset
	cp pkg/Info.plist $(APP).app/Contents/
	echo "APPL????" > $(APP).app/Contents/PkgInfo
	cp -r database $(APP).app/Contents/Resources
	cp -r assets $(APP).app/Contents/Resources
	sips -z 16 16     assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_16x16.png
	sips -z 32 32     assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_16x16@2x.png
	sips -z 32 32     assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_32x32.png
	sips -z 64 64     assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_32x32@2x.png
	sips -z 128 128   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_128x128.png
	sips -z 256 256   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_128x128@2x.png
	sips -z 256 256   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_256x256.png
	sips -z 512 512   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_256x256@2x.png
	sips -z 512 512   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_512x512.png
	cp ludo $(APP).app/Contents/MacOS
	iconutil -c icns -o $(APP).app/Contents/Resources/$(APP).icns $(APP).app/Contents/Resources/$(APP).iconset
	rm -rf $(APP).app/Contents/Resources/$(APP).iconset

empty.dmg:
	mkdir -p template
	hdiutil create -fs HFSX -layout SPUD -size 200m empty.dmg -srcfolder template -format UDRW -volname $(BUNDLENAME) -quiet
	rmdir template

dmg: empty.dmg $(APP).app
	mkdir -p wc
	hdiutil attach empty.dmg -noautoopen -quiet -mountpoint wc
	rm -rf wc/$(APP).app
	ditto -rsrc $(APP).app wc/$(APP).app
	ln -s /Applications wc/Applications
	WC_DEV=`hdiutil info | grep wc | grep "Apple_HFS" | awk '{print $$1}'` && hdiutil detach $$WC_DEV -quiet -force
	rm -f $(BUNDLENAME)-*.dmg
	hdiutil convert empty.dmg -quiet -format UDZO -imagekey zlib-level=9 -o $(BUNDLENAME).dmg

zip: ludo
	mkdir -p $(BUNDLENAME)/
	cp ludo $(BUNDLENAME)/
	cp -r database $(BUNDLENAME)/
	cp -r assets $(BUNDLENAME)/
	7z a $(BUNDLENAME).zip $(BUNDLENAME)\

tar: ludo
	mkdir -p $(BUNDLENAME)/
	cp ludo $(BUNDLENAME)/
	cp -r database $(BUNDLENAME)/
	cp -r assets $(BUNDLENAME)/
	tar -zcf $(BUNDLENAME).tar.gz $(BUNDLENAME)\

clean:
	rm -rf $(BUNDLENAME).app ludo wc empty.dmg $(BUNDLENAME).dmg $(BUNDLENAME)-*