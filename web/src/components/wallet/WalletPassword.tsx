import React, { ChangeEvent, Component } from "react";
import { Container } from "../ui/Container";
import { Card } from "../ui/Card";
import { Box } from "../ui/Box";
import { ArrowButton, BackButton } from "../ui/Button";
import { Input } from "../ui/Input";
import { LoadingSpinner } from "../ui/LoadingSpinner";
import { Text } from "../ui/Text";

type PasswordUiType = "LOGIN" | "CREATE";

export interface ValidationResult<T> {
  value: T;
  success: boolean;
  validationMessages: string[];
  isError: boolean;
}

export interface WalletPasswordProps {
  onComplete:  (password: string) => void;
  onCancel: () => void;
  password?: string;
  uiType: PasswordUiType;
  errorMessage?: string;
}

export interface WalletPasswordState {
  password?: string;
  confirmPassword?: string;
  isValidating: boolean;
  validationResult?: ValidationResult<string>;
}

export class WalletPassword extends Component<
  WalletPasswordProps,
  WalletPasswordState
> {
  constructor(props: WalletPasswordProps) {
    super(props);
    // bind events
    this.componentDidMount = this.componentDidMount.bind(this);
    this.componentDidUnmount = this.componentDidUnmount.bind(this);
    this.onCompleteButtonClick = this.onCompleteButtonClick.bind(this);
    this.state = {
      isValidating: false
    };
  }

  componentDidMount(): void {}

  componentDidUnmount(): void {}

  handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    let name: string = e.target.name;
    let value: string = e.target.value;
    if (name === "password") {
      this.setState((state) => ({ ...state, password: value }));
    } else if (name === "confirmPassword") {
      this.setState((state) => ({ ...state, confirmPassword: value }));
    }
    //resetValidationForField({
    //  scope: validationScopes.password,
    // name: "password"
    //});
  };

  private onCompleteButtonClick = () => {
    if (this.state && this.state.password) {
      this.props.onComplete(this.state.password)
    }
  }

  render() {
    const { isValidating, validationResult } = this.state;
    const { uiType, errorMessage } = this.props;
    const validationFailed =
      typeof validationResult !== "undefined" && !validationResult.success;
    const showFieldErrors =
      validationFailed &&
      typeof validationResult !== "undefined" &&
      !validationResult.isError;
    return (
      <>
        <Container height="50vh" margin="10% 0 0 0">
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
                padding="2em 12em 2em 8em"
              >
                <Text fontSize="14px">
                  {uiType === "LOGIN" ? "Enter password" : "Create a Password"}
                </Text>
                <Input
                  value={this.state.password}
                  name="password"
                  onChange={this.handleChange}
                  placeholder="Password"
                  type="password"
                  margin="1em 0 1em 0"
                  padding="0 1em 0 1em"
                  autoFocus={true}
                  error={showFieldErrors}
                  disabled={isValidating}
                />
                {uiType === "LOGIN" ? (
                  <></>
                ) : (
                  <>
                    <Text fontSize="14px">Confirm Password</Text>
                    <Input
                      value={this.state.confirmPassword}
                      name="confirmPassword"
                      onChange={this.handleChange}
                      placeholder="Password"
                      type="password"
                      margin="1em 0 1em 0"
                      padding="0 1em 0 1em"
                      error={showFieldErrors}
                      disabled={isValidating}
                    />
                  </>
                )}
                {validationFailed ? (
                  (typeof validationResult !== "undefined"
                    ? validationResult.validationMessages
                    : []
                  ).map((e, i) => (
                    <Text align="center" color="#e30429" key={i}>
                      {e}
                    </Text>
                  ))
                ) : (
                  <></>
                )}
                {errorMessage ? (
                  <Text align="center" color="#e30429">
                    {errorMessage}
                  </Text>
                ) : (
                  <></>
                )}
              </Card>
            </Box>
            <Box
              direction="column"
              width="700px"
              align="right"
              margin="0 auto 0 auto"
            >
              <ArrowButton
                label={uiType === "CREATE" ? "Continue" : "Log in"}
                type="submit"
                disabled={isValidating}
                onClick={() => this.onCompleteButtonClick()}
              />
              {isValidating ? (
                <LoadingSpinner
                  active
                  label={
                    uiType === "LOGIN"
                      ? "Logging in..."
                      : "Encrypting your data..."
                  }
                  size={50}
                />
              ) : (
                <></>
              )}
            </Box>
          </Box>
        </Container>
      </>
    );
  }
}
