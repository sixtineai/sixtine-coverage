function getEnv(name, defaultValue) {
    return window?.configs?.[name] || defaultValue
}

export default {
    coverageUrl: getEnv('SIXTINEAI_COVERAGE_URL', './sample-data/demo-cover.out'),
    sourceRoot: getEnv('SIXTINEAI_SOURCE_ROOT', './sample-data/demo-source/'),
    modulePrefix: getEnv('SIXTINEAI_MODULE_PREFIX', ''),
    routeCoverageUrl: getEnv('SIXTINEAI_ROUTE_COVERAGE_URL', './sample-data/demo-route-coverage.rcov'),

    coverageData: getEnv('SIXTINEAI_COVERAGE_DATA', ''),
    routeCoverageData: getEnv('SIXTINEAI_ROUTE_COVERAGE_DATA', ''),
    codeFilesMap: getEnv('SIXTINEAI_CODE_FILES_MAP', ''),

    thresholdGreen: 80,
    thresholdYellow: 50,
}
