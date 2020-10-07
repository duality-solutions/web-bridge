import React, { Component } from "react";
import { Button, Form, Header } from "semantic-ui-react";
import { ConfigurationWebResponse, RestartResponse, WebServerConfig, WebServerRestartRequest } from "../api/WebServerConfig";
import { ConfigurationIceResponse, IceServerConfig } from "../api/IceServerConfig";
import { ConfigurationResponse, RequestConfig, RestUrl }from "../api/Config";
import axios from 'axios';

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

// TODO: read settings file to get rest web server URL and port instead of using RestUrl constant variable 
export class Settings extends Component<SettingsProps, SettingsState> {
  constructor(props: SettingsProps) {
    super(props);
    let configResponse: ConfigurationResponse = { result: {
      IceServers: [
        {
          URL: props.defaultIceUrl,
          UserName: props.defaultIceUser,
          Credential: props.defaultIcePassword
        }
      ],
      WebServer: {
        BindAddress: props.defaultBind,
        ListenPort: props.defaultPort,
        AllowCIDR: props.defaultAllow
      }
    }};
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
    this.getConfigSettings = this.getConfigSettings.bind(this);
    this.updateIceConfig = this.updateIceConfig.bind(this);
    this.updateWebServerConfig = this.updateWebServerConfig.bind(this);
    this.restartWebServer = this.restartWebServer.bind(this);
    // set state
    this.setState({ config: configResponse} );
  }

  componentDidMount(): void {
    this.getConfigSettings();
  }

  componentWillUnmount(): void {}
  
  private getConfigSettings = async () => {
    var self = this;
    await axios.get<ConfigurationResponse>(RestUrl + "config", RequestConfig).then(function (response) {
      self.setState({ config: response.data });
    }).catch(function (error) {
      console.log("Get Configuration Settings [Get] Error: " + error);
    });
  }

  private updateIceConfig = async () => {
    let IceServer: IceServerConfig = {
      URL: this.state.URL ? this.state.URL : this.state.config.result.IceServers[0].URL,
      UserName: this.state.UserName ? this.state.UserName : this.state.config.result.IceServers[0].UserName,
      Credential: this.state.Credential ? this.state.Credential : this.state.config.result.IceServers[0].Credential,
    };
    let IceServers: IceServerConfig[] = [ IceServer ]
    await axios.post<ConfigurationIceResponse>(RestUrl + "config/ice", IceServers, RequestConfig).then(function (response) {
      console.log(response.data);
    }).catch(function (error) {
      console.log("Update ICE Server Config [Post] Error: " + error);
    });
  }

  private updateWebServerConfig = async () => {
    let webserver: WebServerConfig = {
      BindAddress: this.state.BindAddress ? this.state.BindAddress : this.state.config.result.WebServer.BindAddress,
      ListenPort: this.state.ListenPort ? this.state.ListenPort : this.state.config.result.WebServer.ListenPort,
      AllowCIDR: this.state.AllowCIDR ? this.state.AllowCIDR : this.state.config.result.WebServer.AllowCIDR,
    };
    await axios.post<ConfigurationWebResponse>(RestUrl + "config/web", webserver, RequestConfig).then(function (response) {
      console.log(response.data);
    }).catch(function (error) {
      console.log("Update Web Server Config [Post] Error: " + error);
    });
  }

  private restartWebServer = async () => {
    let postData: WebServerRestartRequest = {
      restart_epoch: 0,
    };
    await axios.put<RestartResponse>(RestUrl + "config/web/restart", postData, RequestConfig).then(function (response) {
      console.log(response.data);
    }).catch(function (error) {
      console.log("Web Server Restart [Put] Error: " + error);
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
                  defaultValue={this.state ? this.state.config.result.IceServers[0].URL : ""}
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
                  defaultValue={this.state ? this.state.config.result.IceServers[0].UserName : ""}
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
                  defaultValue={this.state ? this.state.config.result.IceServers[0].Credential: ""}
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
                  defaultValue={this.state ? this.state.BindAddress ? this.state.BindAddress : this.state.config.result.WebServer.BindAddress : "" }
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
                  defaultValue={this.state ? this.state.ListenPort ? this.state.ListenPort : this.state.config.result.WebServer.ListenPort : "" }
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
                  defaultValue={this.state ? this.state.AllowCIDR ? this.state.AllowCIDR : this.state.config.result.WebServer.AllowCIDR : "" }
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
