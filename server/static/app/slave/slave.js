'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider) {
    $stateProvider
      .state('main.slaves', {
        url: '/slaves',
        templateUrl: 'static/app/slave/slave.html',
        controller: 'SlaveCtrl'
      })
      .state('main.slaveDetail', {
        url: '/slaves/:id',
        templateUrl: 'static/app/slave/slave.detail.html',
        controller: 'SlaveDetailCtrl'
      });
  });
