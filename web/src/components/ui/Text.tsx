import styled from "styled-components";

interface H1Props {
  align?: string;
  margin?: string;
  color?: string;
  colored?: boolean;
  fontWeight?: string;
  theme?: { blue: string };
  minwidth?: string;
}

interface ParaProps {
  align?: string;
  margin?: string;
  color?: string;
  colored?: boolean;
  theme?: { blue: string };
  fontSize?: string;
  notUserSelectable?: boolean;
  fontWeight?: string;
  disabled?: boolean;
  lineHeight?: string;
}

const StyledHeader = styled("h1")<H1Props>`
  text-align: ${(props) => props.align || "start"};
  letter-spacing: 0.03em;
  margin: ${(props) => props.margin || "0 0 0 0"};
  font-weight: ${(props) => props.fontWeight || "bold"};
  color: ${(props) => {
    if (props.colored) return props.theme.blue;
    else if (props.color) return props.color;
    else return "black";
  }};
  min-width: 500px;
`;

const StyledHeader3 = styled("h3")<H1Props>`
  text-align: ${(props) => props.align || "start"};
  letter-spacing: 0.03em;
  margin: ${(props) => props.margin || "0 0 0 0"};
  font-weight: ${(props) => props.fontWeight || "bold"};
  color: ${(props) => {
    if (props.colored) return props.theme.blue;
    else if (props.color) return props.color;
    else return "black";
  }};
  min-width: ${(props) => props.minwidth || "500px"};
  word-wrap: break-word;
`;

const StyledText = styled("p")<ParaProps>`
  user-select: ${(props) => (props.notUserSelectable ? "none" : "auto")};
  text-align: ${(props) => (props.align ? props.align : "center")};
  margin: ${(props) => (props.margin ? props.margin : "1em 0 0 0")};
  line-height: 1.4em;
  font-size: ${(props) => (props.fontSize ? props.fontSize : "9")};
  font-weight: ${(props) => (props.fontWeight ? props.fontWeight : "normal")};
  opacity: ${(props) => (props.disabled ? 0.4 : 1)};
  color: ${(props) => {
    if (props.colored) return props.theme.blue;
    else if (props.color) return props.color;
    else return "black";
  }};
  line-height: ${(props) => (props.lineHeight ? props.lineHeight : "")};
`;

export default StyledText;

export { StyledHeader as H1, StyledHeader3 as H3, StyledText as Text };
