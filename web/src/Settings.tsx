import React, { Component } from "react";
import { Button, Form, Header } from "semantic-ui-react";
import { RestUrl, RequestConfig, ConfigurationResponse, ConfigurationIceResponse } from "./WebServiceConfig";
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
  config: ConfigurationResponse
}

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
      console.log("Get Settings Post Error: " + error);
    });
  }

  private saveIceSettings = async () => {
    //var self = this;
    //todo: get current ice server values.
    await axios.post<ConfigurationIceResponse>(RestUrl + "config/ice", this.state.config.result.IceServers[0], RequestConfig).then(function (response) {
      console.log(response.data);
    }).catch(function (error) {
      console.log("Get Settings Post Error: " + error);
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
                <label>ICE Settings URL</label>
                <input placeholder="URL" value={this.state && this.state.config.result.IceServers[0] ? this.state.config.result.IceServers[0].URL : this.props.defaultIceUrl} />
              </Form.Field>
              <Form.Field>
                <label>UserName</label>
                <input placeholder="UserName" value={this.state && this.state.config.result.IceServers[0] ? this.state.config.result.IceServers[0].UserName : this.props.defaultIceUser } />
              </Form.Field>
              <Form.Field>
                <label>Password</label>
                <input placeholder="Password" value={this.state && this.state.config.result.IceServers[0] ? this.state.config.result.IceServers[0].Credential : this.props.defaultIcePassword } />
              </Form.Field>
            </Form.Group>
            <Button onClick={() => this.saveIceSettings()} type="submit">Update ICE</Button>
            <Header as="h4">Web Server Settings</Header>
            <Form.Group className="leftAlign field" widths="equal">
              <Form.Field>
                <label>Web Server Bind Address</label>
                <input placeholder="Bind Address" value={this.state && this.state.config.result.WebServer.BindAddress } />
              </Form.Field>
              <Form.Field>
                <label>Web Server Port</label>
                <input placeholder="Allow CIDR" value={this.state && this.state.config.result.WebServer.ListenPort } />
              </Form.Field>
              <Form.Field>
                <label>Web Server Allow CIDR</label>
                <input placeholder="Allow CIDR" value={this.state && this.state.config.result.WebServer.AllowCIDR } />
              </Form.Field>
            </Form.Group>
            <Button type="submit">Update Web Server</Button>
          </div>
        </Form>
      </div>
    );
  }
}
