.PHONY: voskopus vic-cloud vic-gateway

all: vic-cloud vic-gateway

voskopus:
	./build-voskopus.sh

go_deps:
	echo `/usr/local/go/bin/go version` && cd $(PWD) && /usr/local/go/bin/go mod download

vic-cloud: voskopus go_deps
	CGO_ENABLED=1 GOARM=7 GOARCH=arm \
	CC=/home/kerigan/.anki/vicos-sdk/dist/4.0.0-r05/prebuilt/bin/arm-oe-linux-gnueabi-clang \
	CXX=/home/kerigan/.anki/vicos-sdk/dist/4.0.0-r05/prebuilt/bin/arm-oe-linux-gnueabi-clang++ \
	PKG_CONFIG_PATH="$(PWD)/voskopus/built/armel/lib/pkgconfig" \
	CGO_CFLAGS="-Wno-implicit-function-declaration -I$(PWD)/voskopus/built/armel/include -I$(PWD)/voskopus/built/armel/include/opus" \
	CGO_CXXFLAGS="-stdlib=libc++ -std=c++11" \
	CGO_LDFLAGS="-L$(PWD)/voskopus/built/armel/lib -L$(PWD)/armlibs/lib/arm-linux-gnueabi/android" \
	/usr/local/go/bin/go build \
	-tags nolibopusfile,vicos \
	-ldflags '-w -s -linkmode internal -extldflags "-static" -r /anki/lib' \
	-o build/vic-cloud \
	cloud/main.go

	upx build/vic-cloud


vic-gateway: go_deps
	CGO_ENABLED=1 GOARM=7 GOARCH=arm CC=/home/kerigan/.anki/vicos-sdk/dist/4.0.0-r05/prebuilt/bin/arm-oe-linux-gnueabi-clang CXX=/home/kerigan/.anki/vicos-sdk/dist/4.0.0-r05/prebuilt/bin/arm-oe-linux-gnueabi-clang++ PKG_CONFIG_PATH="$(PWD)/voskopus/lib/pkgconfig" CGO_CFLAGS="-I$(PWD)/voskopus/include -I$(PWD)/voskopus/include/opus -I$(PWD)/voskopus/include/ogg" CGO_CXXFLAGS="-stdlib=libc++ -std=c++11" CGO_LDFLAGS="-L$(PWD)/voskopus/lib -L$(PWD)/armlibs/lib/arm-linux-gnueabi/android" /usr/local/go/bin/go build -tags nolibopusfile,vicos -ldflags '-w -s -linkmode internal -extldflags "-static" -r /anki/lib' -o build/vic-gateway gateway/*.go

	upx build/vic-gateway

