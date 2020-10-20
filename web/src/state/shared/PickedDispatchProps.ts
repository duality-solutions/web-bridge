export type PickedDispatchProps<T, K extends keyof T=keyof T> = {
    [P in K]: T[P] extends (...args: any[]) => any ? (...args: Parameters<T[P]>) => void : never;
};
