# Go Scrcpy Client

This package allows you to view and control android device in realtime.

Note: This gif is compressed and experience lower quality than actual.


## Contribution & Development
Already implemented all functions in scrcpy server 1.20.  
Please check scrcpy server 1.20 source code: [Link](https://github.com/Genymobile/scrcpy/tree/v1.20/server)

## Reference & Appreciation
- Core: [scrcpy](https://github.com/Genymobile/scrcpy)
- Idea: [py-android-viewer](https://github.com/razumeiko/py-android-viewer)
- CI: [index.py](https://github.com/index-py/index.py)

## env 
```bash
require ffmpeg < 5

export FFMPEG_ROOT=/usr/local/opt/ffmpeg@4
export CGO_LDFLAGS="-L$FFMPEG_ROOT/lib/ -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
export CGO_CFLAGS="-I$FFMPEG_ROOT/include"
export LDFLAGS="-L/usr/local/opt/ffmpeg@4/lib"
export CPPFLAGS="-I/usr/local/opt/ffmpeg@4/include"
export LD_LIBRARY_PATH=/usr/local/opt/ffmpeg@4/lib
export PKG_CONFIG_PATH="/usr/local/opt/ffmpeg@4/lib/pkgconfig"

```