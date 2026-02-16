function getEnv(name, defaultValue) {
    return window?.configs?.[name] || defaultValue
}

export default {
    coverageUrl: getEnv('SIXTINEAI_COVERAGE_URL', './cover.out'),
    sourceRoot: getEnv('SIXTINEAI_SOURCE_ROOT', ''),
    modulePrefix: getEnv('SIXTINEAI_MODULE_PREFIX', ''),
    routeCoverageUrl: getEnv('SIXTINEAI_ROUTE_COVERAGE_URL', ''),

    coverageData: getEnv('SIXTINEAI_COVERAGE_DATA', ''),
    routeCoverageData: getEnv('SIXTINEAI_ROUTE_COVERAGE_DATA', ''),
    codeFilesMap: getEnv('SIXTINEAI_CODE_FILES_MAP', ''),

    thresholdGreen: 80,
    thresholdYellow: 50,
}
