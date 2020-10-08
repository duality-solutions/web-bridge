import React, { Component, useState } from "react";
import { Box } from "./ui/Box";
import { BackButton } from "./ui/Button";
import { Card } from "./ui/Card";
import { Container } from "./ui/Container";
import { Dropzone, DropzoneError } from "./ui/Dropzone";
import { SecureFileIcon } from "./ui/Images";
import { H1, H3, Text } from "./ui/Text";
import { FilePathInfo } from "../shared/FilePathInfo";

export interface WalletFileRestoreProps {
  onComplete: () => void;
  onCancel: () => void;
}

export interface WalletFileRestoreState {
    stateError?: string
}

export class WalletFileRestore extends Component<WalletFileRestoreProps, WalletFileRestoreState> {
  private newWallet?: boolean;
  constructor(props: WalletFileRestoreProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
  }

  componentDidMount(): void {}

  componentWillUnmount(): void {}

  render() {
    const [error, setError] = useState<DropzoneError | undefined>(undefined);
    const filesSelectedHandler = (files: FilePathInfo[]) => {
        if (files.length !== 1) {
            setError({ title: "More that one file selected", message: "Please select only one file" });
            return;
        }
        const file: FilePathInfo = files[0];
        if (file.size > 131072) { //128KiB
            setError({ title: "File is too large", message: "Please select a mnemonic recovery file" });
            return;
        }
        //mnemonicRestoreFilePathSubmitted(file.path);
        //secureFilePassword();
        this.props.onComplete();
    }

    return (
      <>
        <H1 align="center" colored fontWeight="600">Restore Account</H1>
        <Container height="50vh" margin="10% 0 0 0">
            <form onSubmit={e => {
                e.preventDefault()
                e.stopPropagation()
                return false;
            }}>
                <Box direction="column" align="center" width="100%">
                    <Box direction="column" width="700px" align="start" margin="0 auto 0 auto">
                        <BackButton onClick={() => this.props.onCancel() } margin="150px 0 0 -80px" />
                        <Card width="100%" align="center" minHeight="225px" padding="2em 4em 2em 2em">
                            <Box display="flex" direction="row" margin="0">
                                <Box width="60px" margin="0">
                                    <SecureFileIcon width="60px" height="60px" />
                                </Box>
                                <Box direction="column" width="500px" align="center" margin="0 auto 0 auto">
                                    <H3 margin="0 0 1em 0">Restore using Secure Restore File </H3>
                                    <Dropzone multiple={false} accept={".psh.json"} filesSelected={filesSelectedHandler} error={error}></Dropzone>
                                    {this.state && this.state.stateError ? <Text align="center" color="#e30429">{this.state.stateError}</Text> : <></>}
                                </Box>
                            </Box>
                        </Card>
                    </Box>
                </Box>
            </form>
        </Container>
      </>
    );
  }
}
