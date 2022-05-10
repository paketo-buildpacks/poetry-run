# Poetry Run Cloud Native Buildpack
## `gcr.io/paketo-buildpacks/poetry-run`

The Paketo Poetry Run CNB sets the start command for a given [poetry](https://python-poetry.org/) application.

This buildpack detects when one of the following conditions is met:

1. ### `BP_POETRY_RUN_TARGET` is set
Example: `BP_POETRY_RUN_TARGET=default_app.server:run`.
The resulting start command for this example would be `poetry run default_app.server:run`.

1. ### `pyproject.toml` exists and contains **exactly one** poetry script
More specifically, the buildpack will detect if `pyproject.toml` looks like the following:

```
[tool.poetry.scripts]
some-script = "some.module:some_method"
```

The resulting start command for this example would be `poetry run some-script`.

See the [`poetry run` documentation](https://python-poetry.org/docs/cli/#run) for more information.

## Integration

This buildpacks writes a start command, so currently there is no conceivable
reason to require it as a dependency.

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh --version <version-number>
```

This will create a `buildpackage.cnb` file under the `build` directory which you
can use to build your app as follows:
```
pack build <app-name> -p <path-to-app> -b <path/to/cpython.cnb> -b <path/to/pip.cnb> -b <path/to/poetry.cnb> -b <path/to/poetry-install.cnb> -b build/buildpackage.cnb
```

### Configuration

#### Custom run command
This buildpack will set a start command that begins with `poetry run`. You can configure a custom command by using `BP_POETRY_RUN_TARGET`.
Example: If `BP_POETRY_RUN_TARGET=default_app.server:run`, this buildpack will set a start command of `poetry run default_app.server:run`.

#### Enabling reloadable process types
You can configure this buildpack to wrap the entrypoint process of your app such that it kills and restarts the process whenever files change in the app's working directory in the container. With this feature enabled, copying new versions of source code into the running container will trigger your app's process to restart. Set the environment variable `BP_LIVE_RELOAD_ENABLED=true` at build time to enable this feature.

## Run Tests

To run all unit tests, run:
```
./scripts/unit.sh
```

To run all integration tests, run:
```
/scripts/integration.sh
```

## Known issues and limitations

* When `BP_POETRY_RUN_TARGET` is not set, only one (and exactly one) script may be defined in the `pyproject.toml` file.
  Zero scripts, or multiple scripts, will result in the buildpack failing detection and therefore not participating in the order group.
