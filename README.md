# go-sso-demo

A demo that shows how SSO works with WorkOS and Go.

## Install

```sh
go get -u github.com/workos-inc/go-sso-demo
```

## How to try

1. Configure a sso connection on [WorkOS SSO connections](https://dashboard.workos.com/sso/connections)
2. Setup a redirect uri on [WorkOS SSO configuration](https://dashboard.workos.com/sso/configuration)*
3. Launch the demo

    ```sh
    go-sso-demo \
        -api-key <workos_api_key> \
        -project-id <workos_project_id> \
        -redirect-uri <redirect_uri>
    ```

*(\*) For local environment, you will need to provide a redirect uri that points to your machine. You can use [Ngrok](https://ngrok.com) to achieve this.*
