import type { Record } from "./public/record";
import { Namespace$result, Namespace$input } from "../artifacts/Namespace";
import { NamespaceStore } from "../plugins/houdini-svelte/stores/Namespace";

export declare type CacheTypeDef = {
    types: {
        Namespace: {
            idFields: never;
            fields: {
                name: {
                    type: string;
                    args: never;
                };
            };
            fragments: [];
        };
        __ROOT__: {
            idFields: {};
            fields: {
                namespace: {
                    type: Record<CacheTypeDef, "Namespace">;
                    args: {
                        id: string;
                    };
                };
                namespaces: {
                    type: (Record<CacheTypeDef, "Namespace">)[] | null;
                    args: never;
                };
            };
            fragments: [];
        };
    };
    lists: {};
    queries: [[NamespaceStore, Namespace$result, Namespace$input]];
};