import React, { Component, FormEvent } from "react";
import { Box } from "../ui/Box";
import { ArrowButton, BackButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { H3, Text } from "../ui/Text";

export interface MnemonicWarningProps {
  onComplete: () => void;
  onCancel: () => void;
}

export interface MnemonicWarningState {}

export class MnemonicWarning extends Component<
  MnemonicWarningProps,
  MnemonicWarningState
> {
  constructor(props: MnemonicWarningProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount(): void {}

  componentDidUnmount(): void {}

  handleSubmit = (e: FormEvent) => {
    //if we don't prevent form submission, causes a browser reload
    e.preventDefault();
    this.props.onComplete();
  };

  render() {
    return (
      <>
        <Container height="50vh" margin="10% 0 0 0">
          <form onSubmit={this.handleSubmit}>
            <Box direction="column" align="center" width="100%">
              <Box
                direction="column"
                width="700px"
                align="start"
                margin="0 auto 0 auto"
              >
                <BackButton
                  onClick={() => this.props.onCancel()}
                  margin="130px 0 0 -100px"
                />
                <Card
                  width="100%"
                  align="center"
                  minHeight="225px"
                  padding="2em 8em 2em 8em"
                  border="solid 1px #e30429"
                >
                  <H3 color="#e30429">Mnemonic Warning !</H3>
                  <Text>
                    Please ensure that you are in a private location and no one
                    can eavesdrop or see your screen.
                  </Text>
                  <Text>
                    By clicking Continue, you will be shown your new wallet
                    mnemonic, which if comprised will grant anyone access to
                    your account.
                  </Text>
                  <Text>Write it down and keep it safe.</Text>
                </Card>
              </Box>
              <Box
                direction="column"
                width="690px"
                align="right"
                margin="0 auto 0 auto"
              >
                <ArrowButton focus label="Continue" type="submit" />
              </Box>
            </Box>
          </form>
        </Container>
      </>
    );
  }
}
