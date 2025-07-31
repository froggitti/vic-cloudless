#!/bin/bash

mkdir -p voskopus
cd voskopus

ORIGPATH="$(pwd)"

#ARMT="$(pwd)/vic-toolchain/arm-linux-gnueabi/bin/arm-linux-gnueabi-"
ARMT="$HOME/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-"

#if [[ ! -f vic-toolchain ]]; then
#    git clone https://github.com/kercre123/vic-toolchain --depth=1
#fi

set -e

function prepareVOSKbuild_ARMARM64() {
    cd $ORIGPATH
    ARCH=$1
    if [[ ${ARCH} == "amd64" ]]; then
        echo "prepareVOSKbuild_ARMARM64: this function is for armhf and arm64 only."
        exit 1
    fi
    mkdir -p build/${ARCH}
    KALDIROOT="$(pwd)/build/${ARCH}/kaldi"
    BPREFIX="$(pwd)/built/${ARCH}"
    cd build/${ARCH}
    expToolchain ${ARCH}
    if [[ ! -f ${KALDIROOT}/KALDIBUILT ]]; then
        git clone -b vosk --single-branch https://github.com/alphacep/kaldi
        cd kaldi/tools
        git clone -b v0.3.20 --single-branch https://github.com/xianyi/OpenBLAS
        git clone -b v3.2.1  --single-branch https://github.com/alphacep/clapack
	    sed -i 's/-mfloat-abi=hard -mfpu=neon/-mfloat-abi=softfp -mfpu=neon-vfpv4/g' ${KALDIROOT}/src/makefiles/*.mk
        echo ${OPENBLAS_ARGS}
        make -C OpenBLAS ONLY_CBLAS=1 TARGET=ARMV7 ${OPENBLAS_ARGS} HOSTCC="gcc -Wno-error" USE_LOCKING=1 ARM_SOFTFP_ABI=1 USE_THREAD=0 NUM_THREADS=2 -j 12
        make -C OpenBLAS ${OPENBLAS_ARGS} HOSTCC="gcc -Wno-error" USE_LOCKING=1 USE_THREAD=0 PREFIX=$(pwd)/OpenBLAS/install install
        rm -rf clapack/BUILD
        mkdir -p clapack/BUILD && cd clapack/BUILD
        cmake -DCMAKE_C_FLAGS="$ARCHFLAGS" -DCMAKE_C_COMPILER_TARGET=$PODHOST \
            -DCMAKE_C_COMPILER=$CC -DCMAKE_SYSTEM_NAME=Generic -DCMAKE_AR=$AR \
            -DCMAKE_TRY_COMPILE_TARGET_TYPE=STATIC_LIBRARY -DCMAKE_POLICY_VERSION_MINIMUM=3.5 \
            -DCMAKE_CROSSCOMPILING=True ..
        make HOSTCC=gcc -j 12 -C F2CLIBS/libf2c
        make  HOSTCC=gcc -j 12 -C BLAS/SRC
        make HOSTCC=gcc  -j 12 -C SRC
        find . -name "*.a" | xargs cp -t ../../OpenBLAS/install/lib
        cd ${KALDIROOT}/tools
        git clone --single-branch https://github.com/alphacep/openfst openfst
	#make openfst
        cd openfst
        autoreconf -i
       CFLAGS="-g -O3" ./configure --prefix=${KALDIROOT}/tools/openfst --enable-static --enable-shared --enable-far --enable-ngram-fsts --enable-lookahead-fsts --with-pic --disable-bin --host=${CROSS_TRIPLE} --build=x86-linux-gnu
       make -j 12 && make install
        cd ${KALDIROOT}/src
        sed -i "s:TARGET_ARCH=\"\`uname -m\`\":TARGET_ARCH=$(echo $CROSS_TRIPLE|cut -d - -f 1):g" configure
        sed -i "s: -O1 : -O3 :g" makefiles/linux_openblas_arm.mk
        ./configure --mathlib=OPENBLAS_CLAPACK --shared --use-cuda=no
        make -j 12 online2 lm rnnlm
        find ${KALDIROOT} -name "*.o" -exec rm {} \;
        touch ${KALDIROOT}/KALDIBUILT
    else
        echo "VOSK dependencies already built for $ARCH"
    fi
    cd $ORIGPATH
}

function expToolchain() {
    export CC="${ARMT}gcc -Wno-error -Wno-implicit-function-declaration"
    export CXX="${ARMT}g++ -Wno-error -Wno-implicit-function-declaration"
    export CPP="${ARMT}cpp"
#    export CFLAGS="-Wno-error"
#    export CXXFLAGS="-Wno-error"
    export LD=${ARMT}ld
    export AR=${ARMT}ar
    #export FC=${ARMT}gfortran
    export RANLIB=${ARMT}ranlib
    export AS=${ARMT}as
    export PODHOST=arm-linux-gnueabi 
    export CROSS_TRIPLE=${PODHOST}
    export CROSS_COMPILE=${ARMT}
    export GOARCH=arm
    export GOARM=7
    export GOOS=linux
    export ARCHFLAGS="-mfloat-abi=softfp -mfpu=neon-vfpv4"
}

function doVOSKbuild() {
    ARCH=$1
    cd $ORIGPATH
    KALDIROOT="$(pwd)/build/${ARCH}/kaldi"
    BPREFIX="$(pwd)/built/${ARCH}"
    if [[ ! -f ${BPREFIX}/lib/libvosk.so ]]; then
        cd build/${ARCH}
        expToolchain $ARCH
        if [[ ! -d vosk-api ]]; then
            git clone https://github.com/alphacep/vosk-api
            cd vosk-api
            git checkout eabd80a848de53e87e5943937146025d42ae570d
            cd ..
        fi
        cd vosk-api/src
        KALDI_ROOT=$KALDIROOT make EXTRA_LDFLAGS=" -lpthread -Wl,-Bstatic -lc++ -Wl,-Bdynamic -lpthread -ldl -lm" -j8
	cd "${ORIGPATH}/build/${ARCH}"
        mkdir -p "${BPREFIX}/lib"
        mkdir -p "${BPREFIX}/include"
        cp vosk-api/src/libvosk.so "${BPREFIX}/lib/"
        cp vosk-api/src/vosk_api.h "${BPREFIX}/include/"
        mkdir -p "${ORIGPATH}/../build"
        cp "${BPREFIX}/lib/libvosk.so" "${ORIGPATH}/../build/"
    else
        echo "VOSK already built for $ARCH"
    fi
    cd $ORIGPATH
}

function buildOPUS() {
    ARCH=$1
    cd $ORIGPATH
    BPREFIX="$(pwd)/built/${ARCH}"
    expToolchain $ARCH

    if [[ ! -f built/${ARCH}/lib/libopus.so.0.10.1 ]]; then
        cd build/${ARCH}
        rm -rf opus
        git clone https://github.com/xiph/opus
        cd opus
        git checkout 08bcc6e46227fca01aa3de3f3512f8b692d8d36b
        ./autogen.sh
        ./configure --host=${PODHOST} --prefix=$BPREFIX
        make -j8
        make install
        cd $ORIGPATH
        touch built/${ARCH}/opus_built
        cp -r built/armel/lib/libopus.so.0.10.1 ../build/libopus.so.0
    else
        echo "OPUS already built for $ARCH"
    fi
}


arch=armel
if [[ ! -f "${ORIGPATH}/built/$arch/lib/libvosk.so" ]]; then
    echo "Compiling VOSK dependencies for $arch"
    prepareVOSKbuild_ARMARM64 "$arch"
    doVOSKbuild "$arch"
fi
buildOPUS "$arch"

cd "${ORIGPATH}"
cp -r built/armel/lib/libopus.so.0.10.1 ../build/libopus.so.0
cp -r built/armel/lib/libvosk.so ../build/

echo "Dependencies complete for $arch."
