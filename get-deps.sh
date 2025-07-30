#!/usr/bin/env bash

if [[ ! -f ~/.anki/vicos-sdk/dist/5.2.1-r06 ]]; then
  echo Getting deps...
  mkdir ~/.anki/vicos-sdk/dist/5.2.1-r06
  cd ~/.anki/vicos-sdk/dist/5.2.1-r06
  wget https://froggitti.net/5.2.1-r06.tar.gz
  gunzip 5.2.1-r06.tar.gz
  tar -xvf 5.2.1-r06.tar
else
  echo You already have the 5.2.1-r06 toolchain installed. Exiting...
  exit
fi
