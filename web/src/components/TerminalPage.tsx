import React, { Component } from "react";
import { Grid } from "semantic-ui-react";
import Terminal from "terminal-in-react";
import { ExecCommand } from "../api/Terminal";

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
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
  }

  componentDidMount(): void {
    this.setState({ mode: this.props.mode });
  }

  componentDidUnmount(): void {}

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
                  ExecCommand(String(cmd));
                }
            }
          />
        </Grid.Row>
      </div>
    );
  }
}
