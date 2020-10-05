import React, { Component } from "react";
import { Grid } from "semantic-ui-react";
import Terminal from "terminal-in-react";

export interface TerminalProps {
  mode: number;
}

export interface TerminalState {
  mode: number;
}

export class TerminalPage extends Component<TerminalProps, TerminalState> {
  /*constructor(props: TerminalProps) {
    super(props);
  }*/

  componentDidMount(): void {
    this.setState({ mode: this.props.mode });
  }

  componentWillUnmount(): void {}

  render() {
    return (
      <div>
        <Grid.Row>
          <Terminal
            color="green"
            backgroundColor="black"
            barColor="black"
            style={{ fontWeight: "bold", fontSize: "1em" }}
            msg="Welcome to the WebBridge terminal."
            startState="maximised"
            hideTopBar
          />
        </Grid.Row>
      </div>
    );
  }
}
