<h1 align="center"><img src="assets/icon.svg" width="100" /> <br />Genshin Real-Time Notes</h1>

Add your real-time notes to your system tray!

<p align="center">
    <img src="./assets/genshin.png" />&nbsp;
    <img src="./assets/hsr.png" />
</p>
<p align="center">
    <img width=300 src="./assets/both.png" />
</p>

# 🛠️ Installing (pre-built binaries)
- Download the [latest .zip release](https://github.com/Beastwick18/genshin-real-time-notes/releases/latest) from the releases tab
- Ensure you have [WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2?form=MA13LH#download) installed
  - Select "Evergreen Standalone Installer"
  - WebView2 comes pre-installed on Windows 11, so you may not have to install it.
- Extract this to wherever you would like it to be installed
- Login to the Hoyolab popup
- Fill out your UID, and set the refresh interval
- Click "done," and everything should be working!

# 🍪 Logging in to Hoyolab
1. Wait for the web page to load, then login with your email and password.
2. Wait for the Battle Chronicle page to load, and copy your UID
3. Paste your UID into the input box at the bottom
4. *(optional)* Change the refresh interval to match how often you would like your data to refresh.
5. Click "Done"

# 🏃 Run on startup
- Create a shortcut to the .exe
- Press `Win + R` and type in `shell:startup` and hit Enter.
- Copy the shortcut to this location

# 🚧 Building from source
## Windows
- Clone the repo:
```
git clone https://github.com/Beastwick18/genshin-real-time-notes
cd genshin-real-time-notes
```
- Run the following command:
```
make
```
- Which will generate `resin.exe` and `stamina.exe` for Genshin and Honkai: Star-Rail respectively.
