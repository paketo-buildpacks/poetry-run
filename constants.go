package poetryrun

// PoetryVenv is the name of the dependency provided by the Poetry Install
// buildpack.
const PoetryVenv = "poetry-venv"

// CPython is the name of the python runtime dependency provided by the CPython
// buildpack: https://github.com/paketo-buildpacks/cpython.
const CPython = "cpython"

// Poetry is the name of the dependency provided by the Poetry buildpack:
// https://github.com/paketo-buildpacks/poetry.
const Poetry = "poetry"

// VenvLayerName is the name of the layer where the venv dependencies are
// installed to.
const VenvLayerName = "poetry-venv"

// CacheLayerName holds the poetry cache.
const CacheLayerName = "cache"
