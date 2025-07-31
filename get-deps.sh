#!/usr/bin/env bash

if [[ ! -d ~/.anki/vicos-sdk/dist/vic-toolchain ]]; then
  echo Getting deps...
  cd ~/.anki/vicos-sdk/dist/ 
  git clone https://github.com/kercre123/vic-toolchain
else
  echo You already have vic-toolchain installed. Exiting...
  exit
fi
