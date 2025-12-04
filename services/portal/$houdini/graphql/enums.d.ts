
type ValuesOf<T> = T[keyof T]
	
export declare const DedupeMatchMode: {
    readonly Variables: "Variables";
    readonly Operation: "Operation";
    readonly None: "None";
}

export type DedupeMatchMode$options = ValuesOf<typeof DedupeMatchMode>
 
export declare const Role: {
    readonly ANONYMOUS: "ANONYMOUS";
    readonly USER: "USER";
}

export type Role$options = ValuesOf<typeof Role>
 