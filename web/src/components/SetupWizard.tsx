import React, { Component } from "react";
import { Grid, Step, StepGroup } from "semantic-ui-react";
import { WalletSetup } from "./wallet/Setup";

export interface SetupWizardProps {
  currentStep: number;
  onComplete: () => void;
}

export interface SetupWizardState {
  currentStep?: number;
}

export class SetupWizard extends Component<SetupWizardProps, SetupWizardState> {
  constructor(props: SetupWizardProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
  }

  componentDidMount(): void {
    this.setState({ currentStep: this.props.currentStep });
  }

  componentDidUnmount(): void {}

  render() {
    return (
      <div className="ui stackable steps">
        <Grid>
          <Grid.Row>
            <StepGroup widths="seven">
              <Step
                icon="privacy"
                href="#"
                active={
                  this.state && this.state.currentStep === 1 ? true : false
                }
                disabled={
                  this.state && this.state.currentStep === 1 ? false : true
                }
              />
              <Step
                icon="user"
                href="#"
                active={
                  this.state && this.state.currentStep === 2 ? true : false
                }
                disabled={
                  this.state && this.state.currentStep === 2 ? false : true
                }
              />
              <Step
                icon="handshake"
                href="#"
                active={
                  this.state && this.state.currentStep === 3 ? true : false
                }
                disabled={
                  this.state && this.state.currentStep === 3 ? false : true
                }
              />
              <Step
                icon="checkmark"
                href="#"
                active={
                  this.state && this.state.currentStep === 4 ? true : false
                }
                disabled={
                  this.state && this.state.currentStep === 4 ? false : true
                }
              />
            </StepGroup>
          </Grid.Row>
          <Grid.Row textAlign="center">
            {this.state && this.state.currentStep === 1 && (
              <Grid.Column>
                <WalletSetup onComplete={() => this.props.onComplete()} />
              </Grid.Column>
            )}
          </Grid.Row>
        </Grid>
      </div>
    );
  }
}
