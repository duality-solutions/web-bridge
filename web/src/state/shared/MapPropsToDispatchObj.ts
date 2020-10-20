export type MapPropsToDispatchObj<T> = {
    [P in keyof T]:T[P]
}
