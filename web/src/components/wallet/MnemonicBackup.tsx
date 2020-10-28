import React, { Component } from "react";
import { Box } from "../ui/Box";
import { ArrowButton, BackButton, LightButton } from "../ui/Button";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { Divider } from "../ui/Divider";
import { PasswordEntry, SafeImage, SecureFileIcon } from "../ui/Images";
import { H3, Text } from "../ui/Text";
import { GetMnemonic } from "../../api/Wallet";

export interface MnemonicBackupProps {
  onCancel: () => void;
  onComplete: () => void;
  enterPassword: () => void;
  onBackupSecureFile: (mnemonic: string) => void;
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
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.backupSecureFile = this.backupSecureFile.bind(this);
    this.getMnemonic = this.getMnemonic.bind(this);
  }

  componentDidMount(): void {
    this.getMnemonic();
  }

  componentDidUnmount(): void {}

  backupSecureFile = async () => {
    if (this.state && this.state.mnemonic) {
      this.props.onBackupSecureFile(this.state.mnemonic);
    }
  };

  private getMnemonic = async () => {
    var self = this;
    GetMnemonic().then((data) => {
      if (data.hdseed !== "") {
        this.setState({ mnemonic: data.mnemonic });
      } else {
        console.log("wallet locked");
        self.props.enterPassword();
      }
    }).catch(function (error) {
      console.log("Rest error. wallet locked");
      self.props.enterPassword();
    });
  };

  render() {
    return (
      <>
        {this.state && this.state.mnemonic && (
          <Container height="50vh" margin="5% 0 0 0">
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
                        {this.state.mnemonic}
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
                          <LightButton onClick={this.backupSecureFile}>
                              Create a secure file
                          </LightButton>
                          <Box display="flex" width="100%" margin="2em 0 0 2em">
                              <SecureFileIcon
                              width="80px"
                              height="80px"
                              style={{ color: "blue" }}
                              />
                          </Box>
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
                    <ArrowButton label="Continue" type="button" onClick={() => this.props.onComplete()}/>
                </Box>
            </Box>
          </Container>
        )}
      </>
    );
  }
}
