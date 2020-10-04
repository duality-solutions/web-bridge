import React, { Component } from "react";
import { Button, Form, Header } from "semantic-ui-react";
import { RestUrl, RequestConfig } from "./WebServiceConfig";
import axios from 'axios';

export interface SettingsProps {
  defaultIceUrl: string;
  defaultIceUser: string;
  defaultIcePassword: string;
  defaultBind: string;
  defaultAllow: string;
}

export interface SettingsState {}

export class Settings extends Component<SettingsProps, SettingsState> {
  private defaultIceUrl: string;
  private defaultIceUser: string;
  private defaultIcePassword: string;
  private defaultBind: string;
  private defaultAllow: string;
  constructor(props: SettingsProps) {
    super(props);
    this.defaultIceUrl = props.defaultIceUrl;
    this.defaultIceUser = props.defaultIceUser;
    this.defaultIcePassword = props.defaultIcePassword;
    this.defaultAllow = props.defaultAllow;
    this.defaultBind = props.defaultBind;
  }

  componentDidMount(): void {}

  componentWillUnmount(): void {}
  
  private getConfigSettings = async () => {
    await axios.get(RestUrl + "config", RequestConfig).then(function (response) {
      let config: string = response.data;
      console.log("Config: " + config);
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
                <input placeholder="URL" value={this.defaultIceUrl} />
              </Form.Field>
              <Form.Field>
                <label>UserName</label>
                <input placeholder="UserName" value={this.defaultIceUser} />
              </Form.Field>
              <Form.Field>
                <label>Password</label>
                <input placeholder="Password" value={this.defaultIcePassword} />
              </Form.Field>
            </Form.Group>
            <Button onClick={() => this.getConfigSettings()} type="submit">Update ICE</Button>
            <Header as="h4">Web Server Settings</Header>
            <Form.Group className="leftAlign field" widths="equal">
              <Form.Field>
                <label>Web Server Bind Address</label>
                <input placeholder="Bind Address" value={this.defaultBind} />
              </Form.Field>
              <Form.Field>
                <label>Web Server Allow CIDR</label>
                <input placeholder="Allow CIDR" value={this.defaultAllow} />
              </Form.Field>
            </Form.Group>
            <Button type="submit">Update Web Server</Button>
          </div>
        </Form>
      </div>
    );
  }
}
