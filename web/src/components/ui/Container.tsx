import styled from "styled-components";

interface ContainerProps {
  height?: string;
  margin?: string;
  minWidth?: string;
  padding?: string;
}

export const Container = styled("div")<ContainerProps>`
  height: ${(props) => (props.height ? props.height : "90vh")};
  margin: ${(props) => props.margin || "5% 0 0 0"};
  min-width: ${(props) => props.minWidth || "0"};
  padding: ${(props) => props.padding || "0"};
`;

export default Container;
