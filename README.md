# cross-medit

Cross Platform memory analysis tool, inspired by apk-medit.

<img width="600" alt="screenshot" src="screenshots/screenshot.png">

# Usage

## Windows

```sh
medit.exe
```

## iOS

Jailbreaking of iphone is required.  
Place medit and Entitlements.plist in /usr/bin.

Connect to the iphone via ssh.

```sh
cd /usr/bin
ldid -SEntitlements.plist medit
./medit
```

## License

MIT License
