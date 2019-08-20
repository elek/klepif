# Klepif

Klepif is a lightweight, simple github build robot. It's inspired by [prow](https://github.com/kubernetes/test-infra/tree/master/prow) which is used to manage kubernetes builds but more lightweight:

 1. It polls the github repository instead of listening to the event (no admin permission is required)
 2. It can trigger new build when required (build is defined as a CLI command)
 3. It can handle the `/label` commands (add label)


## Usage

Please see the `klepif.yaml.example` as an example file. `KLEPIF_GITHUB_TOKEN` can be set as an environment variable.

Most of the time you should use:

```
klepif check
```

It checks the latest changes (compared to a local cache which is saved to a dir).

`klepif run [PRNUM]` command can be used to debug the execution (it calls the build plugins even if the pr is not changed).



