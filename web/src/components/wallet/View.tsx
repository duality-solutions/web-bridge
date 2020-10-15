import React, { Component } from "react";
import { RequestConfig } from "../../api/Config";
import { WalletAddressResponse } from "../../api/Wallet";
import { Box } from "../ui/Box";
import { Card } from "../ui/Card";
import { Container } from "../ui/Container";
import { QRCode } from "../ui/QRCode";
import { H1, H3 } from "../ui/Text";
import axios from 'axios';

export interface WalletViewProps {}

export interface WalletViewState {
  walletAddress?: string;
}

export class WalletView extends Component<WalletViewProps, WalletViewState> {
  constructor(props: WalletViewProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentWillUnmount = this.componentWillUnmount.bind(this);
    this.getWalletAddresses = this.getWalletAddresses.bind(this);
    this.getWalletAddresses();
  }

  componentDidMount(): void {}

  componentWillUnmount(): void {}

  private getWalletAddresses = () => {
    var self = this;
    axios.get<WalletAddressResponse>("/wallet/defaultaddress", RequestConfig).then(function (response) {
      self.setState( { walletAddress: response.data.address });
    }).catch(function (error) {
      console.log("Execute dynamic-cli JSON RCP [Get] Error: " + error);
    });
  };

  render() {
    return (
      <>
        <H1 align="center" color="black">
          Wallet Info
        </H1>
        <Container height="300vh">
        {this.state &&  (
          <Box direction="column" align="center" width="100%">
            <Box display="flex" direction="row" align="center" width="100%">
              <Card
                padding="2em 1em 1em 1em"
                min-height="140px"
                min-width="220px"
              >
                {this.state.walletAddress && this.state.walletAddress.length > 0 && (
                  <>
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
                      value={this.state.walletAddress}
                      size={240}
                    />
                  </>
                )}
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
                  <span>{this.state.walletAddress} </span>
                </div>
              </Card>
            </Box>
          </Box>
        )}
        </Container>
      </>
    );
  }
}
