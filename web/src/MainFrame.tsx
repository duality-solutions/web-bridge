import React from 'react';
import { Component } from "react";
import { Header, Icon, /*Image,*/ Menu, Segment, Sidebar } from 'semantic-ui-react';
//import logo from './logo.svg';

export interface MainFrameProps {
    currentPage?: string
};

export interface MainFrameState {
    setupComplete?: boolean
};

export class MainFrame extends Component<MainFrameProps, MainFrameState> {
    private currentPage: string;
    constructor(props: MainFrameProps) {
        super(props);
        this.currentPage = props.currentPage ? props.currentPage : "home";
    }

    componentDidMount(): void {
    }

    componentWillUnmount(): void {
    }
    
    changePage(): void {
        console.log("MainFrame.changePage" + this.currentPage);
    }

    render() {
        return (
            <div className="ui sidebar">
                {/*<img src={logo} className="App-logo" alt="logo" />*/}
                <Sidebar.Pushable as={Segment}>
                    <Sidebar
                        as={Menu}
                        icon='labeled'
                        inverted
                        vertical
                        visible
                    >
                        <Menu.Item as='a'>
                            <Icon name='home' />
                            Home
                        </Menu.Item>
                        <Menu.Item as='a'>
                            <Icon name='chain' />
                            Chain
                        </Menu.Item>
                        <Menu.Item as='a'>
                            <Icon name='user secret' />
                            Accounts
                        </Menu.Item>
                        <Menu.Item as='a'>
                            <Icon name='connectdevelop' />
                            Bridges
                        </Menu.Item>
                    </Sidebar>
                    <Sidebar.Pusher>
                        <Segment basic>
                            <Header as='h3'>Home</Header>
                        </Segment>
                    </Sidebar.Pusher>
                </Sidebar.Pushable>
            </div>
        )
    }
}