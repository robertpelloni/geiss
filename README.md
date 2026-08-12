[![Build for Windows using MSVC](https://github.com/geissomatik/geiss/actions/workflows/windows-msvc.yml/badge.svg)](https://github.com/geissomatik/geiss/actions/workflows/windows-msvc.yml)


# Geiss Screensaver and Winamp plug-in

This repository holds the latest code of the Geiss Screensaver and Winamp plug-in;
for general information about Geiss please see:

- [At a glance: Geiss](https://www.geisswerks.com/about_geiss.html)
- [Official Homepage of the Geiss Screensaver and Winamp plug-in](https://www.geisswerks.com/geiss/)
- [Screenshots from Geiss](https://www.geisswerks.com/geiss/shots.html)
- [How 'Geiss' Worked](https://www.geisswerks.com/geiss/secrets.html)
- [Geiss @ WinampHeritage.com](https://winampheritage.com/visualizations/geiss)
- [Geisswerks — home page of Ryan M. Geiss](https://www.geisswerks.com/)

For a YouTube video showing Geiss in action, please click:

[![Video "Windows 95 + Winamp + Geiss = Aesthetic"](https://img.youtube.com/vi/MfFcRsUZCxw/hqdefault.jpg)](https://youtu.be/MfFcRsUZCxw?t=15)


# Compiling locally

In order to compile the Geiss Winamp plug-in and/or the Geiss Screensaver,
you need:

- A Git clone or snapshot of this repository
- CMake >=3.15
- Microsoft Visual Studio >=9/2008 (due to the current use of MSVC style assembly)

Then, the code can be built using:

```console
# cmake -G "Visual Studio 17 2022" -A Win32 -DCMAKE_BUILD_TYPE=RelWithDebInfo .
# cmake --build .
```

## What is Geiss?

GEISS is sound-reactive eye candy for your computer. It comes in two forms: one version is a screensaver for Windows 95/98/ME/2000/XP/Vista and the other, a visualization plugin for Winamp. Geiss uses hand-tuned x86 assembly language to reach blazing speeds, generating smooth and graceful realtime graphics... and it can hear whatever your computer hears and synchronize the graphics it generates to the sound you're hearing.

## License

3-Clause BSD License
Copyright (c) 1998-2022 Ryan Geiss (@geissomatik)
