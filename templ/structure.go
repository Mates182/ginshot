package templates

func GetStructure() string {
	return `{
    "structure": [
        "cmd",
        {
            "config": ["cors"]
        },
        {
            "internal": [
                "controller",
                "service",
                "repository",
                "models",
                {
                    "data": ["requests", "responses"]
                },
                "event",
                "middleware"
            ]
        },
        "pkg",
        "router",
        "docs",
        "test",
        {
            "deployments": ["k8s"]
        }
    ]
}`
}
