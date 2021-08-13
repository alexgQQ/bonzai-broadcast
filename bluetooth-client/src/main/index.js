import path from 'path';
import {app, BrowserWindow} from 'electron';
import DownloadManager from 'electron-download-manager';

const entryUrl = process.env.NODE_ENV === 'development'
  ? 'http://localhost:8080/index.html'
  : `file://${path.join(__dirname, 'index.html')}`;

DownloadManager.register();
let window = null;

app.commandLine.appendSwitch('enable-web-bluetooth', true);

app.on('ready', () => {
  window = new BrowserWindow({
    width: 650, height: 1000, resizable: false});
  // window.webContents.openDevTools()
  window.loadURL(entryUrl);
  window.on('closed', () => window = null);
});

app.on('window-all-closed', () => {
  if(process.platform !== 'darwin') {
    app.quit();
  }
});