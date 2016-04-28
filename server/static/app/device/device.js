'use strict';

angular.module('aptWebApp')
    .config(function($stateProvider) {
        $stateProvider
            .state('main.deviceDetail', {
                url: '/devices/:ip/:id',
                templateUrl: 'static/app/device/device.detail.html',
                controller: 'DeviceDetailCtrl'
            })
    });
