# DeckBuilder

## How to build
Clone repository:
```
git clone https://github.com/HardDie/DeckBuilder --recursive
```

Check that all necessary packages are installed
```
./deployment/check_binary.sh
```

Build web
```
make web-build
```

Build binary
```
cd deployment
./build_linux.sh
./build_darwin.sh
./build_windows.sh
```

The resulting files can be found in the deployment/out folder
