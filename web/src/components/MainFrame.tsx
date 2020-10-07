import React, { Component } from "react";
import { Settings } from "./Settings";
import { TerminalPage } from "./TerminalPage";
import {
  Button,
  Dropdown,
  Grid,
  Header,
  Icon,
  Image,
  Menu,
  Modal,
  Segment
} from "semantic-ui-react";
// TODO: Add logo to page
// Add subpages that display
//import logo from './logo.svg';

export interface MainFrameProps {
  currentPage?: string;
}

export interface MainFrameState {
  setupComplete?: boolean;
  currentPage?: string;
  open?: boolean;
}

export class MainFrame extends Component<MainFrameProps, MainFrameState> {
  private currentPage: string;
  constructor(props: MainFrameProps) {
    super(props);
    this.currentPage = props.currentPage ? props.currentPage : "home";
    this.state = {
      setupComplete: false,
      currentPage: this.currentPage,
      open: true
    };
  }

  componentDidMount(): void {}

  componentWillUnmount(): void {}

  private changePage(pageName: string): void {
    this.currentPage = pageName;
    console.log("MainFrame.changePage " + this.currentPage);
    this.setState({ currentPage: pageName, open: true });
  }

  render() {
    return (
      <div>
        <Grid>
          <Grid.Column width={2}>
            <div className="menu">
              <div className="toc">
                <Menu className="inverted vertical thin left fixed">
                  <Menu.Item onClick={() => this.changePage("home")} as="a">
                    <Icon name="home" />
                    Home
                  </Menu.Item>
                  <Menu.Item onClick={() => this.changePage("chain")} as="a">
                    <Icon name="chain" />
                    Chain
                  </Menu.Item>
                  <Menu.Item onClick={() => this.changePage("accounts")} as="a">
                    <Icon name="user secret" />
                    Accounts
                  </Menu.Item>
                  <Menu.Item onClick={() => this.changePage("bridges")} as="a">
                    <Icon name="connectdevelop" />
                    Bridges
                  </Menu.Item>
                  <Menu.Item onClick={() => this.changePage("terminal")} as="a">
                    <Icon name="terminal" />
                    Terminal
                  </Menu.Item>
                  <Dropdown item text="More">
                    <Dropdown.Menu>
                      <Dropdown.Item
                        onClick={() => this.changePage("profile")}
                        icon="edit"
                        text="Edit Profile"
                      />
                      <Dropdown.Item
                        onClick={() => this.changePage("language")}
                        icon="globe"
                        text="Choose Language"
                      />
                      <Dropdown.Item
                        onClick={() => this.changePage("settings")}
                        icon="settings"
                        text="Settings"
                      />
                    </Dropdown.Menu>
                  </Dropdown>
                </Menu>
              </div>
            </div>
          </Grid.Column>
          <Grid.Column stretched width={12}>
            <div>
              <div className="article">
                {this.currentPage === "home" && (
                  <Segment basic raised textAlign="center">
                    <Header as="h3">Home</Header>
                    <Image src="https://react.semantic-ui.com/images/wireframe/paragraph.png" />
                  </Segment>
                )}
                {this.currentPage === "chain" && (
                  <Segment basic raised textAlign="center">
                    <Modal
                      onClose={() => this.setState({ open: false })}
                      onOpen={() => this.setState({ open: true })}
                      open={this.state.open}
                      trigger={<Button>Show Modal</Button>}
                    >
                      <Modal.Header>Select a Photo</Modal.Header>
                      <Modal.Content image>
                        <Image
                          size="medium"
                          src="https://react.semantic-ui.com/images/wireframe/paragraph.png"
                          wrapped
                        />
                        <Modal.Description>
                          <Header>Default Profile Image</Header>
                          <p>
                            We've found the following gravatar image associated
                            with your e-mail address.
                          </p>
                          <p>Is it okay to use this photo?</p>
                        </Modal.Description>
                      </Modal.Content>
                      <Modal.Actions>
                        <Button
                          color="black"
                          onClick={() => this.setState({ open: false })}
                        >
                          Nope
                        </Button>
                        <Button
                          content="Yep, that's me"
                          labelPosition="right"
                          icon="checkmark"
                          onClick={() => this.setState({ open: false })}
                          positive
                        />
                      </Modal.Actions>
                    </Modal>
                  </Segment>
                )}
                {this.currentPage === "accounts" && (
                  <Segment basic raised textAlign="center">
                    <Header as="h3">Accounts</Header>
                    <Image src="https://react.semantic-ui.com/images/wireframe/paragraph.png" />
                  </Segment>
                )}
                {this.currentPage === "bridges" && (
                  <Segment basic raised textAlign="center">
                    <Header as="h3">Bridges</Header>
                    <Image src="https://react.semantic-ui.com/images/wireframe/paragraph.png" />
                  </Segment>
                )}
                {this.currentPage === "terminal" && (
                  <Segment basic raised textAlign="center">
                    <Header as="h3">Terminal</Header>
                    <TerminalPage mode={0} />
                  </Segment>
                )}
                {this.currentPage === "settings" && (
                  <Segment basic raised textAlign="center">
                    <Settings
                      defaultIceUrl="turn:ice.bdap.io:3478"
                      defaultIceUser="test"
                      defaultIcePassword="Admin@123"
                      defaultBind="0.0.0.0"
                      defaultAllow="127.0.0.1/0"
                      defaultPort={35350}
                    />
                  </Segment>
                )}
              </div>
            </div>
          </Grid.Column>
        </Grid>
      </div>
    );
  }
}
