import React, { Component } from "react";
import { Button, Form, Header } from "semantic-ui-react";
import { RestartWebServer, UpdateWebServerConfig, WebServerConfig, WebServerRestartRequest } from "../api/WebServerConfig";
import { UpdateIceConfig, IceServerConfig } from "../api/IceServerConfig";
import { GetConfigSettings, ConfigurationResponse }from "../api/Config";

export interface SettingsProps {
  defaultIceUrl: string;
  defaultIceUser: string;
  defaultIcePassword: string;
  defaultBind: string;
  defaultAllow: string;
  defaultPort: number;
}

export interface SettingsState {
  config: ConfigurationResponse;
  URL?: string;
  UserName?: string;
  Credential?: string;
  BindAddress?: string;
  ListenPort?: number;
  AllowCIDR?: string;
}

export class Settings extends Component<SettingsProps, SettingsState> {
  constructor(props: SettingsProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.getConfigSettings = this.getConfigSettings.bind(this);
    this.updateIceConfig = this.updateIceConfig.bind(this);
    this.updateWebServerConfig = this.updateWebServerConfig.bind(this);
    this.restartWebServer = this.restartWebServer.bind(this);
  }

  componentDidMount(): void {
    this.getConfigSettings();
  }

  componentDidUnmount(): void {}

  private getConfigSettings = () => {
    GetConfigSettings().then((data) => {
      this.setState( { config: data} );
    });
  };

  private updateIceConfig = async () => {
    let IceServer: IceServerConfig = {
      URL: this.state.URL ? this.state.URL : this.state.config.IceServers[0].URL,
      UserName: this.state.UserName ? this.state.UserName : this.state.config.IceServers[0].UserName,
      Credential: this.state.Credential ? this.state.Credential : this.state.config.IceServers[0].Credential,
    };
    let IceServers: IceServerConfig[] = [ IceServer ]
    UpdateIceConfig(IceServers).then((data) => {
      console.log(JSON.stringify(data, null, 2));
    });
  }

  private updateWebServerConfig = async () => {
    let webserver: WebServerConfig = {
      BindAddress: this.state.BindAddress ? this.state.BindAddress : this.state.config.WebServer.BindAddress,
      ListenPort: this.state.ListenPort ? this.state.ListenPort : this.state.config.WebServer.ListenPort,
      AllowCIDR: this.state.AllowCIDR ? this.state.AllowCIDR : this.state.config.WebServer.AllowCIDR,
    };
    UpdateWebServerConfig(webserver).then((data) => {
      console.log(JSON.stringify(data, null, 2));
    });
  }

  private restartWebServer = async () => {
    let postData: WebServerRestartRequest = {
      RestartEpoch: 0,
    };
    RestartWebServer(postData).then((data) => {
      console.log(JSON.stringify(data, null, 2));
    });
  }

  render() {
    return (
      <div>
        <Header as="h3">WebBridge Configuration Settings</Header>
        <Form>
          <Header as="h4">ICE Server</Header>
          <div className="ui form">
            <Form.Group className="leftAlign field" widths="equal">
              <Form.Field>
                <Form.Input
                  label="ICE Settings URL"
                  placeholder="URL"
                  data-tooltip="Please enter the ICE server URL"
                  data-position="bottom right"
                  data-inverted
                  onChange={(e, data) => this.setState( { URL: String(data.value) })}
                  defaultValue={this.state ? this.state.config.IceServers[0].URL : ""}
                />
              </Form.Field>
              <Form.Field>
                <Form.Input
                  label="User Name"
                  placeholder="User Name"
                  data-tooltip="Please enter the ICE server user name"
                  data-position="bottom right"
                  data-inverted
                  onChange={(e, data) => this.setState( { UserName: String(data.value) })}
                  defaultValue={this.state ? this.state.config.IceServers[0].UserName : ""}
                />
              </Form.Field>
              <Form.Field>
                <Form.Input
                  label="Credential"
                  placeholder="ICE Server Credential"
                  data-tooltip="Please enter the ICE server password/credential"
                  data-position="bottom right"
                  data-inverted
                  onChange={(e, data) => this.setState( { Credential: String(data.value) })}
                  defaultValue={this.state ? this.state.config.IceServers[0].Credential: ""}
                />
              </Form.Field>
            </Form.Group>
            <Button onClick={() => this.updateIceConfig()} type="submit">Update ICE Settings</Button>
            <Header as="h4">Web Server Settings</Header>
            <Form.Group className="leftAlign field" widths="equal">
              <Form.Field>
                <Form.Input
                  label="Web Server Bind Address"
                  placeholder="Bind Address"
                  data-tooltip="Please enter the web server local bind address"
                  data-position="bottom right"
                  data-inverted
                  onChange={(e, data) => this.setState( { BindAddress: String(data.value) })}
                  defaultValue={this.state ? this.state.BindAddress ? this.state.BindAddress : this.state.config.WebServer.BindAddress : "" }
                />
              </Form.Field>
              <Form.Field>
                <Form.Input
                  label="Web Server Listen Port"
                  placeholder="Listen Port Number"
                  data-tooltip="Please enter the web server listen port number"
                  data-position="bottom right"
                  data-inverted
                  onChange={(e, data) => this.setState( { ListenPort: Number(data.value) })}
                  defaultValue={this.state ? this.state.ListenPort ? this.state.ListenPort : this.state.config.WebServer.ListenPort : "" }
                />
              </Form.Field>
              <Form.Field>
                <Form.Input
                  label="Web Server Allow CIDR"
                  placeholder="Allow Address CIDR List"
                  data-tooltip="Please enter the web server allow CIDR list"
                  data-position="bottom right"
                  data-inverted
                  onChange={(e, data) => this.setState( { AllowCIDR: String(data.value) })}
                  defaultValue={this.state ? this.state.AllowCIDR ? this.state.AllowCIDR : this.state.config.WebServer.AllowCIDR : "" }
                />
              </Form.Field>
            </Form.Group>
            <Button onClick={() => this.updateWebServerConfig()} type="submit">Update Web Server</Button>
            <Button onClick={() => this.restartWebServer()} type="submit">Restart Web Server</Button>
          </div>
        </Form>
      </div>
    );
  }
}
