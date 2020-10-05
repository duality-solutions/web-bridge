# WebBridge Admin Console

## Build Admin Console

### NVM Ubuntu install

Install the latest release of Node Virtual Machine (NVM) [Linux](https://github.com/nvm-sh/nvm/releases)

```bash
curl -sL https://raw.githubusercontent.com/creationix/nvm/v0.36.0/install.sh -o install_nvm.sh && chmod +x ./install_nvm.sh && ./install_nvm.sh
```

Use NodeJS v12.18.4

```bash
nvm install 12.18.4 && nvm use 12.18.4 && npm install && npm install -g yarn && yarn build
```

### NVM Windows install

Install [Chocolatey](https://chocolatey.org)

```shell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
```

Install the latest release of Node Virtual Machine (NVM) [Windows](https://github.com/coreybutler/nvm-windows/releases)

Use NodeJS v12.18.4

```shell
nvm install 12.18.4;nvm use 12.18.4;
npm install;yarn build
```

### Build website

```shell
yarn build
```

## Access admin console locally

```http
http://localhost:35350/admin
```
