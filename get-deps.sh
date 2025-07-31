#!/usr/bin/env bash

if [[ ! -d ~/.anki/vicos-sdk/dist/1.1.0-r04 ]]; then
  echo Getting deps...
  mkdir ~/.anki/vicos-sdk/dist/1.1.0-r04
  cd ~/.anki/vicos-sdk/dist/1.1.0-r04
  wget https://github.com/kercre123/anki-deps/releases/download/1.1.0-r04/vicos-sdk-1.1.0-r04-x86_64-ubuntu-16.04.tar.gz
  gunzip 1.1.0-r04.tar.gz
  tar -xvf 1.1.0-r04.tar
else
  echo You already have the 1.1.0-r04 toolchain installed. Exiting...
  exit
fi
