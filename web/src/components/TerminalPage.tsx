import React, { Component } from "react";
import { Grid } from "semantic-ui-react";
import Terminal from "terminal-in-react";
import { JsonRpc } from "../api/Terminal";
import { RequestConfig, RestUrl } from "../api/Config";
import axios from 'axios';

export interface TerminalProps {
  mode: number;
}

export interface TerminalState {
  mode: number;
  command?: string;
}

export class TerminalPage extends Component<TerminalProps, TerminalState> {
  constructor(props: TerminalProps) {
    super(props);

    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.execCommand = this.execCommand.bind(this);
  }

  componentDidMount(): void {
    this.setState({ mode: this.props.mode });
  }

  componentWillUnmount(): void {}

  private isBoolean(value: string): boolean {
    return ((value != null) && (value !== '') && (value.toLocaleLowerCase() === 'false' || value.toLocaleLowerCase() === 'true'));
  }

  private isNumber(value: string): boolean {
    return ((value != null) && (value !== '') && !isNaN(Number(value.toString())));
  }

  private stringToParamsArray(s: string[]): (string | number | boolean)[] {
    var objArray: (string | number | boolean)[] = [];
    s.forEach(element => {
        if (this.isNumber(element)){
            objArray.push(Number(element));
        } else if (this.isBoolean(element)) {
            if (element.toLocaleLowerCase() === 'true') {
                objArray.push(true);
            } else {
                objArray.push(false);
            }
        } else {
            objArray.push(element);
        }
    });
    return objArray;
  }

  private execCommand = async (cmd: string) => {
    var parsed: string[] = cmd.split(',');
    if (parsed.length > 0) {
        let method: string = parsed[0];
        let params: string[] = parsed.slice(1, parsed.length);
        let paramsObj: (string | number | boolean)[] = this.stringToParamsArray(params);
        let command: JsonRpc = {
            jsonrpc: "2.0",
            method: method,
            params: paramsObj,
            id: "123" // TODO: create unique id
        };
        await axios.post<object>(RestUrl + "blockchain/jsonrpc", command, RequestConfig).then(function (response) {
            console.log(JSON.stringify(response.data, null, 2));
        }).catch(function (error) {
            console.log("Execute dynamic-cli JSON RCP [Put] Error: " + error);
        });
    }
  }

  render() {
    return (
      <div>
        <Grid.Row>
          <Terminal
            watchConsoleLogging
            color="green"
            backgroundColor="black"
            barColor="black"
            style={{ fontWeight: "bold", fontSize: "1em" }}
            msg="Welcome to the WebBridge terminal."
            startState="maximised"
            hideTopBar
            commandPassThrough={
                (cmd) => {
                    this.execCommand(String(cmd));
                }
            }
          />
        </Grid.Row>
      </div>
    );
  }
}
