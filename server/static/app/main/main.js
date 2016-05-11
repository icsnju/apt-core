'use strict';

var translationsEN = {
    NAV_DASHBOARD: 'Dashboard',
    NAV_JOBS: 'Jobs',
    NAV_ADD: 'Add',

    DEVICE_ID: 'Device ID',
    DEVICE_MANU: 'Manufacturers',
    DEVICE_MODEL: 'Model',
    DEVICE_API: 'API',
    DEVICE_BV: 'Build Version',
    DEVICE_CPU: 'CPU Abi',
    DEVICE_RESOLUTION: 'Resolution',
    DEVICE_STATUS: 'Status',

    NODE: 'Node',
    NODES: 'Nodes',
    DEVICES: 'Devices',
    JOBS: 'Jobs',

    SLAVE_IP: 'IP',
    SLAVE_DEVICES: 'Devices',
    SLAVE_TASKS: 'Tasks',

    JOB_ID: 'Job ID',
    JOB_START_TIME: 'Start Time',
    JOB_FINISH_TIME: 'Finish Time',
    JOB_STATUS: 'Status',
    JOB_DETAIL: 'Job Detail',
    JOB_FRAME: 'Framework',
    JOB_SELECT: 'Device Select',
    JOB_DEVICE: 'Device Number',
    JOB_OPTIONS: 'Options',
    TASK_STATUS: 'Status',
    TASK_RESULT: 'Result',
    ADD_JOB: 'Add Job',

    FRAMEWORK_SELECT: 'Test Framework Select',
    APP_INPUT: 'Application Input',
    FILE: 'File',
    PACKAGE_NAME: 'Package Name',
    ARG: 'Arguments',
    TEST_INPUT: 'Test File Input',
    NEXT_STEP: 'Next Step',
    ALERT_FRAME: 'Please fill out the testing framework.',
    ALERT_DEVICE: 'Please select some devices.',
    SUBMIT: 'Submit'
};

var translationsCN = {
    NAV_DASHBOARD: '节点列表',
    NAV_JOBS: '测试列表',
    NAV_ADD: '添加测试',

    DEVICE_ID: '设备ID',
    DEVICE_MANU: '生产商',
    DEVICE_MODEL: '型号',
    DEVICE_API: 'API',
    DEVICE_BV: '版本',
    DEVICE_CPU: 'CPU型号',
    DEVICE_RESOLUTION: '分辨率',
    DEVICE_STATUS: '状态',

    NODE: '节点',
    NODES: '节点列表',
    DEVICES: '设备列表',
    JOBS: '测试列表',

    SLAVE_IP: 'IP',
    SLAVE_DEVICES: '设备数量',
    SLAVE_TASKS: '任务数量',

    JOB_ID: '测试 ID',
    JOB_START_TIME: '开始时间',
    JOB_FINISH_TIME: '结束时间',
    JOB_STATUS: '状态',
    JOB_DETAIL: '测试详情',
    JOB_FRAME: '测试框架',
    JOB_SELECT: '设备选择',
    JOB_DEVICE: '设备数量',
    JOB_OPTIONS: '操作',
    TASK_STATUS: '状态',
    TASK_RESULT: '结果',
    ADD_JOB: '添加测试',

    FRAMEWORK_SELECT: '测试框架选择',
    APP_INPUT: '被测应用选择',
    FILE: '文件',
    PACKAGE_NAME: '应用包名',
    ARG: '测试参数',
    TEST_INPUT: '测试脚本选择',
    NEXT_STEP: '下一步',
    ALERT_FRAME: '请填写完整的测试框架信息.',
    ALERT_DEVICE: '请选择设备.',
    SUBMIT: '提交'
};

angular.module('aptWebApp')
    .config(function($stateProvider, $translateProvider) {
        $stateProvider
            .state('main', {
                templateUrl: 'static/app/main/main.html',
                controller: 'MainCtrl'
            });
        $translateProvider.translations('en', translationsEN);
        $translateProvider.translations('cn', translationsCN);
        $translateProvider.preferredLanguage('cn');
        $translateProvider.useSanitizeValueStrategy('escape');
    });
