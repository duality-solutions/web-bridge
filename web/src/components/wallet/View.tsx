import React, { Component } from "react";
import { Box } from "../ui/Box";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { QRCode } from "../ui/QRCode";
import { H1, H3 } from "../ui/Text";

export interface WalletViewProps {}

export interface WalletViewState {}

export class WalletView extends Component<WalletViewProps, WalletViewState> {
  private walletStealthAddress: string =
    "L3D2VnognoqcZAFTtYVhFN4UNjyspBKbn1GxEEtr25uhynPGk6FWoGG5zhkBbhWnbuzvKpArwawqadGJeEzb83P7da9hudK2gvYZbw";
  private walletAddress: string = "DFbNCTtvWx8eohBEoVJEFbBR47BFkqYSJz";
  constructor(props: WalletViewProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
  }

  componentDidMount(): void {}

  componentWillUnmount(): void {}

  render() {
    return (
      <>
        <H1 align="center" color="black">
          Wallet Info
        </H1>
        <Container height="300vh">
          <Box direction="column" align="center" width="100%">
            <Box display="flex" direction="row" align="center" width="100%">
              <Card
                //onClick={() => this.props.restoreUsingMnemonic()}
                padding="2em 1em 1em 1em"
                min-height="140px"
                min-width="220px"
              >
                <H3 align="center" color="black">
                  Stealth Address
                </H3>
                <QRCode
                  bgColor="#00000000"
                  fgColor="#4a4a4aff"
                  ecLevel="M"
                  minPadding={30}
                  minimumCellSize={4}
                  qrStyle="dots"
                  value={this.walletStealthAddress}
                  size={205}
                />
                <div
                  style={{
                    wordBreak: "break-all",
                    maxWidth: "500px",
                    fontFamily: '"Courier New", Courier, monospace',
                    textAlign: "center",
                    color: "#4a4a4a",
                    margin: "0 10px 10px",
                    position: "relative"
                  }}
                >
                  <span>{this.walletStealthAddress} </span>
                </div>
              </Card>
              <Card
                //onClick={() => this.props.restoreUsingMnemonic()}
                padding="2em 1em 1em 1em"
                min-height="140px"
                min-width="220px"
              >
                <H3 align="center" color="black">
                  Wallet Address
                </H3>
                <QRCode
                  bgColor="#00000000"
                  fgColor="#4a4a4aff"
                  ecLevel="M"
                  minPadding={30}
                  minimumCellSize={4}
                  qrStyle="dots"
                  value={this.walletAddress}
                  size={240}
                />
                <div
                  style={{
                    wordBreak: "break-all",
                    maxWidth: "500px",
                    fontFamily: '"Courier New", Courier, monospace',
                    textAlign: "center",
                    color: "#4a4a4a",
                    margin: "0 10px 10px",
                    position: "relative"
                  }}
                >
                  <span>{this.walletAddress} </span>
                </div>
              </Card>
            </Box>
          </Box>
        </Container>
      </>
    );
  }
}
