APP ?= Ludo
ARCH ?= x86_64
VERSION ?= dev
BUNDLENAME = $(APP)-$(OS)-$(ARCH)-$(VERSION)

CORES = atari800 bluemsx swanstation fbneo fceumm gambatte genesis_plus_gx handy lutro mednafen_ngp mednafen_pce mednafen_pce_fast mednafen_pcfx mednafen_psx mednafen_saturn mednafen_supergrafx mednafen_vb mednafen_wswan mgba melonds np2kai o2em pcsx_rearmed picodrive pokemini prosystem snes9x stella2014 vecx virtualjaguar

ifeq ($(ARCH), arm)
	CORES := $(filter-out swanstation,$(CORES))
	CORES := $(filter-out mednafen_pcfx,$(CORES))
	CORES := $(filter-out mednafen_saturn,$(CORES))
	CORES := $(filter-out melonds,$(CORES))
endif

ifeq ($(OS), Windows)
	CORES += mupen64plus_next
endif

DYLIBS = $(addprefix cores/, $(addsuffix _libretro.dylib,$(CORES)))
DLLS = $(addprefix cores/, $(addsuffix _libretro.dll,$(CORES)))
SOBJS = $(addprefix cores/, $(addsuffix _libretro.so,$(CORES)))

ifeq ($(OS), OSX)
	BUILDBOTURL=http://buildbot.libretro.com/nightly/apple/osx/$(ARCH)/latest
endif
ifeq ($(OS), Linux)
	ifeq ($(ARCH), arm)
		BUILDBOTURL=http://buildbot.libretro.com/nightly/linux/armv7-neon-hf/latest
	else
		BUILDBOTURL=http://buildbot.libretro.com/nightly/linux/$(ARCH)/latest
	endif
endif
ifeq ($(OS), Windows)
	BUILDBOTURL=http://buildbot.libretro.com/nightly/windows/$(ARCH)/latest
endif

ludo:
	go build

ludo.exe:
	go build -ldflags '-H=windowsgui'

cores/%_libretro.dylib cores/%_libretro.dll cores/%_libretro.so:
	mkdir -p cores
	@if [ $(basename $(@F)) = swanstation_libretro ]; then\
		wget -c https://github.com/kivutar/swanstation/releases/download/v2021-08-01/swanstation-$(OS)-$(ARCH).zip -O $@.zip;\
	else\
		wget -c $(BUILDBOTURL)/$(@F).zip -O $@.zip;\
	fi
	unzip $@.zip -d cores
	rm $@.zip

$(APP).app: ludo $(DYLIBS)
	mkdir -p $(APP).app/Contents/MacOS
	mkdir -p $(APP).app/Contents/Resources/$(APP).iconset
	cp Info.plist $(APP).app/Contents/
	sed -i.bak 's/0.1.0/$(VERSION)/' $(APP).app/Contents/Info.plist
	rm $(APP).app/Contents/Info.plist.bak
	echo "APPL????" > $(APP).app/Contents/PkgInfo
	cp -r database $(APP).app/Contents/Resources
	cp -r assets $(APP).app/Contents/Resources
	cp -r cores $(APP).app/Contents/Resources
	codesign --force --options runtime --verbose --timestamp --sign "7069CC8A4AE9AFF0493CC539BBA4FA345F0A668B" \
		--entitlements entitlements.xml $(APP).app/Contents/Resources/cores/*.dylib
	rm -rf $(APP).app/Contents/Resources/database/.git
	rm -rf $(APP).app/Contents/Resources/assets/.git
	sips -z 16 16   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_16x16.png
	sips -z 32 32   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_16x16@2x.png
	sips -z 32 32   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_32x32.png
	sips -z 64 64   assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_32x32@2x.png
	sips -z 128 128 assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_128x128.png
	sips -z 256 256 assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_128x128@2x.png
	sips -z 256 256 assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_256x256.png
	sips -z 512 512 assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_256x256@2x.png
	sips -z 512 512 assets/icon.png --out $(APP).app/Contents/Resources/$(APP).iconset/icon_512x512.png
	cp ludo $(APP).app/Contents/MacOS
	codesign --force --options runtime --verbose --timestamp --sign "7069CC8A4AE9AFF0493CC539BBA4FA345F0A668B" \
		--entitlements entitlements.xml $(APP).app/Contents/MacOS/ludo
	iconutil -c icns -o $(APP).app/Contents/Resources/$(APP).icns $(APP).app/Contents/Resources/$(APP).iconset
	rm -rf $(APP).app/Contents/Resources/$(APP).iconset
	codesign --force --options runtime --verbose --timestamp --sign "7069CC8A4AE9AFF0493CC539BBA4FA345F0A668B" \
		--entitlements entitlements.xml $(APP).app

empty.dmg:
	mkdir -p template
	hdiutil create -fs HFSX -layout SPUD -size 200m empty.dmg -srcfolder template -format UDRW -volname $(BUNDLENAME) -quiet
	rmdir template

# For OSX
dmg: empty.dmg $(APP).app
	mkdir -p wc
	hdiutil attach empty.dmg -noautoopen -quiet -mountpoint wc
	rm -rf wc/$(APP).app
	ditto -rsrc $(APP).app wc/$(APP).app
	ln -s /Applications wc/Applications
	WC_DEV=`hdiutil info | grep wc | grep "Apple_HFS" | awk '{print $$1}'` && hdiutil detach $$WC_DEV -quiet -force
	rm -f $(BUNDLENAME)-*.dmg
	hdiutil convert empty.dmg -quiet -format UDZO -imagekey zlib-level=9 -o $(BUNDLENAME).dmg
	codesign --force --options runtime --verbose --timestamp --sign "7069CC8A4AE9AFF0493CC539BBA4FA345F0A668B" \
		--entitlements entitlements.xml $(BUNDLENAME).dmg

# For Windows
zip: ludo.exe $(DLLS)
	mkdir -p $(BUNDLENAME)/
	./rcedit-x64 ludo.exe --set-icon assets/icon.ico
	cp ludo.exe $(BUNDLENAME)/
	cp -r database $(BUNDLENAME)/
	cp -r assets $(BUNDLENAME)/
	cp -r cores $(BUNDLENAME)/
	7z a $(BUNDLENAME).zip $(BUNDLENAME)\

# For Linux
tar: ludo $(SOBJS)
	mkdir -p $(BUNDLENAME)/
	cp ludo $(BUNDLENAME)/
	cp -r database $(BUNDLENAME)/
	cp -r assets $(BUNDLENAME)/
	cp -r cores $(BUNDLENAME)/
	tar -zcf $(BUNDLENAME).tar.gz $(BUNDLENAME)\

# For Debian

DEB_ARCH = amd64
ifeq ($(ARCH), arm)
	DEB_ARCH = armhf
endif
DEB_ROOT = ludo_$(VERSION)-1_$(DEB_ARCH)

deb: ludo $(SOBJS)
	mkdir -p $(DEB_ROOT)/DEBIAN
	mkdir -p $(DEB_ROOT)/etc
	mkdir -p $(DEB_ROOT)/usr/bin
	mkdir -p $(DEB_ROOT)/usr/lib/ludo
	mkdir -p $(DEB_ROOT)/usr/share/ludo
	mkdir -p $(DEB_ROOT)/usr/share/applications
	mkdir -p $(DEB_ROOT)/usr/share/icons/hicolor/1024x1024/apps/
	touch $(DEB_ROOT)/etc/ludo.toml
	echo "cores_dir = \"/usr/lib/ludo\"" >> $(DEB_ROOT)/etc/ludo.toml
	echo "assets_dir = \"/usr/share/ludo/assets\"" >> $(DEB_ROOT)/etc/ludo.toml
	echo "database_dir = \"/usr/share/ludo/databbase\"" >> $(DEB_ROOT)/etc/ludo.toml
	cp ludo $(DEB_ROOT)/usr/bin
	cp cores/* $(DEB_ROOT)/usr/lib/ludo
	cp -r assets $(DEB_ROOT)/usr/share/ludo
	cp -r database $(DEB_ROOT)/usr/share/ludo
	cp assets/icon.png $(DEB_ROOT)/usr/share/icons/hicolor/1024x1024/apps/ludo.png
	cp ludo.desktop $(DEB_ROOT)/usr/share/applications
	cp control $(DEB_ROOT)/DEBIAN
	sed -i.bak 's/VERSION/$(VERSION)/' $(DEB_ROOT)/DEBIAN/control
	sed -i.bak 's/ARCH/$(DEB_ARCH)/' $(DEB_ROOT)/DEBIAN/control
	rm $(DEB_ROOT)/DEBIAN/control.bak
	dpkg-deb --build $(DEB_ROOT)

clean:
	rm -rf Ludo.app ludo wc *.dmg *.deb $(BUNDLENAME)-* cores/
