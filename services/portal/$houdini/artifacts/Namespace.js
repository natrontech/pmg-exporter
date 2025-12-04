export default {
    "name": "Namespace",
    "kind": "HoudiniQuery",
    "hash": "8469b126ffa0f00b38aa6b5f01f7ae2aac1d4ad4067c24b6c11e02b3eea1f4e6",

    "raw": `query Namespace {
  namespaces {
    name
  }
}
`,

    "rootType": "Query",
    "stripVariables": [],

    "selection": {
        "fields": {
            "namespaces": {
                "type": "Namespace",
                "keyRaw": "namespaces",
                "nullable": true,

                "selection": {
                    "fields": {
                        "name": {
                            "type": "String",
                            "keyRaw": "name",
                            "visible": true
                        }
                    }
                },

                "visible": true
            }
        }
    },

    "pluginData": {
        "houdini-svelte": {}
    },

    "policy": "CacheOrNetwork",
    "partial": false
};

"HoudiniHash=8e8e893b95d1a9f2085a8da3428352f501cc3fab1992d017e8086fb9f118b2f3";