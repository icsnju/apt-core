'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider) {
    $stateProvider
      .state('table', {
        url: '/table',
        templateUrl: 'static/app/table/table.html',
      })
      .state('table.all', {
        url: '/all',
        templateUrl: 'static/app/table/table-all.html',
        controller: 'TableCtrl'
      })
      .state('table.tasks', {
        url: '/{jobid:int}',
        templateUrl: 'static/app/table/table-task.html',
        controller: 'TasksCtrl'
      });
  });
