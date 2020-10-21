export interface ValidationResult<T> {
    value: T;
    success: boolean;
    validationMessages: string[];
    isError: boolean
}

interface NameIndictator {
    name: string;
    scope: string;
}

export interface NameIndicatorWithValue<T> extends NameIndictator {
    value: T;
}

export type NamedValue<T> = T extends void ? NameIndictator : NameIndicatorWithValue<T>;

export const createValidatedFailurePayload = <T>(fieldScope: string, fieldName: string, message: string, fieldValue: T, isError = false): NamedValue<ValidationResult<T>> => ({
    scope: fieldScope,
    name: fieldName,
    value: {
        success: false,
        validationMessages: [message],
        value: fieldValue,
        isError: isError
    }
})
