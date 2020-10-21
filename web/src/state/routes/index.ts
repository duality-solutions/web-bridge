import React from "react";
import { push } from "connected-react-router";
import { RouteComponentProps } from "react-router";
import mananageWallet from "../store/manageWallet";

export interface RouteInfo {
  path: string;
  component:
    | React.ComponentType<RouteComponentProps<any>>
    | React.ComponentType<any>;
  exact: boolean;
}

const route = (
  path: string,
  component:
    | React.ComponentType<RouteComponentProps<any>>
    | React.ComponentType<any>,
  exact: boolean = true
): RouteInfo => ({ path, component, exact });

const routingTable = {
  mananageWallet: route("/mananageWallet", mananageWallet)
};

export const pushRoute = (route: RouteInfo) => push(route.path);

//deepFreeze(routingTable)

export const appRoutes = routingTable;
