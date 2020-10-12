import React, { Component, FormEvent } from "react";
import { Box } from "../ui/Box";
import { ArrowButton, BackButton, LightButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { Divider } from "../ui/Divider";
import { PasswordEntry, SafeImage, SecureFileIcon } from "../ui/Images";
import { LoadingSpinner } from "../ui/LoadingSpinner";
import { H3, Text } from "../ui/Text";

export interface MnemonicBackupProps {
  mnemonic?: string;
  onCancel: () => void;
}

export interface MnemonicBackupState {
  mnemonic?: string;
}

export class MnemonicBackup extends Component<
  MnemonicBackupProps,
  MnemonicBackupState
> {
  constructor(props: MnemonicBackupProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleFileCreation = this.handleFileCreation.bind(this);
  }

  componentDidMount(): void {
    this.setState({ mnemonic: this.props.mnemonic });
  }

  componentWillUnmount(): void {}

  handleSubmit = (e: FormEvent) => {
    //if we don't prevent form submission, causes a browser reload
    this.setState({ mnemonic: undefined });
    e.preventDefault();
  };

  handleFileCreation = (e: FormEvent) => {
    //if we don't prevent form submission, causes a browser reload
    e.preventDefault();
  };

  render() {
    return (
      <>
        {this.state && this.state.mnemonic && (
          <Container height="50vh" margin="5% 0 0 0">
            <form onSubmit={this.handleSubmit}>
              <Box direction="column" align="center" width="100%">
                <Box
                  direction="column"
                  width="700px"
                  align="start"
                  margin="0 auto 0 auto"
                >
                  <Card
                    width="100%"
                    align="center"
                    minHeight="225px"
                    padding="2em 4em 2em 4em"
                  >
                    <H3>Your Mnemonic Pass Phrase</H3>
                    <Card
                      width="100%"
                      align="center"
                      padding="1em"
                      border="solid 1px grey"
                      background="#fafafa"
                    >
                      <BackButton
                        onClick={() => this.props.onCancel()}
                        margin="130px 0 0 -100px"
                      />
                      <Text
                        color="grey"
                        align="center"
                        margin="0"
                        //notUserSelectable
                      >
                        {this.props.mnemonic}
                      </Text>
                    </Card>
                    <Box display="flex" direction="row">
                      <Box width="50%" margin="0 1em 0 0">
                        <Text align="center">
                          Write or print this phrase and
                        </Text>
                        <Text margin="0" align="center">
                          keep it somewhere safe.
                        </Text>
                        <Box display="flex" width="100%" margin="2em 0 0 0">
                          <PasswordEntry width="80px" height="80px" />
                          <span
                            style={{
                              color: "#2e77d0",
                              lineHeight: "1.2em",
                              fontSize: "300%"
                            }}
                          >
                            &#8594;
                          </span>
                          <SafeImage width="80px" height="80px" />
                        </Box>
                      </Box>
                      <Box
                        width="14px"
                        direction="column"
                        alignContents="center"
                      >
                        <Divider />
                        <Text margin="0">or</Text>
                        <Divider />
                      </Box>
                      <Box width="30%" margin="0 0 0 3em">
                        <Box display="flex" width="100%" margin="2em 0 0 2em">
                          <SecureFileIcon
                            width="80px"
                            height="80px"
                            style={{ color: "blue" }}
                          />
                        </Box>
                        <LightButton onClick={this.handleFileCreation}>
                          Create a secure file
                        </LightButton>
                      </Box>
                    </Box>
                  </Card>
                </Box>
                <Box
                  direction="column"
                  width="700px"
                  align="right"
                  margin="0 auto 0 auto"
                >
                  <ArrowButton label="Continue" type="submit" />
                </Box>
              </Box>
            </form>
          </Container>
        )}
        {this.state && !this.state.mnemonic && (
          <LoadingSpinner
            active={typeof this.props.mnemonic === "undefined"}
            label="Generating your mnemonic passphrase"
            size={50}
            opaque
          />
        )}
      </>
    );
  }
}