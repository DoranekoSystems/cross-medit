# cross-medit

Cross Platform memory analysis tool, inspired by apk-medit.

<img width="600" alt="screenshot" src="https://github.com/DoranekoSystems/cross-medit/assets/96031346/65727e79-c3cd-41d6-9083-8fa4f270bdf8">


# Usage

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
