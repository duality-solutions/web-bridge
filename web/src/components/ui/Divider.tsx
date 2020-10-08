import styled from "styled-components";

interface DividerProps {
  borderLeft?: string;
  background?: string;
  width?: string;
  height?: string;
  opacity?: string;
  margin?: string;
}

const StyledDivider = styled("div")<DividerProps>`
  border-left: ${(props) =>
    props.borderLeft ? props.borderLeft : "1px solid #38546d"};
  background: ${(props) => (props.background ? props.background : "#606060")};
  width: ${(props) => (props.width ? props.width : "1px")};
  height: ${(props) => (props.height ? props.height : "80px")};
  opacity: ${(props) => (props.opacity ? props.opacity : "0.2")};
  margin: ${(props) => (props.margin ? props.margin : "0 auto")};
`;

export { StyledDivider as Divider };
